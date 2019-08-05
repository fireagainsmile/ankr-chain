package account

import (
	"sync"

	"github.com/Ankr-network/ankr-chain/common"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	"github.com/Ankr-network/ankr-chain/types"
)

var (
	onceAM     sync.Once
	instanceAM *AccountManager
)

type AccountManager struct {

}

func (am *AccountManager) InitBalance(appStore appstore.AppStore){
	appStore.Set(types.PrefixBalanceKey([]byte(am.InitAccountAddress())),
		[]byte("10000000000000000000000000000:1")) // 10000000000,000,000,000,000,000,000, 10 billion tokens
}

func (am *AccountManager) InitAccountAddress() string {
    if common.RM == common.RunModeTesting {
		return "B508ED0D54597D516A680E7951F18CAD24C7EC9FCFCD67"
	}else if common.RM == common.RunModeProd {
		return "52E90523B5262E3AC2582F08A23068EE898D445EDF4D18"
	}

    return ""
}

func (am *AccountManager) GasAccountAddress() string {
	if common.RM == common.RunModeTesting {
		return "64BC85F08C03F42B17EAAF5AFFAF9BFAF96CFCB85CA2F3"
	}else if common.RM == common.RunModeProd {
		return "47A65FBF3FADD12B81959AA3D8DF5E300E8C9FBFF98770"
	}

	return ""
}

func AccountManagerInstance() *AccountManager{
	onceAM.Do(func(){
		instanceAM = new(AccountManager)
	})

	return instanceAM
}


