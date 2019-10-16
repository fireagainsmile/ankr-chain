/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	outputDir  = "outputDir"
	outputFlag = "output"
)

type CompileOptions interface {
	Options() []string
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Args:  cobra.MinimumNArgs(1),
	Use:   "compile",
	Short: "ankr smart contract compile tool",
	Long:  `This is used to compile C/C++ source file into wasm format`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: compile,
}

func compile(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("error: no input file")
		return
	}
	err := execClang(args)
	if err != nil {
		fmt.Println(err)
		return
	}
	r := NewRegexpParser()
	sourceFile := filterSrcFile(args).name
	r.ParseFile(sourceFile)
	if !r.ValidContract() {
		fmt.Println("[Error]: Contract compilation stopped!")
		return
	}
	err = GenAbi(sourceFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = execWasmLd(args)
	if err != nil {
		fmt.Println(err)
		return
	}

}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().String(outputFlag, "", "output file directory")
	viper.BindPFlag(outputDir, rootCmd.Flags().Lookup(outputFlag))
}
