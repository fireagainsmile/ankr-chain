package tx

import (
	"fmt"
	"math/big"

	"github.com/Ankr-network/ankr-chain/account"
	ankrcmm "github.com/Ankr-network/ankr-chain/common"
	"github.com/Ankr-network/ankr-chain/common/code"
	ankrcrypto "github.com/Ankr-network/ankr-chain/crypto"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	"github.com/Ankr-network/wagon/exec/gas"
	"github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)


type ImplTxMsgCDCV0 interface {
	SignerAddr() []string
	Type() string
	Bytes(txSerializer TxSerializer) []byte
	SetSecretKey(sk ankrcrypto.SecretKey)
	SecretKey() ankrcrypto.SecretKey
	PermitKey(store appstore.AppStore, pubKey []byte) bool
	ProcessTx(context ContextTx, metric gas.GasMetric, isOnlyCheck bool) (uint32, string, []cmn.KVPair)
}

type TxMsgCDCV0 struct {
	ChID        ankrcmm.ChainID         `json:"chainid"`
	Nonce       uint64                  `json:"nonce"`
	GasLimit    []byte                  `json:"gaslimit"`
	GasPrice    ankrcmm.Amount          `json:"gasprice"`
	GasUsed     *big.Int                `json:"gasused"`
	Signs       []ankrcrypto.Signature  `json:"signs"`
	Memo        string                  `json:"memo"`
	Version     string                  `json:"version"`
	ImplTxMsgCDCV0                       `json:"data"`
}

func (tx *TxMsgCDCV0) SpendGas(gas *big.Int) bool {
	if tx.GasUsed == nil {
		tx.GasUsed = new(big.Int).SetUint64(0)
	}

	gasUsedT := new(big.Int).SetUint64(tx.GasUsed.Uint64())
	gasUsedT = new(big.Int).Add(gasUsedT, gas)

	subGas := new(big.Int).Sub(gasUsedT, new(big.Int).SetBytes(tx.GasLimit))

	if subGas.Cmp(big.NewInt(0)) == 1 || subGas.Cmp(big.NewInt(0)) == 0 {
		return false
	}

	tx.GasUsed.SetUint64(gasUsedT.Uint64())

	return true
}

func (tx *TxMsgCDCV0) DeliverTx(context ContextTx) types.ResponseDeliverTx {
	context.AppStore().IncTotalTx()

	codeT, log, tags := tx.ProcessTx(context, tx, false)
	if codeT != code.CodeTypeOK {
		return types.ResponseDeliverTx{Code: codeT, Log: log}
	}

	if tx.GasUsed == nil || tx.GasUsed.Cmp(big.NewInt(0)) == 0 {
		return types.ResponseDeliverTx{Code: code.CodeTypeOK, GasWanted: 0, GasUsed: 0, Tags: tags}
	}

	subGas := new(big.Int).Sub(tx.GasUsed, new(big.Int).SetBytes(tx.GasLimit))
	if subGas.Cmp(big.NewInt(0)) == 1 || subGas.Cmp(big.NewInt(0)) == 0 {
		return types.ResponseDeliverTx{Code: code.CodeTypeGasNotEnough, Log: fmt.Sprintf("TxMsg DeliverTx, gas not enough, got %s", tx.GasUsed.String())}
	}

	usedFee := new(big.Int).Mul(tx.GasUsed, new(big.Int).SetBytes(tx.GasPrice.Value))
	balFrom, _, _, _, err := context.AppStore().Balance(tx.SignerAddr()[0], tx.GasPrice.Cur.Symbol, 0, false)
	if err != nil {
		return types.ResponseDeliverTx{Code: code.CodeTypeLoadBalError, Log: fmt.Sprintf("TxMsg DeliverTx, get bal err=%s， addr=%s", err.Error(), tx.SignerAddr()[0])}
	}
	if usedFee.Cmp(balFrom) == 1 || usedFee.Cmp(balFrom) == 0 {
		return types.ResponseDeliverTx{Code: code.CodeTypeFeeNotEnough, Log: fmt.Sprintf("TxMsg DeliverTx, fee not enough, got %s, expected %s", usedFee.String(), balFrom.String())}
	}

	balFrom = new(big.Int).Sub(balFrom, usedFee)

	context.AppStore().SetBalance(tx.SignerAddr()[0], ankrcmm.Amount{ankrcmm.Currency{tx.GasPrice.Cur.Symbol, 18}, balFrom.Bytes()})

	foundBal, _, _, _, err := context.AppStore().Balance(account.AccountManagerInstance().FoundAccountAddress(), tx.GasPrice.Cur.Symbol, 0, false)
	if err != nil {
		return types.ResponseDeliverTx{Code: code.CodeTypeLoadBalError, Log: fmt.Sprintf("TxMsg DeliverTx, get bal err=%s， addr=%s", err.Error(), account.AccountManagerInstance().FoundAccountAddress())}
	}
	foundBal = new(big.Int).Add(foundBal, usedFee)
	context.AppStore().SetBalance(account.AccountManagerInstance().FoundAccountAddress(), ankrcmm.Amount{ankrcmm.Currency{tx.GasPrice.Cur.Symbol, 18}, foundBal.Bytes()})

	return types.ResponseDeliverTx{Code: code.CodeTypeOK, Log: log, GasWanted: new(big.Int).SetBytes(tx.GasLimit).Int64(), GasUsed: tx.GasUsed.Int64(), Tags: tags}
}
