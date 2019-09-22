package contract

import (
	"github.com/Ankr-network/ankr-chain/context"
	"github.com/Ankr-network/ankr-chain/types"
	"github.com/tendermint/tendermint/libs/log"
)

type ContractType int
const (
	_ ContractType = iota
	ContractTypeNative
	ContractTypeRuntime
)

type Contract struct {
   contractType ContractType  `json:"contracttype"`
   address 		string        `json:"address"`
   name         string        `json:"name"`
   code         []byte        `json:"code"`
}

func (c *Contract) Call(context context.ContextContract, log log.Logger, contractName string, method string, params []*types.Param) (interface{}, error) {
	if c.contractType == ContractTypeNative {
		return NewInvoker(context, log).Invoke(contractName, method, params)
	}

	return nil, nil
}