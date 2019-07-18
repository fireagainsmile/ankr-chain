package commands

import (
	"fmt"

	"github.com/Ankr-network/ankr-chain/node"
	"github.com/spf13/cobra"
	tmcorecmd "github.com/tendermint/tendermint/cmd/tendermint/commands"
	tmcorecmn "github.com/tendermint/tendermint/libs/common"
)

func NewRunNodeCmd(nodeProvider node.AnkrNodeProvider) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node",
		Short: "Run the ankrchain node",
		RunE: func(cmd *cobra.Command, args []string) error {
			n, err := nodeProvider(config, logger)
			if err != nil {
				return fmt.Errorf("Failed to create node: %v", err)
			}

			// Stop upon receiving SIGTERM or CTRL-C.
			tmcorecmn.TrapSignal(logger, func() {
				if n.Node.IsRunning() {
					n.Node.Stop()
				}
			})

			if err := n.Node.Start(); err != nil {
				return fmt.Errorf("Failed to start node: %v", err)
			}
			logger.Info("Started node", "nodeInfo", n.Node.Switch().NodeInfo())

			// Run forever.
			select {}
		},
	}

	tmcorecmd.AddNodeFlags(cmd)
	AddHistoryStorageNodeFlags(cmd, config.HistoryDB.Type, config.HistoryDB.Host, config.HistoryDB.Name)
	return cmd
}
