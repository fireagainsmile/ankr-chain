package module

import (
	"github.com/go-interpreter/wagon/wasm"
	"reflect"
)

type ModuleEnv struct {
	wasmModule    *wasm.Module
	importedFuncs map[string]*wasm.Function
}

func NewModuleEnv() *ModuleEnv {
	mEnv :=  &ModuleEnv{
		wasm.NewModule(),
		make(map[string]*wasm.Function),
	}

	mEnv.RegisterImportedFunc(PrintIFunc, &wasm.Function{
		Sig: &wasm.FunctionSig{ParamTypes: []wasm.ValueType{wasm.ValueTypeI32}},
		Body: &wasm.FunctionBody{},
		Host: reflect.ValueOf(Print_i),
		Name: PrintIFunc,
	})

	mEnv.RegisterImportedFunc(PrintSFunc, &wasm.Function{
		Sig: &wasm.FunctionSig{ParamTypes: []wasm.ValueType{wasm.ValueTypeI32}},
		Body: &wasm.FunctionBody{},
		Host: reflect.ValueOf(Print_s),
		Name: PrintSFunc,
	})

	mEnv.RegisterImportedFunc(StrlenFunc, &wasm.Function{
		Sig: &wasm.FunctionSig{ParamTypes: []wasm.ValueType{wasm.ValueTypeI32}},
		Body: &wasm.FunctionBody{},
		Host: reflect.ValueOf(Strlen),
		Name: StrlenFunc,
	})

	mEnv.RegisterImportedFunc(StrcmpFunc, &wasm.Function{
		Sig: &wasm.FunctionSig{ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32}, ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32}},
		Body: &wasm.FunctionBody{},
		Host: reflect.ValueOf(Strcmp),
		Name: StrcmpFunc,
	})

	return mEnv
}

func (m *ModuleEnv) ImportedFuncTable() map[string]*wasm.Function {
	return m.importedFuncs
}

func (m *ModuleEnv) WasmModule() *wasm.Module {
	return m.wasmModule
}

func (m *ModuleEnv) RegisterImportedFunc(name string, importedFunc *wasm.Function) {
	m.importedFuncs[name] = importedFunc
}
