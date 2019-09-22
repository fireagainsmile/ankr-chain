package module

import (
	"encoding/json"

	"github.com/Ankr-network/ankr-chain/context"
	"github.com/Ankr-network/ankr-chain/types"
	"github.com/go-interpreter/wagon/exec"
)

const (
	PrintSFunc = "print_s"
	PrintIFunc = "print_i"
	StrlenFunc = "strlen"
	StrcmpFunc = "strcmp"
)

func Print_s(proc *exec.Process, strIdx int32) {
	str, err := proc.ReadString(int64(strIdx))
	if err != nil {
		proc.VM().Logger().Error("Print_s", "err", err)
	}
	proc.VM().Logger().Info("Print_s", "str", str)
}

func Print_i(proc *exec.Process, v int32) {
	proc.VM().Logger().Info("Print_i", "v", v)
}

func Strlen(proc *exec.Process, strIdx int32) int {
	len, err := proc.VM().Strlen(uint(strIdx))
	if err != nil {
		return -1
	}

	return len
}

func Strcmp(proc *exec.Process, strIdx1 int32, strIdx2 int32) int32 {
	cmpR, _ := proc.VM().Strcmp(uint(strIdx1), uint(strIdx2))
	return int32(cmpR)
}

func ContractCall(proc *exec.Process, contractIndex int32, methodIndex int32, paramJsonIndex int32, rtnType int32) int64 {
	toReadContractName, err := proc.ReadString(int64(contractIndex))
	if err != nil {
		proc.VM().Logger().Error("ContractCall read ContractName err", "err", err)
		return -1
	}

	toReadMethodName, err := proc.ReadString(int64(methodIndex))
	if err != nil {
		proc.VM().Logger().Error("ContractCall read MethodName err", "err", err)
		return -1
	}

	toReadJsonParam, err := proc.ReadString(int64(paramJsonIndex))
	if err != nil {
		proc.VM().Logger().Error("ContractCall read jsonParam err", "err", err)
		return -1
	}

	toReadRTNType, err := proc.ReadString(int64(rtnType))
	if err != nil {
		proc.VM().Logger().Error("ContractCall read rtnType err", "err", err)
		return -1
	}

	code, err := context.GetBCContext().LoadContract([]byte(toReadContractName))

	params := make([]*types.Param, 0)
    err =  json.Unmarshal([]byte(toReadJsonParam), params)
    if err != nil {
		proc.VM().Logger().Error("ContractCall json.Unmarshal err", "JsonParam", toReadJsonParam, "err", err)
		return -1
	}

    contrInvoker := proc.VM().ContrInvoker()
    if contrInvoker == nil {
		proc.VMContext().PushVM(proc.VM())
		rtnIndex, _ := proc.VM().ContrInvoker().InvokeInternal(proc.VMContext(), code, toReadContractName, toReadMethodName, params, toReadRTNType)
		lastVM, _:= proc.VMContext().PopVM()
		proc.VMContext().SetRunningVM(lastVM)
		switch rtnIndex.(type) {
		case int32:
			return int64(rtnIndex.(int32))
		case int64:
			return rtnIndex.(int64)
		case string:
			lastVM.SetBytes([]byte(rtnIndex.(string)))
		default:
			return -1
		}
	}else {
		proc.VM().Logger().Error("there is no contrInvoker set")
	}

    return -1
}

func ContractDelegateCall(proc *exec.Process, contractIndex int32, methodIndex int32, paramJsonIndex int32, rtnType int32) int64 {
	toReadContractName, err := proc.ReadString(int64(contractIndex))
	if err != nil {
		proc.VM().Logger().Error("ContractCall read ContractName err", "err", err)
		return -1
	}

	toReadMethodName, err := proc.ReadString(int64(methodIndex))
	if err != nil {
		proc.VM().Logger().Error("ContractCall read MethodName err", "err", err)
		return -1
	}

	toReadJsonParam, err := proc.ReadString(int64(paramJsonIndex))
	if err != nil {
		proc.VM().Logger().Error("ContractCall read jsonParam err", "err", err)
		return -1
	}

	toReadRTNType, err := proc.ReadString(int64(rtnType))
	if err != nil {
		proc.VM().Logger().Error("ContractCall read rtnType err", "err", err)
		return -1
	}

	code, err := context.GetBCContext().LoadContract([]byte(toReadContractName))

	params := make([]*types.Param, 0)
	err =  json.Unmarshal([]byte(toReadJsonParam), params)
	if err != nil {
		proc.VM().Logger().Error("ContractCall json.Unmarshal err", "JsonParam", toReadJsonParam, "err", err)
		return -1
	}

	contrInvoker := proc.VM().ContrInvoker()
	if contrInvoker == nil {
		proc.VMContext().PushVM(proc.VM())
		rtnIndex, _ := proc.VM().ContrInvoker().InvokeInternal(proc.VMContext(), code, toReadContractName, toReadMethodName, params, toReadRTNType)
		lastVM, _:= proc.VMContext().PopVM()
		proc.VMContext().SetRunningVM(lastVM)
		switch rtnIndex.(type) {
		case int32:
			return int64(rtnIndex.(int32))
		case int64:
			return rtnIndex.(int64)
		case string:
			lastVM.SetBytes([]byte(rtnIndex.(string)))
		default:
			return -1
		}
	}else {
		proc.VM().Logger().Error("there is no contrInvoker set")
	}

	return -1
}
