package parser

type ContractState struct {
	className string
	requiredRule map[int]bool
}

func NewContractClass() *ContractState {
	return &ContractState{
		requiredRule: make(map[int]bool),
	}
}
