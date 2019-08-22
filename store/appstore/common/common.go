package common

import (
	"encoding/json"

	dbm "github.com/tendermint/tendermint/libs/db"
)

var (
	StateKey = []byte("stateKey")
)

type State struct {
	DB      dbm.DB
	Size    int64  `json:"size"`
	Height  int64  `json:"height"`
	AppHash []byte `json:"app_hash"`
}

func LoadState(db dbm.DB) State {
	stateBytes := db.Get(StateKey)
	var state State
	if len(stateBytes) != 0 {
		err := json.Unmarshal(stateBytes, &state)
		if err != nil {
			panic(err)
		}
	}
	state.DB = db

	return state
}
