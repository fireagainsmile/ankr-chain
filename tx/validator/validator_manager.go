package validator

import (
	"bytes"
	"fmt"
	"math/big"
	"sort"
	"sync"

	"github.com/Ankr-network/ankr-chain/account"
	ankrcrypto "github.com/Ankr-network/ankr-chain/crypto"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	"github.com/tendermint/tendermint/abci/types"
)

var (
	onceVM     sync.Once
	instanceVM *ValidatorManager
)

type validatorPair struct {
	validatorAddress string
	validatorInfo    *ankrtypes.ValidatorInfo
}

type valList []validatorPair

func (l valList) Swap(i, j int) {l[i], l[j] = l[j], l[i]}
func (l valList) Len() int {return len(l)}
func (l valList) Less(i, j int) bool {
	valIStakeV := new(big.Int).SetBytes(l[i].validatorInfo.StakeAmount.Value)
	valJStakeV := new(big.Int).SetBytes(l[j].validatorInfo.StakeAmount.Value)
	return valIStakeV.Cmp(valJStakeV) == -1
}

type ValidatorManager struct {
	requiredValCnt int
	valMap         map[string]*ankrtypes.ValidatorInfo
}

func (vm *ValidatorManager) Power(stakeAmount *account.Amount) int64{
	return new(big.Int).SetBytes(stakeAmount.Value).Int64()
}

func (vm *ValidatorManager) InitValidator(valUp *types.ValidatorUpdate, appStore appstore.AppStore) error {
	pubKeyHandler, err := ankrcrypto.GetValPubKeyHandler(&ankrtypes.ValPubKey{valUp.PubKey.Type, valUp.PubKey.Data})
	if err != nil {
		return fmt.Errorf("can't find the respond crypto pubkey handler:type=%s, err=%v", valUp.PubKey.Type, err)
	}

	valInfo := &ankrtypes.ValidatorInfo {
		ValAddress: pubKeyHandler.Address().String(),
	    PubKey: ankrtypes.ValPubKey{valUp.PubKey.Type, valUp.PubKey.Data},
		Power: valUp.Power,
	}

	appStore.SetValidator(valInfo)

	vm.valMap[valInfo.ValAddress] = valInfo

	return nil
}

func (vm *ValidatorManager) CreateValidator(valInfo *ankrtypes.ValidatorInfo, appStore appstore.AppStore) {
	appStore.SetValidator(valInfo)

	vm.valMap[valInfo.ValAddress] = valInfo
}

func (vm *ValidatorManager) UpdateValidator(valInfo *ankrtypes.ValidatorInfo, setFlag ankrtypes.ValidatorInfoSetFlag, appStore appstore.AppStore) error {
	valInfoS, err := appStore.Validator(valInfo.ValAddress)
	if err != nil {
		return err
	}

	if setFlag & ankrtypes.ValidatorInfoSetName == 1 {
		valInfoS.Name = valInfo.Name
	}else if setFlag & ankrtypes.ValidatorInfoSetValAddress == 1 {
		valInfoS.ValAddress = valInfo.ValAddress
	}else if setFlag & ankrtypes.ValidatorInfoSetPubKey == 1 {
		valInfoS.PubKey.Type = valInfo.PubKey.Type
		valInfoS.PubKey.Data = valInfo.PubKey.Data
	}else if setFlag & ankrtypes.ValidatorInfoSetStakeAddress == 1 {
		valInfoS.StakeAddress = valInfo.StakeAddress
	}else if setFlag & ankrtypes.ValidatorInfoSetStakeAmount == 1 {
		if valInfoS.StakeAmount.Cur.Symbol != valInfo.StakeAmount.Cur.Symbol {
			return fmt.Errorf("ValidatorManager UpdateValidator currency conflict: orignal symbol=%s, set symbol=%s", valInfoS.StakeAmount.Cur.Symbol, valInfo.StakeAmount.Cur.Symbol)
		}
		valInfoS.StakeAmount.Cur   = valInfo.StakeAmount.Cur
		valInfoS.StakeAmount.Value = valInfo.StakeAmount.Value
	} else if setFlag & ankrtypes.ValidatorInfoSetValidHeight == 1 {
		valInfoS.ValidHeight = valInfo.ValidHeight
	} else {
		return fmt.Errorf("ValidatorManager UpdateValidator invalid set setFlag=%s", setFlag)
	}

	appStore.SetValidator(valInfoS)

	vm.valMap[valInfoS.ValAddress] = valInfoS

	return nil
}

func (vm *ValidatorManager) RemoveValidator(valAddr string, appStore appstore.AppStore) {
	appStore.RemoveValidator(valAddr)
	delete(vm.valMap, valAddr)
}

func (vm *ValidatorManager) Reset() {
	vm.valMap = make(map[string]*ankrtypes.ValidatorInfo)
}

func (vm *ValidatorManager) ValBeginBlock(req types.RequestBeginBlock, appStore appstore.AppStore) {
	for _, byzantineVal := range req.ByzantineValidators {
		appStore.RemoveValidator(string(byzantineVal.Validator.Address))
		delete(vm.valMap, string(byzantineVal.Validator.Address))
	}
}

func (vm *ValidatorManager) ValUpdates() []types.ValidatorUpdate {
	var requiredVals valList
	for k, v := range vm.valMap {
		if len(requiredVals) < vm.requiredValCnt {
			requiredVals = append(requiredVals, validatorPair{k, v})
		} else {
			sort.Sort(requiredVals)
			valInfo := requiredVals[len(requiredVals) -1]
			if new(big.Int).SetBytes(v.StakeAmount.Value).Cmp(new(big.Int).SetBytes(valInfo.validatorInfo.StakeAmount.Value)) < -1 {
				requiredVals[len(requiredVals) -1] = validatorPair{k, v}
			}
		}
	}

	vp := make([]types.ValidatorUpdate, len(requiredVals))
	for k, valPair := range requiredVals {
		vp[k].PubKey.Type = valPair.validatorInfo.PubKey.Type
		vp[k].PubKey.Data = valPair.validatorInfo.PubKey.Data
		vp[k].Power       = valPair.validatorInfo.Power
	}

	return vp
}

func (vm *ValidatorManager) Validators(appStore appstore.AppStore) (validators []types.ValidatorUpdate) {
	itr := appStore.DB().Iterator(nil, nil)
	for ; itr.Valid(); itr.Next() {
		if isValidatorTx(itr.Key()) {
			validator := new(types.ValidatorUpdate)
			err := types.ReadMessage(bytes.NewBuffer(itr.Value()), validator)
			if err != nil {
				panic(err)
			}
			validators = append(validators, *validator)
		}
	}
	return
}

func (vm *ValidatorManager) TotalValidatorPowers(appStore appstore.AppStore) int64 {
	return appStore.TotalValidatorPowers()
}

func ValidatorManagerInstance() *ValidatorManager {
	onceVM.Do(func(){
		instanceVM = &ValidatorManager{4, make(map[string]*ankrtypes.ValidatorInfo)}
	})

	return instanceVM
}

