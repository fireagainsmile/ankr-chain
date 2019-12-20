package serializer

import (
	"errors"
	"github.com/Ankr-network/ankr-chain/tx/cdcv0"

	"github.com/Ankr-network/ankr-chain/tx"
	"github.com/tendermint/go-amino"
)

type TxSerializerCDC struct {
	txCDCV1 *amino.Codec
	txCDCV0 *amino.Codec
}

func NewTxSerializerCDC() *TxSerializerCDC {
	return &TxSerializerCDC {CreateTxCDC(), cdcv0.CreateTxCDC()}
}

func (txs *TxSerializerCDC) Serialize(txMsg *tx.TxMsg) ([]byte, error) {
	return txs.txCDCV1.MarshalBinaryLengthPrefixed(txMsg)
}

func (txs *TxSerializerCDC) MarshalJSON(msg interface{}) ([]byte, error) {
	return txs.txCDCV1.MarshalJSON(msg)
}

func (txs *TxSerializerCDC) DeserializeCDCV1(txBytes []byte) (*tx.TxMsg, error) {
	var txMsg tx.TxMsg

	if len(txBytes) == 0 {
		return nil, errors.New("nil tx")
	}

	err := txs.txCDCV1.UnmarshalBinaryLengthPrefixed(txBytes, &txMsg)
	if err == nil {
		return &txMsg, nil
	} else {
		return nil, err
	}
}

func (txs *TxSerializerCDC) DeserializeCDCV0(txBytes []byte) (*tx.TxMsgCDCV0, error) {
	var txMsg tx.TxMsgCDCV0

	if len(txBytes) == 0 {
		return nil, errors.New("nil tx")
	}

	err := txs.txCDCV0.UnmarshalBinaryLengthPrefixed(txBytes, &txMsg)
	if err == nil {
		return &txMsg, nil
	} else {
		return nil, err
	}
}