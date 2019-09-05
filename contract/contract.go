package contract

import (
	"fmt"
	"github.com/Ankr-network/ankr-chain/akvm/exec"
	"github.com/Ankr-network/ankr-chain/log"
)

type Contract struct {

}

func (c *Contract) Invoke(code []byte, contractName string, method string, param []*Param) error {
	akvm := exec.NewWASMVirtualMachine(code, log.DefaultRootLogger.With("conract", contractName))
	if akvm == nil {
		return fmt.Errorf("can't creat vitual machiane: contractName=%s, method=%s", contractName, method)
	}

	fnIndex := akvm.ExportFnIndex(method)
	if fnIndex == -1 {
		return fmt.Errorf("can't get valid fnIndex: method=%s", method)
	}

	fSig := akvm.FuncSig(fnIndex)
	if len(fSig.Sig.ParamTypes) != len(param){
		return fmt.Errorf("input params' len invlid: len=%d", len(param))
	}

	var args []uint64

	for _, p := range param {
		if p.ParamType == "string" {
			val := p.Value.(string)
			arg, err := akvm.SetBytes([]byte(val))
			if err != nil {
				return fmt.Errorf("param err: index=%d, type=string, val=%s", p.Index, val)
			}

			args = append(args, arg)
		} else if  p.ParamType == "int32" {
			val := p.Value.(int32)
			args = append(args, uint64(val))
		} else if  p.ParamType == "int64" {
			val := p.Value.(int64)
			args = append(args, uint64(val))
		}else {
			return fmt.Errorf("param err: index=%d, type=%s", p.Index, p.ParamType)
		}
	}

	_, err := akvm.Execute(fnIndex, args...)

	return err
}