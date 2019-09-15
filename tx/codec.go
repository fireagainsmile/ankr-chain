package tx

import (
	"github.com/tendermint/go-amino"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
)

var TxCdc *amino.Codec

func init() {
	TxCdc = amino.NewCodec()
	cryptoAmino.RegisterAmino(TxCdc)
	TxCdc.RegisterInterface((*ImplTxMsg)(nil), nil)
	TxCdc.RegisterConcrete(&TxMsg{}, "ankr-chain/tx/txMsg", nil)
	TxCdc.RegisterConcrete(&TxMsgTesting{}, "ankr-chain/tx/txMsgTesting", nil)
}
