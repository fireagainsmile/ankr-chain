package client

import (
	ankrcmm "github.com/Ankr-network/ankr-chain/common"
	ankrcrypto "github.com/Ankr-network/ankr-chain/crypto"
	"github.com/Ankr-network/ankr-chain/tx"
)

type TxMsgHeader struct {
	ChID     ankrcmm.ChainID
	Fee      tx.TxFee
	GasPrice ankrcmm.Amount
	Memo     string
	Version  string
}

type TxMsgBuilder struct {
	msgHeader  TxMsgHeader
	msgData    tx.ImplTxMsg
	serializer tx.TxSerializer
	key        ankrcrypto.SecretKey
}

func NewTxMsgBuilder(msgHeader TxMsgHeader, msgData tx.ImplTxMsg, serializer tx.TxSerializer, key ankrcrypto.SecretKey) *TxMsgBuilder {
	return &TxMsgBuilder{msgHeader, msgData, serializer, key}
}

func (builder *TxMsgBuilder) BuildOnly(nonce uint64) ([]byte, error) {
	txMsg := &tx.TxMsg{ChID: builder.msgHeader.ChID, Nonce: nonce, Fee: builder.msgHeader.Fee, GasPrice: builder.msgHeader.GasPrice, Memo: builder.msgHeader.Memo, Version: builder.msgHeader.Version, ImplTxMsg: builder.msgData}

	return  txMsg.SignAndMarshal(builder.serializer, builder.key)
}

func (builder *TxMsgBuilder) BuildAndCommit(c *Client) (txHash string, commitHeight int64, err error){
	signer := builder.msgData.SignerAddr()
	resp, err := c.Query("/store/account", &ankrcmm.NonceQueryReq{signer[0]})
	if err != nil {
		return "", -1, err
    }

	nonce := resp.(ankrcmm.NonceQueryResp).Nonce

	txMsg := &tx.TxMsg{ChID: builder.msgHeader.ChID, Nonce: nonce, Fee: builder.msgHeader.Fee, GasPrice: builder.msgHeader.GasPrice, Memo: builder.msgHeader.Memo, Version: builder.msgHeader.Version, ImplTxMsg: builder.msgData}

	txBytes, err := txMsg.SignAndMarshal(builder.serializer, builder.key)
	if err != nil {
		return "", -1, err
	}

	return c.BroadcastTxCommit(txBytes)
}


