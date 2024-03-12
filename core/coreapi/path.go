package coreapi

import (
	"context"
	"fmt"
	gopath "path"

	"github.com/bittorrent/go-btfs/namesys/resolve"
	coreiface "github.com/bittorrent/interface-go-btfs-core"
	path "github.com/bittorrent/interface-go-btfs-core/path"
	"github.com/ipfs/boxo/fetcher"
	ipfspath "github.com/ipfs/boxo/path"
	ipfspathresolver "github.com/ipfs/boxo/path/resolver"
	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
)

// ResolveNode resolves the path `p` using Unixfs resolver, gets and returns the
// resolved Node.
func (api *CoreAPI) ResolveNode(ctx context.Context, p path.Path) (ipld.Node, error) {
	rp, err := api.ResolvePath(ctx, p)
	if err != nil {
		return nil, err
	}

	node, err := api.dag.Get(ctx, rp.Cid())
	if err != nil {
		return nil, err
	}
	return node, nil
}

// ResolvePath resolves the path `p` using Unixfs resolver, returns the
// resolved path.
func (api *CoreAPI) ResolvePath(ctx context.Context, p path.Path) (path.Resolved, error) {
	if _, ok := p.(path.Resolved); ok {
		return p.(path.Resolved), nil
	}
	ipath, _ := ipfspath.NewPath(p.String())
	ipath, err := resolve.ResolveIPNS(ctx, api.namesys, ipath)
	if err == resolve.ErrNoNamesys {
		return nil, coreiface.ErrOffline
	} else if err != nil {
		return nil, err
	}

	var dataFetcher fetcher.Factory
	switch ipath.Segments()[0] {
	case "btfs":
		dataFetcher = api.unixFSFetcherFactory
	case "ipld":
		dataFetcher = api.ipldFetcherFactory
	default:
		return nil, fmt.Errorf("unsupported path namespace: %s", p.Namespace())
	}

	r := ipfspathresolver.NewBasicResolver(dataFetcher)
	ip, _ := ipfspath.NewImmutablePath(ipath)
	node, rest, err := r.ResolveToLastNode(ctx, ip)
	if err != nil {
		return nil, err
	}

	root, err := cid.Parse(ipath.Segments()[1])
	if err != nil {
		return nil, err
	}

	return path.NewResolvedPath(ipath, node, root, gopath.Join(rest...)), nil
}

// ResolveIpnsPath resolves only the IPNS path `p` using the Unixfs resolver and returns the
// resolved path.
func (api *CoreAPI) ResolveIpnsPath(ctx context.Context, p path.Path) (*ipfspath.Path, error) {
	if err := p.IsValid(); err != nil {
		return nil, err
	}

	ipath, _ := ipfspath.NewPath(p.String())
	ipath, err := resolve.ResolveIPNS(ctx, api.namesys, ipath)
	if err == resolve.ErrNoNamesys {
		return nil, coreiface.ErrOffline
	} else if err != nil {
		return nil, err
	}
	return &ipath, nil
}
