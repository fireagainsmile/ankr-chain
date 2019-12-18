package handler

import (
	"fmt"
	"github.com/spf13/viper"
	"math/big"
	"net/http"

	"github.com/Ankr-network/ankr-chain/client"
	ankrcmm "github.com/Ankr-network/ankr-chain/common"
	"github.com/Ankr-network/ankr-chain/common/code"
	"github.com/Ankr-network/ankr-chain/crypto"
	lascmm "github.com/Ankr-network/ankr-chain/service/las/common"
	"github.com/Ankr-network/ankr-chain/tx/serializer"
	"github.com/Ankr-network/ankr-chain/tx/token"
	"github.com/shopspring/decimal"
)

func convertToDevimalValFromFloat64str(val string) (*big.Int, error) {
	valF,  err := decimal.NewFromString(val)
	if err != nil {
		return nil, err
	}

	if !valF.GreaterThan(decimal.Zero) {
		return nil, fmt.Errorf("invalid val: %s", val)
	}

	valDString := valF.Shift(18).String()

	valBigInt, isSucc := new(big.Int).SetString(valDString, 10)
	if !isSucc {
		return nil, fmt.Errorf("invalid val: %s", val)
	}

	return valBigInt, nil
}

func convertToFloat64strFromDevimalVal(val *big.Int) string {
	valD := decimal.NewFromBigInt(val, 0)
	return valD.Shift(-18).String()
}

func TxTransferHandler(c *client.Client) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		var transferInfo TxMsgTransferInfo
		err := ReadPostBody(resp, req, Cdc, &transferInfo)
		if err != nil {
			WriteErrorResponse(resp, http.StatusBadRequest, err.Error())
			return
		}

		gasPriceBigValWithDecimal, err := convertToDevimalValFromFloat64str(transferInfo.Header.GasPrice.Value)
		if err != nil {
			WriteErrorResponse(resp, http.StatusBadRequest, err.Error())
			return
		}

		accountBigValWithDecimal, err := convertToDevimalValFromFloat64str(transferInfo.Data.Amount.Value)

		msgHeader := client.TxMsgHeader{
			ChID:     ankrcmm.ChainID(viper.GetString(lascmm.FlagChainID)),
			GasLimit: new(big.Int).SetUint64(uint64(transferInfo.Header.GasLimit)).Bytes(),
			GasPrice: ankrcmm.Amount{ankrcmm.Currency{transferInfo.Header.GasPrice.Symbol, 18}, gasPriceBigValWithDecimal.Bytes()},
			Memo:     transferInfo.Header.Memo,
			Version:  "1.0",
		}

		tfMsg := &token.TransferMsg{
			FromAddr: transferInfo.Data.FromAddr,
			ToAddr:   transferInfo.Data.ToAddr,
			Amounts: []ankrcmm.Amount{ankrcmm.Amount{ankrcmm.Currency{transferInfo.Data.Amount.Symbol, 18}, accountBigValWithDecimal.Bytes()}},
		}

		txSerializer := serializer.NewTxSerializerCDC()

		key := crypto.NewSecretKeyEd25519(transferInfo.Data.PriKey)

		builder := client.NewTxMsgBuilder(msgHeader, tfMsg,  txSerializer, key)

		cResult, err := builder.BuildAndCommitWithRawResult(c)
		if err != nil{
			WriteErrorResponse(resp, http.StatusInternalServerError, err.Error())
			return
		}

		if cResult == nil {
			WriteErrorResponse(resp, http.StatusInternalServerError, "Invalid Commit Result")
			return
		}

		commitR := &TxCommitResult{}

		commitR.TxHash  = cResult.Hash.String()
		commitR.Height  = uint32(cResult.Height)
		commitR.Log     = cResult.DeliverTx.Log
		commitR.GasUsed = uint32(cResult.DeliverTx.GasUsed)

		commitR.Status = "Success"
		if cResult.DeliverTx.Code != code.CodeTypeOK {
			commitR.Status = "Fail"
		}

		commitR.TimeStamp = TagValue("app.timestamp", cResult.DeliverTx.Tags)

		if cResult.DeliverTx.Code != code.CodeTypeOK {
			commitR.Err = fmt.Sprintf("BuildAndCommitWithRawResult DeliverTx response code not ok, code=%d, log=%s", cResult.DeliverTx.Code, cResult.DeliverTx.Log)
		}

		rtnJson, err := Cdc.MarshalJSON(commitR)
		if err != nil {
			WriteErrorResponse(resp, http.StatusInternalServerError, err.Error())
			return
		}

		resp.Write(rtnJson)

		return
	}
}

