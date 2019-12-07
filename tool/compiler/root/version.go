package root

import (
	"fmt"

	"github.com/spf13/cobra"
)

var compilerVersion = "v0.0.1"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the contract compiler version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version:",compilerVersion)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}