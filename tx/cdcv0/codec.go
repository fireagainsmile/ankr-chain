package cdcv0

import (
	"github.com/Ankr-network/ankr-chain/tx"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/encoding/amino"
)

func CreateTxCDC() *amino.Codec {
	txCdc := amino.NewCodec()
	cryptoAmino.RegisterAmino(txCdc)
	txCdc.RegisterInterface((*tx.ImplTxMsgCDCV0)(nil), nil)
	txCdc.RegisterConcrete(&tx.TxMsgCDCV0{}, "ankr-chain/tx/txMsgCDCV0", nil)
	txCdc.RegisterConcrete(&TransferMsg{}, "ankr-chain/tx/cdcv0/tranferTxMsg", nil)

	return txCdc
}
