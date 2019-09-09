package exec

import (
	"io/ioutil"
	"testing"

	"github.com/Ankr-network/ankr-chain/log"
	"github.com/stretchr/testify/assert"
)

func TestExecuteWithNoReturn(t *testing.T) {
	rawBytes, err := ioutil.ReadFile("F:/GoPath/src/github.com/Ankr-network/ankr-chain/contract/cpp/example/TestContract.wasm")
	if err != nil {
		t.Errorf("can read wasm file: %s", err.Error())
	}

	wasmVM := NewWASMVirtualMachine(rawBytes, log.DefaultRootLogger.With("contract", "test"))
	assert.NotEqual(t, wasmVM, nil)

	arg1, _ := wasmVM.wasmVM.SetBytes([]byte("Test"))
	fnIndex := wasmVM.ExportFnIndex("testFunc")
	assert.NotEqual(t, fnIndex, -1)
	_, err = wasmVM.Execute(fnIndex, "", arg1)
	if err != nil {
		t.Fatalf("could not execute Main: %v", err)
	}
}

func TestExecuteWithIntReturn(t *testing.T) {
	rawBytes, err := ioutil.ReadFile("F:/GoPath/src/github.com/Ankr-network/ankr-chain/contract/cpp/example/TestContract.wasm")
	if err != nil {
		t.Errorf("can read wasm file: %s", err.Error())
	}

	wasmVM := NewWASMVirtualMachine(rawBytes, log.DefaultRootLogger.With("contract", "test"))
	assert.NotEqual(t, wasmVM, nil)

	arg1, _ := wasmVM.wasmVM.SetBytes([]byte("Test"))
	fnIndex := wasmVM.ExportFnIndex("testFuncWithInt")
	assert.NotEqual(t, fnIndex, -1)
	rtnIndex, err := wasmVM.Execute(fnIndex, "", arg1)
	if err != nil {
		t.Fatalf("could not execute Main: %v", err)
	}

	t.Logf("testFuncWithInt rtn=%d", rtnIndex)
}

func TestExecuteWithStringReturn(t *testing.T) {
	rawBytes, err := ioutil.ReadFile("F:/GoPath/src/github.com/Ankr-network/ankr-chain/contract/cpp/example/TestContract.wasm")
	if err != nil {
		t.Errorf("can read wasm file: %s", err.Error())
	}

	wasmVM := NewWASMVirtualMachine(rawBytes, log.DefaultRootLogger.With("contract", "test"))
	assert.NotEqual(t, wasmVM, nil)

	arg1, _ := wasmVM.wasmVM.SetBytes([]byte("Test"))
	fnIndex := wasmVM.ExportFnIndex("testFuncWithString")
	assert.NotEqual(t, fnIndex, -1)
	rtnIndex, err := wasmVM.Execute(fnIndex, "string", arg1)
	if err != nil {
		t.Fatalf("could not execute Main: %v", err)
	}

	rtnStr, _ := rtnIndex.(string)

	t.Logf("testFuncWithInt rtn=%s", rtnStr)
}

