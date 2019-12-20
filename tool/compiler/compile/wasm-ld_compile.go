package compile

import (
	"fmt"
	"github.com/Ankr-network/ankr-chain/tool/compiler/abi"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var OutPutDir string

type WasmOptions struct {
	Compiler string
	Flags    []string
}

var DefaultWasmOptions = WasmOptions{
	Flags: []string{"--no-entry", "--strip-all", "--allow-undefined", "--no-merge-data-segments", "-zstack-size=16384", "--stack-first"},
}

//used to add flags that without parameters
func (wo *WasmOptions)addFlags(args []string) *WasmOptions {
	wo.Flags = append(wo.Flags, args...)
	return wo
}

func NewDefaultWasmOptions() *WasmOptions {
	return &DefaultWasmOptions
}

func (wo *WasmOptions) Options() []string {
	return wo.Flags
}
// by executing wasm-ld to transfer object into wasm format
// and remove

func (wasmOp *WasmOptions)Execute(args []string) error {
	distFileName := filepath.Base(args[0])
	srcFileName := filepath.Base(abi.ContractMainFile)
	srcFileSlice := strings.Split(srcFileName, ".")
	srcFile := fmt.Sprintf("%s.o", srcFileSlice[0])
	distFile := distFileName
	distSlice := strings.Split(distFile, ".")
	distFile = fmt.Sprintf("%s.wasm", distSlice[0])
	wasmArgs := wasmOp.Options()
	wasmArgs = append(wasmArgs, srcFile, "-o", distFile)
	out, err := exec.Command("wasm-ld", wasmArgs...).Output()
	if err != nil {
		return handleExeError( "wasm-ld", err)
	}
	if len(out) != 0 {
		fmt.Println(string(out))
	}
	err = addPrefixToFile(distFile, *abi.ABIPrefix)
	if err != nil {
		return err
	}
	//remove intermediate file
	err = os.Remove(srcFile)
	if err != nil {
		return err
	}
	if abi.ContractMainFile == abi.TempCppFile {
		err = os.Remove(abi.TempCppFile)
		if err != nil {
			return err
		}
	}
	renameFile := filepath.Join(OutPutDir, distFile)
	if _, err = os.Stat(OutPutDir); os.IsNotExist(err) {
		err = os.MkdirAll(OutPutDir, 0600)
		if err != nil {
			return err
		}
	}
	return os.Rename(distFile, renameFile)
}

func addPrefixToFile(fileName string, prefix abi.BinPrefix) error {
	srcByte, err := ioutil.ReadFile(fileName)
	if err != nil{
		return err
	}
	prefixArray := prefix.Byte()
	newByte := append(prefixArray[:], srcByte...)
	f, err := os.OpenFile(fileName, os.O_TRUNC|os.O_CREATE, 0600)
	defer f.Close()
	if err != nil {
		return err
	}
	_, err = f.Write(newByte)
	if err != nil {
		return err
	}
	return nil
}
