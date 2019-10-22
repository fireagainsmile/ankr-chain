package cmd

import (
  "fmt"
  "os"
  "github.com/Ankr-network/ankr-chain/tool/ankrc/cmd/compiler"
  "github.com/Ankr-network/ankr-chain/tool/ankrc/cmd/decompiler"
  "github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
  Use:   "ankrc",
  Short: "A brief description of your application",
  //	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func init()  {
  rootCmd.AddCommand(compiler.CompileCmd, decompiler.DecompileCmd)
}
