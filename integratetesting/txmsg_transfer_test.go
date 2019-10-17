package integratetesting

import (
	"math/big"
	"testing"

	"github.com/Ankr-network/ankr-chain/client"
	ankrcmm "github.com/Ankr-network/ankr-chain/common"
	"github.com/Ankr-network/ankr-chain/crypto"
	"github.com/Ankr-network/ankr-chain/tx/serializer"
	"github.com/Ankr-network/ankr-chain/tx/token"
	"github.com/stretchr/testify/assert"
)

func TestTxTransferWithNode(t *testing.T) {
	c := client.NewClient("localhost:26657")

	msgHeader := client.TxMsgHeader{
		ChID: "test-chain-fIp7bA",
		GasLimit: new(big.Int).SetUint64(1000).Bytes(),
		GasPrice: ankrcmm.Amount{ankrcmm.Currency{"ANKR", 18}, new(big.Int).SetUint64(10000000000000).Bytes()},
		Memo: "test transfer",
		Version: "1.0",
	}

	tfMsg := &token.TransferMsg{FromAddr: "B508ED0D54597D516A680E7951F18CAD24C7EC9FCFCD67",
		ToAddr:  "454D92DC842F532683E820DF6C3784473AD9CCF222D8FB",
		Amounts: []ankrcmm.Amount{ankrcmm.Amount{ankrcmm.Currency{"ANKR", 18}, new(big.Int).SetUint64(6000000000000000000).Bytes()}},
	}

	txSerializer := serializer.NewTxSerializerCDC()

	key := crypto.NewSecretKeyEd25519("wmyZZoMedWlsPUDVCOy+TiVcrIBPcn3WJN8k5cPQgIvC8cbcR10FtdAdzIlqXQJL9hBw1i0RsVjF6Oep/06Ezg==")

	builder := client.NewTxMsgBuilder(msgHeader, tfMsg,  txSerializer, key)

	txHash, cHeight, err := builder.BuildAndCommit(c)

	assert.Equal(t, err, nil)

	t.Logf("TestTxTransferWithNode sucessful: txHash=%s, cHeight=%d", txHash, cHeight)

	resp := &ankrcmm.BalanceQueryResp{}
	c.Query("/store/balance", &ankrcmm.BalanceQueryReq{"B508ED0D54597D516A680E7951F18CAD24C7EC9FCFCD67", "ANKR"}, resp)

	t.Logf("bal=%s", resp.Amount)
}
