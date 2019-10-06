package router

/*
import (
	"sync"

	"github.com/Ankr-network/ankr-chain/context"
	"github.com/Ankr-network/ankr-chain/tx/serializer"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
)

var (
	onceMR     sync.Once
	instanceMR *MsgRouter
)

type TxMessageHandler interface {
	CheckTx(context context.ContextTx) types.ResponseCheckTx
	DeliverTx(context context.ContextTx) types.ResponseDeliverTx
}

type MsgRouter struct {
	routerMap    map[string]TxMessageHandler
	txSerializer serializer.TxSerializer
	mrLog        log.Logger
}

func(mr *MsgRouter) SetLogger(mrLog log.Logger) {
	mr.mrLog = mrLog
}

func(mr *MsgRouter) AddTxMessageHandler(txMsgName string, txMsgHandler TxMessageHandler) {
	mr.routerMap[txMsgName] = txMsgHandler
}

func(mr *MsgRouter) TxMessageHandler(tx []byte) (TxMessageHandler, interface{}) {
	txMsg, err := mr.txSerializer.Deserialize(tx)
	if err != nil {
		if mr.mrLog != nil {
			mr.mrLog.Error("can't decode tx", "err", err)
		}

		return nil, nil
	}

	if txHandler, ok:= mr.routerMap[txMsg.Type()]; ok {
		return txHandler, txMsg.ImplTxMsg
	}else {
		mr.mrLog.Error("can't find the respond txmsg handler", "txType", txMsg.Type())
		return nil, nil
	}
}

func MsgRouterInstance() *MsgRouter {
	onceMR.Do(func(){
		routerMap := make(map[string]TxMessageHandler)
		instanceMR = &MsgRouter{routerMap: routerMap, txSerializer: serializer.NewTxSerializer()}
	})

	return instanceMR
}

*/









