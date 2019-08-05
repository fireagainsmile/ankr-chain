package common

type RunMode int
const (
	_ RunMode = iota
	RunModeTesting
	RunModeProd
)

var RM = RunModeTesting
