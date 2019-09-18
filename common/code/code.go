package code

import (
	"github.com/tendermint/tendermint/abci/types"
)

const (
	CodeTypeOK                   uint32 = types.CodeTypeOK
	CodeTypeEncodingError        uint32 = 1
	CodeTypeBadNonce             uint32 = 2
	CodeTypeUnauthorized         uint32 = 3
	CodeTypeP2PFilterPathError   uint32 = 4
	CodeTypeP2PFilterNotPass     uint32 = 5
	CodeQueryDataLenZero         uint32 = 6
	CodeQueryNoQueryHandlerFound uint32 = 7
	CodeTypeUnknownError         uint32 = 8
)
