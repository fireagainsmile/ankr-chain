package iavl

import (
	"errors"

	"github.com/Ankr-network/ankr-chain/types"
	"github.com/tendermint/iavl"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

type IavlStore struct {
	tree           *iavl.MutableTree
	keepVersionNum int64
	log            log.Logger
}

func NewIavlStore(db dbm.DB, cacheSize int, keepVersionNum int64, logStore log.Logger) *IavlStore {
	if db == nil {
		panic("can't create IvalStore, db nil")
	}

	if logStore == nil {
		panic("can't create IvalStore, logStore nil")
	}

	tree := iavl.NewMutableTree(db, cacheSize)
	if tree == nil {
		panic("create MutableTree failed")
	}

	return &IavlStore{tree:tree, keepVersionNum: keepVersionNum, log: logStore}
}

func (s *IavlStore) Set(key []byte, value []byte) bool {
	return s.tree.Set(key, value)
}

func (s *IavlStore) Get(key []byte) ([]byte, error) {
	if len(key) == 0 {
		s.log.Error("key is nil")
		return nil, errors.New("key is nil")
	}

	version := int64(s.tree.Version())
	_, val := s.tree.GetVersioned(key, version)

	return val, nil
}

func (s *IavlStore) Has(key []byte) bool {
	return s.tree.Has(key)
}

func (s *IavlStore) Remove(key []byte) ([]byte, bool) {
	return s.tree.Remove(key)
}

func (s *IavlStore) Commit() (types.CommitID, error) {
	rHash, ver, err := s.tree.SaveVersion()
	if err != nil {
		panic(err)
	}

	if ver > s.keepVersionNum  && s.keepVersionNum > 0{
		err = s.tree.DeleteVersion(ver-s.keepVersionNum)
		if err != nil {
			panic(err)
		}
	}

	return types.CommitID{ver, rHash}, nil
}

func (s *IavlStore) LatestVersion() types.CommitID {
	ver   :=  s.tree.Version()
	rHash := s.tree.Hash()

	return types.CommitID{ver, rHash}
}

func (s *IavlStore) LoadVersion(ver int64) (int64, error) {
	return s.tree.LoadVersion(ver)
}

// Load the latest versioned tree from disk.
func (s *IavlStore) Load() (int64, error) {
	return s.tree.Load()
}

func (s *IavlStore) Rollback() {
	s.tree.Rollback()
	s.log.Debug("IavlStore rollback happens")
}