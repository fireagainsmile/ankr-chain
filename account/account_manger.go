package account

import (
	"sync"

	"github.com/Ankr-network/ankr-chain/common"
)

var (
	onceAM     sync.Once
	instanceAM *AccountManager
)

type AdminAccountType int
const (
	_ AdminAccountType = AdminAccountType(AccountAdminOP-1)
	AdminAccountOP
	AdminAccountValidator
	AdminAccountFound
	AdminAccountMetering
)

type AccountManager struct {
	adminAccMap map[AdminAccountType]string
}

func (am *AccountManager) GenesisAccountAddress() string {
    if common.RM == common.RunModeTesting {
		return "B508ED0D54597D516A680E7951F18CAD24C7EC9FCFCD67"
	}else if common.RM == common.RunModeProd {
		return "52E90523B5262E3AC2582F08A23068EE898D445EDF4D18"
	}

    return ""
}

func (am *AccountManager) FoundAccountAddress() string {
	if common.RM == common.RunModeTesting {
		return "64BC85F08C03F42B17EAAF5AFFAF9BFAF96CFCB85CA2F3"
	}else if common.RM == common.RunModeProd {
		return "47A65FBF3FADD12B81959AA3D8DF5E300E8C9FBFF98770"
	}

	return ""
}

func (am *AccountManager) AdminOpAccount(opType AdminAccountType) string {
	return am.adminAccMap[opType]
}

func AccountManagerInstance() *AccountManager{
	onceAM.Do(func(){
		if common.RM == common.RunModeTesting {
			instanceAM = &AccountManager{map[AdminAccountType]string{
				AdminAccountOP : "sxhP4F6OLZKNPQ2lG13WcHzitNX9++h56cppBDhwMlI=",
				AdminAccountValidator : "trwr09Y8sqIdg2H7vhJFsf4aBowBzqkMOjzAGu2ZF6E=",
				AdminAccountFound : "dBCzB+l/WYxqk+i54a4addy1XhiIK5t0IAZ5OKtegWY=",
				AdminAccountMetering : "wvHG3EddBbXQHcyJal0CS/YQcNYtEbFYxejnqf9OhM4=",

			}}
		} else if common.RM == common.RunModeProd {
			instanceAM = &AccountManager{map[AdminAccountType]string{
				AdminAccountOP : "j90knB4tx3d6xi9KefyCl2FwS/hd/jpEj+cbHdzFcqM=",
				AdminAccountValidator : "cGSgVIfAsXWbuWImGxJlNzfqruzuGA+4JXv5gfB0FyY=",
				AdminAccountFound : "sasRoTNPFzpJIHkTILaJaBnhcoC78zJk1Jy3s1/xvAE=",
				AdminAccountMetering : "cOKct2+weTftBpTvhvFKqzg9tBkN7gG/gtFVuoE53e0=",

			}}
		}
	})

	return instanceAM
}


