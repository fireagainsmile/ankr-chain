package p2p

import (
	"github.com/Ankr-network/ankr-chain/common/code"
	"strings"

	"github.com/tendermint/tendermint/abci/types"
)

type PeerFilter struct {
	peerIDs   []string
	peerAddrs []string
}

func NewPeerFilter() *PeerFilter {
	return new(PeerFilter)
}

func (pf *PeerFilter) PeerIDSet() []string {
	return pf.peerIDs
}

func (pf *PeerFilter) PeerAddrSet() []string {
	return pf.peerAddrs
}

func (pf *PeerFilter) Config(pfConf string) error {
	pfSegs := strings.Split(pfConf, ",")
	for _, pfItem := range pfSegs {
		pfItemSegs := strings.Split(pfItem, ":")
		if len(pfItemSegs) == 2 {
			pf.peerAddrs = append(pf.peerAddrs, pfItem)
		}

		if len(pfItemSegs) == 1 {
			pf.peerIDs = append(pf.peerIDs, pfItem)
		}
	}

	return nil
}

func (pf *PeerFilter) Query(filterType string, peer string) (resQuery types.ResponseQuery) {
	switch filterType {
	case "id":
		for _, peerID := range pf.peerIDs {
			if peerID == peer {
				resQuery.Code = code.CodeTypeOK
				return
			}
		}
	case "addr":
		for _, peerAddr := range pf.peerAddrs {
			if peerAddr == peer {
				resQuery.Code = code.CodeTypeOK
				return
			}
		}
	}

	resQuery.Code = code.CodeTypeP2PFilterNotPass

	return
}
