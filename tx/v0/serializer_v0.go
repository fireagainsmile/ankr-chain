package v0

import (
	"errors"
	"strings"

	ankrcmm "github.com/Ankr-network/ankr-chain/common"
)

func ParseTxPrefix(tx []byte) (string, error) {
	if strings.HasPrefix(string(tx), ankrcmm.ValidatorSetChangePrefix) {
		return ankrcmm.ValidatorSetChangePrefix, nil
	}else if strings.HasPrefix(string(tx), ankrcmm.TrxSendPrefix) {
		return ankrcmm.TrxSendPrefix, nil
	}else if strings.HasPrefix(string(tx), ankrcmm.SetMeteringPrefix) {
		return ankrcmm.SetMeteringPrefix, nil
	}else if strings.HasPrefix(string(tx), ankrcmm.SetCertPrefix) {
		return ankrcmm.SetCertPrefix, nil
	}else if strings.HasPrefix(string(tx), ankrcmm.RemoveCertPrefix) {
		return ankrcmm.RemoveCertPrefix, nil
	}else if strings.HasPrefix(string(tx), ankrcmm.SetBalancePrefix) {
		return ankrcmm.SetBalancePrefix, nil
	}else if strings.HasPrefix(string(tx), ankrcmm.SetOpPrefix) {
		return ankrcmm.SetOpPrefix, nil
	}else if strings.HasPrefix(string(tx), ankrcmm.SetStakePrefix) {
		return ankrcmm.SetOpPrefix, nil
	}else {
		return "", errors.New("unknown tx")
	}

	return "", nil
}

type TxSerializerV0 struct {
}

func (txs *TxSerializerV0) Deserialize(txBytes []byte) (txType string, data interface{}, err error) {
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
