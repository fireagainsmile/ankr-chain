package p2p

import (
	"github.com/Ankr-network/ankr-chain/common/code"
	"github.com/tendermint/tendermint/abci/types"
)

type Seeds struct {
	seedList string
}

func NewSeeds() *Seeds {
	return new(Seeds)
}

func (sds *Seeds) Config(sConf string) error{
	sds.seedList = sConf

	return nil
}

func (sds *Seeds) Query() (resQuery types.ResponseQuery) {
	resQuery.Code  = code.CodeTypeOK
	resQuery.Value = []byte(sds.seedList)
	return
}
