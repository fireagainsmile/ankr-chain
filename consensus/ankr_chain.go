package ankrchain

import (
	"fmt"
	"strings"

	"github.com/Ankr-network/ankr-chain/common/code"
	act "github.com/Ankr-network/ankr-chain/module/account"
	_ "github.com/Ankr-network/ankr-chain/module/metering"
	_ "github.com/Ankr-network/ankr-chain/module/token"
	val "github.com/Ankr-network/ankr-chain/module/validator"
	"github.com/Ankr-network/ankr-chain/router"
	"github.com/Ankr-network/ankr-chain/store/appstore"
    akver "github.com/Ankr-network/ankr-chain/version"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmCoreTypes "github.com/tendermint/tendermint/types"
	"github.com/tendermint/tendermint/version"
)

var _ types.Application = (*AnkrChainApplication)(nil)

type AnkrChainApplication struct {
	appName string
	app     appstore.AppStore
	logger  log.Logger
}

func NewAnkrChainApplication(dbDir string, appName string, l log.Logger) *AnkrChainApplication {
	appStore := appstore.NewAppStore(dbDir, l.With("module", "AppStore"))

	return &AnkrChainApplication{
		appName: appName,
		app:     appStore,
		logger:  l,
	}

	router.MsgRouterInstance().SetLogger(l.With("module", "AnkrChainRouter"))
}

func (app *AnkrChainApplication) SetLogger(l log.Logger) {
	app.logger = l
}

func (app *AnkrChainApplication) Info(req types.RequestInfo) types.ResponseInfo {
	return types.ResponseInfo{
		Data:             app.appName,
		Version:          version.ABCIVersion,
		AppVersion:       akver.APPVersion,
		LastBlockHeight:  app.app.Height(),
		LastBlockAppHash: app.app.APPHash(),
	}
}

func (app *AnkrChainApplication) SetOption(req types.RequestSetOption) types.ResponseSetOption {
	return types.ResponseSetOption{}
}

// tx is either "val:pubkey/power" or "key=value" or just arbitrary bytes
func (app *AnkrChainApplication) DeliverTx(tx []byte) types.ResponseDeliverTx {
	txMsgHandler := router.MsgRouterInstance().TxMessageHandler(tx)
	if txMsgHandler != nil {
		return txMsgHandler.DeliverTx(tx, app.app)
	}

	return types.ResponseDeliverTx{
                        Code: code.CodeTypeEncodingError,
                        Log:  fmt.Sprintf("Unexpected command. Got %v", tx)}
}

func (app *AnkrChainApplication) CheckTx(tx []byte) types.ResponseCheckTx {
	txMsgHandler := router.MsgRouterInstance().TxMessageHandler(tx)
	if txMsgHandler != nil {
		return txMsgHandler.CheckTx(tx, app.app)
	}

	return types.ResponseCheckTx{
		Code: code.CodeTypeEncodingError,
		Log:  fmt.Sprintf("Unexpected. Got %v", tx)}
}

// Commit will panic if InitChain was not called
func (app *AnkrChainApplication) Commit() types.ResponseCommit {
	return app.app.Commit()
}

func (app *AnkrChainApplication) Query(reqQuery types.RequestQuery) types.ResponseQuery {
	return app.app.Query(reqQuery)
}

// Save the validators in the merkle tree
func (app *AnkrChainApplication) InitChain(req types.RequestInitChain) types.ResponseInitChain {
	var initTotalPowers int64
	for _, v := range req.Validators {
		codeT, _, _ := val.ValidatorManagerInstance().UpdateValidator(v, app.app)
		if codeT != code.CodeTypeOK {
			app.logger.Error("Error updating validators", "code", codeT)
		}

		initTotalPowers += v.Power
	}

	if initTotalPowers > tmCoreTypes.MaxTotalVotingPower {
		app.logger.Error("The init total validator powers reach max %d", "maxtotalvalidatorpower", tmCoreTypes.MaxTotalVotingPower)
		return types.ResponseInitChain{}
	}

	sbytes := string(req.AppStateBytes)
	if len(sbytes) > 0 {
		sbytes = sbytes[1 : len(sbytes)-1]
		addressAndBalance := strings.Split(sbytes, ":")
		if len(addressAndBalance) != 2 {
			app.logger.Error("Error read app states", "appstate", sbytes)
			return types.ResponseInitChain{}
		}
		addressS, balanceS := addressAndBalance[0], addressAndBalance[1]
		fmt.Println(addressS)
		fmt.Println(balanceS)
		//app.app.state.db.Set(prefixBalanceKey([]byte(addressS)), []byte(balanceS+":1"))
		//app.app.state.Size += 1
		//app.app.Commit()
	}

	act.AccountManagerInstance().InitBalance(app.app)
	val.ValidatorManagerInstance().InitValidator(app.app)

	return types.ResponseInitChain{}
}

// Track the block hash and header information
func (app *AnkrChainApplication) BeginBlock(req types.RequestBeginBlock) types.ResponseBeginBlock {
	// reset valset changes
	val.ValidatorManagerInstance().Reset()
	return types.ResponseBeginBlock{}
}

// Update the validator set
func (app *AnkrChainApplication) EndBlock(req types.RequestEndBlock) types.ResponseEndBlock {
	return types.ResponseEndBlock{ValidatorUpdates: val.ValidatorManagerInstance().ValUpdates()}
}