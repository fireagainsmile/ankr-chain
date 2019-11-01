package decompile

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

const (
	WasmPrefixLength = 10
	tmpFileName = "temp.wasm"
)

var outPut string
type Wasm2WatOptions struct {
	Compiler   string
	input string
	outPut string
	Extension []string
}

var DefaultWasm2Wat = Wasm2WatOptions{
	Compiler: "wasm2wat",
	outPut:outPut,
}

func NewWasm2WatOp() *Wasm2WatOptions {
	return &DefaultWasm2Wat
}

func (w *Wasm2WatOptions)AddExtension(arg string) *Wasm2WatOptions {
	w.Extension = append(w.Extension,arg)
	return w
}

func (w *Wasm2WatOptions)Options() []string {
	var options []string
	options = append(options, w.input)
	options = append(options, "-o", outPut)
	options = append(options, w.Extension...)
	return options
}

func (w *Wasm2WatOptions)addInput (file string) *Wasm2WatOptions {
	w.input = file
	return w
}

// execute wasm2wat hello.wasm -o hello.wat
// Execute read from hello.wasm and generate a temp file after trim heading prefix.
// delete temp file after decompiled wasm into wat
func (w *Wasm2WatOptions)Execute(args[]string) error {
	srcFile := filterWasm(args)
	if srcFile == ""{
		return errors.New("no wasm file found! ")
	}

	fileByte, err := ioutil.ReadFile(srcFile)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(tmpFileName, os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		return err
	}

	_, err = f.Write(fileByte[WasmPrefixLength:])
	if err != nil {
		return err
	}

	w.addInput(tmpFileName)

	//gen
	out, err := exec.Command(w.Compiler, w.Options()...).Output()
	if err != nil {
		return err
	}
	if string(out) != "" {
		return errors.New(string(out))
	}
	return nil
}

func filterWasm(args []string) string {
	for _, v := range args {
		oldLength := len(v)
		if oldLength != len(strings.TrimRight(v, "wasm")) {
			return v
		}
	}
	return  ""
}
