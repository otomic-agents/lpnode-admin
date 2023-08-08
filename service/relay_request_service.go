package service

import (
	"admin-panel/types"
	"admin-panel/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
)

type RelayRequestService struct {
}

func NewRelayRequestService() *RelayRequestService {
	return &RelayRequestService{}
}

func (*RelayRequestService) RegisterAccount(name string, profile string) (res types.RelayRegisterResponse, err error) {
	res = types.RelayRegisterResponse{}
	relayUrl := os.Getenv("RELAY_ACCESS_URL")
	if relayUrl == "" {
		err = errors.New("Relay url无法获取")
		return
		// relayUrl = "http://localhost:18009"
	}
	type registerPayload = struct {
		Name         string `json:"name"`
		Profile      string `json:"profile"`
		LpnodeApiKey string `json:"lpnode_api_key"`
	}
	strTime := fmt.Sprintf("%d", (int64(time.Now().UnixNano())))
	sEnc := base64.StdEncoding.EncodeToString([]byte(strTime))
	sendPayload := registerPayload{
		Name:         name,
		Profile:      profile,
		LpnodeApiKey: sEnc,
	}
	log.Println(sendPayload)
	url := fmt.Sprintf("%s/relay-admin-panel/lpnode_admin_panel/register_lp", relayUrl)
	log.Println(url)
	
	resp, body, errs := gorequest.New().Post(url).
		Send(sendPayload).
		End()
	if len(errs) > 0 {
		err = errors.New(errs[0].Error())
		return
	}
	if resp.StatusCode != 200 {
		err = utils.GetNoEmptyError(err)
		err = errors.WithMessage(err, fmt.Sprintf("不正确的服务器状态响应码:%d", resp.StatusCode))
		return
	}

	err = json.Unmarshal([]byte(body), &res)
	if res.Code != 200 {
		err = errors.WithMessage(utils.GetNoEmptyError(err), "relay注册失败, relay code !=200")
		return
	}
	if res.RelayApiKey == "" || res.LpnodeApiKey == "" {
		err = errors.WithMessage(utils.GetNoEmptyError(err), "relay注册失败, RelayApiKey LpnodeApiKey empty")
		return
	}
	if err != nil {
		return
	}
	return
}
