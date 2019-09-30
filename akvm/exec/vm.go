package exec

import (
	"bytes"

	"github.com/Ankr-network/ankr-chain/akvm/memory"
	"github.com/Ankr-network/ankr-chain/akvm/module"
	"github.com/go-interpreter/wagon/exec"
	"github.com/go-interpreter/wagon/log"
	"github.com/go-interpreter/wagon/wasm"
)

type WASMVirtualMachine struct {
	wasmVM       *exec.VM
	envModule    *module.ModuleEnv
	log          log.Logger
}

func NewWASMVirtualMachine(code []byte, log log.Logger) *WASMVirtualMachine {
	wasmVM :=  &WASMVirtualMachine{ envModule: module.NewModuleEnv()}
	wasmVM.log = log
	wasmVM.loadAndInstantiateModule(code)

	return wasmVM
}

func (wvm *WASMVirtualMachine) loadAndInstantiateModule(code []byte) {
	if wvm.envModule == nil {
		panic("WASMVirtualMachine envModle nil")
	}

	importResolver := module.NewImportResolver(wvm.envModule)
	m, err := wasm.ReadModule(bytes.NewReader(code), importResolver.Resolve)
	if err != nil {
		panic(err)
	}

	m.HeapMem = memory.NewHeapMemory()

	/*err = validate.VerifyModule(m)
	if err != nil {
		panic(err)
	}*/

	vm, err := exec.NewVM(m)
	if err != nil {
		panic(err)
	}

	vm.SetLogger(wvm.log)

	wvm.wasmVM = vm
}

func (wvm *WASMVirtualMachine) ExportFnIndex(fnName string) int64 {
	if wvm.wasmVM.Module() == nil || wvm.wasmVM.Module().Export == nil {
		return -1
	}

	exportEntry, ok := wvm.wasmVM.Module().Export.Entries[fnName]
	if ok && exportEntry.Kind == wasm.ExternalFunction{
		return int64(exportEntry.Index)
	}

	return -1
}

func (wvm *WASMVirtualMachine) FuncSig(fnIndex int64) wasm.Function {
	return wvm.wasmVM.Module().FunctionIndexSpace[fnIndex]
}

func (wvm *WASMVirtualMachine) SetBytes(bytes []byte) (uint64, error) {
	return wvm.SetBytes(bytes)
}

func (wvm *WASMVirtualMachine) ReadString(off int64) (string, error) {
	return wvm.wasmVM.ReadString(off)
}

func (wvm *WASMVirtualMachine) SetContrInvoker(contrInvoker exec.ContractInvoker){
	wvm.wasmVM.SetContrInvoker(contrInvoker)
}

func (wvm *WASMVirtualMachine) ContrInvoker() exec.ContractInvoker{
	return wvm.wasmVM.ContrInvoker()
}

func (wvm *WASMVirtualMachine) Execute(fnIndex int64, rtnType string, args ...uint64)(interface{}, error) {
	return wvm.wasmVM.ExecCode(fnIndex, rtnType, args...)
}


