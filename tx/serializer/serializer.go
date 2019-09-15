package serializer

import (
	"github.com/Ankr-network/ankr-chain/tx"
)

type TxSerializer interface {
	Serialize(txMsg *tx.TxMsg) ([]byte, error)
	Deserialize(txBytes []byte) (*tx.TxMsg, error)
}

func NewTxSerializer() TxSerializer {
	return new(TxSerializerCDC)
}
