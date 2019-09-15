package iavl

import (
	"fmt"
	"github.com/Ankr-network/ankr-chain/account"
	"math/big"
)

const (
	StoreAccountPrefix = "accstore:"
)

func ContainAccountPrefix(address string) string {
	return containPrefix(address, StoreAccountPrefix)
}

func stripAccountKeyPrefix(key string) (string, error) {
	return stripKeyPrefix(key, StoreAccountPrefix)
}

func (sp *IavlStoreApp) InitGenesisAccount() {
	addr := account.AccountManagerInstance().GenesisAccountAddress()

	var accInfo account.AccountInfo
	accInfo.Nonce   = 0
	accInfo.Address = addr
	accInfo.PubKey  = ""
	accInfo.Asserts = []account.Assert{{"ANKR", "10000000000000000000000000000"}}

	sp.addAccount(&accInfo)
}

func (sp *IavlStoreApp) InitFoundAccount() {
	addr := account.AccountManagerInstance().FoundAccountAddress()

	var accInfo account.AccountInfo
	accInfo.Nonce   = 0
	accInfo.Address = addr
	accInfo.PubKey  = ""
	accInfo.Asserts = []account.Assert{{"ANKR", "0"}}

	sp.addAccount(&accInfo)
}

func (sp *IavlStoreApp) addAccount(accInfo *account.AccountInfo) {
	if sp.iavlSM.storeMap[IavlStoreAccountKey].Has([]byte(ContainAccountPrefix(accInfo.Address))) {
		return
	}

	bytes := account.EncodeAccount(sp.cdc, accInfo)

	sp.iavlSM.storeMap[IavlStoreAccountKey].Set([]byte(ContainAccountPrefix(accInfo.Address)), bytes)
}

func (sp *IavlStoreApp) updateOnce(address string, nonce uint64) error {
	if !sp.iavlSM.storeMap[IavlStoreAccountKey].Has([]byte(ContainAccountPrefix(address))) {
		return fmt.Errorf("can't find the respond account from store: address=%s", address)
	}

	accBytes, _ := sp.iavlSM.storeMap[IavlStoreAccountKey].Get([]byte(address))
	accInfo := account.DecodeAccount(sp.cdc, accBytes)

	accInfo.Nonce = nonce

	bytes := account.EncodeAccount(sp.cdc, &accInfo)

	updated := sp.iavlSM.storeMap[IavlStoreAccountKey].Set([]byte(ContainAccountPrefix(accInfo.Address)), bytes)
	if !updated {
		return fmt.Errorf("update account's nonce fail: address=%s", address)
	}

	return nil
}

func (sp *IavlStoreApp) updatePubKey(address string, pubKey string) error {
	if !sp.iavlSM.storeMap[IavlStoreAccountKey].Has([]byte(address)) {
		return fmt.Errorf("can't find the respond account from store: address=%s", address)
	}

	accBytes, _ := sp.iavlSM.storeMap[IavlStoreAccountKey].Get([]byte(ContainAccountPrefix(address)))
	accInfo := account.DecodeAccount(sp.cdc, accBytes)

	accInfo.PubKey = pubKey

	bytes := account.EncodeAccount(sp.cdc, &accInfo)

	updated := sp.iavlSM.storeMap[IavlStoreAccountKey].Set([]byte(ContainAccountPrefix(accInfo.Address)), bytes)
	if !updated {
		return fmt.Errorf("update account's nonce fail: address=%s", address)
	}

	return nil
}

