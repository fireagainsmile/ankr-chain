package types

import (
	"github.com/tendermint/go-amino"
)

const (
	ContractTokenStorePrefix = "TOKEN:"
)

const (
	CodePrefixLen = 10
)

type ContractType int
const (
	_ ContractType = iota
	ContractTypeNative  = 0x01
	ContractTypeRuntime = 0x02
)

type Param struct {
	Index     int          `json:"index"`
	ParamType string       `json:"paramType"`
	Value     interface{} `json:"value"`
}

type ParamSort []*Param

func (p ParamSort) Len() int {return len(p)}
func (p ParamSort) Swap(i, j int) { p[i], p[j] = p[j], p[i]}
func (p ParamSort) Less(i, j int) bool {
	return p[i].Index < p[j].Index
}

type ContractResult struct {
	IsSuccess bool          `json:"issuccess"`
	ResultType string       `json:"resultType"`
	Value      interface{} `json:"value"`
}

type ContractInfo struct {
	Addr      string   `json:"addr"`
	Name      string   `json:"name"`
	Codes     []byte   `json:"codes"`
	CodesDesc string   `json:"codesdesc"`
	//TxHashs   []string `json:"txhashs"`
}

func EncodeContractInfo(cdc *amino.Codec, cInfo *ContractInfo) []byte {
	jsonBytes, _ := cdc.MarshalJSON(cInfo)

	return jsonBytes
}

func DecodeContractInfo(cdc *amino.Codec, bytes []byte) (cInfo ContractInfo) {
	cdc.UnmarshalJSON(bytes, &cInfo)
	return
}
