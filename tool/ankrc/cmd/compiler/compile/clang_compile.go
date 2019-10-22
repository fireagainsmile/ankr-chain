package compile

import (
	"errors"
	abi2 "github.com/Ankr-network/ankr-chain/tool/ankrc/cmd/compiler/abi"
	"os/exec"
)

var (
	cPlusType = "c++"
	cType = "c"
)

type ClangOptions struct {
	Compiler   string
	Compile    string
	Optimise   string
	Target     string
	Standard   string
	Extensions []string
}

var DefaultClangOptions = ClangOptions{
	Compiler: "clang++",
	Optimise: "-O3",
	Compile:  "-c",
	Target:   "--target=wasm32",
	Standard: "-std=c++14",
}

func (co *ClangOptions)Options() (args []string) {
	args = append(args, co.Optimise, co.Compile, co.Target)
	args = append(args, co.Extensions...)
	return
}

// set clang compile options
//func (*ClangOptions)setOption *ClangOptions{...}
func (cp *ClangOptions) withC()  *ClangOptions {
	cp.Compiler = "clang"
	return cp
}

func (cp *ClangOptions) withCpp() *ClangOptions  {
	cp.Compiler = "clang++"
	return cp
}

type srcContract struct {
	name string
	fileType string // the contract whether is c or c++ type
}

func NewClangOption() *ClangOptions {
	return &DefaultClangOptions
}

// compile c/c++ file into object
func (co *ClangOptions)Execute(args []string) error  {

	clangArgs := co.Options()
	clangArgs = append(clangArgs, abi2.TempCppFile)

	out, err := exec.Command(co.Compiler, clangArgs...).Output()
	if err != nil {
		return err
	}
	if string(out) != "" {
		return errors.New(string(out))
	}
	return nil
}
