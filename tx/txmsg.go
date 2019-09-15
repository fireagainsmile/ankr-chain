package tx

import (
	"fmt"
	"github.com/Ankr-network/ankr-chain/common"
	"github.com/Ankr-network/ankr-chain/common/code"
	ankrcrypto "github.com/Ankr-network/ankr-chain/crypto"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	"github.com/Ankr-network/ankr-chain/tx/serializer"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	"github.com/Workiva/go-datastructures/threadsafe/err"
	"github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

type ImplTxMsg interface {
	SignerAddr() []string
	GasWanted() int64
	GasUsed() int64
	Type() string
	Bytes() []byte
	SetSecretKey(sk ankrcrypto.SecretKey)
	SecretKey() ankrcrypto.SecretKey
	ProcessTx(appStore appstore.AppStore, isOnlyCheck bool) (uint32, string, []cmn.KVPair)
}

type TxMsg struct {
	ChID        common.ChainID          `json:"chainid"`
	Nonce       uint64                  `json:"nonce"`
    Fee         TxFee                   `json:"fee"`
	Signs       []ankrcrypto.Signature  `json:"signs"`
	Memo        string                  `json:"memo"`
    ImplTxMsg                           `json:"data"`
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

func (tx *TxMsg) verifySignature() (uint32, string) {
	txMsgT := &TxMsg{tx.ChID, tx.Nonce, tx.Fee, nil, tx.Memo, tx.ImplTxMsg}
	toVerifyBytes, err := serializer.NewTxSerializer().Serialize(txMsgT)
	if err != nil {
		return code.CodeTypeVerifySignaError, err.Error()
	}

	for i, signerAddr := range tx.SignerAddr() {
		if len(signerAddr) != ankrtypes.KeyAddressLen {
			return  code.CodeTypeInvalidAddress, fmt.Sprintf("Unexpected signer address. Got %v, len=%d", signerAddr, len(signerAddr))
		}

		addr := tx.Signs[i].PubKey.Address()
		if len(addr) != ankrtypes.KeyAddressLen {
			return  code.CodeTypeInvalidAddress, fmt.Sprintf("Unexpected signer. Got %v, addr len=%d", addr, len(addr))
		}

		isOk := tx.Signs[i].PubKey.VerifyBytes(toVerifyBytes, tx.Signs[i].Signed)
		if !isOk {
			return code.CodeTypeVerifySignaError, fmt.Sprintf("can't pass sign verifying for signer: pubKey=%s", string(sign.PubKey.Bytes()))
		}
	}

	return code.CodeTypeOK, ""

}

func (tx *TxMsg) BasicVerify(appStore appstore.AppStore) (uint32, string) {
	codeV, log := tx.verifySignature()
	if codeV != code.CodeTypeOK {
		return codeV, log
	}

	onceStore, err := appStore.Nonce(tx.SignerAddr()[0])
	if err != nil {
		return code.CodeTypeGetStoreNonceError, err.Error()
	}
	if onceStore != tx.Nonce {
		return code.CodeTypeBadNonce, fmt.Sprintf("bad nonce: address=%s", tx.SignerAddr()[0])
	}

    return code.CodeTypeOK, ""
}

func (b *TxMsg) CheckTx(appStore appstore.AppStore) types.ResponseCheckTx {
	codeT, log := b.BasicVerify(appStore)
	if codeT != code.CodeTypeOK {
		return types.ResponseCheckTx{Code: codeT, Log: log}
	}

	codeT, log, _ = b.ProcessTx(appStore, true)
	if codeT != code.CodeTypeOK {
		return types.ResponseCheckTx{Code: codeT, Log: log}
	}

	return types.ResponseCheckTx{Code: code.CodeTypeOK, GasWanted: 1}
}

func (b *TxMsg) DeliverTx(appStore appstore.AppStore) types.ResponseDeliverTx {
	codeT, log, tags := b.ProcessTx(appStore, false)
	if codeT != code.CodeTypeOK {
		return types.ResponseDeliverTx{Code: codeT, Log: log}
	}

	return types.ResponseDeliverTx{Code: code.CodeTypeOK, GasWanted: b.GasWanted(), GasUsed: b.GasUsed(), Tags: tags}
}
