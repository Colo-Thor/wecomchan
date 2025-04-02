package service

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"huaweicloud.com/go-runtime/events/apig"
	"io/ioutil"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/riba2534/wecomchan/go-scf/consts"
	"github.com/riba2534/wecomchan/go-scf/dal"
	"github.com/riba2534/wecomchan/go-scf/model"
	"github.com/riba2534/wecomchan/go-scf/utils"
	"huaweicloud.com/go-runtime/go-api/context"
)

func WeComChanService(ctx context.RuntimeContext, event apig.APIGTriggerEvent) map[string]interface{} {
	sendKey := getQuery("sendkey", event)
	msgType := getQuery("msg_type", event)
	msg := getQuery("msg", event)
	if msgType == "" || msg == "" {
		return utils.MakeResp(-1, "param error")
	}
	if sendKey != consts.SEND_KEY {
		return utils.MakeResp(-1, "sendkey error")
	}
	toUser := getQuery("to_user", event)
	if toUser == "" {
		toUser = consts.WECOM_TOUID
	}
	if err := postWechatMsg(dal.AccessToken, msg, msgType, toUser); err != nil {
		return utils.MakeResp(0, err.Error())
	}
	return utils.MakeResp(0, "success")
}

func postWechatMsg(accessToken, msg, msgType, toUser string) error {
	content := &model.WechatMsg{
		ToUser:                 toUser,
		AgentId:                consts.WECOM_AID,
		MsgType:                msgType,
		DuplicateCheckInterval: 600,
		Text: &model.MsgText{
			Content: msg,
		},
	}
	b, _ := jsoniter.Marshal(content)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{Timeout: 10 * time.Second, Transport: tr}
	req, _ := http.NewRequest("POST", fmt.Sprintf(consts.WeComMsgSendURL, accessToken), bytes.NewBuffer(b))
	req.Header.Set("Content-type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("[postWechatMsg] failed, err=", err)
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println("postWechatMsg statusCode is not 200")
		return errors.New("statusCode is not 200")
	}
	respBodyBytes, _ := ioutil.ReadAll(resp.Body)
	postResp := &model.PostResp{}
	if err := jsoniter.Unmarshal(respBodyBytes, postResp); err != nil {
		fmt.Println("postWechatMsg json Unmarshal failed, err=", err)
		return err
	}
	if postResp.Errcode != 0 {
		fmt.Println("postWechatMsg postResp.Errcode != 0, err=", postResp.Errmsg)
		return errors.New(postResp.Errmsg)
	}
	return nil
}

func getQuery(key string, event apig.APIGTriggerEvent) string {
	switch event.HttpMethod {
	case "GET":
		value := event.QueryStringParameters[key]
		if len(value) > 0 && value != "" {
			return value
		}
		return ""
	case "POST":
		if event.IsBase64Encoded {
			return jsoniter.Get([]byte(event.GetRawBody()), key).ToString()
		}
		return jsoniter.Get([]byte(event.Body), key).ToString()
	default:
		return ""
	}
}
