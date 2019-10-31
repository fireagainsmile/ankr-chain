package contract

import (
	ankrcmm "github.com/Ankr-network/ankr-chain/common"
	"github.com/Ankr-network/ankr-chain/context"
	"github.com/Ankr-network/ankr-chain/contract/native"
	"github.com/Ankr-network/ankr-chain/contract/runtime"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	"github.com/tendermint/tendermint/libs/log"
)

type Contract interface {
	Call(context     context.ContextContract,
		appStore     appstore.AppStore,
		conType      ankrcmm.ContractType,
		conPatt      ankrcmm.ContractPatternType,
		code         []byte,
		contractName string,
		method       string,
		params       []*ankrcmm.Param,
		rtnType      string) (*ankrcmm.ContractResult, error)
}

type ContractImpl struct {
	invokerMap map[ankrcmm.ContractType]Invoker
}

func NewContract(store appstore.AppStore,  log log.Logger) Contract {
	cImpl := &ContractImpl{make(map[ankrcmm.ContractType]Invoker)}
	cImpl.init(store, log)

	return cImpl
}

func (c *ContractImpl) init(store appstore.AppStore, log log.Logger) {
	c.registerInvoker(store, log)
}

func (c *ContractImpl) registerInvoker(store appstore.AppStore, log log.Logger){
	nativeInvoker := native.NewNativeInvoker(store, log)
	c.invokerMap[ankrcmm.ContractTypeNative] = nativeInvoker

	runtimeInvoker := runtime.NewRuntimeInvoke(log)
	c.invokerMap[ankrcmm.ContractTypeRuntime] = runtimeInvoker
}

func (c *ContractImpl) Call(context context.ContextContract, appStore appstore.AppStore, conType ankrcmm.ContractType, conPatt ankrcmm.ContractPatternType, code []byte, contractName string, method string, params []*ankrcmm.Param, rtnType string) (*ankrcmm.ContractResult, error) {
	return c.invokerMap[conType].Invoke(context, conPatt, appStore, code, contractName, method, params, rtnType)
}