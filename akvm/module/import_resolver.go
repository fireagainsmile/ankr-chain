package module

import (
	"errors"

	"github.com/Ankr-network/wagon/wasm"
)

type ImportResolver interface {
	Resolve(name string) (*wasm.Module, error)
}

type ImplImportResolver struct {
	envModle *ModuleEnv
}

func NewImportResolver(envModle *ModuleEnv) ImportResolver {
	return &ImplImportResolver{envModle}
}

func (r *ImplImportResolver) resolverEnvModule() (*wasm.Module, error) {
	if r.envModle == nil || r.envModle.wasmModule == nil{
		panic("envModle nil")
	}

	m := r.envModle.wasmModule

	m.Export = &wasm.SectionExports {
		Entries: make(map[string]wasm.ExportEntry),
	}

	var funcIndex uint32 = 0
	funcTable := r.envModle.ImportedFuncTable()
	for funcName, importedFunc := range funcTable {
		m.FunctionIndexSpace = append(m.FunctionIndexSpace, *importedFunc)
		m.Export.Entries[funcName] = wasm.ExportEntry {
			FieldStr: funcName,
			Kind:     wasm.ExternalFunction,
			Index:    funcIndex,
		}

		funcIndex++
	}

	return r.envModle.wasmModule, nil
}

func (r *ImplImportResolver) Resolve(name string) (*wasm.Module, error) {
	if name == "env"  {
		return r.resolverEnvModule()
	}else {
		return nil, errors.New("not support importing other's module at present")
	}

}


