package native

import (
	"fmt"
	"reflect"

	"github.com/Ankr-network/ankr-chain/context"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	"github.com/tendermint/tendermint/libs/log"
)

type NativeInvoker struct {
	nativeConracts map[string]interface{}
	context       context.ContextContract
	log           log.Logger
}

func NewNativeInvoker(context context.ContextContract, log log.Logger) *NativeInvoker {
	nativeConracts := map[string]interface{}{"AnkrCoin" : NewAnkrCoin(context, log)}

	return &NativeInvoker{nativeConracts, context, log}
}

func (invoker *NativeInvoker)Invoke(code []byte, contractName string, method string, params []*ankrtypes.Param, rtnType string) (interface{}, error) {
	natiContractI, ok := invoker.nativeConracts[contractName]
	if !ok {
		invoker.log.Error("NativeInvoker Invoke, can't find the responding contract", "contractName", contractName)
		return nil, fmt.Errorf("NativeInvoker Invoke, can't find the responding contract, contractName=%d", contractName)
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
		return rtnValues[0].Interface(), nil
	}

	return nil, nil
}


