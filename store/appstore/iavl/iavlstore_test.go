package iavl

import (
	"github.com/stretchr/testify/assert"
	"testing"

	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func TestSaveWithLevelDB(t *testing.T) {
	db, err := dbm.NewGoLevelDB("teststore", "./teststore")
	if err != nil {
		panic(err)
	}

	iavalStore := NewIavlStore(db, 100, 10, log.NewNopLogger())

	iavalStore.Set(ankrtypes.PrefixBalanceKey([]byte("7246BCE86AC2BA9CAC1B00D229B0AE08F58E3A4A1F8BD4")), []byte("1000000000"+":"+"10"))

	iavalStore.Set([]byte("key1"), []byte("value1"))
	iavalStore.Set([]byte("key2"), []byte("value2"))

	val, _ := iavalStore.Get(ankrtypes.PrefixBalanceKey([]byte("7246BCE86AC2BA9CAC1B00D229B0AE08F58E3A4A1F8BD4")))
	name1, err := iavalStore.Get([]byte("key1"))

	assert.Equal(t, string(val), "1000000000:10")
	assert.Equal(t, err, nil)
	assert.Equal(t, string(name1), "value1")

	db.Close()
}

func TestLoadWithLevelDB(t *testing.T) {
	db, err := dbm.NewGoLevelDB("teststore", "./teststore")
	if err != nil {
		panic(err)
	}

	iavalStore := NewIavlStore(db, 100, 10, log.NewNopLogger())

	iavalStore.Load()

	name1, err := iavalStore.Get([]byte("key1"))

	assert.Equal(t, err, nil)
	assert.Equal(t, string(name1), "value1")

	db.Close()
}

