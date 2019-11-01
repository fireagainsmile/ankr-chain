package commands

import "github.com/spf13/cobra"

func AddPeerFilterNodeFlags(cmd *cobra.Command, peers  string) {
	cmd.Flags().String("allowedpeers", peers, "The allowed peers' id or addr by ',' seperation")
}
