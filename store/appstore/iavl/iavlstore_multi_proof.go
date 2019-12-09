package iavl

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/tendermint/iavl"

	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/merkle"
)

const (
	ProofOpMultiStore = "multistore"
)

type IavlStoreMultiProof struct {
	CommInfo commitInfo
}

func (proof *IavlStoreMultiProof) RootHash() []byte {
	return proof.CommInfo.AppHash
}

type IavlStoreMultiOp struct {
	key    []byte
	Proof  *IavlStoreMultiProof `json:"proof"`
}

func NewIavlStoreMultiOp(key []byte, proof *IavlStoreMultiProof) *IavlStoreMultiOp {
	return &IavlStoreMultiOp{
		key,
		proof,
	}
}

func (op IavlStoreMultiOp) Run(args [][]byte) ([][]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("IavlStoreMultiOp Value size is not 1")
	}
	value := args[0]

	root  := op.Proof.RootHash()

	for _, commI :=  range op.Proof.CommInfo.Commits {
		if commI.Name == string(op.key) {
			if bytes.Equal(value, commI.CID.Hash) {
				return [][]byte{root}, nil
			}else {
				return nil, fmt.Errorf("IavlStoreMultiOp %s's hash mismatch: Expected %X, Got %X", commI.Name,value, commI.CID.Hash)
			}
		}
	}

	return nil, fmt.Errorf("IavlStoreMultiOp can't find responding hash: key=%s", string(op.key))
}

func (op IavlStoreMultiOp) GetKey() []byte {
	return op.key
}

func (op IavlStoreMultiOp) ProofOp() merkle.ProofOp {
	dataBytes := amino.NewCodec().MustMarshalBinaryLengthPrefixed(op)
	return merkle.ProofOp{
		Type: ProofOpMultiStore,
		Key:  op.key,
		Data: dataBytes,
	}
}

func IavlStoreMultiOpDecoder(pop merkle.ProofOp) (merkle.ProofOperator, error) {
	if pop.Type != ProofOpMultiStore {
		return nil, fmt.Errorf("unexpected ProofOp.Type; got %v, want %v", pop.Type, ProofOpMultiStore)
	}

	var op IavlStoreMultiOp
	err := amino.NewCodec().UnmarshalBinaryLengthPrefixed(pop.Data, &op)
	if err != nil {
		return nil, fmt.Errorf("decoding ProofOp.Data into IavlStoreMultiOp: %w", err)
	}
	return NewIavlStoreMultiOp(pop.Key, op.Proof), nil
}

func DefaultProofRuntime() ( *merkle.ProofRuntime) {
	proofRT := merkle.NewProofRuntime()
	proofRT.RegisterOpDecoder(merkle.ProofOpSimpleValue, merkle.SimpleValueOpDecoder)
	proofRT.RegisterOpDecoder(iavl.ProofOpIAVLValue, iavl.IAVLValueOpDecoder)
	proofRT.RegisterOpDecoder(iavl.ProofOpIAVLAbsence, iavl.IAVLAbsenceOpDecoder)
	proofRT.RegisterOpDecoder(ProofOpMultiStore, IavlStoreMultiOpDecoder)

	return proofRT
}






