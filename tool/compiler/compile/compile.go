package compile

import (
	"fmt"
	"github.com/Ankr-network/ankr-chain/tool/compiler/abi"
	"github.com/Ankr-network/ankr-chain/tool/compiler/decompile"
	"github.com/Ankr-network/ankr-chain/tool/compiler/parser"
	"github.com/spf13/cobra"
	"github.com/cheggaaa/pb"
	"os"
	"strings"
)

var (
	outputFlag = "output"
)

type CompileOptions interface {
	Options() []string
}

type Executable interface {
	Execute(args []string) error
}

// rootCmd represents the base command when called without any subcommands
var CompileCmd= &cobra.Command{
	Args:  cobra.MinimumNArgs(1),
	Use:   "compile",
	Short: "ankr smart contract compile tool",
	Long:  `Compile ankr smart contract and generate WebAssembly binary format`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: compile,
}

func init()  {
	CompileCmd.AddCommand(decompile.DecompileCmd)
}

func compile(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("expected filename argument.")
		cmd.Help()
		return
	}
	fmt.Println("compiling ", args)
	bar := pb.StartNew(100)
	bar.SetWidth(100)
	bar.ShowFinalTime = false

	err := exeCommand(abi.NewContractClass(), args)
	if err != nil {
		fmt.Println(err)
		return
	}
	bar.Add(25)
	bar.AlwaysUpdate = true

	//exec clang commands
	err = exeCommand(NewClangOption(), args)
	if err != nil {
		fmt.Println(err)
		return
	}
	bar.Add(25)

	//exec smart contract rule check
	sourceFile := getSrcFile(args)
	err = exeCommand(parser.NewRegexpParser(), []string{sourceFile})
	if err != nil {
		fmt.Println(err)
		return
	}
	bar.Add(25)

	//exec wasm-ld to generate binary file
	err = exeCommand(NewDefaultWasmOptions(), args)
	if err != nil {
		fmt.Println(err)
		return
	}
	bar.Add(25)
	bar.Finish()

	fmt.Println("Compile smart contract finished.")
}

//helper function to execute commands
func exeCommand(cmd Executable, args []string) error {
	return cmd.Execute(args)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := CompileCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	CompileCmd.Flags().StringVar(&OutPutDir, outputFlag, "./", "output file directory")
	CompileCmd.Flags().BoolVar(&abi.GenerateAbi, "gen-abi", false, "generate abi")
}

func getSrcFile(args []string) string {
	for _, arg := range args {
		argSlice := strings.Split(arg, ".")
		if len(argSlice) == 2{
			switch argSlice[1] {
			case "cpp", "cc","c":
				return arg
			}
		}
	}
	return ""
}