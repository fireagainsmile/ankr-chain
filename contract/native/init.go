package native

import (
	"github.com/Ankr-network/ankr-chain/store/appstore"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
)

func Init(store appstore.AppStore) {
	store.SaveContract([]byte(ankrtypes.ContractTokenStorePrefix + "ANKR"), []byte{ankrtypes.ContractTypeNative})
}
