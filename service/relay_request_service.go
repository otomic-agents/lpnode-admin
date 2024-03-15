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
		err = errors.New("cannot find relay url")
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
		err = errors.WithMessage(err, fmt.Sprintf("incorrect server status code response:%d", resp.StatusCode))
		return
	}

	err = json.Unmarshal([]byte(body), &res)
	if err != nil {
		err = errors.WithMessage(err, "json unmarshal error:")
		return
	}
	if res.Code != 200 && res.Code != 30210 {
		err = errors.WithMessage(utils.GetNoEmptyError(err), "relay registration failed, relay code != 200 and != 30210")
		return
	}
	if res.RelayApiKey == "" || res.LpnodeApiKey == "" {
		err = errors.WithMessage(utils.GetNoEmptyError(err), "relay registration failed, relayapikey or lpnodeapikey empty")
		return
	}
	return
}
