package decoder

import (
	"errors"
	"strings"

	ankrtypes "github.com/Ankr-network/ankr-chain/types"
)

func ParseTxPrefix(tx []byte) (string, error) {
	if strings.HasPrefix(string(tx), ankrtypes.ValidatorSetChangePrefix) {
		return ankrtypes.ValidatorSetChangePrefix, nil
	}else if strings.HasPrefix(string(tx), ankrtypes.TrxSendPrefix) {
		return ankrtypes.TrxSendPrefix, nil
	}else if strings.HasPrefix(string(tx), ankrtypes.SetMeteringPrefix) {
		return ankrtypes.SetMeteringPrefix, nil
	}else if strings.HasPrefix(string(tx), ankrtypes.SetCertPrefix) {
		return ankrtypes.SetCertPrefix, nil
	}else if strings.HasPrefix(string(tx), ankrtypes.RemoveCertPrefix) {
		return ankrtypes.RemoveCertPrefix, nil
	}else if strings.HasPrefix(string(tx), ankrtypes.SetBalancePrefix) {
		return ankrtypes.SetBalancePrefix, nil
	}else if strings.HasPrefix(string(tx), ankrtypes.SetOpPrefix) {
		return ankrtypes.SetOpPrefix, nil
	}else if strings.HasPrefix(string(tx), ankrtypes.SetStakePrefix) {
		return ankrtypes.SetOpPrefix, nil
	}else {
		return "", errors.New("unknown tx")
	}

	return "", nil
}

type TxDecoderV0 struct {
}

func (txdv1 *TxDecoderV0) Decode(txBytes []byte) (txType string, data interface{}, err error) {
	if len(txBytes) == 0 {
		txType = ""
		data   = nil
		err    = errors.New("nil tx")
		return
	}

	txPrefix, err := ParseTxPrefix(txBytes)
	if err != nil {
		txType = ""
		data   = nil
		return "", nil, err
	}

	txType = txPrefix
	data   = strings.Split(string(txBytes[len(txPrefix):]), ":")

	return
}
