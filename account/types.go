package account

import "math/big"

type Code struct {
	Name      string   `json:"name"`
	CodeBytes []byte   `json:"codebytes"`
	CodeDescs string   `json:"codedescs"`
}


type Currency struct {
	Symbol  string `json:"symbol"`
	Decimal int64  `json:"decimal"`
}

type Amount struct {
	Cur   Currency `json:"currency"`
	Value *big.Int  `json:"value"`
}

type AccountType int
const (
	_ AccountType = iota
	AccountGeneral
	AccountGenesis
	AccountAdminOP
	AccountAdminValidator
	AccountAdminFound
	AccountAdminMetering
	AccountContract
)

type AccountInfo struct {
	AccType  AccountType `json:"accounttype"`
	Nonce    uint64      `json:"nonce"`
	Address  string      `json:"address"`
	PubKey   string      `json:"pubkey"`
	Amounts  []Amount    `json:"asserts"`
}

type AllowanceInfo struct {
	sender  string
	spender string
	amount  Amount
}

type ContractInfo struct {
	Address  string  `json:"address"`
	Codes    []Code  `json:"codes"`
}