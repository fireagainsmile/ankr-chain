package cmd

import (
	"fmt"
	"github.com/Ankr-network/ankr-chain/tool/compiler/abi"
	"github.com/Ankr-network/ankr-chain/tool/compiler/parser"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	outputDir  = "outputDir"
	outputFlag = "output"
	genAbi bool
	genAbiFlag = "gen-abi"
)

type CompileOptions interface {
	Options() []string
}

type Executable interface {
	Execute(args []string) error
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

	//exec clang commands
	err := exeCommand(NewClangOption(), args)
	if err != nil {
		fmt.Println(err)
		return
	}

	//exec smart contract rule check
	sourceFile := filterSrcFile(args).name
	err = exeCommand(parser.NewRegexpParser(), []string{sourceFile})
	if err != nil {
		fmt.Println(err)
		return
	}

	//exec smart contract abi gen
	if genAbi{
		err = abi.GenAbi(sourceFile)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	//exec wasm-ld to generate binary file
	err = exeCommand(NewDefaultWasmOptions(), args)
	if err != nil {
		fmt.Println(err)
		return
	}

}

//helper function to execute commands
func exeCommand(cmd Executable, args []string) error {
	return cmd.Execute(args)
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
	//rootCmd.Flags().Bool(genAbiFlag,false, "generate abi file")
	rootCmd.Flags().BoolVar(&genAbi, "gen-abi", false, "generate abi")
	viper.BindPFlag(outputDir, rootCmd.Flags().Lookup(outputFlag))
}
