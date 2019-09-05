package exec

import (
	"github.com/Ankr-network/ankr-chain/log"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	rawBytes, err := ioutil.ReadFile("F:/GoPath/src/github.com/Ankr-network/ankr-chain/contract/cpp/example/TestContract.wasm")
	if err != nil {
		t.Errorf("can read wasm file: %s", err.Error())
	}

	wasmVM := NewWASMVirtualMachine(rawBytes, log.DefaultRootLogger.With("contract", "test"))
	assert.NotEqual(t, wasmVM, nil)

	arg1, _ := wasmVM.wasmVM.SetBytes([]byte("Test"))
	fnIndex := wasmVM.ExportFnIndex("testFunc")
	assert.NotEqual(t, fnIndex, -1)
	_, err = wasmVM.Execute(fnIndex, arg1)
	if err != nil {
		t.Fatalf("could not execute Main: %v", err)
	}
}