func (sp *IavlStoreApp) updateBalance(address string, assert account.Assert, nonce uint64) error {
	if !sp.iavlSM.storeMap[IavlStoreAccountKey].Has([]byte(ContainAccountPrefix(address))) {
		return fmt.Errorf("can't find the respond account from store: address=%s", address)
	}

	accBytes, _ := sp.iavlSM.storeMap[IavlStoreAccountKey].Get([]byte(address))
	accInfo := account.DecodeAccount(sp.cdc, accBytes)

	findAcc := false
	for _, ass := range accInfo.Asserts {
		if ass.Symbol == assert.Symbol {
			ass.Amount = assert.Amount
			findAcc = true
		}
	}

	if !findAcc {
		accInfo.Asserts = append(accInfo.Asserts, assert)
	}

	bytes := account.EncodeAccount(sp.cdc, &accInfo)

	updated := sp.iavlSM.storeMap[IavlStoreAccountKey].Set([]byte(ContainAccountPrefix(accInfo.Address)), bytes)
	if !updated {
		return fmt.Errorf("update account's nonce fail: address=%s", address)
	}

	return nil
}

func (sp *IavlStoreApp) GetAssert(address string, symbol string) (*account.Assert, error) {
	if !sp.iavlSM.storeMap[IavlStoreAccountKey].Has([]byte(ContainAccountPrefix(address))) {
		return nil, fmt.Errorf("can't find the respond account from store: address=%s", address)
	}

	accBytes, _ := sp.iavlSM.storeMap[IavlStoreAccountKey].Get([]byte(ContainAccountPrefix(address)))
	accInfo := account.DecodeAccount(sp.cdc, accBytes)
	for _, ass := range accInfo.Asserts {
		if ass.Symbol == symbol {
			return &ass, nil
		}
	}

	return nil, fmt.Errorf("can't find the respond assert from store: symbol=%s", symbol)
}

func (sp *IavlStoreApp) Nonce(address string) (uint64, error) {
	if !sp.iavlSM.storeMap[IavlStoreAccountKey].Has([]byte(ContainAccountPrefix(address))) {
		return 0, fmt.Errorf("can't find the respond account from store: address=%s", address)
	}

	accBytes, _ := sp.iavlSM.storeMap[IavlStoreAccountKey].Get([]byte(ContainAccountPrefix(address)))
	accInfo := account.DecodeAccount(sp.cdc, accBytes)


	return accInfo.Nonce, nil
}

func (sp *IavlStoreApp) SetBalance(address string, amount account.Assert, nonce uint64) {
	if !sp.iavlSM.IavlStore(IavlStoreAccountKey).Has([]byte(ContainAccountPrefix(address))) {
		var accInfo account.AccountInfo
		accInfo.Nonce   = nonce
		accInfo.Address = address
		accInfo.PubKey  = ""
		accInfo.Asserts = []account.Assert{{"ANKR", "0"}}

		sp.addAccount(&accInfo)
	}else {
		sp.updateBalance(address, amount, nonce)
	}
}

func (sp *IavlStoreApp) Balance(address string, symbol string) (*big.Int, error) {
	assert, err := sp.GetAssert(address, symbol)
	if err != nil {
		return nil, err
	}

	balInt, _ :=  new(big.Int).SetString(assert.Amount, 10)

	return balInt, nil
}

func (sp *IavlStoreApp) AccountList() ([]byte, uint64) {
	addrCount := uint64(0)
	addressList := ""

	endBytes := prefixEndBytes([]byte(StoreAccountPrefix))

	sp.iavlSM.storeMap[IavlStoreAccountKey].tree.IterateRange([]byte(StoreAccountPrefix), endBytes, true, func(key []byte, value []byte) bool{
		if len(key) >= len(StoreAccountPrefix) && string(key[0:len(StoreAccountPrefix)]) == StoreAccountPrefix {
			accAddr, err := stripAccountKeyPrefix(string(key))
			if err != nil {
				sp.storeLog.Error("stripAccountKeyPrefix error", "err", err)
			}else {
				addrCount++
				addressList += addressList + ";" + accAddr
			}
		}

		return false
	})

	if addrCount > 0 {
		addressList = addressList[1:]
		return []byte(addressList), addrCount
	}else {
		return nil, addrCount
	}
}
