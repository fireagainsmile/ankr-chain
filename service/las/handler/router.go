package handler

import (
	"github.com/Ankr-network/ankr-chain/client"
	lascmm "github.com/Ankr-network/ankr-chain/service/las/common"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func RegisterHandler() *mux.Router {
	r := mux.NewRouter()

	nodeAddr := viper.GetString(lascmm.FlagNode)

	c := client.NewClient(nodeAddr)
	r.HandleFunc("/v1/version", QueryVersionHandler()).Methods("GET")
	r.HandleFunc("/v1/node/info", QueryNodeInfoHandler(c)).Methods("GET")
	r.HandleFunc("/v1/account/generate/{count}", GenerateAccounts(c)).Methods("GET")
	r.HandleFunc("/v1/account/nonce/{address}", QueryAccountNonceHandler(c)).Methods("GET")
	r.HandleFunc("/v1/account/balance", QueryAccountBalanceHandler(c)).Methods("GET")
	r.HandleFunc("/v1/tx/transfer", TxTransferHandler(c)).Methods("POST")
	r.HandleFunc("/v1/block/tx/transfers", QueryBlockTxTransfersHandler(c)).Methods("GET")
	r.HandleFunc("/v1/block/syncing",QueryBlockSyncing(c)).Methods("GET")

	return r
}
