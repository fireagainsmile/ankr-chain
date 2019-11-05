package iavl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func TestIavlStoreVersion(t *testing.T) {
	db := dbm.NewMemDB()
	storeLog := log.NewNopLogger()

	iavlStore := NewIavlStore(db, 100, 10, storeLog)

	iavlStore.Set([]byte("testkey1"), []byte("testvalue1"))
	iavlStore.Commit()
	t.Logf("version=%d", iavlStore.tree.Version())

	val1, err := iavlStore.Get([]byte("testkey1"))
	assert.Equal(t, nil, err)
	assert.Equal(t, "testvalue1", string(val1))

	iavlStore.Set([]byte("testkey1"), []byte("testvalue2"))
	//iavlStore.Commit()
	t.Logf("version=%d", iavlStore.tree.Version())

	val1, err = iavlStore.Get([]byte("testkey1"))
	assert.Equal(t, nil, err)
	assert.Equal(t, "testvalue2", string(val1))

	iavlStore.Rollback()
	t.Logf("version=%d", iavlStore.tree.Version())

	val1, err = iavlStore.Get([]byte("testkey1"))
	assert.Equal(t, nil, err)
	assert.Equal(t, "testvalue1", string(val1))
}

func TestIavlStoreVersionCommitRollback(t *testing.T) {
	db := dbm.NewMemDB()
	storeLog := log.NewNopLogger()

	iavlStore := NewIavlStore(db, 100, 10, storeLog)

	iavlStore.Set([]byte("testkey1"), []byte("testvalue1"))
	iavlStore.Commit()
	t.Logf("version=%d", iavlStore.tree.Version())

	val1, err := iavlStore.Get([]byte("testkey1"))
	assert.Equal(t, nil, err)
	assert.Equal(t, "testvalue1", string(val1))

	iavlStore.Set([]byte("testkey1"), []byte("testvalue2"))
	iavlStore.Commit()
	t.Logf("version=%d", iavlStore.tree.Version())

	val1, err = iavlStore.Get([]byte("testkey1"))
	assert.Equal(t, nil, err)
	assert.Equal(t, "testvalue2", string(val1))

	iavlStore.Rollback()
	t.Logf("version=%d", iavlStore.tree.Version())

	val1, err = iavlStore.Get([]byte("testkey1"))
	assert.Equal(t, nil, err)
	assert.Equal(t, "testvalue2", string(val1))
}

