package handler

import (
	"github.com/Ankr-network/ankr-chain/version"
	"net/http"

	"github.com/Ankr-network/ankr-chain/client"
)

func QueryNodeInfoHandler(c *client.Client) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		rsStatus, err := c.Status()
		if err != nil {
			WriteErrorResponse(resp, http.StatusInternalServerError, err.Error())
			return
		}

		respData := &NodeInfo{}

		respData.ID           = string(rsStatus.NodeInfo.ID())
		respData.ListenAddr   = rsStatus.NodeInfo.ListenAddr
		respData.ChainID      = rsStatus.NodeInfo.Network
		respData.Moniker      = rsStatus.NodeInfo.Moniker
		respData.TxIndex      = rsStatus.NodeInfo.Other.TxIndex
		respData.RPCAddr      = rsStatus.NodeInfo.Other.RPCAddress
		respData.LatestHeight = uint32(rsStatus.SyncInfo.LatestBlockHeight)
		respData.Version      =  version.NodeVersion

		respJson, err :=Cdc.MarshalJSON(respData)
		if err != nil {
			WriteErrorResponse(resp, http.StatusInternalServerError, err.Error())
			return
		}
		resp.Header().Set("Content-Type", "application/json")
		resp.Write(respJson)

		return
	}
}
