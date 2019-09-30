package exec

import (
	"io/ioutil"
	"testing"

	"github.com/Ankr-network/ankr-chain/log"
	"github.com/stretchr/testify/assert"
)

func TestExecuteWithNoReturn(t *testing.T) {
	rawBytes, err := ioutil.ReadFile("F:/GoPath/src/github.com/Ankr-network/ankr-chain/contract/example/cpp/TestContract.wasm")
	if err != nil {
		t.Errorf("can't read wasm file: %s", err.Error())
	}

	wasmVM := NewWASMVirtualMachine(rawBytes, log.DefaultRootLogger.With("contract", "test"))
	assert.NotEqual(t, wasmVM, nil)

	jsonArg := "{\"testStr\":\"testFunc arg\"}"

	method, _ := wasmVM.wasmVM.SetBytes([]byte("testFunc"))
	arg, _ := wasmVM.wasmVM.SetBytes([]byte(jsonArg))
	fnIndex := wasmVM.ExportFnIndex("ContractEntry")
	assert.NotEqual(t, fnIndex, -1)
	_, err = wasmVM.Execute(fnIndex, "", []uint64{method, arg}...)
	if err != nil {
		t.Fatalf("could not execute Main: %v", err)
	}
}

func TestExecuteWithIntReturn(t *testing.T) {
	rawBytes, err := ioutil.ReadFile("F:/GoPath/src/github.com/Ankr-network/ankr-chain/contract/example/cpp/TestContract.wasm")
	if err != nil {
		t.Errorf("can't read wasm file: %s", err.Error())
	}

	wasmVM := NewWASMVirtualMachine(rawBytes, log.DefaultRootLogger.With("contract", "test"))
	assert.NotEqual(t, wasmVM, nil)

	jsonArg := "{\"testStr\":\"testFuncWithInt arg\"}"

	method, _ := wasmVM.wasmVM.SetBytes([]byte("testFuncWithInt"))
	arg, _ := wasmVM.wasmVM.SetBytes([]byte(jsonArg))
	fnIndex := wasmVM.ExportFnIndex("ContractEntry")

	assert.NotEqual(t, fnIndex, -1)
	rtnStr, err := wasmVM.Execute(fnIndex, "string", []uint64{method, arg}...)
	if err != nil {
		t.Fatalf("could not execute Main: %v", err)
	}

	t.Logf("testFuncWithInt rtn=%s", rtnStr)
}

func TestExecuteWithStringReturn(t *testing.T) {
	rawBytes, err := ioutil.ReadFile("F:/GoPath/src/github.com/Ankr-network/ankr-chain/contract/example/cpp/TestContract.wasm")
	if err != nil {
		t.Errorf("can't read wasm file: %s", err.Error())
	}

	wasmVM := NewWASMVirtualMachine(rawBytes, log.DefaultRootLogger.With("contract", "test"))
	assert.NotEqual(t, wasmVM, nil)

	jsonArg := "{\"testStr\":\"testFuncWithString arg\"}"

	method, _ := wasmVM.wasmVM.SetBytes([]byte("testFuncWithString"))
	arg, _ := wasmVM.wasmVM.SetBytes([]byte(jsonArg))
	fnIndex := wasmVM.ExportFnIndex("ContractEntry")

	assert.NotEqual(t, fnIndex, -1)
	rtnStr, err := wasmVM.Execute(fnIndex, "string", []uint64{method, arg}...)
	if err != nil {
		t.Fatalf("could not execute Main: %v", err)
	}

	t.Logf("testFuncWithInt rtn=%s", rtnStr)
}

