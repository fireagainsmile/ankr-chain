package native

import (
	"fmt"
	"math/big"

	"github.com/Ankr-network/ankr-chain/account"
	ankrcmm "github.com/Ankr-network/ankr-chain/common"
	"github.com/Ankr-network/ankr-chain/context"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	"github.com/go-interpreter/wagon/exec/gas"
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
	addrBytes := make([]byte, ankrcmm.KeyAddressLen/2)
	addrBytes[ankrcmm.KeyAddressLen/2-1] = 0x01
	codePrefixBytes := ankrcmm.GenerateContractCodePrefix(ankrcmm.ContractTypeNative, ankrcmm.ContractVMTypeUnknown, ankrcmm.ContractPatternTypeUnknown)
	store.CreateCurrency("ANKR", &ankrcmm.Currency{"ANKR", 18})
	store.BuildCurrencyCAddrMap("ANKR", string(addrBytes))
	store.SaveContract(string(addrBytes), &ankrcmm.ContractInfo{string(addrBytes), "ANKR", account.AccountManagerInstance().GenesisAccountAddress(), codePrefixBytes, ""})
	totalSup, _ := new(big.Int).SetString("10000000000000000000000000000", 10)
	store.SetBalance(account.AccountManagerInstance().GenesisAccountAddress(), ankrcmm.Amount{ankrcmm.Currency{"ANKR", 18},totalSup.Bytes()})
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

	balSender = new(big.Int).Sub(balSender, value)
	balTo     = new(big.Int).Add(balTo, value)

	stepGas := gas.GasSlowStep * 2
	ac.context.SpendGas(new(big.Int).SetUint64(stepGas))

	ac.context.SetBalance(ac.context.SenderAddr(), ankrcmm.Amount{ankrcmm.Currency{ac.symbol, 18}, balSender.Bytes()})
	ac.context.SetBalance(toAddr, ankrcmm.Amount{ankrcmm.Currency{ac.symbol, 18}, balTo.Bytes()})

	gasUsed := uint64(len(balSender.Bytes())) * gas.GasContractByte
	ac.context.SpendGas(new(big.Int).SetUint64(gasUsed))

	gasUsed = uint64(len(balTo.Bytes())) * gas.GasContractByte
	ac.context.SpendGas(new(big.Int).SetUint64(gasUsed))

	toAddrParam   := fmt.Sprintf("\"%s\"", toAddr)
	amountParam   := fmt.Sprintf("\"%s\"", amount)

	jsonArgFromat := "[{\"index\":1,\"Name\":\"toAddr\",\"ParamType\":\"string\",\"Value\":%s}," +
		              "{\"index\":2,\"Name\":\"amount\",\"ParamType\":\"string\",\"Value\":%s}]"

	jsonArg := fmt.Sprintf(jsonArgFromat, toAddrParam, amountParam)

	TrigEvent("Transfer(string, string, string))", jsonArg, ac.log, ac.context)

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

	balFrom = new(big.Int).Sub(balFrom, value)
	balTo   = new(big.Int).Add(balTo, value)

	stepGas := gas.GasSlowStep * 2
	ac.context.SpendGas(new(big.Int).SetUint64(stepGas))

	ac.context.SetBalance(ac.context.SenderAddr(), ankrcmm.Amount{ankrcmm.Currency{ac.symbol,18}, balFrom.Bytes()})
	ac.context.SetBalance(toAddr, ankrcmm.Amount{ankrcmm.Currency{ac.symbol, 18}, balTo.Bytes()})

	gasUsed := uint64(len(balFrom.Bytes())) * gas.GasContractByte
	ac.context.SpendGas(new(big.Int).SetUint64(gasUsed))

	gasUsed = uint64(len(balTo.Bytes())) * gas.GasContractByte
	ac.context.SpendGas(new(big.Int).SetUint64(gasUsed))

	fromAddrParam := fmt.Sprintf("\"%s\"", fromAddr)
	toAddrParam   := fmt.Sprintf("\"%s\"", toAddr)
	amountParam   := fmt.Sprintf("\"%s\"", amount)

	jsonArgFromat := "[{\"index\":1,\"Name\":\"fromAddr\",\"ParamType\":\"string\",\"Value\":%s}," +
		              "{\"index\":2,\"Name\":\"toAddr\",\"ParamType\":\"string\",\"Value\":%s}," +
		              "{\"index\":3,\"Name\":\"amount\",\"ParamType\":\"string\",\"Value\":%s}]"

	jsonArg := fmt.Sprintf(jsonArgFromat, fromAddrParam, toAddrParam, amountParam)

	TrigEvent("TransferFrom(string, string, string))", jsonArg, ac.log, ac.context)

	return true
}

