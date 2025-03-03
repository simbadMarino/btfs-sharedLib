package test

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/bittorrent/go-btfs/core"
	"github.com/bittorrent/go-btfs/core/bootstrap"
	"github.com/bittorrent/go-btfs/core/coreapi"
	mock "github.com/bittorrent/go-btfs/core/mock"
	"github.com/bittorrent/go-btfs/core/node/libp2p"
	"github.com/bittorrent/go-btfs/keystore"
	"github.com/bittorrent/go-btfs/repo"

	config "github.com/bittorrent/go-btfs-config"
	coreiface "github.com/bittorrent/interface-go-btfs-core"
	"github.com/bittorrent/interface-go-btfs-core/tests"
	"github.com/ipfs/go-datastore"
	syncds "github.com/ipfs/go-datastore/sync"
	"github.com/ipfs/go-filestore"
	ci "github.com/libp2p/go-libp2p/core/crypto"
	peer "github.com/libp2p/go-libp2p/core/peer"
	mocknet "github.com/libp2p/go-libp2p/p2p/net/mock"
)

const testPeerID = "QmTFauExutTsy4XP6JbMFcw2Wa9645HJt2bTqL6qYDCKfe"

type NodeProvider struct{}

func (NodeProvider) MakeAPISwarm(ctx context.Context, fullIdentity bool, n int) ([]coreiface.CoreAPI, error) {
	mn := mocknet.New()

	nodes := make([]*core.IpfsNode, n)
	apis := make([]coreiface.CoreAPI, n)

	for i := 0; i < n; i++ {
		var ident config.Identity
		if fullIdentity {
			sk, pk, err := ci.GenerateKeyPair(ci.RSA, 2048)
			if err != nil {
				return nil, err
			}

			id, err := peer.IDFromPublicKey(pk)
			if err != nil {
				return nil, err
			}
			kbytes, err := ci.MarshalPrivateKey(sk)
			if err != nil {
				return nil, err
			}

			ident = config.Identity{
				PeerID:  id.String(),
				PrivKey: base64.StdEncoding.EncodeToString(kbytes),
			}
		} else {
			ident = config.Identity{
				PeerID: testPeerID,
			}
		}

		c := config.Config{}
		c.Addresses.Swarm = []string{fmt.Sprintf("/ip4/18.0.%d.1/tcp/4001", i)}
		c.Identity = ident
		c.Experimental.FilestoreEnabled = true

		ds := syncds.MutexWrap(datastore.NewMapDatastore())
		r := &repo.Mock{
			C: c,
			D: ds,
			K: keystore.NewMemKeystore(),
			F: filestore.NewFileManager(ds, filepath.Dir(os.TempDir())),
		}

		node, err := core.NewNode(ctx, &core.BuildCfg{
			Routing: libp2p.DHTServerOption,
			Repo:    r,
			Host:    mock.MockHostOption(mn),
			Online:  fullIdentity,
			ExtraOpts: map[string]bool{
				"pubsub": true,
			},
		})
		if err != nil {
			return nil, err
		}
		nodes[i] = node
		apis[i], err = coreapi.NewCoreAPI(node)
		if err != nil {
			return nil, err
		}
	}

	err := mn.LinkAll()
	if err != nil {
		return nil, err
	}

	bsinf := bootstrap.BootstrapConfigWithPeers(
		[]peer.AddrInfo{
			nodes[0].Peerstore.PeerInfo(nodes[0].Identity),
		},
	)

	for _, n := range nodes[1:] {
		if err := n.Bootstrap(bsinf); err != nil {
			return nil, err
		}
	}

	return apis, nil
}

func TestIface(t *testing.T) {
	tests.TestApi(&NodeProvider{})(t)
}
