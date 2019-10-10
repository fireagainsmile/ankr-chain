package account

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/go-amino"
)

func TestEncodeAccount(t *testing.T) {
	cdc := amino.NewCodec()

	//RegisterCodec(cdc)

	var accInfo AccountInfo
	accInfo.Nonce = 1
	accInfo.Address = "5AEBA6EB8BC51DA277CCF1EF229F0C05D9535FA36CC872"
	accInfo.PubKey = "430EF1CCC9C24958779C84E6EDE748F8B743E60475D5098837B8A4971C3468FA"
	accInfo.Amounts = []Amount{{Currency{"ANKR", 18}, new(big.Int).SetUint64(100).Bytes()}}

	accBytes := EncodeAccount(cdc, &accInfo)
	accDe := DecodeAccount(cdc, accBytes)
	assert.Equal(t, accDe.Nonce, uint64(1))
	assert.Equal(t, accDe.Address, "5AEBA6EB8BC51DA277CCF1EF229F0C05D9535FA36CC872")
	assert.Equal(t, accDe.PubKey, "430EF1CCC9C24958779C84E6EDE748F8B743E60475D5098837B8A4971C3468FA")
	assert.Equal(t, accDe.Amounts[0].Cur.Symbol, "ANKR")
	assert.Equal(t, new(big.Int).SetBytes(accDe.Amounts[0].Value).String(), "100")
}
