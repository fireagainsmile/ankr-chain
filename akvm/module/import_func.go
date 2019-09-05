package module

import (
	"github.com/go-interpreter/wagon/exec"
)

const (
	PrintSFunc = "print_s"
	PrintIFunc = "print_i"
	StrlenFunc = "strlen"
	StrcmpFunc = "strcmp"
)

func Print_s(proc *exec.Process, strIdx int32) {
	toReads := make([]byte, 100)
	len, _ := proc.ReadAt(toReads, int64(strIdx))
	str := string(toReads[:len])
	proc.VM().Logger().Info("Print_s", "str", str)
}

func Print_i(proc *exec.Process, v int32) {
	proc.VM().Logger().Info("Print_i", "v", v)
}

func Strlen(proc *exec.Process, strIdx int32) int {
	len, err := proc.VM().Strlen(uint(strIdx))
	if err != nil {
		return -1
	}

	return len
}

func Strcmp(proc *exec.Process, strIdx1 int32, strIdx2 int32) int32 {
	cmpR, _ := proc.VM().Strcmp(uint(strIdx1), uint(strIdx2))
	return int32(cmpR)
}
