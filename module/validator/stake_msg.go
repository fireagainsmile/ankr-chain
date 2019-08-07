package validator

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/Ankr-network/ankr-chain/common/code"
	apm "github.com/Ankr-network/ankr-chain/module"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

type StakeMsg struct {
	apm.BaseTxMsg
}

func (s *StakeMsg) GasWanted() int64 {
	return 0
}

func (s *StakeMsg) GasUsed() int64 {
	return 0
}

func (s *StakeMsg) ProcessTx(tx []byte, appStore appstore.AppStore, isOnlyCheck bool) (uint32, string, []cmn.KVPair) {
	tx = tx[len(ankrtypes.SetStakePrefix):]
	trxSetStakeSlices := strings.Split(string(tx), ":")
	if len(trxSetStakeSlices) != 4 {
		return code.CodeTypeEncodingError, fmt.Sprintf("Expected trx set balance. Got %v", trxSetStakeSlices), nil
	}

	amountS := trxSetStakeSlices[0]
	//nonceS := trxSetStakeSlices[1]

	amountSet, err := new(big.Int).SetString(string(amountS), 10)
	if !err {
		return code.CodeTypeEncodingError, fmt.Sprintf("Unexpected amount. Got %v", amountS), nil
	} else { // amountSet < 0
		zeroN, _ := new(big.Int).SetString("0", 10)
		if amountSet.Cmp(zeroN) == -1 {
			return code.CodeTypeEncodingError, fmt.Sprintf("Unexpected amount, negative num. Got %v", amountS), nil
		}
	}

	if isOnlyCheck {
		return code.CodeTypeOK, "", nil
	}

	//app.app.state.db.Set(prefixStakeKey([]byte("")), []byte(amountS +":"+ nonceS))
	//app.app.state.Size += 1

	tags := []cmn.KVPair{
		{Key: []byte("app.type"), Value: []byte("SetStake")},
	}
	return code.CodeTypeOK, "", tags
}

