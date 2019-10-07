package tx

import (
	"fmt"
	"math/big"

	"github.com/Ankr-network/ankr-chain/account"
	"github.com/Ankr-network/ankr-chain/common"
	"github.com/Ankr-network/ankr-chain/common/code"
	ankrcrypto "github.com/Ankr-network/ankr-chain/crypto"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	"github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

type ImplTxMsg interface {
	SignerAddr() []string
	GasWanted() int64
	GasUsed() int64
	Type() string
	Bytes(txSerializer TxSerializer) []byte
	SetSecretKey(sk ankrcrypto.SecretKey)
	SecretKey() ankrcrypto.SecretKey
	PermitKey(store appstore.AppStore, pubKey []byte) bool
	ProcessTx(context ContextTx, isOnlyCheck bool) (uint32, string, []cmn.KVPair)
}

type TxFee struct {
	Amount account.Amount `json:"amount"`
	Gas    *big.Int       `json:"gas"`
}

type TxMsg struct {
	ChID        common.ChainID          `json:"chainid"`
	Nonce       uint64                  `json:"nonce"`
    Fee         TxFee                   `json:"fee"`
	GasPrice    account.Amount          `json:"gasprice"`
	Signs       []ankrcrypto.Signature  `json:"signs"`
	Memo        string                  `json:"memo"`
	Version     string                  `json:"version"`
    ImplTxMsg                           `json:"data"`
}

type txSignMsg struct {
	ChID     common.ChainID   `json:"chainid"`
	Nonce    uint64           `json:"nonce"`
	Fee      TxFee            `json:"fee"`
	GasPrice account.Amount   `json:"gasprice"`
	Memo     string           `json:"memo"`
	Version  string           `json:"version"`
	Data     []byte           `json:"data"`
}

func (ts txSignMsg) Bytes(txSerializer TxSerializer) []byte {
	msgBytes, err := txSerializer.MarshalJSON(&ts)
	if err != nil {
		panic(err)
	}

	return msgBytes
}

func (tx *TxMsg) signMsg(txSerializer TxSerializer) *txSignMsg {
	return &txSignMsg{
		ChID:     tx.ChID,
		Nonce:    tx.Nonce,
		Fee:      tx.Fee,
		GasPrice: tx.GasPrice,
		Memo:     tx.Memo,
		Version:  tx.Version,
		Data:     tx.ImplTxMsg.Bytes(txSerializer),
	}
}

func (tx *TxMsg) SignAndMarshal(txSerializer TxSerializer, key ankrcrypto.SecretKey) ([]byte, error) {
	signMsg := tx.signMsg(txSerializer)
	if signMsg != nil {
		signMsgBytes := signMsg.Bytes(txSerializer)
		signature, err := key.Sign(signMsgBytes)
		if err != nil {
			panic(err)
		}

		tx.Signs = []ankrcrypto.Signature{*signature}

		return txSerializer.Serialize(tx)
	}

	return nil, nil
}

func (tx *TxMsg) verifySignature(store appstore.AppStore, txSerializer TxSerializer) (uint32, string) {
	signMsg := tx.signMsg(txSerializer)
	toVerifyBytes := signMsg.Bytes(txSerializer)
	for i, signerAddr := range tx.SignerAddr() {
		if len(signerAddr) != ankrtypes.KeyAddressLen {
			return  code.CodeTypeInvalidAddress, fmt.Sprintf("Unexpected signer address. Got %v, len=%d", signerAddr, len(signerAddr))
		}

		if !tx.PermitKey(store, tx.Signs[i].PubKey.Bytes()) {
			return code.CodeTypeNotPermitPubKey, fmt.Sprintf("not permit public key: %v", tx.Signs[i].PubKey.Bytes())
		}

		isOk := tx.SecretKey().Verify(toVerifyBytes, &tx.Signs[i])
		if !isOk {
			return code.CodeTypeVerifySignaError, fmt.Sprintf("can't pass sign verifying for signer: pubKey=%v", tx.Signs[i].PubKey.Bytes())
		}
	}

	return code.CodeTypeOK, ""
}

func (tx *TxMsg) verifyMinGasPrice(context ContextTx) (uint32, string) {
	minGasPrice := context.MinGasPrice()
	gasPriceVal := new(big.Int).SetBytes(tx.GasPrice.Value)
	minGasPriceVal :=  new(big.Int).SetBytes(minGasPrice.Value)
	if tx.GasPrice.Cur.Symbol != minGasPrice.Cur.Symbol || gasPriceVal.Cmp(minGasPriceVal) < -1 || gasPriceVal.Cmp(minGasPriceVal) == 0{
		return code.CodeTypeGasPriceIrregular, fmt.Sprintf("irregular tx gas price: txGasSymbol=%s, txGasPriceVal=%s, minGasSymbol=%s, minGasPriceVal=%s",
			tx.GasPrice.Cur.Symbol, gasPriceVal.String(),
			minGasPrice.Cur.Symbol, minGasPriceVal.String())
	}

	return code.CodeTypeOK, ""
}

func (tx *TxMsg) BasicVerify(context ContextTx) (uint32, string) {
	codeV, log := tx.verifySignature(context.AppStore(), context.TxSerializer())
	if codeV != code.CodeTypeOK {
		return codeV, log
	}

	onceStore, err := context.AppStore().Nonce(tx.SignerAddr()[0])
	if err != nil {
		return code.CodeTypeGetStoreNonceError, err.Error()
	}
	if onceStore != tx.Nonce {
		return code.CodeTypeBadNonce, fmt.Sprintf("bad nonce: address=%s", tx.SignerAddr()[0])
	}

    return code.CodeTypeOK, ""
}

func (tx *TxMsg) CheckTx(context ContextTx) types.ResponseCheckTx {
	codeT, log := tx.BasicVerify(context)
	if codeT != code.CodeTypeOK {
		return types.ResponseCheckTx{Code: codeT, Log: log}
	}

	codeT, log = tx.verifyMinGasPrice(context)
	if codeT != code.CodeTypeOK {
		return types.ResponseCheckTx{Code: codeT, Log: log}
	}

	codeT, log, _ = tx.ProcessTx(context, true)
	if codeT != code.CodeTypeOK {
		return types.ResponseCheckTx{Code: codeT, Log: log}
	}

	return types.ResponseCheckTx{Code: code.CodeTypeOK, GasWanted: tx.GasWanted()}
}

func (tx *TxMsg) DeliverTx(context ContextTx) types.ResponseDeliverTx {
	codeT, log, tags := tx.ProcessTx(context, false)
	if codeT != code.CodeTypeOK {
		return types.ResponseDeliverTx{Code: codeT, Log: log}
	}

	return types.ResponseDeliverTx{Code: code.CodeTypeOK, GasWanted: tx.GasWanted(), GasUsed: tx.GasUsed(), Tags: tags}
}
