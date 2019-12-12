package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/rpc/core/types"
	"time"
)

var (
	// parameter names
	subscribeContent = "subscribeContent"
	subscribeTitle = "subscribeTitle"
	subscribeUrl = "subscribeUrl"
	subscribeTimeOut = "subscribeTimeOut"
	subscribeMaxCap = "subscribeMaxCap"
)
// subscribeCmd represents the subscribe command
var subscribeCmd = &cobra.Command{
	Use:   "subscribe",
	Short: "subscribe events from ankr chain",
	Run: subscribeFromAnkr,
}

func init() {
	err := addStringFlag(subscribeCmd, subscribeContent, queryParam, "", "","subscription query string", required)
	if err != nil {
		panic(err)
	}
	err = addStringFlag(subscribeCmd, subscribeTitle, titlePara, "", "subscription","subscription title of the query to be named", notRequired)
	err = addStringFlag(subscribeCmd, subscribeUrl, urlParam, "", "","url of the validator node", required)
	err = addIntFlag(subscribeCmd, subscribeTimeOut, timeoutParam, "", 30,"the time of seconds to wait before can not receive any response", notRequired)
	err = addIntFlag(subscribeCmd, subscribeMaxCap, capParam, "", 100,"maximum subscriptions to be subscribed from ankr chain", notRequired)
}

func subscribeFromAnkr(cmd *cobra.Command, args []string) {
	client := newAnkrHttpClient(viper.GetString(queryUrl))
	maxCap := viper.GetInt(subscribeMaxCap)
	outChan := make(chan core_types.ResultEvent, maxCap)

	timeOut := viper.GetInt(subscribeTimeOut)
	timeDuration := time.Duration(timeOut)
	subTitle := viper.GetString(subscribeTitle)
	subContent := viper.GetString(subscribeContent)
	defer  close(outChan)
	err := client.SubscribeAndWait(subTitle, subContent, timeDuration*time.Second, maxCap, outChan)
	if err != nil {
		fmt.Println("Failed to subscribe from ankr chain:", err.Error())
		return
	}
	for {
		select {
		case out := <-outChan:
			displayStruct(out)
		}
	}
}