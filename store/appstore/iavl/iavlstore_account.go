package iavl

import (
	"fmt"
	"math/big"

	"github.com/Ankr-network/ankr-chain/account"
	ankrcmm "github.com/Ankr-network/ankr-chain/common"
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

	var accInfo ankrcmm.AccountInfo
	accInfo.AccType = ankrcmm.AccountGenesis
	accInfo.Nonce   = 0
	accInfo.Address = addr
	accInfo.PubKey  = ""
	accInfo.Amounts = []ankrcmm.Amount{ankrcmm.Amount{ankrcmm.Currency{"ANKR",18}, totalSupply.Bytes()}}


	sp.addAccountInfo(&accInfo)
}

func (sp *IavlStoreApp) InitFoundAccount() {
	addr := account.AccountManagerInstance().FoundAccountAddress()

	var accInfo ankrcmm.AccountInfo
	accInfo.AccType = ankrcmm.AccountFound
	accInfo.Nonce   = 0
	accInfo.Address = addr
	accInfo.PubKey  = ""
	accInfo.Amounts = []ankrcmm.Amount{{ankrcmm.Currency{Symbol: "ANKR"}, new(big.Int).SetUint64(0).Bytes()}}

	sp.addAccountInfo(&accInfo)
}

func (sp *IavlStoreApp) addAccountInfo(accInfo *ankrcmm.AccountInfo) {
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

	accBytes, _ := sp.iavlSM.storeMap[IavlStoreAccountKey].Get([]byte(containAccountPrefix(address)))
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

func (sp *IavlStoreApp) updateBalance(address string, assert ankrcmm.Amount) error {
	if !sp.iavlSM.storeMap[IavlStoreAccountKey].Has([]byte(containAccountPrefix(address))) {
		return fmt.Errorf("can't find the respond account from store: address=%s", address)
	}

	accBytes, _ := sp.iavlSM.storeMap[IavlStoreAccountKey].Get([]byte(containAccountPrefix(address)))
	accInfo := account.DecodeAccount(sp.cdc, accBytes)

	findAcc := false
	for i, ass := range accInfo.Amounts {
		if ass.Cur.Symbol == assert.Cur.Symbol {
			accInfo.Amounts[i].Value = make([]byte, len(assert.Value))
			copy(accInfo.Amounts[i].Value, assert.Value)
			findAcc = true
			break
		}
	}

	if !findAcc {
		accInfo.Amounts = append(accInfo.Amounts, assert)
	}

	bytes := account.EncodeAccount(sp.cdc, &accInfo)

	updated := sp.iavlSM.storeMap[IavlStoreAccountKey].Set([]byte(containAccountPrefix(address)), bytes)
	if !updated {
		return fmt.Errorf("update account's nonce fail: address=%s", address)
	}

	return nil
}

func (sp *IavlStoreApp) GetAssert(address string, symbol string) (*ankrcmm.Amount, error) {
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

func (sp *IavlStoreApp)  NonceQuery(address string) (*ankrcmm.NonceQueryResp, error) {
	nonce, err := sp.Nonce(address)

	return &ankrcmm.NonceQueryResp{nonce}, err
}

func (sp *IavlStoreApp) Nonce(address string) (uint64, error) {
	if !sp.iavlSM.storeMap[IavlStoreAccountKey].Has([]byte(containAccountPrefix(address))) {
		sp.AddAccount(address, ankrcmm.AccountGenesis)
		return 0, nil
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

func (sp *IavlStoreApp) AddAccount(address string, accType ankrcmm.AccountType) {
	if !sp.iavlSM.IavlStore(IavlStoreAccountKey).Has([]byte(containAccountPrefix(address))) {
		var accInfo ankrcmm.AccountInfo
		accInfo.AccType = accType
		accInfo.Nonce   = 0
		accInfo.Address = address
		accInfo.PubKey  = ""
		accInfo.Amounts= []ankrcmm.Amount{{ankrcmm.Currency{"ANKR", 18}, new(big.Int).SetUint64(0).Bytes()}}

		sp.addAccountInfo(&accInfo)
	}
}

func (sp *IavlStoreApp) AccountQuery(address string) (*ankrcmm.AccountQueryResp, error) {
	if !sp.iavlSM.IavlStore(IavlStoreAccountKey).Has([]byte(containAccountPrefix(address))) {
		accBytes, _ := sp.iavlSM.storeMap[IavlStoreAccountKey].Get([]byte(containAccountPrefix(address)))
		accInfo := account.DecodeAccount(sp.cdc, accBytes)

		return &ankrcmm.AccountQueryResp{
			accInfo.AccType,
			accInfo.Nonce,
			accInfo.Address,
			accInfo.PubKey,
			accInfo.Amounts,
		}, nil
	}

	return nil, fmt.Errorf("there is no responding account info: addr=%s", address)
}

func (sp *IavlStoreApp) SetBalance(address string, amount ankrcmm.Amount) {
	if !sp.iavlSM.IavlStore(IavlStoreAccountKey).Has([]byte(containAccountPrefix(address))) {
		var accInfo ankrcmm.AccountInfo
		accInfo.AccType = ankrcmm.AccountGeneral
		accInfo.Nonce   = 0
		accInfo.Address = address
		accInfo.PubKey  = ""
		accInfo.Amounts = []ankrcmm.Amount{amount}

		sp.addAccountInfo(&accInfo)
	}else {
		sp.updateBalance(address, amount)
	}
}

func (sp *IavlStoreApp) Balance(address string, symbol string) (*big.Int, error) {
	assert, err := sp.GetAssert(address, symbol)
	if err != nil {
		return nil, err
	}

	return new(big.Int).SetBytes(assert.Value), nil
}

func (sp *IavlStoreApp) BalanceQuery(address string, symbol string) (*ankrcmm.BalanceQueryResp, error) {
	bal, err := sp.Balance(address, symbol)
	if err != nil || bal == nil {
		return nil, err
	}

	return &ankrcmm.BalanceQueryResp{bal.String()}, err
}

func (sp *IavlStoreApp) SetAllowance(addrSender string, addrSpender string, amount ankrcmm.Amount) {
	key := containAccountAllowPrefix(addrSender, addrSpender, amount.Cur.Symbol)
	if !sp.iavlSM.IavlStore(IavlStoreAccountKey).Has([]byte(key)) {
		sp.iavlSM.IavlStore(IavlStoreAccountKey).Set([]byte(key), amount.Value)
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

func (sp *IavlStoreApp) AccountList() ([]string, uint64) {
	addrCount := uint64(0)
	var addressList []string

	endBytes := prefixEndBytes([]byte(StoreAccountPrefix))

	sp.iavlSM.storeMap[IavlStoreAccountKey].tree.IterateRange([]byte(StoreAccountPrefix), endBytes, true, func(key []byte, value []byte) bool{
		if len(key) >= len(StoreAccountPrefix) && string(key[0:len(StoreAccountPrefix)]) == StoreAccountPrefix {
			accAddr, err := stripAccountKeyPrefix(string(key))
			if err != nil {
				sp.storeLog.Error("stripAccountKeyPrefix error", "err", err)
			}else {
				addrCount++
				addressList = append(addressList, accAddr)
			}
		}

		return false
	})

	if addrCount > 0 {
		addressList = addressList[1:]
		return addressList, addrCount
	}else {
		return nil, addrCount
	}
}
