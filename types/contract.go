package types

type Param struct {
	Index     uint8        `json:"index"`
	ParamType string       `json:"paramType"`
	Value     interface{} `json:"value"`
}

type ParamSort []*Param

func (p ParamSort) Len() int {return len(p)}
func (p ParamSort) Swap(i, j int) { p[i], p[j] = p[j], p[i]}
func (p ParamSort) Less(i, j int) bool {
	return p[i].Index < p[j].Index
}
