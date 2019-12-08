package commands

import (
	"fmt"

	"github.com/Ankr-network/ankr-chain/log"
	"github.com/Ankr-network/ankr-chain/version"
	"github.com/spf13/cobra"
	tmcorecommands "github.com/tendermint/tendermint/cmd/tendermint/commands"
	"github.com/tendermint/tendermint/p2p"
)

func resetAll(cmd *cobra.Command, args []string) {
	tmcorecommands.ResetAll(config.DBDir(), config.P2P.AddrBookFile(), config.PrivValidatorKeyFile(),
		config.PrivValidatorStateFile(), log.DefaultRootLogger)
}

var ResetAllCmd = &cobra.Command{
	Use:   "unsafe_reset_all",
	Short: "(unsafe) Remove all the data and WAL, reset this node's validator to genesis state",
	Run:   resetAll,
}

// ShowNodeIDCmd dumps node's ID to the standard output.
var ShowNodeIDCmd = &cobra.Command{
	Use:   "show_node_id",
	Short: "Show this node's ID",
	RunE:  showNodeID,
}

func showNodeID(cmd *cobra.Command, args []string) error {
	nodeKey, err := p2p.LoadNodeKey(config.NodeKeyFile())
	if err != nil {
		return err
	}

	fmt.Println(nodeKey.ID())
	return nil
}

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version info",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.APPVersion)
	},
}

func AddTendermintCoreCommands(cmd *cobra.Command) {
	cmd.AddCommand(
		ResetAllCmd,
		ShowNodeIDCmd,
		VersionCmd)
}


