package types

import (
	"time"
)

type TransactionType int

const (
	TransacrionTypeSendTx TransactionType = iota
	TransacrionTypeMetering
	TransacrionTypeSetBalance
	TransacrionTypeSetStake
	TransacrionTypeValidatorTx
)

type TransactionHead struct {
	TxHash string
	TxType string
	Height int64
	Index  uint32
	Time   time.Time
}

type TransactionSendTx struct {
    TransactionHead
	FromAddress string `json:"fromAddress"`
	ToAddress   string `json:"toAddress"`
	Amount      string `json:"amount"`
}

type TransactionMetering struct {
	TransactionHead
	DC string       `json:"dc"`
	NS string       `json:"ns"`
	Value string    `json:"value"`
}

type TransactionSetBalanceTx struct {
	TransactionHead
	Address  string `json:"address"`
	Amount   string `json:"amount"`
}

type TransactionSetStakeTx struct {
	TransactionHead
	Amount string `json:"amount"`
}

type TransactionSetValidatorTx struct {
	TransactionHead
	ValidatorPubkey string  `json:"validatorPubKey"`
	Power           string   `json:"power"`
}


