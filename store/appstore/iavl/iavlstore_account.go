package iavl

import (
	"fmt"
	"github.com/Ankr-network/ankr-chain/account"
	"math/big"
)

const (
	StoreAccountPrefix      = "accstore:"
	StoreAccountAllowPrefix = "accallowstore:"
	StoreValidatorPrefix    = "valstore:"
)

func containAccountPrefix(address string) string {
	return containPrefix(address, StoreAccountPrefix)
}

func containAccountAllowPrefix(addrSender string, addrSpender string, symbol string) string {
	return containPrefix(addrSender+"_"+addrSpender+"_"+symbol, StoreAccountAllowPrefix)
}

func containValidatorPrefix(address string) string {
	return containPrefix(address, StoreValidatorPrefix)
}

func stripAccountKeyPrefix(key string) (string, error) {
	return stripKeyPrefix(key, StoreAccountPrefix)
}

func stripValidatorKeyPrefix(key string) (string, error) {
	return stripKeyPrefix(key, StoreValidatorPrefix)
}

func (sp *IavlStoreApp) InitGenesisAccount() {
	addr := account.AccountManagerInstance().GenesisAccountAddress()

	totalSupply, _ := new(big.Int).SetString("100000000000000000000000000000", 10)

	var accInfo account.AccountInfo
	accInfo.Nonce   = 0
	accInfo.Address = addr
	accInfo.PubKey  = ""
	accInfo.Amounts = []account.Amount{account.Amount{account.Currency{"ANKR",18}, totalSupply}}


	sp.addAccount(&accInfo)
}

func (sp *IavlStoreApp) InitFoundAccount() {
	addr := account.AccountManagerInstance().FoundAccountAddress()

	var accInfo account.AccountInfo
	accInfo.Nonce   = 0
	accInfo.Address = addr
	accInfo.PubKey  = ""
	accInfo.Amounts = []account.Amount{{account.Currency{Symbol: "ANKR"}, new(big.Int).SetUint64(0)}}

	sp.addAccount(&accInfo)
}

func (sp *IavlStoreApp) addAccount(accInfo *account.AccountInfo) {
	if sp.iavlSM.storeMap[IavlStoreAccountKey].Has([]byte(containAccountPrefix(accInfo.Address))) {
		return
	}

	bytes := account.EncodeAccount(sp.cdc, accInfo)

	sp.iavlSM.storeMap[IavlStoreAccountKey].Set([]byte(containAccountPrefix(accInfo.Address)), bytes)
}

func (sp *IavlStoreApp) updateOnce(address string, nonce uint64) error {
	if !sp.iavlSM.storeMap[IavlStoreAccountKey].Has([]byte(containAccountPrefix(address))) {
		return fmt.Errorf("can't find the respond account from store: address=%s", address)
	}

	accBytes, _ := sp.iavlSM.storeMap[IavlStoreAccountKey].Get([]byte(address))
	accInfo := account.DecodeAccount(sp.cdc, accBytes)

	accInfo.Nonce = nonce

	bytes := account.EncodeAccount(sp.cdc, &accInfo)

	updated := sp.iavlSM.storeMap[IavlStoreAccountKey].Set([]byte(containAccountPrefix(accInfo.Address)), bytes)
	if !updated {
		return fmt.Errorf("update account's nonce fail: address=%s", address)
	}

	return nil
}

func (sp *IavlStoreApp) updatePubKey(address string, pubKey string) error {
	if !sp.iavlSM.storeMap[IavlStoreAccountKey].Has([]byte(address)) {
		return fmt.Errorf("can't find the respond account from store: address=%s", address)
	}

	accBytes, _ := sp.iavlSM.storeMap[IavlStoreAccountKey].Get([]byte(containAccountPrefix(address)))
	accInfo := account.DecodeAccount(sp.cdc, accBytes)

	accInfo.PubKey = pubKey

	bytes := account.EncodeAccount(sp.cdc, &accInfo)

	updated := sp.iavlSM.storeMap[IavlStoreAccountKey].Set([]byte(containAccountPrefix(accInfo.Address)), bytes)
	if !updated {
		return fmt.Errorf("update account's nonce fail: address=%s", address)
	}

	return nil
}

