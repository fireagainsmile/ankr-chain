package serializer

import (
	"github.com/Ankr-network/ankr-chain/tx"
	"github.com/Ankr-network/ankr-chain/tx/contract"
	"github.com/Ankr-network/ankr-chain/tx/metering"
	"github.com/Ankr-network/ankr-chain/tx/token"
	"github.com/Ankr-network/ankr-chain/tx/validator"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/encoding/amino"
)

func CreateTxCDC() *amino.Codec {
	txCdc := amino.NewCodec()
	cryptoAmino.RegisterAmino(txCdc)
	txCdc.RegisterInterface((*tx.ImplTxMsg)(nil), nil)
	txCdc.RegisterConcrete(&tx.TxMsg{}, "ankr-chain/tx/txMsg", nil)
	txCdc.RegisterConcrete(&token.TransferMsg{}, "ankr-chain/tx/token/tranferTxMsg", nil)
	txCdc.RegisterConcrete(&validator.ValidatorMsg{}, "ankr-chain/tx/validator/validatorMsg", nil)
	txCdc.RegisterConcrete(&metering.SetCertMsg{}, "ankr-chain/tx/metering/setCertMsg", nil)
	txCdc.RegisterConcrete(&metering.RemoveCertMsg{}, "ankr-chain/tx/metering/removeCertMsg", nil)
	txCdc.RegisterConcrete(&metering.MeteringMsg{}, "ankr-chain/tx/metering/meteringMsg", nil)
	txCdc.RegisterConcrete(&contract.ContractDeployMsg{}, "ankr-chain/tx/contract/ContractDeployMsg", nil)
	txCdc.RegisterConcrete(&contract.ContractInvokeMsg{}, "ankr-chain/tx/contract/ContractInvokeMsg", nil)
	txCdc.RegisterConcrete(&tx.TxMsgTesting{}, "ankr-chain/tx/txMsgTesting", nil)

	return txCdc
}
