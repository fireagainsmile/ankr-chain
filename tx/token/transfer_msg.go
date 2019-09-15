package token

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/Ankr-network/ankr-chain/account"
	"github.com/Ankr-network/ankr-chain/common/code"
	ankrcrypto "github.com/Ankr-network/ankr-chain/crypto"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	"github.com/Ankr-network/ankr-chain/tx"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

const (
	MIN_TOKEN_SEND = "5000000000000000000" // 5 tokens
)

func NewTransferTxMsg() *tx.TxMsg {
	return &tx.TxMsg{ImplTxMsg: new(TransferMsg)}
}

type TransferMsg struct {
	FromAddr string           `json:"fromaddr"`
	ToAddr   string           `json:"toaddr"`
	Asserts []account.Assert  `json:"asserts"`
}

func (tr *TransferMsg) SignerAddr() []string {
	return []string {tr.FromAddr}
}

func (tr *TransferMsg) GasWanted() int64 {
	return 0
}

func (tr *TransferMsg) GasUsed() int64 {
	gasUsed, _ := strconv.ParseInt(MIN_TOKEN_SEND, 0, 64)

	return gasUsed
}

func (tr *TransferMsg) Type() string {
	return ankrtypes.TrxSendPrefix
}

func (tr *TransferMsg) Bytes() []byte {
	return nil
}
func (tr *TransferMsg) SetSecretKey(sk ankrcrypto.SecretKey) {

}

func (tr *TransferMsg) SecretKey() ankrcrypto.SecretKey {
	return nil
}

func (tr *TransferMsg) ProcessTx(appStore appstore.AppStore, isOnlyCheck bool) (uint32, string, []cmn.KVPair){
	if len(tr.FromAddr) != ankrtypes.KeyAddressLen {
		return  code.CodeTypeEncodingError, fmt.Sprintf("Unexpected from address. Got %v", tr.FromAddr), nil
	}

	if len(tr.ToAddr) != ankrtypes.KeyAddressLen {
		return  code.CodeTypeEncodingError, fmt.Sprintf("Unexpected to address. Got %v", tr.FromAddr), nil
	}

	ankrBal, err := appStore.Balance(tr.FromAddr, "ANKR")
	if err != nil {
		return code.CodeTypeBalError, err.Error(), nil
	}
	minGas, _ := new(big.Int).SetString(MIN_TOKEN_SEND, 10)
	if ankrBal.Cmp(minGas) == -1 || ankrBal.Cmp(minGas) == 0 {
		return code.CodeTypeGasNotEnough, fmt.Sprintf("not enough gas: address=%s, bal=%s", tr.FromAddr, ankrBal.String()), nil
	}

	var fromBalI *big.Int
	for _, assert := range  tr.Asserts {
		amountSend, ret := new(big.Int).SetString(assert.Amount, 10)
		if !ret {
			return code.CodeTypeEncodingError, fmt.Sprintf("Unexpected amount. Got %v", assert.Amount), nil
		} else {
			zeroN, _ := new(big.Int).SetString("0", 10)
			if amountSend.Cmp(zeroN) == -1 {
				return code.CodeTypeEncodingError, fmt.Sprintf("Unexpected amount, negative num. Got %v, assert=%s", assert.Amount, assert.Symbol), nil
			}


			if assert.Symbol == "ANKR" {
				minN, _ := new(big.Int).SetString(MIN_TOKEN_SEND, 10)
				if amountSend.Cmp(minN) == -1 || amountSend.Cmp(minN) == 0 {
					return code.CodeTypeEncodingError, fmt.Sprintf("Unexpected amount, not enough amount. Got %v", assert.Amount), nil
				}

				// check stake here. If from balance is less than stake, let it fail.
				var isSucess bool
				cstake, _ := new(big.Int).SetString("0", 10)
				value := appStore.Get([]byte(ankrtypes.AccountStakePrefix))
				if value == nil || string(value) == "" {
					// do nothing for now
				} else {
					stakeNonceSlices := strings.Split(string(value), ":")
					cstake, isSucess = new(big.Int).SetString(string(stakeNonceSlices[0]), 10)
					if !isSucess {
						return code.CodeTypeEncodingError, fmt.Sprintf("stake format error, %v", stakeNonceSlices[0]), nil
					}
				}

				fromBalI = ankrBal
				amountSendTemp, _ := new(big.Int).SetString(amountSend.String(), 10)
				amountSendTemp.Add(amountSendTemp, minGas)
				amountSendTemp.Add(amountSendTemp, cstake)
				if fromBalI.Cmp(amountSendTemp) == -1 {
					return code.CodeTypeTransferNotEnough, fmt.Sprintf("Not enough balance amount: %s, tr amount=%s, assert=%s", fromBalI.String(), amountSendTemp.String(), assert.Symbol), nil
				}
			} else {
				fromBalI, err = appStore.Balance(tr.FromAddr, assert.Symbol)
				if err != nil {
					return code.CodeTypeBalError, err.Error(), nil
				}

				if fromBalI.Cmp(amountSend) == -1 {
					return code.CodeTypeTransferNotEnough, fmt.Sprintf("Not enough balance amount: %s, tr amount=%s, assert=%s", fromBalI.String(), assert.Amount, assert.Symbol), nil
				}
			}

			if !isOnlyCheck {
				fromBalI.Sub(fromBalI, amountSend)

				toBalI, err := appStore.Balance(tr.ToAddr, assert.Symbol)
				if err != nil {
					return code.CodeTypeBalError, err.Error(), nil
				}
				toBalI.Add(toBalI, amountSend)

				fundBalI, err := appStore.Balance(account.AccountManagerInstance().FoundAccountAddress(), assert.Symbol)
				if err != nil {
					return code.CodeTypeBalError, err.Error(), nil
				}
				fundBalI.Add(toBalI, minGas)

				appStore.SetBalance(tr.FromAddr, account.Assert{assert.Symbol, fromBalI.String()})
				appStore.SetBalance(tr.ToAddr, account.Assert{assert.Symbol, toBalI.String()})
				appStore.SetBalance(account.AccountManagerInstance().FoundAccountAddress(), account.Assert{assert.Symbol, fundBalI.String()})
			}
		}
	}

	if isOnlyCheck {
		return code.CodeTypeOK, "", nil
	}

	ankrBal.Sub(ankrBal, minGas)

	tvalue := time.Now().UnixNano()
	tags := []cmn.KVPair{
		{Key: []byte("app.fromaddress"), Value: []byte(tr.FromAddr)},
		{Key: []byte("app.toaddress"), Value: []byte(tr.ToAddr)},
		{Key: []byte("app.timestamp"), Value: []byte(strconv.FormatInt(tvalue, 10))},
		{Key: []byte("app.type"), Value: []byte("Send")},
	}

	return code.CodeTypeOK, "", tags
}
