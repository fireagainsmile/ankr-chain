package code

import (
	"github.com/tendermint/tendermint/abci/types"
)

const (
	CodeTypeOK                       uint32 = types.CodeTypeOK
	CodeTypeEncodingError            uint32 = 1
	CodeTypeDecodingError            uint32 = 2
	CodeTypeBadNonce                 uint32 = 3
	CodeTypeUnauthorized             uint32 = 4
	CodeTypeP2PFilterPathError       uint32 = 5
	CodeTypeP2PFilterNotPass         uint32 = 6
	CodeQueryDataLenZero             uint32 = 7
	CodeTypeUnknownError             uint32 = 8
	CodeTypeMismatchChainID          uint32 = 9
	CodeTypeGetStoreNonceError       uint32 = 10
	CodeTypeInvalidAddress           uint32 = 11
	CodeTypeVerifySignaError         uint32 = 12
	CodeTypeLoadBalError             uint32 = 13
	CodeTypeBalNotEnough             uint32 = 14
	CodeTypeGasNotEnough             uint32 = 15
	CodeTypeFeeNotEnough             uint32 = 16
	CodeTypeGasPriceIrregular        uint32 = 17
	CodeTypeLoadContractErr          uint32 = 18
	CodeTypeCallContractErr          uint32 = 19
	CodeTypeInvalidStakeCurrency     uint32 = 20
	CodeTypeInvalidValidatorPubKey   uint32 = 21
	CodeTypeNotPermitPubKey          uint32 = 22
	CodeTypeContractInvalidCodeSize  uint32 = 23
	CodeTypeContractAddrTakenUp      uint32 = 24
	CodeTypeContractInvalidAddr      uint32 = 25
	CodeTypeContractCantFound        uint32 = 26
	CodeTypeQueryInvalidStoreName    uint32 = 27
	CodeTypeQueryInvalidQueryReqData uint32 = 28
	CodeTypeNoV0TxHandler            uint32 = 29
	CodeTypeInvalidFromAddr          uint32 = 30
)
