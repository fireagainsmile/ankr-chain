package handler

type GasPriceLas struct {
	Symbol  string `json:"symbol"`
	Value   string `json:"value"`
}

type AmountLas struct {
	Symbol  string  `json:"symbol"`
	Value   string  `json:"value"`
}

type TxMsgHeaderLas struct {
	GasLimit uint32           `json:"gasLimit"`
	GasPrice GasPriceLas      `json:"gasPrice"`
	Memo     string           `json:"memo"`
}

type TxMsgTransferDataLas struct {
	PriKey   string     `json:"priKey"`
	FromAddr string     `json:"fromAddr"`
	ToAddr   string     `json:"toAddr"`
	Amount   AmountLas  `json:"amount"`
}

type TxMsgTransferInfo struct {
	Header TxMsgHeaderLas        `json:"header"`
	Data   TxMsgTransferDataLas  `json:"data"`
}

type TxCommitResult struct {
	TxHash    string  `json:"hash"`
	Height    uint32  `json:"height"`
	Status    string  `json:"status"`
	GasUsed   uint32  `json:"gasUsed"`
	TimeStamp string  `json:"timeStamp"`
	Log       string  `json:"log"`
	Err       string  `json:"err"`
}

type TranserResultItem struct {
	TxHash    string       `json:"hash"`
	Height    uint32       `json:"height"`
	FromAddr  string       `json:"fromAddr"`
	ToAddr    string       `json:"toAddr"`
	Amounts   []*AmountLas `json:"amounts"`
	Status    string       `json:"status"`
	GasLimit  uint32       `json:"gasLimit"`
	GasPrice  GasPriceLas  `json:"gasPrice"`
	GasUsed   uint32       `json:"gasUsed"`
	TimeStamp string       `json:"timeStamp"`
	Memo      string       `json:"memo"`
}

type TranserResultsOfOneBlock struct {
	TotalTx         uint32               `json:"totalTxs"`
	TransferResults []*TranserResultItem `json:"transferResults"`
}

type blockSyncing struct {
	Syncing bool  `json:"syncing"`
}

type AccountInfoItem struct {
	PriKey  string  `json:"priKey"`
	PubKey  string  `json:"pubKey"`
	Address string  `json:"address"`
}

type NodeInfo struct {
	ID            string  `json:"id"`
	ListenAddr    string  `json:"listenAddr"`
	ChainID       string  `json:"chainID"`
	Moniker       string  `json:"moniker"`
	TxIndex       string  `json:"txIndex"`
	RPCAddr       string  `json:"rpcAddr"`
	LatestHeight  uint32  `json:"latestHeight"`
	Version       string  `json:"version"`
}

