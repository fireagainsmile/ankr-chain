package validator

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/Ankr-network/ankr-chain/common/code"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	"github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

type StakeMsg struct {
	
}

func (s *StakeMsg)CheckTx(tx []byte, appStore appstore.AppStore) types.ResponseCheckTx {
	tx = tx[len(ankrtypes.SetStakePrefix):]
		trxSetStakeSlices := strings.Split(string(tx), ":")
		if len(trxSetStakeSlices) != 4 {
			return types.ResponseCheckTx{
				Code: code.CodeTypeEncodingError,
				Log:  fmt.Sprintf("Expected trx set stake. Got %v", trxSetStakeSlices)}
		}

		return types.ResponseCheckTx{Code: code.CodeTypeOK, GasWanted: 1}
}

/* this function is disabled for now */
func (s *StakeMsg) DeliverTx(tx []byte, appStore appstore.AppStore) types.ResponseDeliverTx {
	tx = tx[len(ankrtypes.SetStakePrefix):]
	trxSetStakeSlices := strings.Split(string(tx), ":")
	if len(trxSetStakeSlices) != 4 {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Expected trx set balance. Got %v", trxSetStakeSlices)}
	}

	amountS := trxSetStakeSlices[0]
	//nonceS := trxSetStakeSlices[1]

	amountSet, err := new(big.Int).SetString(string(amountS), 10)
	if !err {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected amount. Got %v", amountS)}
	} else { // amountSet < 0
		zeroN, _ := new(big.Int).SetString("0", 10)
		if amountSet.Cmp(zeroN) == -1 {
			return types.ResponseDeliverTx{
				Code: code.CodeTypeEncodingError,
				Log:  fmt.Sprintf("Unexpected amount, negative num. Got %v", amountS)}
		}
	}

	//app.app.state.db.Set(prefixStakeKey([]byte("")), []byte(amountS +":"+ nonceS))
	//app.app.state.Size += 1

	tags := []cmn.KVPair{
		{Key: []byte("app.type"), Value: []byte("SetStake")},
	}
	return types.ResponseDeliverTx{Code: code.CodeTypeOK,  GasUsed: 0, Tags: tags}
}
