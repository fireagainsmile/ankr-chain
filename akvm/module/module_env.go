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

	mEnv.RegisterImportedFunc(JsonObjectIndexFunc, &wasm.Function{
		Sig: &wasm.FunctionSig{ParamTypes: []wasm.ValueType{wasm.ValueTypeI32}, ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32}},
		Body: &wasm.FunctionBody{},
		Host: reflect.ValueOf(JsonObjectIndex),
		Name: JsonObjectIndexFunc,
	})

	mEnv.RegisterImportedFunc(JsonCreateObjectFunc, &wasm.Function{
		Sig: &wasm.FunctionSig{ParamTypes: []wasm.ValueType{}, ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32}},
		Body: &wasm.FunctionBody{},
		Host: reflect.ValueOf(JsonCreateObject),
		Name: JsonCreateObjectFunc,
	})

	mEnv.RegisterImportedFunc(JsonGetIntFunc, &wasm.Function{
		Sig: &wasm.FunctionSig{ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32}, ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32}},
		Body: &wasm.FunctionBody{},
		Host: reflect.ValueOf(JsonGetInt),
		Name: JsonGetIntFunc,
	})

	mEnv.RegisterImportedFunc(JsonGetStringFunc, &wasm.Function{
		Sig: &wasm.FunctionSig{ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32}, ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32}},
		Body: &wasm.FunctionBody{},
		Host: reflect.ValueOf(JsonGetString),
		Name: JsonGetStringFunc,
	})

	mEnv.RegisterImportedFunc(JsonPutIntFunc, &wasm.Function{
		Sig: &wasm.FunctionSig{ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32}, ReturnTypes: []wasm.ValueType{}},
		Body: &wasm.FunctionBody{},
		Host: reflect.ValueOf(JsonPutInt),
		Name: JsonPutIntFunc,
	})

	mEnv.RegisterImportedFunc(JsonPutStringFunc, &wasm.Function{
		Sig: &wasm.FunctionSig{ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32}, ReturnTypes: []wasm.ValueType{}},
		Body: &wasm.FunctionBody{},
		Host: reflect.ValueOf(JsonPutString),
		Name: JsonPutStringFunc,
	})

	mEnv.RegisterImportedFunc(JsonToStringFunc, &wasm.Function{
		Sig: &wasm.FunctionSig{ParamTypes: []wasm.ValueType{wasm.ValueTypeI32}, ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32}},
		Body: &wasm.FunctionBody{},
		Host: reflect.ValueOf(JsonToString),
		Name: JsonToStringFunc,
	})

	mEnv.RegisterImportedFunc(ContractCallFunc, &wasm.Function{
		Sig: &wasm.FunctionSig{ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32}, ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32}},
		Body: &wasm.FunctionBody{},
		Host: reflect.ValueOf(ContractCall),
		Name: JsonToStringFunc,
	})

	mEnv.RegisterImportedFunc(ContractDelegateCallFunc, &wasm.Function{
		Sig: &wasm.FunctionSig{ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32}, ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32}},
		Body: &wasm.FunctionBody{},
		Host: reflect.ValueOf(ContractDelegateCall),
		Name: ContractDelegateCallFunc,
	})

	mEnv.RegisterImportedFunc(TrigEventFunc, &wasm.Function{
		Sig: &wasm.FunctionSig{ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32}, ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32}},
		Body: &wasm.FunctionBody{},
		Host: reflect.ValueOf(TrigEvent),
		Name: TrigEventFunc,
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
