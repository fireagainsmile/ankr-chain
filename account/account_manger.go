package account

import (
	"sync"

	"github.com/Ankr-network/ankr-chain/common"
	"github.com/Ankr-network/ankr-chain/store/appstore"
)

var (
	onceAM     sync.Once
	instanceAM *AccountManager
)

type AccountManager struct {
	adminAccMap map[common.AccountType]string
}

func (am *AccountManager) Init(store appstore.AppStore) {
	for k, v := range am.adminAccMap {
		store.AddAccount(v, k)
	}
}

func (am *AccountManager) GenesisAccountAddress() string {
    return am.adminAccMap[common.AccountGenesis]
}

func (am *AccountManager) FoundAccountAddress() string {
	return am.adminAccMap[common.AccountFound]
}

func (am *AccountManager) AdminOpAccount(opType common.AccountType) string {
	return am.adminAccMap[opType]
}

func AccountManagerInstance() *AccountManager{
	onceAM.Do(func(){
		if common.RM == common.RunModeTesting {
			instanceAM = &AccountManager{map[common.AccountType]string{
				common.AccountGenesis : "B508ED0D54597D516A680E7951F18CAD24C7EC9FCFCD67",
				common.AccountFound : "64BC85F08C03F42B17EAAF5AFFAF9BFAF96CFCB85CA2F3",
				common.AccountAdminOP : "sxhP4F6OLZKNPQ2lG13WcHzitNX9++h56cppBDhwMlI=",
				common.AccountAdminValidator : "trwr09Y8sqIdg2H7vhJFsf4aBowBzqkMOjzAGu2ZF6E=",
				common.AccountAdminFound : "dBCzB+l/WYxqk+i54a4addy1XhiIK5t0IAZ5OKtegWY=",
				common.AccountAdminMetering : "wvHG3EddBbXQHcyJal0CS/YQcNYtEbFYxejnqf9OhM4=",
			}}
		} else if common.RM == common.RunModeProd {
			instanceAM = &AccountManager{map[common.AccountType]string{
				common.AccountGenesis : "52E90523B5262E3AC2582F08A23068EE898D445EDF4D18",
				common.AccountFound : "47A65FBF3FADD12B81959AA3D8DF5E300E8C9FBFF98770",
				common.AccountAdminOP : "j90knB4tx3d6xi9KefyCl2FwS/hd/jpEj+cbHdzFcqM=",
				common.AccountAdminValidator : "cGSgVIfAsXWbuWImGxJlNzfqruzuGA+4JXv5gfB0FyY=",
				common.AccountAdminFound : "sasRoTNPFzpJIHkTILaJaBnhcoC78zJk1Jy3s1/xvAE=",
				common.AccountAdminMetering : "cOKct2+weTftBpTvhvFKqzg9tBkN7gG/gtFVuoE53e0=",
			}}
		}
	})

	return instanceAM
}


