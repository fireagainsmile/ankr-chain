package p2p



import (
	"fmt"
	"github.com/Ankr-network/ankr-chain/common/code"
	"github.com/Ankr-network/ankr-chain/router"
	"github.com/tendermint/tendermint/abci/types"
	"strings"
)

type P2PQueryHandler struct {
	pf *PeerFilter
	sds *Seeds
}

func Init(pf *PeerFilter, sds *Seeds){
	qh := &P2PQueryHandler{pf, sds}
	router.QueryRouterInstance().AddQueryHandler("p2p", qh)
}

func (pqh *P2PQueryHandler) Query(reqQuery types.RequestQuery) (resQuery types.ResponseQuery) {
	path := reqQuery.Path
	pathSegs := strings.Split(path[1:], "/")
	if len(pathSegs) == 0 {
		resQuery = types.ResponseQuery{Code: code.CodeTypeP2PFilterPathError, Log: fmt.Sprintf("p2p filter query pah error: %s", reqQuery.Path)}
		return
	}
	switch pathSegs[0] {
	case "filter":
		if len(pathSegs) != 3 {
			resQuery = types.ResponseQuery{Code: code.CodeTypeP2PFilterPathError, Log: fmt.Sprintf("p2p filter query path error: %s", reqQuery.Path)}
			return
		}
		if pqh.pf != nil{
			return pqh.pf.Query(pathSegs[1], pathSegs[2])
		}
	case "seeds":
       if pqh.sds != nil {
       		return pqh.sds.Query()
	   }

	}
}


