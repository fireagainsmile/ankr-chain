package handler

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"

	"github.com/shopspring/decimal"
	tmcmm "github.com/tendermint/tendermint/libs/common"
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

func WriteErrorResponse(w http.ResponseWriter, status int, err string) {
	w.WriteHeader(status)
	w.Write([]byte(err))
}

func ReadPostBody(w http.ResponseWriter, r *http.Request, cdc *Codec, req interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("invalid post body")
			WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}
	}()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return err
	}

	err = cdc.UnmarshalJSON(body, req)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return err
	}

	return nil
}

func TagValue(tagKey string, tags tmcmm.KVPairs) string {
	for _, kv := range tags {
		if string(kv.Key) == tagKey {
			return string(kv.Value)
		}
	}

	return ""
}