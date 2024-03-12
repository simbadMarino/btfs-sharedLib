package coreapi

import (
	"context"
	"fmt"

	coreiface "github.com/bittorrent/interface-go-btfs-core"
	caopts "github.com/bittorrent/interface-go-btfs-core/options"
	"github.com/bittorrent/interface-go-btfs-core/path"
	pin "github.com/ipfs/boxo/pinning/pinner"

	bserv "github.com/ipfs/boxo/blockservice"
	offline "github.com/ipfs/boxo/exchange/offline"
	"github.com/ipfs/boxo/ipld/merkledag"
	"github.com/ipfs/go-cid"
)

type PinAPI CoreAPI

func (api *PinAPI) Add(ctx context.Context, p path.Path, opts ...caopts.PinAddOption) error {
	dagNode, err := api.core().ResolveNode(ctx, p)
	if err != nil {
		return fmt.Errorf("pin: %s", err)
	}

	settings, err := caopts.PinAddOptions(opts...)
	if err != nil {
		return err
	}

	defer api.blockstore.PinLock(ctx).Unlock(ctx)

	err = api.pinning.Pin(ctx, dagNode, settings.Recursive, "")
	if err != nil {
		return fmt.Errorf("pin: %s", err)
	}

	if err := api.provider.Provide(dagNode.Cid()); err != nil {
		return err
	}

	return api.pinning.Flush(ctx)
}

func (api *PinAPI) Ls(ctx context.Context, opts ...caopts.PinLsOption) (<-chan coreiface.Pin, error) {
	settings, err := caopts.PinLsOptions(opts...)
	if err != nil {
		return nil, err
	}

	switch settings.Type {
	case "all", "direct", "indirect", "recursive":
	default:
		return nil, fmt.Errorf("invalid type '%s', must be one of {direct, indirect, recursive, all}", settings.Type)
	}

	return api.pinLsAll(ctx, settings.Type, settings.Detailed), nil
}

func (api *PinAPI) IsPinned(ctx context.Context, p path.Path, opts ...caopts.PinIsPinnedOption) (string, bool, error) {
	resolved, err := api.core().ResolvePath(ctx, p)
	if err != nil {
		return "", false, fmt.Errorf("error resolving path: %s", err)
	}

	settings, err := caopts.PinIsPinnedOptions(opts...)
	if err != nil {
		return "", false, err
	}

	mode, ok := pin.StringToMode(settings.WithType)
	if !ok {
		return "", false, fmt.Errorf("invalid type '%s', must be one of {direct, indirect, recursive, all}", settings.WithType)
	}

	return api.pinning.IsPinnedWithType(ctx, resolved.Cid(), mode)
}

// Rm pin rm api
func (api *PinAPI) Rm(ctx context.Context, p path.Path, opts ...caopts.PinRmOption) error {
	rp, err := api.core().ResolvePath(ctx, p)
	if err != nil {
		return err
	}

	settings, err := caopts.PinRmOptions(opts...)
	if err != nil {
		return err
	}

	// Note: after unpin the pin sets are flushed to the blockstore, so we need
	// to take a lock to prevent a concurrent garbage collection
	defer api.blockstore.PinLock(ctx).Unlock(ctx)

	if err = api.pinning.Unpin(ctx, rp.Cid(), settings.Recursive); err != nil {
		return err
	}

	return api.pinning.Flush(ctx)
}

func (api *PinAPI) Update(ctx context.Context, from path.Path, to path.Path, opts ...caopts.PinUpdateOption) error {
	settings, err := caopts.PinUpdateOptions(opts...)
	if err != nil {
		return err
	}

	fp, err := api.core().ResolvePath(ctx, from)
	if err != nil {
		return err
	}

	tp, err := api.core().ResolvePath(ctx, to)
	if err != nil {
		return err
	}

	defer api.blockstore.PinLock(ctx).Unlock(ctx)

	err = api.pinning.Update(ctx, fp.Cid(), tp.Cid(), settings.Unpin)
	if err != nil {
		return err
	}

	return api.pinning.Flush(ctx)
}

type pinStatus struct {
	err      error
	cid      cid.Cid
	ok       bool
	badNodes []coreiface.BadPinNode
}

// BadNode is used in PinVerifyRes
type badNode struct {
	path path.Resolved
	err  error
}

func (s *pinStatus) Ok() bool {
	return s.ok
}

func (s *pinStatus) BadNodes() []coreiface.BadPinNode {
	return s.badNodes
}

func (n *badNode) Path() path.Resolved {
	return n.path
}

func (n *badNode) Err() error {
	return n.err
}