func (sp *IavlStoreApp) updateBalance(address string, assert account.Amount) error {
	if !sp.iavlSM.storeMap[IavlStoreAccountKey].Has([]byte(containAccountPrefix(address))) {
		return fmt.Errorf("can't find the respond account from store: address=%s", address)
	}

	accBytes, _ := sp.iavlSM.storeMap[IavlStoreAccountKey].Get([]byte(address))
	accInfo := account.DecodeAccount(sp.cdc, accBytes)

	findAcc := false
	for _, ass := range accInfo.Amounts {
		if ass.Cur.Symbol == assert.Cur.Symbol {
			ass.Value= assert.Value
			findAcc = true
		}
	}

	if !findAcc {
		accInfo.Amounts = append(accInfo.Amounts, assert)
	}

	bytes := account.EncodeAccount(sp.cdc, &accInfo)

	updated := sp.iavlSM.storeMap[IavlStoreAccountKey].Set([]byte(containAccountPrefix(accInfo.Address)), bytes)
	if !updated {
		return fmt.Errorf("update account's nonce fail: address=%s", address)
	}

	return nil
}

func (sp *IavlStoreApp) GetAssert(address string, symbol string) (*account.Amount, error) {
	if !sp.iavlSM.storeMap[IavlStoreAccountKey].Has([]byte(containAccountPrefix(address))) {
		return nil, fmt.Errorf("can't find the respond account from store: address=%s", address)
	}

	accBytes, _ := sp.iavlSM.storeMap[IavlStoreAccountKey].Get([]byte(containAccountPrefix(address)))
	accInfo := account.DecodeAccount(sp.cdc, accBytes)
	for _, ass := range accInfo.Amounts {
		if ass.Cur.Symbol == symbol {
			return &ass, nil
		}
	}

	return nil, fmt.Errorf("can't find the respond assert from store: symbol=%s", symbol)
}

func (sp *IavlStoreApp) Nonce(address string) (uint64, error) {
	if !sp.iavlSM.storeMap[IavlStoreAccountKey].Has([]byte(containAccountPrefix(address))) {
		return 0, fmt.Errorf("can't find the respond account from store: address=%s", address)
	}

	accBytes, _ := sp.iavlSM.storeMap[IavlStoreAccountKey].Get([]byte(containAccountPrefix(address)))
	accInfo := account.DecodeAccount(sp.cdc, accBytes)

	return accInfo.Nonce, nil
}

func (sp *IavlStoreApp) IncNonce(address string) (uint64, error) {
	if !sp.iavlSM.storeMap[IavlStoreAccountKey].Has([]byte(containAccountPrefix(address))) {
		return 0, fmt.Errorf("can't find the respond account from store: address=%s", address)
	}

	accBytes, _ := sp.iavlSM.storeMap[IavlStoreAccountKey].Get([]byte(containAccountPrefix(address)))
	accInfo := account.DecodeAccount(sp.cdc, accBytes)
	accInfo.Nonce++

	bytes := account.EncodeAccount(sp.cdc, &accInfo)

	sp.iavlSM.storeMap[IavlStoreAccountKey].Set([]byte(containAccountPrefix(accInfo.Address)), bytes)

	return accInfo.Nonce, nil
}

func (sp *IavlStoreApp) SetBalance(address string, amount account.Amount) {
	if !sp.iavlSM.IavlStore(IavlStoreAccountKey).Has([]byte(containAccountPrefix(address))) {
		var accInfo account.AccountInfo
		accInfo.Nonce   = 0
		accInfo.Address = address
		accInfo.PubKey  = ""
		accInfo.Amounts= []account.Amount{amount}

		sp.addAccount(&accInfo)
	}else {
		sp.updateBalance(address, amount)
	}
}

func (sp *IavlStoreApp) Balance(address string, symbol string) (*big.Int, error) {
	assert, err := sp.GetAssert(address, symbol)
	if err != nil {
		return nil, err
	}

	return assert.Value, nil
}

func (sp *IavlStoreApp) SetAllowance(addrSender string, addrSpender string, amount account.Amount) {
	key := containAccountAllowPrefix(addrSender, addrSpender, amount.Cur.Symbol)
	if !sp.iavlSM.IavlStore(IavlStoreAccountKey).Has([]byte(key)) {
		sp.iavlSM.IavlStore(IavlStoreAccountKey).Set([]byte(key), []byte(amount.Value.String()))
	}
}

func (sp *IavlStoreApp) Allowance(addrSender string, addrSpender string, symbol string) (*big.Int, error){
	key := containAccountAllowPrefix(addrSender, addrSpender, symbol)
	if !sp.iavlSM.IavlStore(IavlStoreAccountKey).Has([]byte(key)) {
		return nil, fmt.Errorf("IavlStoreApp Allowance not exist key: key=%s", key)
	}

	val, err := sp.iavlSM.IavlStore(IavlStoreAccountKey).Get([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("IavlStoreApp Allowance get key err: key=%s, err=%v", key, err)
	}

	rtnI, _ :=  new(big.Int).SetString(string(val), 10)

	return rtnI, nil

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
