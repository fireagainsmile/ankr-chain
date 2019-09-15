package code

import (
	"github.com/tendermint/tendermint/abci/types"
)

const (
	CodeTypeOK                 uint32 = types.CodeTypeOK
	CodeTypeEncodingError      uint32 = 1
	CodeTypeDecodingError      uint32 = 2
	CodeTypeBadNonce           uint32 = 3
	CodeTypeUnauthorized       uint32 = 4
	CodeTypeP2PFilterPathError uint32 = 5
	CodeTypeP2PFilterNotPass   uint32 = 6
	CodeQueryDataLenZero       uint32 = 7
	CodeTypeUnknownError       uint32 = 8
	CodeTypeMismatchChainID    uint32 = 9
	CodeTypeGetStoreNonceError uint32 = 10
	CodeTypeInvalidAddress     uint32 = 11
	CodeTypeVerifySignaError   uint32 = 12
	CodeTypeBalError           uint32 = 13
	CodeTypeGasNotEnough       uint32 = 14
	CodeTypeTransferNotEnough  uint32 = 15
)
