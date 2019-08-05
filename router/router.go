package router

import (
	"errors"
	"strings"
	"sync"

	"github.com/Ankr-network/ankr-chain/store/appstore"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
)


var (
	onceMR     sync.Once
	instanceMR *MsgRouter
)

type TxMessageHandler interface {
	CheckTx(tx []byte, appStore appstore.AppStore) types.ResponseCheckTx
	DeliverTx(tx []byte, appStore appstore.AppStore) types.ResponseDeliverTx
}

func ParseTxPrefix(tx []byte) (string, error) {
	if strings.HasPrefix(string(tx), ankrtypes.ValidatorSetChangePrefix) {
		return ankrtypes.ValidatorSetChangePrefix, nil
	}else if strings.HasPrefix(string(tx), ankrtypes.TrxSendPrefix) {
		return ankrtypes.TrxSendPrefix, nil
	}else if strings.HasPrefix(string(tx), ankrtypes.SetMeteringPrefix) {
		return ankrtypes.SetMeteringPrefix, nil
	}else if strings.HasPrefix(string(tx), ankrtypes.SetCertPrefix) {
		return ankrtypes.SetCertPrefix, nil
	}else if strings.HasPrefix(string(tx), ankrtypes.RemoveCertPrefix) {
		return ankrtypes.RemoveCertPrefix, nil
	}else if strings.HasPrefix(string(tx), ankrtypes.SetBalancePrefix) {
		return ankrtypes.SetBalancePrefix, nil
	}else if strings.HasPrefix(string(tx), ankrtypes.SetOpPrefix) {
		return ankrtypes.SetOpPrefix, nil
	}else if strings.HasPrefix(string(tx), ankrtypes.SetStakePrefix) {
		return ankrtypes.SetOpPrefix, nil
	}else {
		return "", errors.New("unknown tx")
	}

	return "", nil
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

func(mr *MsgRouter) TxMessageHandler(tx []byte) TxMessageHandler {
	txMsgName, err := ParseTxPrefix(tx)
	if err != nil {
		if mr.mrLog != nil {
			mr.mrLog.Error("unknown tx", "err", err)
		}

		return nil
	}

	if txHandler, ok:= mr.routerMap[txMsgName]; ok {
		return txHandler
	}else {
		mr.mrLog.Error("can't find the respond txmsg handler", "txmsgname", txMsgName)
		return nil
	}
}

func MsgRouterInstance() *MsgRouter {
	onceMR.Do(func(){
		routerMap := make(map[string]TxMessageHandler)
		instanceMR = &MsgRouter{routerMap: routerMap}
	})

	return instanceMR
}










