package root

import (
	"github.com/Ankr-network/ankr-chain/tool/compiler/compile"
	"github.com/Ankr-network/ankr-chain/tool/compiler/decompile"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Args:cobra.MinimumNArgs(1),
	Use:   "contract-compiler",
	Short: "ankr smart contract compiler",
}

func init(){
	RootCmd.AddCommand(compile.CompileCmd, decompile.DecompileCmd)
}