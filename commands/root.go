package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	ankrconfg "github.com/Ankr-network/ankr-chain/config"
	"github.com/Ankr-network/ankr-chain/log"
	tmcorecmd "github.com/tendermint/tendermint/cmd/tendermint/commands"
	tmcorecfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	tmflags "github.com/tendermint/tendermint/libs/cli/flags"
	tmcorelog "github.com/tendermint/tendermint/libs/log"
)

var (
	config = ankrconfg.DefaultAnkrConfig()
)

func init() {
	registerFlagsRootCmd(RootCmd)
}

func registerFlagsRootCmd(cmd *cobra.Command) {
	cmd.PersistentFlags().String("log_level", config.LogLevel, "Log level")
}

// ParseConfig retrieves the default environment configuration,
// sets up the Tendermint root and ensures that the root exists
func ParseConfig() (*ankrconfg.AnkrConfig, error) {
	conf := ankrconfg.DefaultAnkrConfig()
	err := viper.Unmarshal(conf)
	if err != nil {
		return nil, err
	}
	conf.SetRoot(conf.RootDir)
	tmcorecfg.EnsureRoot(conf.RootDir)
	if err = conf.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("Error in config file: %v", err)
	}
	return conf, err
}

// RootCmd is the root command for Tendermint core.
var RootCmd = &cobra.Command{
	Use:   "ankrchain",
	Short: "Ankr chain for distributed cloud compute network",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if cmd.Name() == tmcorecmd.VersionCmd.Name() {
			return nil
		}
		config, err = ParseConfig()
		if err != nil {
			return err
		}
		if config.LogFormat == tmcorecfg.LogFormatJSON {
			log.DefaultRootLogger = tmcorelog.NewTMJSONLogger(tmcorelog.NewSyncWriter(os.Stdout))
		}
		log.DefaultRootLogger, err = tmflags.ParseLogLevel(config.LogLevel, log.DefaultRootLogger, ankrconfg.DefaultLogLevel())
		if err != nil {
			return err
		}
		if viper.GetBool(cli.TraceFlag) {
			log.DefaultRootLogger = tmcorelog.NewTracingLogger(log.DefaultRootLogger)
		}
		log.DefaultRootLogger = log.DefaultRootLogger.With("tx", "main")
		return nil
	},
}
