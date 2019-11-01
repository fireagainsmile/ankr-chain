package main

import (
	"fmt"

	"github.com/Ankr-network/ankr-chain/store/historystore/db/mysql"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	ms := mysql.NewMySql("indexer:ankr1234@tcp(devnet-index.cj9yqyrk0eui.us-west-1.rds.amazonaws.com:3306)", "dccntenderminthistory", nil)
	if ms != nil{
		fmt.Sprintf("sucess")
	}

	ms.UpdateAccount("100", "5AEBA6EB8BC51DA277CCF1EF229F0C05D9535FA36CC872")
}



