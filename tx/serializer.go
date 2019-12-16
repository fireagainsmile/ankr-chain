package tx

type TxSerializer interface {
	Serialize(txMsg *TxMsg) ([]byte, error)
	MarshalJSON(interface{}) ([]byte, error)
	DeserializeCDCV1(txBytes []byte) (*TxMsg, error)
	DeserializeCDCV0(txBytes []byte) (*TxMsgCDCV0, error)
}
