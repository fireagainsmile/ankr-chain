package handler

import (
	"net/http"

	"github.com/Ankr-network/ankr-chain/version"
)

func QueryVersionHandler() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {

		respData := &LASVersion{version.LasVersion}

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