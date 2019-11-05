package native

import (
	"fmt"
	"math/big"
	"reflect"

	ankrcmm "github.com/Ankr-network/ankr-chain/common"
	"github.com/Ankr-network/ankr-chain/context"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	"github.com/tendermint/tendermint/libs/log"
)

type NativeInvoker struct {
	nativeConracts map[string]interface{}
	context        context.ContextContract
	log            log.Logger
}

func NewNativeInvoker(store appstore.AppStore, log log.Logger) *NativeInvoker {
	nativeConracts := map[string]interface{}{"ANKR" : NewAnkrCoin(store, log)}

	return &NativeInvoker{nativeConracts, nil, log}
}

func (invoker *NativeInvoker) SetContextContract(context context.ContextContract) {
	invoker.context = context
	for _, v := range invoker.nativeConracts {
		methodValue := reflect.ValueOf(v).MethodByName("SetContextContract")
		if methodValue.IsValid() {
			methodValue.Call([]reflect.Value{reflect.ValueOf(context)})
		}
	}
}

func (invoker *NativeInvoker) Invoke(context context.ContextContract, conPatt ankrcmm.ContractPatternType, appStore appstore.AppStore, code []byte, contractName string, method string, params []*ankrcmm.Param, rtnType string) (*ankrcmm.ContractResult, error) {
	invoker.SetContextContract(context)
	natiContractI, ok := invoker.nativeConracts[contractName]
	if !ok {
		invoker.log.Error("NativeInvoker Invoke, can't find the responding contract", "contractName", contractName)
		return nil, fmt.Errorf("NativeInvoker Invoke, can't find the responding contract, contractName=%s", contractName)
	}

	natiContract := reflect.ValueOf(natiContractI)
	methodValue := natiContract.MethodByName(method)
	args := make([]reflect.Value, len(params))
	for pIndex, param := range params {
		if pIndex != param.Index {
			invoker.log.Error("NativeInvoker Invoke, method param order invalid", "contractName", contractName, "method", method)
			return nil, fmt.Errorf("NativeInvoker Invoke, method param order invalid, contractName=%s, method=%s", contractName, method)
		}
		args[pIndex] = reflect.ValueOf(param.Value)
	}

	rtnValues := methodValue.Call(args)
	if len(rtnValues) > 0 {
		if rtnValues[0].Type().Name() == rtnType {
			context.SpendGas(new(big.Int).SetUint64(0))
		  	return &ankrcmm.ContractResult{true, rtnValues[0].Type().Name(), rtnValues[0].Interface()}, nil
		}else {
			return &ankrcmm.ContractResult{false, rtnValues[0].Type().Name(), rtnValues[0].Interface()}, nil
		}
	}

	return nil, fmt.Errorf("invalid native contract call: contractName=%s, method=%s", contractName, method)
}


