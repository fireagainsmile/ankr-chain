package tx

import (
	"testing"

	ankrcmm "github.com/Ankr-network/ankr-chain/common"
	"github.com/stretchr/testify/assert"
)


func TestTxCodec(t *testing.T) {

	txansferTxMsg := &TxMsg{ImplTxMsg: &TxMsgTesting{ToAddr: "ToAddr", Asserts: []ankrcmm.Assert{{"ANKR","100"}}}}

	txansferTxMsg.ChID = "ankrchain"

	txBytes, err := TxCdc.MarshalBinaryLengthPrefixed(txansferTxMsg)
	assert.Equal(t, err, nil)

	var txansferTxMsg1 TxMsg

	err = TxCdc.UnmarshalBinaryLengthPrefixed(txBytes, &txansferTxMsg1)
	assert.Equal(t, err, nil)
	txTrMsg := txansferTxMsg1.ImplTxMsg.(*TxMsgTesting)
	assert.Equal(t, txTrMsg.Asserts[0].Symbol, "ANKR")
}
