package types

import (
	"fmt"
	"testing"

	"github.com/tendermint/go-amino"
)

func TestContractInfoEncode(t *testing.T) {
	txCdc := amino.NewCodec()

	cInfo := &ContractInfo{Name: "TestCont", Codes: []byte("9012678910")/*, TxHashs: []string{"0x58902720", "0x58902723"}*/}
	cInfoBytes := EncodeContractInfo(txCdc, cInfo)

	cInfoR := DecodeContractInfo(txCdc, cInfoBytes)

	//fmt.Printf("%v, code=%s, txHashs=%v\n", cInfoR, string(cInfoR.Codes), cInfoR.TxHashs)
	fmt.Printf("%v, code=%s\n", cInfoR, string(cInfoR.Codes))
}