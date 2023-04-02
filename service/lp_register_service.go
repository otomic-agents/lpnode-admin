package service

import (
	"admin-panel/logger"
	database "admin-panel/mongo_database"
	"admin-panel/utils"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"strings"
	"time"
)

type LpRegisterLogicService struct {
}

func NewLpRegisterLogicService() *LpRegisterLogicService {
	return &LpRegisterLogicService{}
}

func (lprls *LpRegisterLogicService) RegisterItem(installId string, serviceName string, clientName string, chainType string, chainId int64, namespace string) (status bool, err error) {
	if serviceName == "" {
		err = fmt.Errorf("%s空的ServiceName", installId)
		return
	}
	log.Println("开始注册", installId, serviceName, clientName, "🏩🏩")
	//clientName avax|bsc|xrp
	err = NewLpCluster().RestartPod(namespace, "obridge-chain-client-", clientName)
	if err != nil {
		log.Println("重启pod发生了错误")
		return
	}
	log.Println("成功重启了pod", "clientName", clientName, "📱📱📱📱📱")
	log.Println("等待...")
	time.Sleep(time.Second * 10)

	log.Println("准备注册Client.", lprls.GetServiceUrl(serviceName))
	url := fmt.Sprintf("http://%s:9100/%s-client-%d/lpnode/register_lpnode", serviceName, chainType, chainId)
	log.Println("注册的地址是", url, "🥩🥩🥩🥩🥩🥩🥩")
	var dataStr = `{"lpnode_server_url":{"on_transfer_out":"http://lpnode-server:9202/lpnode/chain_client/on_transfer_out","on_transfer_in":"http://lpnode-server:9202/lpnode/chain_client/on_transfer_in","on_confirm":"http://lpnode-server:9202/lpnode/chain_client/on_confirm","on_refunded":"http://lpnode-server:9202/lpnode/chain_client/on_refund"}}`
	dataStr, _ = sjson.Set(dataStr, "chainType", chainType)
	if chainType == "near" {
		nearTokenList := []struct {
			Address string `bson:"address"`
			TokenId string `bson:"tokenId"`
		}{}
		findErr, cursor := database.FindAll("main", "tokens", bson.M{"chainType": "near"})
		if findErr != nil {
			err = errors.WithMessage(findErr, "查询币对发生了错误")
			return
		}
		if err = cursor.All(context.TODO(), &nearTokenList); err != nil {
			return
		}
		for _, nearToken := range nearTokenList {
			cursor.Decode(&nearToken)
			tokenHex, convertErr := utils.Base58ToHexString(nearToken.Address)
			if convertErr != nil {
				err = errors.WithMessage(err, "convert token error")
				return
			}
			key := strings.Replace(nearToken.TokenId, ".", "\\.", 1)
			dataStr, _ = sjson.Set(dataStr, fmt.Sprintf("token_map.%s.receiver_id", key), tokenHex)
		}
		logger.System.Debug("需要传送的值是", dataStr)
	}
	retryer := utils.RetryerNew().SetOption(&utils.RepetOption{
		Interval: 2000,
		MaxCount: 15,
	})
	err = retryer.Repet(func() error {
		_, ok, err := utils.NewHttpCall().PostJsonCall(&utils.HttpCallRequestOption{
			Url:     url,
			Timeout: 5000,
			JsonStr: dataStr,
			TestOKFun: func(bodyStr string) bool {
				log.Println("bodyis:", bodyStr)
				return gjson.Get(bodyStr, "code").Int() == 200
			},
		})
		if err != nil {
			return err
		}
		if ok {
			return nil
		} else {
			return errors.New("暂时未就绪")
		}
	})

	log.Println("sendData:________")
	log.Println(dataStr)
	log.Println("___________________")
	if err != nil {
		err = errors.WithMessage(err, fmt.Sprintf("注册出错%s", installId))
		return
	}
	log.Println("🩳🩳🩳🩳")
	status = true
	return
}
func (lprls *LpRegisterLogicService) GetServiceUrl(serviceName string) string {
	return serviceName
}

func (lprls *LpRegisterLogicService) UnregisterItem(installId string, serviceName string) (status bool, err error) {
	if serviceName == "" {
		err = fmt.Errorf("%s空的ServiceName", installId)
		return
	}
	log.Println("开始卸载", installId, "📗📗📗📗📗📗📗📗📗📗📗")
	return
}
func (lprls *LpRegisterLogicService) IsRegister(installId string, serviceName string) (status bool, err error) {
	if serviceName == "" {
		err = fmt.Errorf("%s空的ServiceName", installId)
		return
	}
	return
}
func (lprls *LpRegisterLogicService) GetRelayApiKey() (apiKey string, err error) {
	v := struct {
		Id          primitive.ObjectID `bson:"_id"`
		RelayApiKey string             `bson:"relayApiKey"`
	}{}
	err = database.FindOne("main", "relayAccounts", bson.M{}, &v)
	if err != nil {
		return
	}
	if v.RelayApiKey == "" {
		err = errors.New("没有找到RelayApiKey,Lp还没有注册到relay")
		return
	}
	apiKey = v.RelayApiKey
	return
}
