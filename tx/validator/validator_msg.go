package validator

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/Ankr-network/ankr-chain/common"
	"github.com/Ankr-network/ankr-chain/common/code"
	ankrcrypto "github.com/Ankr-network/ankr-chain/crypto"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	apm "github.com/Ankr-network/ankr-chain/tx"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	"github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	tmCoreTypes "github.com/tendermint/tendermint/types"
	"strconv"
)

type ValidatorMsg struct {
	apm.TxMsg
}

func (v *ValidatorMsg) GasWanted() int64 {
	return 0
}

func (v *ValidatorMsg) GasUsed() int64 {
	return 0
}

func (v *ValidatorMsg) Type() string {
	return ankrtypes.TrxSendPrefix
}

func (v *ValidatorMsg) Bytes() []byte {
	return nil
}
func (v *ValidatorMsg) SetSecretKey(sk ankrcrypto.SecretKey) {

}

func (v *ValidatorMsg) SecretKey() ankrcrypto.SecretKey {
	return nil
}

func (v *ValidatorMsg) ProcessTx(txMsg interface{}, appStore appstore.AppStore, isOnlyCheck bool) (uint32, string,  []cmn.KVPair) {
	pubKeyAndPower, ok := txMsg.([]string)
	if !ok {
		return  code.CodeTypeEncodingError, fmt.Sprintf("invalid tx set op msg"), nil
	}

	if len(pubKeyAndPower) != 5 {
		return code.CodeTypeEncodingError, fmt.Sprintf("Expected 'pubkey/power'. Got %v", pubKeyAndPower), nil
	}
	pubkeyS   := pubKeyAndPower[0]
	powerS    :=pubKeyAndPower[1]
	nonceS    := pubKeyAndPower[2]
	adminPubS := pubKeyAndPower[3]
	sigS      := pubKeyAndPower[4]

	var admin_pubkey_str string = ""
	admin_pubkey := appStore.Get([]byte(ankrtypes.ADMIN_OP_VAL_PUBKEY_NAME))
	if len(admin_pubkey) == 0 {
		fmt.Println("use default ADMIN_OP_VAL_PUBKEY_NAME")
		admin_pubkey_str = adminPubKeyOfValidator()
	} else {
		admin_pubkey_str = string(admin_pubkey)
	}

	if adminPubS != admin_pubkey_str {
		return code.CodeTypeEncodingError, fmt.Sprintf("Unexpected pubkey. Got %v", adminPubS), nil
	}

	powerInt, err := strconv.ParseInt(string(powerS), 10, 64)
	if err != nil {
		return  code.CodeTypeEncodingError, fmt.Sprintf("Unexpected power. Got %v", powerS), nil
	} else { // power < 0
		if powerInt < 0 {
			return code.CodeTypeEncodingError, fmt.Sprintf("Unexpected amount, negative num. Got %v", powerS), nil
		}
	}

	curValidatorCount := ValidatorManagerInstance().TotalValidatorPowers(appStore)
	if (curValidatorCount + int64(powerInt)) > tmCoreTypes.MaxTotalVotingPower {
		return code.CodeTypeEncodingError, fmt.Sprintf("Total powers %d will reach with the power %d", tmCoreTypes.MaxTotalVotingPower, powerInt), nil
	}

	nonceInt, err_n := strconv.ParseInt(string(nonceS), 10, 64)
	if err_n != nil {
		return code.CodeTypeEncodingError, fmt.Sprintf("Unexpected nonce. Got %v", nonceS), nil
	}

	pDec, _ := base64.StdEncoding.DecodeString(sigS)
	pubKeyObject, err_d := common.DeserilizePubKey(adminPubS)
	if err_d != nil {
		return code.CodeTypeEncodingError, fmt.Sprintf("Deserilize pubkey failure. Got %v", adminPubS), nil
	}

	s256 := common.ConvertBySha256([]byte(fmt.Sprintf("%s%s%s", pubkeyS, powerS, nonceS)))
	bb := pubKeyObject.VerifyBytes(s256[:32], pDec)
	if !bb {
		fmt.Println("Bad signature, transaction failed.", sigS)
		return code.CodeTypeEncodingError, fmt.Sprintf("DeliverTx, Bad signature. Got %v", sigS), nil
	}

	var inNonceInt int64 = 0
	inNonce := appStore.Get(([]byte(ankrtypes.SET_VAL_NONCE)))
	if len(inNonce) == 0 {
		inNonceInt = 0
	} else {
		inNonceIntValue, err_p:= strconv.ParseInt(string(inNonce), 10, 64)
		if err_p != nil || inNonceInt < 0 {
			return  code.CodeTypeEncodingError, fmt.Sprintf("Unexpected from nonce. Got %v", inNonce), nil
		}
		inNonceInt = inNonceIntValue
	}

	if (inNonceInt + 1) != nonceInt {
		return code.CodeTypeEncodingError, fmt.Sprintf("nonce should be one more than last nonce. Got %v", nonceS), nil
	}

	// decode the pubkey
	pubkey, err := hex.DecodeString(pubkeyS)
	if err != nil {
		return code.CodeTypeEncodingError, fmt.Sprintf("Pubkey (%s) is invalid hex", pubkeyS), nil
	}

	// decode the power
	power, err := strconv.ParseInt(powerS, 10, 64)
	if err != nil {
		return code.CodeTypeEncodingError, fmt.Sprintf("Power (%s) is not an int", powerS), nil
	}

	if isOnlyCheck {
		return code.CodeTypeOK, "", nil
	}
	// update
	appStore.Set([]byte(ankrtypes.SET_VAL_NONCE), []byte(nonceS))
	return ValidatorManagerInstance().UpdateValidator(types.Ed25519ValidatorUpdate(pubkey, int64(power)), appStore)
}
