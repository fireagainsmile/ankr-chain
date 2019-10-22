package tx

type TxSerializer interface {
	Serialize(txMsg *TxMsg) ([]byte, error)
	MarshalJSON(interface{}) ([]byte, error)
	Deserialize(txBytes []byte) (*TxMsg, error)
}

