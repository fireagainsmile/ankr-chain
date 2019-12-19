package commands

import (
	"net/http"
	"os"

	"github.com/Ankr-network/ankr-chain/client"
	lascmm "github.com/Ankr-network/ankr-chain/service/las/common"
	lashandler "github.com/Ankr-network/ankr-chain/service/las/handler"
	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	tmserver "github.com/tendermint/tendermint/rpc/lib/server"
)

func Start() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start",
		Short:   "Start ankrchain-las (ankrchain lite access service), a local REST server with swagger-ui: http://localhost:32669/swagger-ui/",
		Example: "ankrchain-las start --chain-id=<chain-id> --proof-verify --node=<node connection address, such as tcp://127.0.0.1:26657, http://127.0.0.1:26657 and https://127.0.0.1:443>",
		RunE: func(cmd *cobra.Command, args []string) error {
			listenAddr := viper.GetString(lascmm.FlagListenAddr)
			router := registerHandler()

			statikFS, err := fs.New()
			if err != nil {
				panic(err)
			}

			staticServer := http.FileServer(statikFS)
			router.PathPrefix("/swagger-ui/").Handler(http.StripPrefix("/swagger-ui/", staticServer))

			logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "ankrchain-las")
			maxOpen := viper.GetInt(lascmm.FlagMaxOpenConnections)

			config := &tmserver.Config{MaxOpenConnections: maxOpen}

			listener, err := tmserver.Listen(
				listenAddr,
				config,
			)
			if err != nil {
				return err
			}

			logger.Info("Starting ankrchain-las service...")

			err = tmserver.StartHTTPServer(listener, router, logger, config)
			if err != nil {
				return err
			}

			logger.Info("ankrchain-las server started")

			// wait forever and cleanup
			cmn.TrapSignal(logger, func() {
				err := listener.Close()
				logger.Error("error closing listener", "err", err)
			})

			return nil
		},
	}

	cmd.Flags().String(lascmm.FlagListenAddr, "tcp://localhost:32669", "The address for the server to listen on")
	cmd.Flags().String(lascmm.FlagCORS, "", "Set the domains that can make CORS requests (* for all)")
	cmd.Flags().String(lascmm.FlagChainID, "", "Chain ID of ankrchain node")
	cmd.Flags().String(lascmm.FlagNode, "tcp://127.0.0.1:26657", "The node connection address")
	cmd.Flags().Int(lascmm.FlagMaxOpenConnections, 1000, "The number of maximum open connections")
	cmd.Flags().Bool(lascmm.FlagProofVerify, false, "Needn't verify proofs of responses")

	return cmd
}

func registerHandler() *mux.Router {
	r := mux.NewRouter()

	nodeAddr := viper.GetString(lascmm.FlagNode)

	c := client.NewClient(nodeAddr)
	r.HandleFunc("/v1/version", lashandler.QueryVersionHandler()).Methods("GET")
	r.HandleFunc("/v1/node/info", lashandler.QueryNodeInfoHandler(c)).Methods("GET")
	r.HandleFunc("/v1/account/generate/{count}", lashandler.GenerateAccounts(c)).Methods("GET")
	r.HandleFunc("/v1/account/nonce/{address}", lashandler.QueryAccountNonceHandler(c)).Methods("GET")
	r.HandleFunc("/v1/account/balance", lashandler.QueryAccountBalanceHandler(c)).Methods("GET")
	r.HandleFunc("/v1/tx/transfer", lashandler.TxTransferHandler(c)).Methods("POST")
	r.HandleFunc("/v1/block/tx/transfers", lashandler.QueryBlockTxTransfersHandler(c)).Methods("GET")
	r.HandleFunc("/v1/block/syncing", lashandler.QueryBlockSyncing(c)).Methods("GET")

	return r
}
