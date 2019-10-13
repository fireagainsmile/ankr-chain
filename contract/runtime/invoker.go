package runtime

import (
	"errors"
	"fmt"
	"reflect"

	akexe "github.com/Ankr-network/ankr-chain/akvm/exec"
	ankrcontext "github.com/Ankr-network/ankr-chain/context"
	"github.com/Ankr-network/ankr-chain/log"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	"github.com/go-interpreter/wagon/exec"
)

const (
	ContractEntry = "ContractEntry"
)

type RuntimeInvoke struct {
	context ankrcontext.ContextAKVM
	log     log.Logger
}

func NewRuntimeInvoke(log log.Logger) *RuntimeInvoke {
	return &RuntimeInvoke{nil, log}
}

func (r *RuntimeInvoke) InvokeInternal(contractAddr string, ownerAddr string, callerAddr string, vmContext *exec.VMContext, code []byte, contractName string, method string, params interface{}, rtnType string) (interface{}, error) {
	paramValues := params.([]*ankrtypes.Param)
	if paramValues == nil && len(paramValues) == 0 {
		return nil, errors.New("invalid params")
	}

	akvm := akexe.NewWASMVirtualMachine(contractAddr, ownerAddr, callerAddr, vmContext.GasMetric(), vmContext.Publisher(), code, r.log)
	if akvm == nil {
		return -1, fmt.Errorf("can't creat vitual machiane: contractName=%s, method=%s", contractName, method)
	}

	akvm.SetContrInvoker(r)

	fnIndex := akvm.ExportFnIndex(ContractEntry)
	if fnIndex == -1 {
		return -1, fmt.Errorf("can't get valid fnIndex: method=%s", method)
	}

	fSig := akvm.FuncSig(fnIndex)
	if len(fSig.Sig.ParamTypes) != len(paramValues){
		return -1, fmt.Errorf("input params' len invlid: len=%d", len(paramValues))
	}

	var args []uint64

	for _, p := range paramValues {
		if p.ParamType == "string" {
			val := p.Value.(string)
			arg, err := akvm.SetBytes([]byte(val))
			if err != nil {
				return -1, fmt.Errorf("param err: index=%d, type=string, val=%s", p.Index, val)
			}

			args = append(args, arg)
		} else if  p.ParamType == "int32" {
			val := p.Value.(int32)
			args = append(args, uint64(val))
		} else if  p.ParamType == "int64" {
			val := p.Value.(int64)
			args = append(args, uint64(val))
		}else {
			return -1, fmt.Errorf("param err: index=%d, type=%s", p.Index, p.ParamType)
		}
	}

	return akvm.Execute(fnIndex, rtnType, args...)
}

func (r *RuntimeInvoke) Invoke(context ankrcontext.ContextContract, appStore appstore.AppStore, code []byte, contractName string, method string, param []*ankrtypes.Param, rtnType string) (*ankrtypes.ContractResult, error) {
	r.context = ankrcontext.CreateContextAKVM(context,appStore)
	akvm := akexe.NewWASMVirtualMachine(context.ContractAddr(), context.OwnerAddr(), context.SenderAddr(), context, context, code, r.log)
	if akvm == nil {
		return &ankrtypes.ContractResult{false, rtnType, nil}, fmt.Errorf("can't creat vitual machiane: contractName=%s, method=%s", contractName, method)
	}

	akvm.SetContrInvoker(r)

	fnIndex := akvm.ExportFnIndex(method)
	if fnIndex == -1 {
		return &ankrtypes.ContractResult{false, rtnType, nil}, fmt.Errorf("can't get valid fnIndex: method=%s", method)
	}

	fSig := akvm.FuncSig(fnIndex)
	if len(fSig.Sig.ParamTypes) != len(param){
		return &ankrtypes.ContractResult{false, rtnType, nil}, fmt.Errorf("input params' len invlid: len=%d", len(param))
	}

	var args []uint64

	for _, p := range param {
		if p.ParamType == "string" {
			val := p.Value.(string)
			arg, err := akvm.SetBytes([]byte(val))
			if err != nil {
				return &ankrtypes.ContractResult{false, rtnType, nil}, fmt.Errorf("param err: index=%d, type=string, val=%s", p.Index, val)
			}

			args = append(args, arg)
		} else if  p.ParamType == "int32" {
			val := p.Value.(int32)
			args = append(args, uint64(val))
		} else if  p.ParamType == "int64" {
			val := p.Value.(int64)
			args = append(args, uint64(val))
		}else {
			return &ankrtypes.ContractResult{false, rtnType, nil}, fmt.Errorf("param err: index=%d, type=%s", p.Index, p.ParamType)
		}
	}

	akvmResult, err := akvm.Execute(fnIndex, rtnType, args...)
	if err != nil {
		return &ankrtypes.ContractResult{false, rtnType, nil}, err
	}

	if reflect.ValueOf(akvmResult).Type().Name() == rtnType  {
		return &ankrtypes.ContractResult{true, reflect.ValueOf(akvmResult).Type().Name(), akvmResult}, err
	}else {
		return &ankrtypes.ContractResult{false, reflect.ValueOf(akvmResult).Type().Name(), akvmResult}, err
	}
}