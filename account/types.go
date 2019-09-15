package account

type Code struct {
	Name      string   `json:"name"`
	CodeBytes []byte   `json:"codebytes"`
	CodeDescs string   `json:"codedescs"`
}

type Assert struct {
	Symbol string    `json:"symbol"`
	Amount string    `json:"amount"`
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
	Asserts  []Assert    `json:"asserts"`
}

type ContractInfo struct {
	Address  string  `json:"address"`
	Codes    []Code  `json:"codes"`
}