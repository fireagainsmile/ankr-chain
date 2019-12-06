package integratetesting

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/Ankr-network/ankr-chain/client"
	ankrcmm "github.com/Ankr-network/ankr-chain/common"
	"github.com/stretchr/testify/assert"
)

func TestBalanceQueryWithProofVerify(t *testing.T) {
	c := client.NewClient("localhost:26657")

	resp := &ankrcmm.BalanceQueryResp{}
	err := c.QueryWithOption("/store/balance", 0, true, "F:\\dccntendermint\\trnode\\", &ankrcmm.BalanceQueryReq{"065E37B3FC243B9FABB1519AB876E7632C510DC9324031", "ANKR"}, resp)

	assert.Equal(t, nil, err)

	t.Logf("resp=%v", resp)
}

func TestCurrencyInfoQuery(t *testing.T) {
	c := client.NewClient("localhost:26657")

	resp := &ankrcmm.CurrencyQueryResp{}
	c.Query("/store/currency", &ankrcmm.CurrencyQueryReq{"ANKR"}, resp)

	fmt.Printf("resp=%v\n", resp)
}

func TestBigInt(t *testing.T) {
	//num1, _ := new(big.Int).SetString("0", 10)
	num1Bytes := []byte{0x00}
	fmt.Printf("num1=%s", new(big.Int).SetBytes(num1Bytes).String())
}

