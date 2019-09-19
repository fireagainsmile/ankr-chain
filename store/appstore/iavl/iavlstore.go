package iavl

import (
	"errors"
	"fmt"
	"github.com/Ankr-network/ankr-chain/common/code"

	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	"github.com/tendermint/iavl"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/merkle"
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

	tree.SaveVersion()

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

func (s *IavlStore) Commit() (ankrtypes.CommitID, error) {
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

	return ankrtypes.CommitID{ver, rHash}, nil
}

func (s *IavlStore) LatestVersion() ankrtypes.CommitID {
	ver   :=  s.tree.Version()
	rHash := s.tree.Hash()

	return ankrtypes.CommitID{ver, rHash}
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

func (s *IavlStore) getHeight(reqHeight int64) int64 {
	height := reqHeight
	if height == 0 {
		latest := s.tree.Version()
		if s.tree.VersionExists(latest - 1) {
			height = latest - 1
		} else {
			height = latest
		}
	}
	return height
}

func (s *IavlStore) Query(reqQuery types.RequestQuery) (resQuery types.ResponseQuery) {
	if len(reqQuery.Data) == 0 {
		resQuery.Code = code.CodeQueryDataLenZero
		resQuery.Log  = "Query cannot be zero length"
		return
	}

	tree := s.tree

	// store the height we chose in the response, with 0 being changed to the
	// latest height
	resQuery.Height = s.getHeight(reqQuery.Height)

	switch reqQuery.Path {
	case "/key":        // get by key
		key := reqQuery.Data // data holds the key bytes

		resQuery.Key = key
		if !s.tree.VersionExists(resQuery.Height) {
			resQuery.Log = iavl.ErrVersionDoesNotExist.Error()
			break
		}

		if reqQuery.Prove {
			value, proof, err := tree.GetVersionedWithProof(key, resQuery.Height)
			if err != nil {
				resQuery.Log = err.Error()
				break
			}
			if proof == nil {
				// Proof == nil implies that the store is empty.
				if value != nil {
					panic("unexpected value for an empty proof")
				}
			}
			if value != nil {
				// value was found
				resQuery.Value = value
				resQuery.Proof = &merkle.Proof{Ops: []merkle.ProofOp{iavl.NewIAVLValueOp(key, proof).ProofOp()}}
			} else {
				// value wasn't found
				resQuery.Value = nil
				resQuery.Proof = &merkle.Proof{Ops: []merkle.ProofOp{iavl.NewIAVLAbsenceOp(key, proof).ProofOp()}}
			}
		} else {
			_, resQuery.Value = tree.GetVersioned(key, resQuery.Height)
		}
	default:
		resQuery.Code = code.CodeTypeUnknownError
		resQuery.Log  =  fmt.Sprintf("Unexpected Query path: %v", reqQuery.Path)
		return
	}

	return
}