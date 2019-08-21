package validator

import (
	"fmt"

	"github.com/Ankr-network/ankr-chain/common/code"
	ankrcrypto "github.com/Ankr-network/ankr-chain/crypto"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	apm "github.com/Ankr-network/ankr-chain/tx"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	"math/big"
)

type StakeMsg struct {
	apm.TxMsg
}

func (s *StakeMsg) GasWanted() int64 {
	return 0
}

func (s *StakeMsg) GasUsed() int64 {
	return 0
}

func (s *StakeMsg) Type() string {
	return ankrtypes.TrxSendPrefix
}

func (s *StakeMsg) Bytes() []byte {
	return nil
}
func (s *StakeMsg) SetSecretKey(sk ankrcrypto.SecretKey) {

}

func (s *StakeMsg) SecretKey() ankrcrypto.SecretKey {
	return nil
}

func (s *StakeMsg) ProcessTx(txMsg interface{}, appStore appstore.AppStore, isOnlyCheck bool) (uint32, string, []cmn.KVPair) {
	trxSetStakeSlices, ok := txMsg.([]string)
	if !ok {
		return  code.CodeTypeEncodingError, fmt.Sprintf("invalid tx set op msg"), nil
	}

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

