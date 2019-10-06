package serializer

import (
	"errors"
	"strings"

	"github.com/Ankr-network/ankr-chain/tx"
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

type TxSerializerV0 struct {
}

func (txdv1 *TxSerializerV0) Deserialize(txBytes []byte) (*tx.TxMsg, error) {
	if len(txBytes) == 0 {
		return nil, errors.New("nil tx")
	}

	_, err := ParseTxPrefix(txBytes)
	if err != nil {
		//txType = ""
		//data   = nil
		return  nil, err
	}

	//txType = txPrefix
	//data   = strings.Split(string(txBytes[len(txPrefix):]), ":")

	return nil, nil
}
