package main

import (
	"os"
	"path/filepath"

	"github.com/Ankr-network/ankr-chain/commands"
	"github.com/Ankr-network/ankr-chain/config"
	ankrnode "github.com/Ankr-network/ankr-chain/node"
	"github.com/tendermint/tendermint/libs/cli"
)

func main() {
	rootCmd := commands.RootCmd
	rootCmd.Use = "ankrchain"
	rootCmd.Short = "ankr chain for distributed cloud compute network"
	commands.AddTendermintCoreCommands(rootCmd)

	rootCmd.AddCommand(commands.InitFilesCmd)

	nodeFunc := ankrnode.NewAnkrNode

	// Create & start node
	rootCmd.AddCommand(commands.NewRunNodeCmd(nodeFunc))

	cmd := cli.PrepareBaseCmd(rootCmd, "TM", os.ExpandEnv(filepath.Join("$HOME", config.DefaultAnkrChainDir())))
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}