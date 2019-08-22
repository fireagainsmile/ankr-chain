package router

import (
	"fmt"
	"strings"
	"sync"

	"github.com/Ankr-network/ankr-chain/store/appstore"
	"github.com/tendermint/tendermint/libs/log"
)

var (
	onceQR     sync.Once
	instanceQR *QueryRouter
)

type QueryRouter struct {
	qrMap map[string]appstore.QueryHandler
	qrLog log.Logger
}

func(qr *QueryRouter) SetLogger(qrLog log.Logger) {
	qr.qrLog = qrLog
}

func (qr *QueryRouter) AddQueryHandler(path string, qHandler appstore.QueryHandler) {
	qr.qrMap[path] = qHandler
}

func (qr *QueryRouter) parseRouterPath(path string) (routerPath string, subPath string, err error) {
	if path == "" {
		routerPath = "store"
		subPath    = ""
		err        = nil
		return
	} else {
		if !strings.HasPrefix(path,"/") || len(path) == 1 {
			qr.qrLog.Error("invalid path for parseRouterPath", "path", path)
			return "", "", fmt.Errorf("invalid path for parseRouterPath, path=%s", path)
		}

		pathSegs := strings.Split(path[1:], "/")
		if len(pathSegs) >= 1 {
			routerPath = pathSegs[0]
		}

		if len(pathSegs) >= 2 {
			subPath = path[len(pathSegs[0])+1:]
		}

		return
	}
}

func (qr* QueryRouter) QueryHandler(path string) (appstore.QueryHandler, string) {
	rPath, subPath, err := qr.parseRouterPath(path)
	if err != nil {
		return nil, ""
	}

	if qHandler, ok := qr.qrMap[rPath]; ok {
		return qHandler, subPath
	}

	qr.qrLog.Error("There is no responding query handler", "path", path)

	return nil, ""
}

func QueryRouterInstance() *QueryRouter {
	onceMR.Do(func(){
		qrMap := make(map[string]appstore.QueryHandler)
		instanceQR = &QueryRouter{qrMap: qrMap}
	})

	return instanceQR
}





