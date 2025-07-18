package service

import (
	"admin-panel/logger"
	database "admin-panel/mongo_database"
	"admin-panel/utils"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LpRegisterLogicService struct {
}

func NewLpRegisterLogicService() *LpRegisterLogicService {
	return &LpRegisterLogicService{}
}

func (lprls *LpRegisterLogicService) RegisterItem(installId string, serviceName string, clientName string, chainType string, chainId int64, namespace string) (status bool, err error) {
	if serviceName == "" {
		err = fmt.Errorf("%s ServiceName is empty", installId)
		return
	}
	log.Println("start registration", installId, serviceName, clientName, "🏩🏩")
	//clientName avax|bsc|xrp
	err = NewLpCluster().RestartPod(namespace, "chain-client-", clientName)
	if err != nil {
		log.Println("error restarting pod")
		return
	}
	log.Println("pod restarted successfully", "clientName", clientName, "📱📱📱📱📱")
	log.Println("waiting...")
	time.Sleep(time.Second * 10)

	log.Println("preparing to register client.", lprls.GetServiceUrl(serviceName))
	url := fmt.Sprintf("http://%s:9100/%s-client-%d/lpnode/register_lpnode", serviceName, chainType, chainId)
	log.Println("registration address is", url)
	var dataStr = `{"lpnode_server_url":{"on_transfer_out":"http://lpnode-server:9202/lpnode/chain_client/on_transfer_out","on_transfer_in":"http://lpnode-server:9202/lpnode/chain_client/on_transfer_in","on_confirm":"http://lpnode-server:9202/lpnode/chain_client/on_confirm","on_refunded":"http://lpnode-server:9202/lpnode/chain_client/on_refund"}}`
	dataStr, _ = sjson.Set(dataStr, "chainType", chainType)
	if chainType == "near" {
		nearTokenList := []struct {
			Address string `bson:"address"`
			TokenId string `bson:"tokenId"`
		}{}
		findErr, cursor := database.FindAll("main", "tokens", bson.M{"chainType": "near"})
		if findErr != nil {
			err = errors.WithMessage(findErr, "error querying token pair")
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
		logger.System.Debug("dataStr:", dataStr)
	}
	retryer := utils.RetryerNew().SetOption(&utils.RepetOption{
		Interval: 2000,
		MaxCount: 15,
	})
	log.Println("sendData:________")
	log.Println(dataStr)
	log.Println("___________________")
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
			return errors.New("not ready temporary")
		}
	})

	if err != nil {
		err = errors.WithMessage(err, fmt.Sprintf("register error,install_id: %s", installId))
		return
	}
	status = true
	return
}
func (lprls *LpRegisterLogicService) RegisterItemWithoutRestart(installId string, serviceName string, clientName string, chainType string, chainId int64, namespace string) (status bool, err error) {
	if serviceName == "" {
		err = fmt.Errorf("%s ServiceName is empty", installId)
		return
	}
	log.Println("start registration", installId, serviceName, clientName, "🏩🏩")
	log.Println("preparing to register client.", lprls.GetServiceUrl(serviceName))
	url := fmt.Sprintf("http://%s:9100/%s-client-%d/lpnode/register_lpnode", serviceName, chainType, chainId)
	checkUrl := fmt.Sprintf("http://%s:9100/%s-client-%d/lpnode/register_lpnode_support_duplication", serviceName, chainType, chainId)
	_, checkOk, checkErr := utils.NewHttpCall().PostJsonCall(&utils.HttpCallRequestOption{
		Url:     checkUrl,
		Timeout: 5000,
		JsonStr: `{}`,
		TestOKFun: func(bodyStr string) bool {
			log.Println("bodyis:", bodyStr)
			return gjson.Get(bodyStr, "code").Int() == 200
		},
	})
	if checkErr != nil {
		err = checkErr
		return
	}
	if !checkOk {
		err = errors.New("the chain client does not support duplicate registration")
		return
	}
	log.Println("registration address is", url)
	var dataStr = `{"lpnode_server_url":{"on_transfer_out":"http://lpnode-server:9202/lpnode/chain_client/on_transfer_out","on_transfer_in":"http://lpnode-server:9202/lpnode/chain_client/on_transfer_in","on_confirm":"http://lpnode-server:9202/lpnode/chain_client/on_confirm","on_refunded":"http://lpnode-server:9202/lpnode/chain_client/on_refund"}}`
	dataStr, _ = sjson.Set(dataStr, "chainType", chainType)
	if chainType == "near" {
		nearTokenList := []struct {
			Address string `bson:"address"`
			TokenId string `bson:"tokenId"`
		}{}
		findErr, cursor := database.FindAll("main", "tokens", bson.M{"chainType": "near"})
		if findErr != nil {
			err = errors.WithMessage(findErr, "error querying token pair")
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
		logger.System.Debug("dataStr:", dataStr)
	}
	retryer := utils.RetryerNew().SetOption(&utils.RepetOption{
		Interval: 2000,
		MaxCount: 15,
	})
	log.Println("sendData:________")
	log.Println("url", url)
	log.Println(dataStr)
	log.Println("___________________")
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
			return errors.New("not ready temporary")
		}
	})


	if err != nil {
		err = errors.WithMessage(err, fmt.Sprintf("register error,install_id: %s", installId))
		return
	}
	status = true
	return
}
func (lprls *LpRegisterLogicService) GetServiceUrl(serviceName string) string {
	return serviceName
}

func (lprls *LpRegisterLogicService) UnregisterItem(installId string, serviceName string) (status bool, err error) {
	if serviceName == "" {
		err = fmt.Errorf("%s empty servicename", installId)
		return
	}
	log.Println("start uninstalling", installId)
	return
}
func (lprls *LpRegisterLogicService) IsRegister(installId string, serviceName string) (status bool, err error) {
	if serviceName == "" {
		err = fmt.Errorf("%s empty servicename", installId)
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
		err = errors.New("cannot find RelayApiKey, lp not register to relay yet")
		return
	}
	apiKey = v.RelayApiKey
	return
}
func (lprls *LpRegisterLogicService) GetLpName() (apiKey string, err error) {
	v := struct {
		Id   primitive.ObjectID `bson:"_id"`
		Name string             `bson:"name"`
	}{}
	err = database.FindOne("main", "relayAccounts", bson.M{}, &v)
	if err != nil {
		return
	}
	if v.Name == "" {
		err = errors.New("cannot find relay account name")
		return
	}
	apiKey = v.Name
	return
}
