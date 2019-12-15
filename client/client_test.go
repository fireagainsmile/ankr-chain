package client

import (
	"fmt"
	"testing"
	"time"

	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

func TestSubscribeAndWait(t *testing.T) {
	c := NewClient("localhost:26657")
	t.Run("TestCon", func(t *testing.T) {
		queryStr := fmt.Sprintf("contract.addr='%s' AND contract.method ='%s'", "0000000000000000000000000000000000000000000001", "TransferFrom")

		outCh := make(chan ctypes.ResultEvent, 100)
		c.SubscribeAndWait("test", queryStr, 30*time.Second, 100, outCh)

		for {
			select {
			 case recvData := <- outCh:
			 	fmt.Printf("Recv data: %v\n", recvData)
			}
		}
	})

}

func TestBlock(t *testing.T) {
	c := NewClient("chain-duke-test.dccn.ankr.com:26657")

	height := int64(1225774)
	rsB, _ := c.Block(&height)

	fmt.Printf("appHash=%s\n", rsB.Block.AppHash.String())

	for _, tx := range rsB.Block.Txs {
		txHash := fmt.Sprintf("%X", tx.Hash())
		fmt.Printf("txHash=%s\n", txHash)
		txEntry, _ := c.Tx(tx.Hash(), false )
		fmt.Printf("code=%d, log=%s\n",  txEntry.TxResult.Code, txEntry.TxResult.Log)
	}
}
