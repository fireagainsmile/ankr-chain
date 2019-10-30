package decompile

import (
  "fmt"
  "github.com/spf13/cobra"
  "os"
)

type Executable interface {
  Execute(args []string) error
}

// rootCmd represents the base command when called without any subcommands
var DecompileCmd= &cobra.Command{
  Args:cobra.MinimumNArgs(1),
  Use:   "decompile",
  Short: "Read a file in the WebAssembly binary format, and convert it to the WebAssembly text format.",
  Run: decompile,
  //	Run: func(cmd *cobra.Command, args []string) { },
}

func decompile(cmd *cobra.Command, args []string) {
  if len(args) != 1 {
    fmt.Println("expected filename argument.")
    cmd.Help()
    return
  }
  fmt.Println("Converting wasm to wast...")
  err := executeCommand(NewWasm2WatOp(), args)
  if err != nil {
    fmt.Println(err)
    return
  }
  fmt.Println("Finished Converting wasm to wast...")
}

func executeCommand(exe Executable, args []string) error {
  return exe.Execute(args)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
  if err := DecompileCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func init() {
  // Here you will define your flags and configuration settings.
  // Cobra supports persistent flags, which, if defined here,
  // will be global for your application.
  DecompileCmd.Flags().StringVarP(&outPut, "output", "o","", "output file name")
}