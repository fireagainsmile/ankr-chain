package native

import (
	"github.com/Ankr-network/ankr-chain/store/appstore"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	"github.com/tendermint/tendermint/libs/log"
)

func Init(store appstore.AppStore, log log.Logger) {
	store.SaveContract([]byte(ankrtypes.ContractTokenStorePrefix + "ANKR"), []byte{ankrtypes.ContractTypeNative})
}