func (ac *AnkrCoin) Approve(spenderAddr string, amount string) bool {
	value, isSucess := new(big.Int).SetString(amount, 10)
	if !isSucess || value == nil{
		ac.log.Error("AnkrCoin Approve invalid amount", "isSucess", isSucess)
	}

	ac.context.SetAllowance(ac.context.SenderAddr(), spenderAddr, ankrcmm.Amount{ankrcmm.Currency{ac.symbol, 18},value.Bytes()})

	gasUsed := uint64(len(value.Bytes())) * gas.GasContractByte
	ac.context.SpendGas(new(big.Int).SetUint64(gasUsed))

	spenderAddrParam   := fmt.Sprintf("\"%s\"", spenderAddr)
	amountParam   := fmt.Sprintf("\"%s\"", amount)

	jsonArgFromat := "[{\"index\":1,\"Name\":\"spenderAddr\",\"ParamType\":\"string\",\"Value\":%s}," +
		              "{\"index\":2,\"Name\":\"amount\",\"ParamType\":\"string\",\"Value\":%s}]"

	jsonArg := fmt.Sprintf(jsonArgFromat, spenderAddrParam, amountParam)

	TrigEvent("TransferFrom(string, string, string))", jsonArg, ac.log, ac.context)

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

	allowVal = new(big.Int).Add(allowVal, addedValue)

	stepGas := gas.GasSlowStep
	ac.context.SpendGas(new(big.Int).SetUint64(stepGas))

	ac.context.SetAllowance(ac.context.SenderAddr(), spenderAddr, ankrcmm.Amount{ankrcmm.Currency{"ANKR", 18},allowVal.Bytes()})

	gasUsed := uint64(len(allowVal.Bytes())) * gas.GasContractByte
	ac.context.SpendGas(new(big.Int).SetUint64(gasUsed))

	spenderAddrParam := fmt.Sprintf("\"%s\"", spenderAddr)
	addedAmountParam := fmt.Sprintf("\"%s\"", addedAmount)

	jsonArgFromat := "[{\"index\":1,\"Name\":\"spenderAddr\",\"ParamType\":\"string\",\"Value\":%s}," +
		"{\"index\":2,\"Name\":\"addedAmount\",\"ParamType\":\"string\",\"Value\":%s}]"

	jsonArg := fmt.Sprintf(jsonArgFromat, spenderAddrParam, addedAmountParam)

	TrigEvent("TransferFrom(string, string, string))", jsonArg, ac.log, ac.context)

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

	allowVal = new(big.Int).Sub(allowVal, subtractedValue)
	stepGas := gas.GasSlowStep
	ac.context.SpendGas(new(big.Int).SetUint64(stepGas))

	ac.context.SetAllowance(ac.context.SenderAddr(), spenderAddr, ankrcmm.Amount{ankrcmm.Currency{"ANKR", 18},allowVal.Bytes()})

	gasUsed := uint64(len(allowVal.Bytes())) * gas.GasContractByte
	ac.context.SpendGas(new(big.Int).SetUint64(gasUsed))

	spenderAddrParam      := fmt.Sprintf("\"%s\"", spenderAddr)
	subtractedAmountParam := fmt.Sprintf("\"%s\"", subtractedAmount)

	jsonArgFromat := "[{\"index\":1,\"Name\":\"spenderAddr\",\"ParamType\":\"string\",\"Value\":%s}," +
		              "{\"index\":2,\"Name\":\"subtractedAmount\",\"ParamType\":\"string\",\"Value\":%s}]"

	jsonArg := fmt.Sprintf(jsonArgFromat, spenderAddrParam, subtractedAmountParam)

	TrigEvent("TransferFrom(string, string, string))", jsonArg, ac.log, ac.context)

	return true
}