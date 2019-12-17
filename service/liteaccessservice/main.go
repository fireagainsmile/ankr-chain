package main

import (
	"os"

	"github.com/Ankr-network/ankr-chain/service/liteaccessservice/commands"
	lascmm "github.com/Ankr-network/ankr-chain/service/liteaccessservice/common"
	_ "github.com/Ankr-network/ankr-chain/service/liteaccessservice/statik"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"
)

var (
	rootCmd = &cobra.Command{
		Use:   "ankrchain-ias",
		Short: "ankrchain lite access service",
	}
)

func main() {
	cobra.EnableCommandSorting = false

	rootCmd.AddCommand(
		commands.Start(),
	)

	executor := cli.PrepareMainCmd(rootCmd, "ANKRCHAINlAS", os.ExpandEnv(lascmm.DefaultLasHome))
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}
