package iavl

import (
	"github.com/stretchr/testify/assert"
	"testing"

	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

)

func TestSaveWithLevelDB(t *testing.T) {
	db, err := dbm.NewGoLevelDB("teststore", "./teststore")
	if err != nil {
		panic(err)
	}

	iavalStore := NewIavlStore(db, 100, 10, log.NewNopLogger())

	iavalStore.Set([]byte("key1"), []byte("value1"))
	iavalStore.Set([]byte("key2"), []byte("value2"))

	iavalStore.Commit()

	name1, err := iavalStore.Get([]byte("key1"))

	assert.Equal(t, err, nil)
	assert.Equal(t, string(name1), "value1")
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
}