func (api *PinAPI) Verify(ctx context.Context) (<-chan coreiface.PinStatus, error) {
	visited := make(map[cid.Cid]*pinStatus)
	bs := api.blockstore
	DAG := merkledag.NewDAGService(bserv.New(bs, offline.Exchange(bs)))
	getLinks := merkledag.GetLinksWithDAG(DAG)
	//recPins := api.pinning.RecursiveKeys(ctx, false)

	var checkPin func(root cid.Cid) *pinStatus
	checkPin = func(root cid.Cid) *pinStatus {
		if status, ok := visited[root]; ok {
			return status
		}

		links, err := getLinks(ctx, root)
		if err != nil {
			status := &pinStatus{ok: false, cid: root}
			status.badNodes = []coreiface.BadPinNode{&badNode{path: path.IpldPath(root), err: err}}
			visited[root] = status
			return status
		}

		status := &pinStatus{ok: true, cid: root}
		for _, lnk := range links {
			res := checkPin(lnk.Cid)
			if !res.ok {
				status.ok = false
				status.badNodes = append(status.badNodes, res.badNodes...)
			}
		}

		visited[root] = status
		return status
	}

	out := make(chan coreiface.PinStatus)
	go func() {
		defer close(out)
		for p := range api.pinning.RecursiveKeys(ctx, false) {
			var res *pinStatus
			if p.Err != nil {
				res = &pinStatus{err: p.Err}
			} else {
				res = checkPin(p.Pin.Key)
			}
			select {
			case <-ctx.Done():
				return
			case out <- res:
			}
		}
	}()

	return out, nil
}

type pinInfo struct {
	pinType string
	path    path.Resolved
	name    string
	err     error
}

func (p *pinInfo) Path() path.Resolved {
	return p.path
}

func (p *pinInfo) Type() string {
	return p.pinType
}

func (p *pinInfo) Err() error {
	return p.err
}

// pinLsAll is an internal function for returning a list of pins
func (api *PinAPI) pinLsAll(ctx context.Context, typeStr string, detailed bool) <-chan coreiface.Pin {
	out := make(chan coreiface.Pin)
	emittedSet := cid.NewSet()

	AddToResultKeys := func(c cid.Cid, name, typeStr string) error {
		if emittedSet.Visit(c) {
			select {
			case out <- &pinInfo{
				pinType: typeStr,
				name:    name,
				path:    path.IpldPath(c),
			}:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		return nil
	}

	go func() {
		defer close(out)

		var rkeys []cid.Cid
		var err error
		if typeStr == "recursive" || typeStr == "all" {
			for streamedCid := range api.pinning.RecursiveKeys(ctx, detailed) {
				if streamedCid.Err != nil {
					out <- &pinInfo{err: streamedCid.Err}
					return
				}
				if err = AddToResultKeys(streamedCid.Pin.Key, streamedCid.Pin.Name, "recursive"); err != nil {
					out <- &pinInfo{err: err}
					return
				}
				rkeys = append(rkeys, streamedCid.Pin.Key)
			}
		}
		if typeStr == "direct" || typeStr == "all" {
			for streamedCid := range api.pinning.DirectKeys(ctx, detailed) {
				if streamedCid.Err != nil {
					out <- &pinInfo{err: streamedCid.Err}
					return
				}
				if err = AddToResultKeys(streamedCid.Pin.Key, streamedCid.Pin.Name, "direct"); err != nil {
					out <- &pinInfo{err: err}
					return
				}
			}
		}
		if typeStr == "indirect" {
			// We need to first visit the direct pins that have priority
			// without emitting them

			for streamedCid := range api.pinning.DirectKeys(ctx, detailed) {
				if streamedCid.Err != nil {
					out <- &pinInfo{err: streamedCid.Err}
					return
				}
				emittedSet.Add(streamedCid.Pin.Key)
			}

			for streamedCid := range api.pinning.RecursiveKeys(ctx, detailed) {
				if streamedCid.Err != nil {
					out <- &pinInfo{err: streamedCid.Err}
					return
				}
				emittedSet.Add(streamedCid.Pin.Key)
				rkeys = append(rkeys, streamedCid.Pin.Key)
			}
		}
		if typeStr == "indirect" || typeStr == "all" {
			walkingSet := cid.NewSet()
			for _, k := range rkeys {
				err = merkledag.Walk(
					ctx, merkledag.GetLinksWithDAG(api.dag), k,
					func(c cid.Cid) bool {
						if !walkingSet.Visit(c) {
							return false
						}
						if emittedSet.Has(c) {
							return true // skipped
						}
						err := AddToResultKeys(c, "", "indirect")
						if err != nil {
							out <- &pinInfo{err: err}
							return false
						}
						return true
					},
					merkledag.SkipRoot(), merkledag.Concurrent(),
				)
				if err != nil {
					out <- &pinInfo{err: err}
					return
				}
			}
		}
	}()

	return out
}

func (api *PinAPI) core() coreiface.CoreAPI {
	return (*CoreAPI)(api)
}
