package router

import (
	"sync"

	"github.com/Ankr-network/ankr-chain/store/appstore"
	"github.com/Ankr-network/ankr-chain/tx/decoder"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
)

var (
	onceMR     sync.Once
	instanceMR *MsgRouter
)

type TxMessageHandler interface {
	CheckTx(txMsg interface{}, appStore appstore.AppStore) types.ResponseCheckTx
	DeliverTx(txMsg interface{}, appStore appstore.AppStore) types.ResponseDeliverTx
}

type MsgRouter struct {
	routerMap map[string]TxMessageHandler
	mrLog log.Logger
}

func(mr *MsgRouter) SetLogger(mrLog log.Logger) {
	mr.mrLog = mrLog
}

func(mr *MsgRouter) AddTxMessageHandler(txMsgName string, txMsgHandler TxMessageHandler) {
	mr.routerMap[txMsgName] = txMsgHandler
}

func(mr *MsgRouter) TxMessageHandler(tx []byte) (TxMessageHandler, interface{}) {
	txType, data, err := new(decoder.TxDecoderAdapter).Decode(tx)
	if err != nil {
		if mr.mrLog != nil {
			mr.mrLog.Error("can't decode tx", "err", err)
		}

		return nil, nil
	}

	if txHandler, ok:= mr.routerMap[txType]; ok {
		return txHandler, data
	}else {
		mr.mrLog.Error("can't find the respond txmsg handler", "txType", txType)
		return nil, nil
	}
}

func MsgRouterInstance() *MsgRouter {
	onceMR.Do(func(){
		routerMap := make(map[string]TxMessageHandler)
		instanceMR = &MsgRouter{routerMap: routerMap}
	})

	return instanceMR
}










