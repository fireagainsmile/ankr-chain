package main

import (
	"fmt"
	"github.com/Ankr-network/ankr-chain/tool/cli/cmd"
)

func main()  {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}