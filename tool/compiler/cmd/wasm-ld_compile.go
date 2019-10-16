package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

const (
	CodePrefixLen = 10
	ExtensionLen = 7
)

type ContractType int
const (
	_ ContractType = iota
	ContractTypeNative  = 0x01
	ContractTypeRuntime = 0x02
	ContractTypeUnknown = 0x03
)

type ContractVMType int
const (
	_ ContractVMType = iota
	ContractVMTypeWASM    = 0x01
	ContractVMTypeUnknown = 0x02
)

type ContractPatternType int
const (
	_  ContractPatternType = iota
	ContractPatternType1       = 0x01
	ContractPatternType2       = 0x02
	ContractPatternTypeUnknown = 0x03
)
type BinPrefix struct {
	TypeBin ContractType
	VMTypeBin ContractVMType
	PattenTypeBin ContractPatternType
	Extension [ExtensionLen]byte
}

func NewBinPrefix(contractType ContractType, contractVmType ContractVMType, contractPattenType ContractPatternType) *BinPrefix {
	return &BinPrefix{
		TypeBin: contractType,
		VMTypeBin: contractVmType,
		PattenTypeBin: contractPattenType,
	}
}

func (b *BinPrefix)SetOption(op [ExtensionLen]byte) *BinPrefix {
	b.Extension = op
	return b
}

func (b BinPrefix)Byte() []byte {
	out := make([]byte, CodePrefixLen, CodePrefixLen)
	out[0] = byte(b.TypeBin)
	out[1] = byte(b.VMTypeBin)
	out[2] = byte(b.VMTypeBin)
	for i, b := range b.Extension {
		out[i + 3] = b
	}
	return out
}

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
	srcFile := filterSrcFile(args).name
	srcFileSlice := strings.Split(srcFile, ".")
	srcFile = fmt.Sprintf("%s.o", srcFileSlice[0])
	distFile := fmt.Sprintf("%s.wasm", srcFileSlice[0])
	//wasmOp := NewDefaultWasmOptions()
	wasmArgs := wasmOp.Options()
	wasmArgs = append(wasmArgs, srcFile, "-o", distFile)
	out, err := exec.Command("wasm-ld.exe", wasmArgs...).Output()
	if err != nil {
		return err
	}
	if string(out) != "" {
		return errors.New(string(out))
	}

	prefix := NewBinPrefix(ContractTypeNative, ContractVMTypeWASM, ContractPatternType1)
	err = addPrefixToFile(distFile, *prefix)
	if err != nil {
		return err
	}
	//remove intermediate file
	err = os.Remove(srcFile)
	if err != nil {
		return err
	}
	renameFile := viper.GetString(outputDir)
	renameFile = path.Join(renameFile, distFile)
	return os.Rename(distFile, renameFile)
}

func addPrefixToFile(fileName string, prefix BinPrefix) error {
	srcByte, err := ioutil.ReadFile(fileName)
	if err != nil{
		return err
	}
	newByte := append(prefix.Byte(), srcByte...)
	return ioutil.WriteFile(fileName, newByte, 0600)
}
