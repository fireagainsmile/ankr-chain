package log

import (
	"os"

	tmcorelog "github.com/tendermint/tendermint/libs/log"
)

type Logger tmcorelog.Logger

var (
	DefaultRootLogger = tmcorelog.NewTMLogger(tmcorelog.NewSyncWriter(os.Stdout))
)