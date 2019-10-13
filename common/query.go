package common

type BalanceQueryReq struct {
	Address string  `json:"address"`
	Symbol  string  `json:"symbol"`
}

type BalanceQueryResp struct {
	Amount string  `json:"amount"`
}

type CertKeyQueryReq struct {
	DCName string  `json:"dcname"`
	NSName string  `json:"nsname"`
}

type CertKeyQueryResp struct {
	PEMBase64 string  `json:"pembase64"`
}

type MeteringQueryReq struct {
	DCName string  `json:"dcname"`
	NSName string  `json:"nsname"`
}

type MeteringQueryResp struct {
	PEMBase64 string  `json:"pembase64"`
}

type ContractQueryReq struct {
	Address string  `json:"address"`
}

type ContractQueryResp struct {
	Addr      string   `json:"addr"`
	Name      string   `json:"name"`
	Owner     string   `json:"owneraddr"`
	Codes     []byte   `json:"codes"`
	CodesDesc string   `json:"codesdesc"`
}

type ValidatorQueryReq struct {
	ValAddr  string  `json:"address"`
}

type ValidatorQueryResp struct {
	Name         string     `json:"name"`
	ValAddress   string     `json:"valaddress"`
	PubKey       ValPubKey  `json:"pubkey"`
	Power        int64      `json:"power"`
	StakeAddress string     `json:"stakeaddress"`
	StakeAmount  Amount     `json:"stakeamount"`
	ValidHeight  uint64     `json:"validheight"`
}