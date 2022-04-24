package main

import (
	"encoding/json"
	"fmt"
	"huaweicloud.com/go-runtime/events/apig"
	"huaweicloud.com/go-runtime/go-api/context"
	"huaweicloud.com/go-runtime/pkg/runtime"
	"strings"

	"github.com/riba2534/wecomchan/go-scf/consts"
	"github.com/riba2534/wecomchan/go-scf/dal"
	"github.com/riba2534/wecomchan/go-scf/service"
	"github.com/riba2534/wecomchan/go-scf/utils"
)

func init() {
	consts.FUNC_NAME = utils.GetEnvDefault("FUNC_NAME", "")
	consts.SEND_KEY = utils.GetEnvDefault("SEND_KEY", "")
	consts.WECOM_CID = utils.GetEnvDefault("WECOM_CID", "")
	consts.WECOM_SECRET = utils.GetEnvDefault("WECOM_SECRET", "")
	consts.WECOM_AID = utils.GetEnvDefault("WECOM_AID", "")
	consts.WECOM_TOUID = utils.GetEnvDefault("WECOM_TOUID", "@all")
	if consts.FUNC_NAME == "" || consts.SEND_KEY == "" || consts.WECOM_CID == "" ||
		consts.WECOM_SECRET == "" || consts.WECOM_AID == "" || consts.WECOM_TOUID == "" {
		fmt.Printf("os.env load Fail, please check your os env.\nFUNC_NAME=%s\nSEND_KEY=%s\nWECOM_CID=%s\nWECOM_SECRET=%s\nWECOM_AID=%s\nWECOM_TOUID=%s\n", consts.FUNC_NAME, consts.SEND_KEY, consts.WECOM_CID, consts.WECOM_SECRET, consts.WECOM_AID, consts.WECOM_TOUID)
		panic("os.env param error")
	}
	fmt.Println("os.env load success!")
}

func HTTPHandler(payload []byte, ctx context.RuntimeContext) (interface{}, error) {
	var event apig.APIGTriggerEvent
	err := json.Unmarshal(payload, &event)
	if err != nil {
		fmt.Println("Unmarshal failed")
		return "invalid data", err
	}

	path := event.Path
	fmt.Println("req->", event.String())
	var result interface{}
	if strings.HasPrefix(path, "/"+consts.FUNC_NAME) {
		result = service.WeComChanService(ctx, event)
	} else {
		// 匹配失败返回原始HTTP请求
		result = event
	}

	return apig.APIGTriggerResponse{
		IsBase64Encoded: false,
		StatusCode:      200,
		Headers: map[string]string{
			"content-type": "application/json",
		},
		Body: utils.MarshalToStringParam(result),
	}, err
}

func main() {
	dal.Init()
	runtime.Register(HTTPHandler)
}
