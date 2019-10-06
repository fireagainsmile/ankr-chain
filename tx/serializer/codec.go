package serializer

import (
	"github.com/Ankr-network/ankr-chain/tx"
	"github.com/Ankr-network/ankr-chain/tx/token"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/encoding/amino"
)

func CreateTxCDC() *amino.Codec {
	txCdc := amino.NewCodec()
	cryptoAmino.RegisterAmino(txCdc)
	txCdc.RegisterInterface((*tx.ImplTxMsg)(nil), nil)
	txCdc.RegisterConcrete(&tx.TxMsg{}, "ankr-chain/tx/txMsg", nil)
	txCdc.RegisterConcrete(&token.TransferMsg{}, "ankr-chain/tx/token/tranferTxMsg", nil)
	txCdc.RegisterConcrete(&tx.TxMsgTesting{}, "ankr-chain/tx/txMsgTesting", nil)

	return txCdc
}
