package client

import (
	"fmt"
	"github.com/Ankr-network/ankr-chain/tx/token"
	"math/big"
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
	c := NewClient("chain.dccn.ankr.com:26657")

	height := int64(1785347)
	rsB, _ := c.Block(&height)

	fmt.Printf("appHash=%s\n", rsB.Block.AppHash.String())

	for _, tx := range rsB.Block.Txs {
		txHash := fmt.Sprintf("%X", tx.Hash())
		fmt.Printf("txHash=%s\n", txHash)
		txEntry, _ := c.Tx(tx.Hash(), false )
		fmt.Printf("code=%d, log=%s\n",  txEntry.TxResult.Code, txEntry.TxResult.Log)
		txMsg, _ := NewTxDecoder().Decode(txEntry.Tx)
		fmt.Printf("txMsg=%v\n", txMsg)
		tf := txMsg.ImplTxMsg.(*token.TransferMsg)
		fmt.Printf("tf=%v\n", tf)
	}
}

func TestTxSearch(t *testing.T) {
	c := NewClient("chain.dccn.ankr.com:26657")
	sRS, _ := c.TxSearch("app.B508ED0D54597D516A680E7951F18CAD24C7EC9FCFCD67=1", false, 1, 100)
	for _, txRs := range sRS.Txs  {
		fmt.Printf("txHash=%s, height=%d\n", txRs.Hash.String(), txRs.Height)
		//fmt.Printf("code=%d, log=%s\n",  txRs.TxResult.Code, txRs.TxResult.Log)
		txMsg, _ := NewTxDecoder().Decode(txRs.Tx)
		//fmt.Printf("txMsg=%v\n", txMsg)
		tf := txMsg.ImplTxMsg.(*token.TransferMsg)
		fmt.Printf("tf.From=%s, tf.To=%s, tf.val=%s, gasUsed=%d, gasLimit=%d\n", tf.FromAddr, tf.ToAddr, new(big.Int).SetBytes(tf.Amounts[0].Value).String(), txRs.TxResult.GasUsed, new(big.Int).SetBytes(txMsg.GasLimit).Uint64())
	}
}
