package cmd

import (
	"errors"
	"os/exec"
	"strings"
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

// compile c/c++ file into object
func execClang(args []string) error {
	srcFile := filterSrcFile(args)
	if srcFile.name == "" {
		return errors.New("no input file")
	}

	var co = &DefaultClangOptions
	switch srcFile.fileType {
	case cType:
		co.withC()
	case cPlusType:
		co.withCpp()
	}


	clangArgs := co.Options()
	clangArgs = append(clangArgs, srcFile.name)

	out, err := exec.Command(co.Compiler, clangArgs...).Output()
	if err != nil {
		return err
	}
	if string(out) != "" {
		return errors.New(string(out))
	}
	return nil
}

func filterSrcFile(args []string) srcContract {
	var src srcContract
	for _, arg := range args {
		argSlice := strings.Split(arg, ".")
		if len(argSlice) == 2{
			switch argSlice[1] {
			case "cpp", "cc":
				src.name = arg
				src.fileType = cPlusType
				return src
			case "c":
				src.name = arg
				src.fileType = cType
				return src
			}
		}
	}
	return src
}
