package contract

import (
	"github.com/Ankr-network/ankr-chain/context"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	"github.com/tendermint/tendermint/libs/log"
)

type Contract struct {
   contractType ankrtypes.ContractType  `json:"contracttype"`
   address 		string                   `json:"address"`
   name         string                  `json:"name"`
   code         []byte                  `json:"code"`
}

func Call(context context.ContextContract, log log.Logger, conType ankrtypes.ContractType, code []byte, contractName string, method string, params []*ankrtypes.Param, rtnType string) (interface{}, error) {
	if conType == ankrtypes.ContractTypeNative {
		return NewInvoker(context, log).Invoke(code, contractName, method, params, rtnType)
	}

	return nil, nil
}