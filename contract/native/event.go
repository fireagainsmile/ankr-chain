package native

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	ankrcmm "github.com/Ankr-network/ankr-chain/common"
	contContr "github.com/Ankr-network/ankr-chain/context"
	"github.com/tendermint/tendermint/libs/log"
)

func TrigEvent(evSrc string, jsonData string, log log.Logger, contextCont contContr.ContextContract) int32 {
	evSrcSegs := strings.Split(evSrc, "(")
	if len(evSrcSegs) != 2 {
		log.Error("TrigEvent event invalid evSrc", "evSrc", evSrc)
		return -1
	}

	method := evSrcSegs[0]
	//argStr := strings.TrimRight(evSrcSegs[1], ")")
	//argTyps := strings.Split(argStr, ",")

	var params []*ankrcmm.Param
	err :=  json.Unmarshal([]byte(jsonData), &params)
	if err != nil {
		log.Error("TrigEvent event json.Unmarshal err", "evData", jsonData, "err", err)
		return -1
	}

	tags := make(map[string]string)
	tags["contract.addr"]   = contextCont.ContractAddr()
	tags["contract.method"] = method
	for _, param := range params {
		tagName  := fmt.Sprintf("contract.method.%s", param.Name)
		tagValue :=  param.Value.(string)
		tags[tagName] = tagValue
	}

	contextCont.PublishWithTags(context.Background(),"contract", tags)

	return 0
}
