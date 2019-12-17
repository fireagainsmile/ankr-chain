package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"

	tmcmm "github.com/tendermint/tendermint/libs/common"
)

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