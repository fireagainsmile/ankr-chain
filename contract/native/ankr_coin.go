package native

import (
	"github.com/Ankr-network/ankr-chain/store/appstore"
	"math/big"

	"github.com/Ankr-network/ankr-chain/account"
	"github.com/Ankr-network/ankr-chain/context"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	"github.com/tendermint/tendermint/libs/log"
)

type AnkrCoin struct {
	name        string
	symbol      string
	decimals    uint8
	totalSupply *big.Int
	context     context.ContextContract
	log         log.Logger
}

func NewAnkrCoin(store appstore.AppStore, log log.Logger) *AnkrCoin {
	addrBytes := make([]byte, ankrtypes.KeyAddressLen/2)
	addrBytes[ankrtypes.KeyAddressLen/2-1] = 0x01
	store.BuildCurrencyCAddrMap("ANKR", string(addrBytes))
	totalSup, _ := new(big.Int).SetString("10000000000000000000000000000", 10)
	store.SetBalance(account.AccountManagerInstance().GenesisAccountAddress(), account.Amount{account.Currency{"ANKR", 18},totalSup.Bytes()})
	return &AnkrCoin{
		"Ankr Network",
		"ANKR", 18,
		totalSup,
		nil,
		log,
	}
}

func (ac *AnkrCoin) SetContextContract(context context.ContextContract) {
	ac.context = context
}

func (ac *AnkrCoin) Name() string {
	return ac.name
}

func (ac *AnkrCoin) Symbol() string {
	return ac.symbol
}

func (ac *AnkrCoin) Decimals() uint8 {
	return ac.decimals
}

func (ac *AnkrCoin) TotalSupply() *big.Int {
	return ac.totalSupply
}

func (ac *AnkrCoin) BalanceOf(addr string) *big.Int {
	bal, err := ac.context.Balance(addr, ac.symbol)
	if err != nil {
		ac.log.Error("AnkrCoin BalanceOf err", "addr", addr, "err", err)
	}

	return bal
}

func (ac *AnkrCoin) Transfer(toAddr string, amount string) bool {
	value, isSucess := new(big.Int).SetString(amount, 10)
	if !isSucess || value == nil{
		ac.log.Error("AnkrCoin Transfer invalid amount", "isSucess", isSucess)
	}

	if toAddr == "" {
		ac.log.Error("AnkrCoin Transfer toAddr blank")
		return false
	}

	balSender := ac.BalanceOf(ac.context.SenderAddr())
	if balSender == nil || balSender.Cmp(value) == -1 || balSender.Cmp(value) == 0 {
		if balSender == nil {
			ac.log.Error("AnkrCoin Transfer sender balance nil", "senderAddr", ac.context.SenderAddr())
		} else {
			ac.log.Error("AnkrCoin Transfer sender balance less than or equal to value", "senderAddr", ac.context.SenderAddr(), "balance", balSender.String(), "value", value.String())
		}

		return false
	}

	balTo := ac.BalanceOf(toAddr)
	if balTo == nil {
		balTo = new(big.Int).SetUint64(0)
	}

	balSender = balSender.Sub(balSender, value)
	balTo     = balTo.Add(balTo, value)
	ac.context.SetBalance(ac.context.SenderAddr(), account.Amount{account.Currency{ac.symbol, 18}, balSender.Bytes()})
	ac.context.SetBalance(toAddr, account.Amount{account.Currency{ac.symbol, 18}, balTo.Bytes()})

	//emit event(to do)

	return true
}

func (ac *AnkrCoin) TransferFrom(fromAddr string, toAddr string, amount string) bool {
	value, isSucess := new(big.Int).SetString(amount, 10)
	if !isSucess || value == nil{
		ac.log.Error("AnkrCoin TransferFrom invalid amount", "isSucess", isSucess)
	}

	if toAddr == "" {
		ac.log.Error("AnkrCoin TransferFrom toAddr blank")
		return false
	}

	balFrom := ac.BalanceOf(fromAddr)
	if balFrom == nil || balFrom.Cmp(value) == -1 || balFrom.Cmp(value) == 0 {
		if balFrom == nil {
			ac.log.Error("AnkrCoin TransferFrom from balance nil", "fromAddr", ac.context.SenderAddr())
		} else {
			ac.log.Error("AnkrCoin TransferFrom from balance less than or equal to value", "fromAddr", ac.context.SenderAddr(), "balance", balFrom.String(), "value", value.String())
		}

		return false
	}

	balTo := ac.BalanceOf(toAddr)
	if balTo == nil {
		balTo = new(big.Int).SetUint64(0)
	}

	balFrom = balFrom.Sub(balFrom, value)
	balTo   = balTo.Add(balTo, value)
	ac.context.SetBalance(ac.context.SenderAddr(), account.Amount{account.Currency{ac.symbol,18}, balFrom.Bytes()})
	ac.context.SetBalance(toAddr, account.Amount{account.Currency{ac.symbol, 18}, balTo.Bytes()})

	//emit event(to do)

	return true
}

func (ac *AnkrCoin) Approve(spenderAddr string, amount string) bool {
	value, isSucess := new(big.Int).SetString(amount, 10)
	if !isSucess || value == nil{
		ac.log.Error("AnkrCoin Approve invalid amount", "isSucess", isSucess)
	}

	ac.context.SetAllowance(ac.context.SenderAddr(), spenderAddr, account.Amount{account.Currency{ac.symbol, 18},value.Bytes()})

	//emit event(to do)

	return true
}

func (ac *AnkrCoin) Allowance(ownerAddr string, spenderAddr string) *big.Int {
	allowVal, err := ac.context.Allowance(ownerAddr, spenderAddr, ac.symbol)
	if err != nil {
		ac.log.Error("AnkrCoin Allowance err", "ownerAddr", ownerAddr, "spenderAddr", spenderAddr, "err", err)
	}

	return allowVal
}

func (ac *AnkrCoin) IncreaseApproval(spenderAddr string, addedAmount string) bool{
	addedValue, isSucess := new(big.Int).SetString(addedAmount, 10)
	if !isSucess || addedValue == nil{
		ac.log.Error("AnkrCoin IncreaseApproval invalid addedAmount", "isSucess", isSucess)
	}

	allowVal := ac.Allowance(ac.context.SenderAddr(), spenderAddr)
	if allowVal == nil {
		ac.log.Error("AnkrCoin IncreaseApproval sender's allowance nil")
		return false
	}

	allowVal = allowVal.Add(allowVal, addedValue)

	//emit event(to do)

	return true
}

func (ac *AnkrCoin) DecreaseApproval(spenderAddr string, subtractedAmount string) bool {
	subtractedValue, isSucess := new(big.Int).SetString(subtractedAmount, 10)
	if !isSucess || subtractedValue == nil{
		ac.log.Error("AnkrCoin DecreaseApproval invalid subtractedAmount", "isSucess", isSucess)
	}

	allowVal := ac.Allowance(ac.context.SenderAddr(), spenderAddr)
	if allowVal == nil {
		ac.log.Error("AnkrCoin IncreaseApproval sender's allowance nil")
		return false
	}

	allowVal = allowVal.Sub(allowVal, subtractedValue)

	//emit event(to do)

	return true
}