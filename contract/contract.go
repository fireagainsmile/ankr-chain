package contract

import (
	"github.com/Ankr-network/ankr-chain/context"
	"github.com/Ankr-network/ankr-chain/contract/native"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	"github.com/tendermint/tendermint/libs/log"
)

var invokerMap map[ankrtypes.ContractType]Invoker

func Init(store appstore.AppStore,  log log.Logger) {
	native.Init(store, log)
	registerInvoker(store, log)
}

func registerInvoker(store appstore.AppStore, log log.Logger){
	invokerMap = make(map[ankrtypes.ContractType]Invoker)

	nativeInvoker := native.NewNativeInvoker(store, log)
	invokerMap[ankrtypes.ContractTypeNative] = nativeInvoker
}

func Call(context context.ContextContract, conType ankrtypes.ContractType, code []byte, contractName string, method string, params []*ankrtypes.Param, rtnType string) (interface{}, error) {
	return invokerMap[conType].Invoke(context, code, contractName, method, params, rtnType)
}