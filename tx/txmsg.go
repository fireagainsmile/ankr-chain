package tx

import (
	"github.com/Ankr-network/ankr-chain/common"
	"github.com/Ankr-network/ankr-chain/common/code"
	ankrcrypto "github.com/Ankr-network/ankr-chain/crypto"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	"github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

type ImplTxMsg interface {
	GasWanted() int64
	GasUsed() int64
	Type() string
	Bytes() []byte
	SetSecretKey(sk ankrcrypto.SecretKey)
	SecretKey() ankrcrypto.SecretKey
	ProcessTx(txMsg interface{}, appStore appstore.AppStore, isOnlyCheck bool) (uint32, string, []cmn.KVPair)
}

type TxMsg struct {
	ChID  common.ChainID         `json:"chainid"`
	Nonce uint64                 `json:"nonce"`
    Fee   TxFee                  `json:"fee"`
	Sign  ankrcrypto.Signature   `json:"signature"`
	Memo  string                 `json:"memo"`
    ImplTxMsg                    `json:"data"`
}

type txSignMsg struct {
	ChID  common.ChainID   `json:"chainid"`
	Nonce uint64           `json:"nonce"`
	Fee   []byte           `json:"fee"`
	Memo  string           `json:"memo"`
	Data  []byte           `json:"data"`
}

func (ts txSignMsg) Bytes() []byte {
	msgBytes, err := TxCdc.MarshalJSON(ts)
	if err != nil {
		panic(err)
	}

	return msgBytes
}

func (tx *TxMsg) signMsg() *txSignMsg {
	return &txSignMsg{
		ChID:  tx.ChID,
		Nonce: tx.Nonce,
		Fee:   tx.Fee.Bytes(),
		Memo:  tx.Memo,
		Data:  tx.ImplTxMsg.Bytes(),
	}
}

func (tx *TxMsg) SignAndMarshal() ([]byte, error) {
	signMsg := tx.signMsg()
	if signMsg != nil {
		signMsgBytes := signMsg.Bytes()
		signature, err := tx.SecretKey().Sign(signMsgBytes)
		if err != nil {
			panic(err)
		}

		tx.Sign = *signature

		return TxCdc.MarshalBinaryLengthPrefixed(tx)
	}

	return nil, nil
}

func (tx *TxMsg) BasicVerify() types.ResponseCheckTx {

}

func (b *TxMsg) CheckTx(txMsg interface{}, appStore appstore.AppStore) types.ResponseCheckTx {
	codeT, log, _ := b.ProcessTx(txMsg, appStore, true)
	if codeT != code.CodeTypeOK {
		return types.ResponseCheckTx{Code: codeT, Log: log}
	}

	return types.ResponseCheckTx{Code: code.CodeTypeOK, GasWanted: 1}
}

func (b *TxMsg) DeliverTx(txMsg interface{}, appStore appstore.AppStore) types.ResponseDeliverTx {
	codeT, log, tags := b.ProcessTx(txMsg, appStore, false)
	if codeT != code.CodeTypeOK {
		return types.ResponseDeliverTx{Code: codeT, Log: log}
	}

	return types.ResponseDeliverTx{Code: code.CodeTypeOK, GasWanted: b.GasWanted(), GasUsed: b.GasUsed(), Tags: tags}
}
