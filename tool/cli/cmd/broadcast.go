package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// broadcastCmd represents the broadcast command

var (
	broadcastUrl = "broadcastUrl"
	broadcastFile = "broadcastFile"
)
var broadcastCmd = &cobra.Command{
	Use:   "broadcast",
	Short: "broadcast signed transaction to ankr chain",
	Run: runBroadcast,
}

func init() {
	err := addStringFlag(broadcastCmd, broadcastUrl, urlParam, "", "", "validator url", required)
	if err != nil {
		panic(err)
	}
	err = addStringFlag(broadcastCmd, broadcastFile, fileParam, "f", "", "transaction signed data", required)
	if err != nil {
		panic(err)
	}
}

func runBroadcast(cmd *cobra.Command, args []string){
	validatorUrl = viper.GetString(broadcastUrl)
	if len(validatorUrl) < 1 {
		fmt.Println("Illegal url is received!")
		return
	}
	cl := newAnkrHttpClient(validatorUrl)
	signedFile := viper.GetString(broadcastFile)
	signedBytes, err := ioutil.ReadFile(signedFile)
	if err != nil {
		fmt.Println("Failed to read from file:", err.Error())
		return
	}
	txHash, commitHeight, _, err := cl.BroadcastTxCommit(signedBytes)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Broadcast transaction success")
	fmt.Println("Hash:", txHash)
	fmt.Println("height:", commitHeight)
}
