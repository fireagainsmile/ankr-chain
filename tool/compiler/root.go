package main

import (
	"github.com/Ankr-network/ankr-chain/tool/compiler/compile"
	"github.com/Ankr-network/ankr-chain/tool/compiler/decompile"
	"github.com/spf13/cobra"
)

var rootCmd= &cobra.Command{
	Args:cobra.MinimumNArgs(1),
	Use:   "compiler",
	Short: "ankr smart contract compiler",
}

func init(){
	rootCmd.AddCommand(compile.CompileCmd, decompile.DecompileCmd)
}