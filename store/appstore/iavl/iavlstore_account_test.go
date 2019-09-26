package iavl

import (
	"math/big"
	"testing"

	"github.com/Ankr-network/ankr-chain/account"
	"github.com/stretchr/testify/assert"
)

func TestIavlStoreAccountCount(t *testing.T) {
	storeApp := NewMockIavlStoreApp()

	var accInfo1 account.AccountInfo
	accInfo1.Nonce = 1
	accInfo1.Address = "5AEBA6EB8BC51DA277CCF1EF229F0C05D9535FA36CC872"
	accInfo1.PubKey = "430EF1CCC9C24958779C84E6EDE748F8B743E60475D5098837B8A4971C3468FA"
	accInfo1.Amounts = []account.Amount{{account.Currency{Symbol: "ANKR"}, new(big.Int).SetUint64(100)}}

	var accInfo2 account.AccountInfo
	accInfo2.Nonce = 1
	accInfo2.Address = "8AEBA6EB8BC51DA277CCF1EF229F0C05D9535FA36CC873"
	accInfo2.PubKey = "930EF1CCC9C24958779C84E6EDE748F8B743E60475D5098837B8A4971C3468FB"
	accInfo2.Amounts = []account.Amount{{account.Currency{Symbol: "ANKR"}, new(big.Int).SetUint64(100)}}

	storeApp.addAccount(&accInfo1)
	storeApp.addAccount(&accInfo2)

	_, accCnt := storeApp.AccountList()

	assert.Equal(t, accCnt, uint64(2))
}
