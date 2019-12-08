package root

import (
	"fmt"

	"github.com/Ankr-network/ankr-chain/version"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the contract compiler version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version:",version.CompilerVersion)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}