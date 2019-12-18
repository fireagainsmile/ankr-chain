package handler

import (
	"fmt"
	"math/big"
	"net/http"
	"strconv"

	"github.com/Ankr-network/ankr-chain/client"
	"github.com/Ankr-network/ankr-chain/common/code"
	ankrtxcmm "github.com/Ankr-network/ankr-chain/tx/common"
	"github.com/Ankr-network/ankr-chain/tx/token"
)

func QueryBlockTxTransfersHandler(c *client.Client) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		if err != nil {
			WriteErrorResponse(resp, http.StatusBadRequest, err.Error())
			return
		}

		if req.Form == nil || len(req.Form) == 0 {
			WriteErrorResponse(resp, http.StatusBadRequest, "QueryTxTransfersHandler, form val's len is 0 or nil")
			return
		}

		heightParam  := req.FormValue("height")
		pageParam    := req.FormValue("page")
		perPageParam := req.FormValue("perPage")

		if heightParam == "" {
			WriteErrorResponse(resp, http.StatusBadRequest, fmt.Sprintf("blank request param height, heightParam=%s", heightParam))
			return
		}
		height, err := strconv.ParseInt(heightParam, 10, 64)
		if err != nil {
			WriteErrorResponse(resp, http.StatusBadRequest, err.Error())
			return
		}

		page := int64(1)
		if pageParam != "" {
			page, err = strconv.ParseInt(pageParam, 10, 64)
			if err != nil {
				WriteErrorResponse(resp, http.StatusBadRequest, err.Error())
				return
			}
		}

		perPage := int64(100)
		if perPageParam != "" {
			perPage, err = strconv.ParseInt(perPageParam, 10, 64)
			if err != nil {
				WriteErrorResponse(resp, http.StatusBadRequest, err.Error())
				return
			}
		}

		queryStr := fmt.Sprintf("tx.height=%d AND app.type='Transfer'", height)
		sResult, err := c.TxSearch(queryStr, true, int(page), int(perPage))
		if err != nil {
			WriteErrorResponse(resp, http.StatusInternalServerError, err.Error())
			return
		}

		respData := &TranserResultsOfOneBlock{
			TotalTx: uint32(sResult.TotalCount),
		}
		for _, txR := range sResult.Txs {

			txMsg, err := client.NewTxDecoder().Decode(txR.Tx)
			if err != nil {
				WriteErrorResponse(resp, http.StatusInternalServerError, err.Error())
				return
			}

			if txMsg.ImplTxMsg.Type() != ankrtxcmm.TxMsgTypeTransfer {
				WriteErrorResponse(resp, http.StatusInternalServerError, err.Error())
				return
			}

			transferMsg := txMsg.ImplTxMsg.(*token.TransferMsg)

			txRItem := new(TranserResultItem)
			txRItem.TxHash   = fmt.Sprintf("%X", txR.Tx.Hash())
			txRItem.Height   = uint32(txR.Height)
			txRItem.FromAddr = transferMsg.FromAddr
			txRItem.ToAddr   = transferMsg.ToAddr

			for _, amount := range transferMsg.Amounts {
				amountLas := new(AmountLas)
				amountLas.Symbol = amount.Cur.Symbol
				amountLas.Value  = convertToFloat64strFromDevimalVal(new(big.Int).SetBytes(amount.Value))
				txRItem.Amounts = append(txRItem.Amounts, amountLas)
			}

			txRItem.Status = "Success"
			if txR.TxResult.Code != code.CodeTypeOK {
				txRItem.Status = "Fail"
			}

			txRItem.GasLimit        = uint32(new(big.Int).SetBytes(txMsg.GasLimit).Uint64())
			txRItem.GasPrice.Symbol = txMsg.GasPrice.Cur.Symbol
			txRItem.GasPrice.Value  = convertToFloat64strFromDevimalVal(new(big.Int).SetBytes(txMsg.GasPrice.Value))
			txRItem.TimeStamp       = TagValue("app.timestamp", txR.TxResult.Tags)
			txRItem.Memo            = txMsg.Memo

			respData.TransferResults = append(respData.TransferResults, txRItem)
		}

		rtnJson, err := Cdc.MarshalJSON(respData)
		if err != nil {
			WriteErrorResponse(resp, http.StatusInternalServerError, err.Error())
			return
		}

		resp.Write(rtnJson)

		return
	}
}

func QueryBlockSyncing(c *client.Client) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		respData := &blockSyncing{}
		rs, err := c.Status()
		if err != nil {
			WriteErrorResponse(resp, http.StatusInternalServerError, err.Error())
			return
		}

		respData.Syncing = rs.SyncInfo.CatchingUp

		respJson, err := Cdc.MarshalJSON(respData)
		if err != nil {
			WriteErrorResponse(resp, http.StatusInternalServerError, err.Error())
			return
		}
		resp.Header().Set("Content-Type", "application/json")
		resp.Write(respJson)
		return
	}
}
