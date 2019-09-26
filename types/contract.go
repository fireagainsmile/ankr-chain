package types

const (
	ContractTokenStorePrefix = "TOKEN:"
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
