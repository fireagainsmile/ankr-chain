package account

import (
	"github.com/Ankr-network/ankr-chain/store/appstore"
	"sync"

	"github.com/Ankr-network/ankr-chain/common"
)

var (
	onceAM     sync.Once
	instanceAM *AccountManager
)

type AdminAccountType int
const (
	_ AdminAccountType = iota//AdminAccountType(AccountAdminOP-1)
	AdminAccountOP
	AdminAccountValidator
	AdminAccountFound
	AdminAccountMetering
)

type AccountManager struct {
	adminAccMap map[AccountType]string
}

func (am *AccountManager) Init(store appstore.AppStore) {
	for k, v := range am.adminAccMap {
		store.AddAccount(v, k)
	}
}

func (am *AccountManager) GenesisAccountAddress() string {
    return am.adminAccMap[AccountGenesis]
}

func (am *AccountManager) FoundAccountAddress() string {
	return am.adminAccMap[AccountFound]
}

func (am *AccountManager) AdminOpAccount(opType AccountType) string {
	return am.adminAccMap[opType]
}

func AccountManagerInstance() *AccountManager{
	onceAM.Do(func(){
		if common.RM == common.RunModeTesting {
			instanceAM = &AccountManager{map[AccountType]string{
				AccountGenesis : "B508ED0D54597D516A680E7951F18CAD24C7EC9FCFCD67",
				AccountFound : "64BC85F08C03F42B17EAAF5AFFAF9BFAF96CFCB85CA2F3",
				AccountAdminOP : "sxhP4F6OLZKNPQ2lG13WcHzitNX9++h56cppBDhwMlI=",
				AccountAdminValidator : "trwr09Y8sqIdg2H7vhJFsf4aBowBzqkMOjzAGu2ZF6E=",
				AccountAdminFound : "dBCzB+l/WYxqk+i54a4addy1XhiIK5t0IAZ5OKtegWY=",
				AccountAdminMetering : "wvHG3EddBbXQHcyJal0CS/YQcNYtEbFYxejnqf9OhM4=",
			}}
		} else if common.RM == common.RunModeProd {
			instanceAM = &AccountManager{map[AccountType]string{
				AccountGenesis : "52E90523B5262E3AC2582F08A23068EE898D445EDF4D18",
				AccountFound : "47A65FBF3FADD12B81959AA3D8DF5E300E8C9FBFF98770",
				AccountAdminOP : "j90knB4tx3d6xi9KefyCl2FwS/hd/jpEj+cbHdzFcqM=",
				AccountAdminValidator : "cGSgVIfAsXWbuWImGxJlNzfqruzuGA+4JXv5gfB0FyY=",
				AccountAdminFound : "sasRoTNPFzpJIHkTILaJaBnhcoC78zJk1Jy3s1/xvAE=",
				AccountAdminMetering : "cOKct2+weTftBpTvhvFKqzg9tBkN7gG/gtFVuoE53e0=",
			}}
		}
	})

	return instanceAM
}


