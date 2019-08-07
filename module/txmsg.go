package module

import (
	"github.com/Ankr-network/ankr-chain/common/code"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	"github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

type ImplTxMsg interface {
	GasWanted() int64
	GasUsed() int64
	ProcessTx(tx []byte, appStore appstore.AppStore, isOnlyCheck bool) (uint32, string, []cmn.KVPair)
}

type BaseTxMsg struct {
	ImplTxMsg
}

func (b *BaseTxMsg) CheckTx(tx []byte, appStore appstore.AppStore) types.ResponseCheckTx {
	codeT, log, _ := b.ProcessTx(tx, appStore, true)
	if codeT != code.CodeTypeOK {
		return types.ResponseCheckTx{Code: codeT, Log: log}
	}

	return types.ResponseCheckTx{Code: code.CodeTypeOK, GasWanted: 1}
}

func (b *BaseTxMsg) DeliverTx(tx []byte, appStore appstore.AppStore) types.ResponseDeliverTx {
	codeT, log, tags := b.ProcessTx(tx, appStore, false)
	if codeT != code.CodeTypeOK {
		return types.ResponseDeliverTx{Code: codeT, Log: log}
	}

	return types.ResponseDeliverTx{Code: code.CodeTypeOK, GasWanted: b.GasWanted(), GasUsed: b.GasUsed(), Tags: tags}
}
