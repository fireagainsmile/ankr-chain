package contract

import (
	"errors"
	"fmt"

	akexe "github.com/Ankr-network/ankr-chain/akvm/exec"
	"github.com/Ankr-network/ankr-chain/log"
	"github.com/Ankr-network/ankr-chain/types"
	"github.com/go-interpreter/wagon/exec"
)

type Contract struct {

}

func (c *Contract) InvokeInternal(vmContext *exec.VMContext, code []byte, contractName string, method string, params interface{}, rtnType string) (interface{}, error) {
	paramValues := params.([]*types.Param)
	if paramValues == nil && len(paramValues) == 0 {
		return nil, errors.New("invalid params")
	}

	akvm := akexe.NewWASMVirtualMachine(code, log.DefaultRootLogger.With("conract", contractName))
	if akvm == nil {
		return -1, fmt.Errorf("can't creat vitual machiane: contractName=%s, method=%s", contractName, method)
	}

	akvm.SetContrInvoker(c)

	fnIndex := akvm.ExportFnIndex(method)
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

func (c *Contract) Invoke(code []byte, contractName string, method string, param []*types.Param, rtnType string) (interface{}, error) {
	akvm := akexe.NewWASMVirtualMachine(code, log.DefaultRootLogger.With("conract", contractName))
	if akvm == nil {
		return -1, fmt.Errorf("can't creat vitual machiane: contractName=%s, method=%s", contractName, method)
	}

	akvm.SetContrInvoker(c)

	fnIndex := akvm.ExportFnIndex(method)
	if fnIndex == -1 {
		return -1, fmt.Errorf("can't get valid fnIndex: method=%s", method)
	}

	fSig := akvm.FuncSig(fnIndex)
	if len(fSig.Sig.ParamTypes) != len(param){
		return -1, fmt.Errorf("input params' len invlid: len=%d", len(param))
	}

	var args []uint64

	for _, p := range param {
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