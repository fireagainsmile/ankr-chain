package native

import (
	"encoding/json"
	"fmt"
	"testing"


	ankrcmm "github.com/Ankr-network/ankr-chain/common"
	"github.com/stretchr/testify/assert"
)

func TestEventJsonArg(t *testing.T) {
	fromAddr := "dfghq"
	toAddr   := "wertr"
	amount  :=  "12712"
	fromAddrParam := fmt.Sprintf("\"%s\"", fromAddr)
	toAddrParam   := fmt.Sprintf("\"%s\"", toAddr)
	amountParam   := fmt.Sprintf("\"%s\"", amount)

	jsonArgFromat := "[{\"index\":1,\"Name\":\"fromAddr\",\"ParamType\":\"string\",\"Value\":%s}," +
		"{\"index\":2,\"Name\":\"toAddr\",\"ParamType\":\"string\",\"Value\":%s}," +
		"{\"index\":3,\"Name\":\"amount\",\"ParamType\":\"string\",\"Value\":%s}]"

	jsonArg := fmt.Sprintf(jsonArgFromat, fromAddrParam, toAddrParam, amountParam)

	t.Logf("jsonArg: %s", jsonArg)

	var params []*ankrcmm.Param
	err :=  json.Unmarshal([]byte(jsonArg), &params)
	assert.Equal(t, nil, err)

	t.Logf("params=%v", params)
}
