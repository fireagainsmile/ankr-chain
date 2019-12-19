package client

import (
	"fmt"

	ankrcmm "github.com/Ankr-network/ankr-chain/common"
	"github.com/Ankr-network/ankr-chain/common/code"
	ankrcrypto "github.com/Ankr-network/ankr-chain/crypto"
	"github.com/Ankr-network/ankr-chain/tx"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

type TxMsgHeader struct {
	ChID      ankrcmm.ChainID
	GasLimit  []byte
	GasPrice  ankrcmm.Amount
	Memo      string
	Version   string
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
	txMsg := &tx.TxMsg{ChID: builder.msgHeader.ChID, Nonce: nonce, GasLimit: builder.msgHeader.GasLimit, GasPrice: builder.msgHeader.GasPrice, Memo: builder.msgHeader.Memo, Version: builder.msgHeader.Version, ImplTxMsg: builder.msgData}

	return  txMsg.SignAndMarshal(builder.serializer, builder.key)
}

func (builder *TxMsgBuilder) BuildAndCommitWithRawResult(c *Client) (*ctypes.ResultBroadcastTxCommit, error){
	signer := builder.msgData.SignerAddr()
	resp := &ankrcmm.NonceQueryResp{}
	err := c.Query("/store/nonce", &ankrcmm.NonceQueryReq{signer[0]}, resp)
	if err != nil {
		return nil, err
	}

	nonce := resp.Nonce

	txMsg := &tx.TxMsg{ChID: builder.msgHeader.ChID, Nonce: nonce, GasLimit: builder.msgHeader.GasLimit, GasPrice: builder.msgHeader.GasPrice, Memo: builder.msgHeader.Memo, Version: builder.msgHeader.Version, ImplTxMsg: builder.msgData}

	txBytes, err := txMsg.SignAndMarshal(builder.serializer, builder.key)
	if err != nil {
		return nil, err
	}

	return c.BroadcastTxCommitWithRawResult(txBytes)
}

func (builder *TxMsgBuilder) BuildAndCommit(c *Client) (txHash string, commitHeight int64, log string, err error){
	signer := builder.msgData.SignerAddr()
	resp := &ankrcmm.NonceQueryResp{}
	err = c.Query("/store/nonce", &ankrcmm.NonceQueryReq{signer[0]}, resp)
	if err != nil {
		return "", -1, "", err
    }

	nonce := resp.Nonce

	txMsg := &tx.TxMsg{ChID: builder.msgHeader.ChID, Nonce: nonce, GasLimit: builder.msgHeader.GasLimit, GasPrice: builder.msgHeader.GasPrice, Memo: builder.msgHeader.Memo, Version: builder.msgHeader.Version, ImplTxMsg: builder.msgData}

	txBytes, err := txMsg.SignAndMarshal(builder.serializer, builder.key)
	if err != nil {
		return "", -1, "", err
	}

	return c.BroadcastTxCommit(txBytes)
}

func (builder *TxMsgBuilder) BuildAndBroadcastSync(c *Client) (data []byte, txHash string, log string, err error){
	signer := builder.msgData.SignerAddr()
	resp := &ankrcmm.NonceQueryResp{}
	err = c.Query("/store/nonce", &ankrcmm.NonceQueryReq{signer[0]}, resp)
	if err != nil {
		return nil, "",  "", err
	}

	nonce := resp.Nonce

	txMsg := &tx.TxMsg{ChID: builder.msgHeader.ChID, Nonce: nonce, GasLimit: builder.msgHeader.GasLimit, GasPrice: builder.msgHeader.GasPrice, Memo: builder.msgHeader.Memo, Version: builder.msgHeader.Version, ImplTxMsg: builder.msgData}

	txBytes, err := txMsg.SignAndMarshal(builder.serializer, builder.key)
	if err != nil {
		return nil, "",  "", err
	}

	btxResult, err := c.BroadcastTxSync(txBytes)
	if err != nil {
		return nil, "",  "", err
	}

	if btxResult.Code != code.CodeTypeOK {
		return nil, "",  "", fmt.Errorf("BroadcastTxSync err: %d, %w", btxResult.Code, err)
	}

	return btxResult.Data.Bytes(), btxResult.Hash.String(), btxResult.Log, nil
}

func (builder *TxMsgBuilder) BuildAndBroadcastAsync(c *Client) (data []byte, txHash string, log string, err error){
	signer := builder.msgData.SignerAddr()
	resp := &ankrcmm.NonceQueryResp{}
	err = c.Query("/store/nonce", &ankrcmm.NonceQueryReq{signer[0]}, resp)
	if err != nil {
		return nil, "",  "", err
	}

	nonce := resp.Nonce

	txMsg := &tx.TxMsg{ChID: builder.msgHeader.ChID, Nonce: nonce, GasLimit: builder.msgHeader.GasLimit, GasPrice: builder.msgHeader.GasPrice, Memo: builder.msgHeader.Memo, Version: builder.msgHeader.Version, ImplTxMsg: builder.msgData}

	txBytes, err := txMsg.SignAndMarshal(builder.serializer, builder.key)
	if err != nil {
		return nil, "",  "", err
	}

	btxResult, err := c.BroadcastTxAsync(txBytes)
	if err != nil {
		return nil, "",  "", err
	}

	if btxResult.Code != code.CodeTypeOK {
		return nil, "",  "", fmt.Errorf("BroadcastTxAsync err: %d, %w", btxResult.Code, err)
	}

	return btxResult.Data.Bytes(), btxResult.Hash.String(), btxResult.Log, nil
}


