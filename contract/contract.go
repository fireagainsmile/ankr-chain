package contract

import (
	"github.com/Ankr-network/ankr-chain/context"
	"github.com/Ankr-network/ankr-chain/contract/native"
	"github.com/Ankr-network/ankr-chain/contract/runtime"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	"github.com/tendermint/tendermint/libs/log"
)

type Contract interface {
	Call(context     context.ContextContract,
		appStore     appstore.AppStore,
		conType      ankrtypes.ContractType,
		code         []byte,
		contractName string,
		method       string,
		params       []*ankrtypes.Param,
		rtnType      string) (*ankrtypes.ContractResult, error)
}

type ContractImpl struct {
	invokerMap map[ankrtypes.ContractType]Invoker
}

func NewContract(store appstore.AppStore,  log log.Logger) Contract {
	cImpl := &ContractImpl{make(map[ankrtypes.ContractType]Invoker)}
	cImpl.init(store, log)

	return cImpl
}

func (c *ContractImpl) init(store appstore.AppStore,  log log.Logger) {
	c.registerInvoker(store, log)
}

func (c *ContractImpl) registerInvoker(store appstore.AppStore, log log.Logger){

	nativeInvoker := native.NewNativeInvoker(store, log)
	c.invokerMap[ankrtypes.ContractTypeNative] = nativeInvoker

	runtimeInvoker := runtime.NewRuntimeInvoke(log)
	c.invokerMap[ankrtypes.ContractTypeRuntime] = runtimeInvoker
}

func (c *ContractImpl) Call(context context.ContextContract, appStore appstore.AppStore, conType ankrtypes.ContractType, code []byte, contractName string, method string, params []*ankrtypes.Param, rtnType string) (*ankrtypes.ContractResult, error) {
	return c.invokerMap[conType].Invoke(context, appStore, code, contractName, method, params, rtnType)
}