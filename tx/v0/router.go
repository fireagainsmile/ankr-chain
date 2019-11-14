package v0

import (
	"fmt"
	"github.com/Ankr-network/ankr-chain/common/code"
	"sync"

	ankrcmm "github.com/Ankr-network/ankr-chain/common"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
)

var (
	onceMR     sync.Once
	instanceMR *MsgRouter
)

type DeliverTxHandler interface {
	ProcessTx(txMsg interface{}, appStore appstore.AppStore) types.ResponseDeliverTx
}

type MsgRouter struct {
	routerMap  map[string]DeliverTxHandler
	serializer TxSerializerV0
	mrLog      log.Logger
}

func(mr *MsgRouter) SetLogger(mrLog log.Logger) {
	mr.mrLog = mrLog
}

func(mr *MsgRouter) AddTxMessageHandler(txMsgName string, txHandler DeliverTxHandler) {
	mr.routerMap[txMsgName] = txHandler
}

func(mr *MsgRouter) txMessageHandler(tx []byte) (DeliverTxHandler, interface{}) {
	txType, data, err := mr.serializer.Deserialize(tx)
	if err != nil {
		if mr.mrLog != nil {
			mr.mrLog.Error("can't decode tx", "err", err, "tx", tx)
		}

		return nil, nil
	}

	if mr.mrLog != nil {
		mr.mrLog.Info("MsgRouter txMessageHandler", "txType", txType)
	}

	if txHandler, ok:= mr.routerMap[txType]; ok {
		return txHandler, data
	}else {
		if mr.mrLog != nil {
			mr.mrLog.Error("can't find the respond txmsg handler", "txType", txType)
		}

		return nil, nil
	}
}

func(mr *MsgRouter) DeliverTx(tx []byte, store appstore.AppStore) types.ResponseDeliverTx {
	handler, msgData := mr.txMessageHandler(tx)
	if handler == nil {
		return types.ResponseDeliverTx{Code: code.CodeTypeNoV0TxHandler, Log: fmt.Sprintf("can't find responding v0 tx handler, tx=%v", tx)}
	}

	resDeliverTx := handler.ProcessTx(msgData, store)

	if mr.mrLog != nil {
		mr.mrLog.Info("MsgRouter DeliverTx ProcessTx", "code", resDeliverTx.Code, "log", resDeliverTx.Log)
	}

	if resDeliverTx.Code == code.CodeTypeOK {
		store.IncTotalTx()
	}

	return resDeliverTx
}

func MsgRouterInstance() *MsgRouter {
	onceMR.Do(func(){
		routerMap := make(map[string]DeliverTxHandler)
		instanceMR = &MsgRouter{routerMap: routerMap}
		instanceMR.AddTxMessageHandler(ankrcmm.TrxSendPrefix, new(transferMsg))
		instanceMR.AddTxMessageHandler(ankrcmm.SetCertPrefix, new(setCertMsg))
		instanceMR.AddTxMessageHandler(ankrcmm.RemoveCertPrefix, new(removeCertMsg))
		instanceMR.AddTxMessageHandler(ankrcmm.SetMeteringPrefix, new(meteringMsg))
	})

	return instanceMR
}











