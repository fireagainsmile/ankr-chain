package validator

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/Ankr-network/ankr-chain/common/code"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	"github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

var (
	onceVM     sync.Once
	instanceVM *ValidatorManager
)

type ValidatorManager struct {
	valUpdates []types.ValidatorUpdate
}

func (vm *ValidatorManager) InitValidator(appStore appstore.AppStore) {
	appStore.Set([]byte(ankrtypes.SET_VAL_NONCE), []byte("0"))
	//appStore.Set(PrefixStakeKey([]byte("")), []byte("0:1"))
	value := []byte("")
	value = appStore.Get([]byte(ankrtypes.AccountStakePrefix))
	if value == nil || string(value) == "" {
		appStore.Set(PrefixStakeKey([]byte("")), []byte("0:1"))
	}
}

// add, update, or remove a validator
func (v *ValidatorManager) UpdateValidator(valUP types.ValidatorUpdate, appStore appstore.AppStore) types.ResponseDeliverTx {
	key := []byte("val:" + string(valUP.PubKey.Data))
	if valUP.Power == 0 {
		// remove validator
		if !appStore.Has(key) {
			return types.ResponseDeliverTx{
				Code: code.CodeTypeUnauthorized,
				Log:  fmt.Sprintf("Cannot remove non-existent validator %X", key)}
		}
		appStore.Delete(key)
	} else {
		// add or update validator
		value := bytes.NewBuffer(make([]byte, 0))
		if err := types.WriteMessage(&valUP, value); err != nil {
			return types.ResponseDeliverTx{
				Code: code.CodeTypeEncodingError,
				Log:  fmt.Sprintf("Error encoding validator: %v", err)}
		}
		appStore.Set(key, value.Bytes())
	}

	// we only update the changes array if we successfully updated the tree
	ValidatorManagerInstance().valUpdates = append(ValidatorManagerInstance().valUpdates, valUP)

	tags := []cmn.KVPair{
		{Key: []byte("app.type"), Value: []byte("UpdateValidator")},
	}
	return types.ResponseDeliverTx{Code: code.CodeTypeOK, GasUsed: 0, Tags: tags}
}

func (vm *ValidatorManager) Reset() {
	vm.valUpdates = make([]types.ValidatorUpdate, 0)
}

func (vm *ValidatorManager) ValUpdates() []types.ValidatorUpdate {
	return vm.valUpdates
}

func ValidatorManagerInstance() *ValidatorManager {
	onceVM.Do(func(){
		instanceVM = &ValidatorManager{make([]types.ValidatorUpdate, 0)}
	})

	return instanceVM
}

