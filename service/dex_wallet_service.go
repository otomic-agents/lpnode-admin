package service

import (
	"admin-panel/logger"
	database "admin-panel/mongo_database"
	"admin-panel/types"
	"admin-panel/utils"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DexWalletLogicService struct {
}

func NewDexWalletLogicService() *DexWalletLogicService {
	return &DexWalletLogicService{}
}

func (*DexWalletLogicService) FindOneByFilter(filter bson.M) (ret types.DBWalletRow, err error) {
	ret = types.DBWalletRow{}
	err = database.FindOne("main", "wallets", filter, &ret)
	if err != nil && strings.Contains(err.Error(), "no documents in result") {
		err = nil
	}
	return
}

func (*DexWalletLogicService) CreateByBsonMap(data interface{}) (err error) {
	err = database.Insert("main", "wallets", data)
	if err != nil {
		return
	}
	return
}
func (*DexWalletLogicService) ListAll(data bson.M) (ret []types.DBWalletRow, err error) {
	emptyList := []types.DBWalletRow{}
	ret = emptyList
	err, cursor := database.FindAll("main", "wallets", data)
	if err != nil {
		return
	}
	var results []types.DBWalletRow
	if err = cursor.All(context.TODO(), &results); err != nil {
		return
	}
	for _, result := range results {
		cursor.Decode(&result)
	}
	ret = results
	return
}
func (*DexWalletLogicService) DeleteById(inputId string) (count int64, err error) {
	objectId, idErr := primitive.ObjectIDFromHex(inputId)
	count = 0
	if idErr != nil {
		err = idErr
		return
	}
	delCount, delErr := database.DeleteOne("main", "wallets", bson.M{
		"_id": objectId,
	})
	if delErr != nil {
		err = delErr
		return
	}
	count = delCount
	return
}
func (dwls *DexWalletLogicService) GetVaultAccessToken() (accessToken string, err error) {
	acctokenUrl := fmt.Sprintf("http://%s/permission/v1alpha1/access", os.Getenv("OS_SYSTEM_SERVER"))
	OS_APP_KEY := os.Getenv("OS_API_KEY")
	OS_APP_SECRET := os.Getenv("OS_API_SECRET")
	timestamp := time.Now().UnixNano() / 1e6 / 1000
	srcText := fmt.Sprintf("%s%d%s", OS_APP_KEY, timestamp, OS_APP_SECRET)
	tokenBytes, err := bcrypt.GenerateFromPassword([]byte(srcText), 10)
	tokenStr := string(tokenBytes)
	tobeSend := "{}"
	tobeSend, _ = sjson.Set(tobeSend, "app_key", OS_APP_KEY)
	tobeSend, _ = sjson.Set(tobeSend, "timestamp", timestamp)
	tobeSend, _ = sjson.Set(tobeSend, "token", tokenStr)
	tobeSend, _ = sjson.SetRaw(tobeSend, "perm", `{"group":"secret.vault","dataType":"key","version":"v1","ops":["Info","Sign"]}`)
	logger.System.Debug("tobeSend is:", "\n", tobeSend)
	resultOption := &utils.HttpCallRequestOption{
		Url:     acctokenUrl,
		Timeout: 1000 * 10,
		JsonStr: tobeSend,
		TestOKFun: func(bodyStr string) bool {
			log.Println("bodyis:", bodyStr)
			code := gjson.Get(bodyStr, "code").Int()
			return code == 0
		},
	}
	accessTokenBody, ok, requestAccessTokenErr := utils.NewHttpCall().PostJsonCall(resultOption)
	if requestAccessTokenErr != nil {
		err = requestAccessTokenErr
		return
	}
	if !ok {
		err = errors.WithMessage(utils.GetNoEmptyError(err), "assert return result error")
		return
	}
	logger.System.Debug("ok is ", ok)
	accessToken = gjson.Get(accessTokenBody, "data.access_token").String()
	return
}
func (dwls *DexWalletLogicService) GetVaultList(accessToken string) (res []types.VaultDataRow, err error) {
	res = make([]types.VaultDataRow, 0)
	listVaultUrl := fmt.Sprintf("http://%s/system-server/v1alpha1/key/secret.vault/v1/Info", os.Getenv("OS_SYSTEM_SERVER"))
	listTobeSend := "{\"t\":1}"
	listVaultRequestOption := &utils.HttpCallRequestOption{
		Header:  map[string]string{"X-Access-Token": accessToken},
		Url:     listVaultUrl,
		Timeout: 1000 * 10,
		JsonStr: listTobeSend,
		TestOKFun: func(bodyStr string) bool {
			log.Println("bodyis:", bodyStr)
			code := gjson.Get(bodyStr, "code").Int()
			return code == 0
		},
	}
	listVaultBody, listOk, listVaultErr := utils.NewHttpCall().PostJsonCall(listVaultRequestOption)
	if listVaultErr != nil {
		err = listVaultErr
		return
	}
	if !listOk {
		err = errors.WithMessage(utils.GetNoEmptyError(err), "get vault list error, assert return error")
		return
	}
	for _, item := range gjson.Get(listVaultBody, "data.data").Array() {
		res = append(res, types.VaultDataRow{
			Address:    item.Get("address").String(),
			HostType:   item.Get("host_type").String(),
			Id:         item.Get("id").String(),
			Name:       item.Get("name").String(),
			SecertType: item.Get("secert_type").String(),
		})
	}
	return
}
func (dwls *DexWalletLogicService) GetVault(storeId string) (res *types.VaultDataRow, err error) {
	res = nil
	accessToken, err := dwls.GetVaultAccessToken()
	if err != nil {
		err = errors.WithMessage(err, "get accessToken error")
		return
	}
	list, err := dwls.GetVaultList(accessToken)
	if err != nil {
		err = errors.WithMessage(err, "GetVaultList error")
		return
	}
	for _, item := range list {
		if item.Id == storeId {
			res = &item
		}
	}
	return
}
func (dwls *DexWalletLogicService) getSecretVaultAccessToken() (accessToken string, err error) {
	acctokenUrl := fmt.Sprintf("http://%s/permission/v1alpha1/access", os.Getenv("OS_SYSTEM_SERVER"))
	OS_APP_KEY := os.Getenv("OS_API_KEY")
	OS_APP_SECRET := os.Getenv("OS_API_SECRET")
	timestamp := time.Now().UnixNano() / 1e6 / 1000
	srcText := fmt.Sprintf("%s%d%s", OS_APP_KEY, timestamp, OS_APP_SECRET)
	tokenBytes, err := bcrypt.GenerateFromPassword([]byte(srcText), 10)
	tokenStr := string(tokenBytes)
	tobeSend := "{}"
	tobeSend, _ = sjson.Set(tobeSend, "app_key", OS_APP_KEY)
	tobeSend, _ = sjson.Set(tobeSend, "timestamp", timestamp)
	tobeSend, _ = sjson.Set(tobeSend, "token", tokenStr)
	tobeSend, _ = sjson.SetRaw(tobeSend, "perm", `{"group":"secret.infisical","dataType":"secret","version":"v1","ops":["CreateSecret?workspace=otmoic","RetrieveSecret?workspace=otmoic"]}`)
	logger.System.Debug("tobeSend is:", "\n", tobeSend)
	resultOption := &utils.HttpCallRequestOption{
		Url:     acctokenUrl,
		Timeout: 1000 * 10,
		JsonStr: tobeSend,
		TestOKFun: func(bodyStr string) bool {
			log.Println("bodyis:", bodyStr)
			code := gjson.Get(bodyStr, "code").Int()
			return code == 0
		},
	}
	accessTokenBody, ok, requestAccessTokenErr := utils.NewHttpCall().PostJsonCall(resultOption)
	if requestAccessTokenErr != nil {
		err = requestAccessTokenErr
		return
	}
	if !ok {
		err = errors.WithMessage(utils.GetNoEmptyError(err), "assert return result error:")
		return
	}
	logger.System.Debug("ok is ", ok)
	accessToken = gjson.Get(accessTokenBody, "data.access_token").String()
	return
}

func (dwls *DexWalletLogicService) SaveToSecretVault(vaultName string, privateKey string, accessToken string) (id string, err error) {
	id = ""
	createSecretUrl := fmt.Sprintf("http://%s/system-server/v1alpha1/secret/secret.infisical/v1/CreateSecret?workspace=otmoic", os.Getenv("OS_SYSTEM_SERVER"))
	createTobeSend := "{}"
	createTobeSend, _ = sjson.Set(createTobeSend, "name", vaultName)
	createTobeSend, _ = sjson.Set(createTobeSend, "value", privateKey)
	createTobeSend, _ = sjson.Set(createTobeSend, "env", "prod")
	createRequestOption := &utils.HttpCallRequestOption{
		Header:  map[string]string{"X-Access-Token": accessToken},
		Url:     createSecretUrl,
		Timeout: 1000 * 10,
		JsonStr: createTobeSend,
		TestOKFun: func(bodyStr string) bool {
			log.Println("bodyis:", bodyStr)
			code := gjson.Get(bodyStr, "code").Int()
			sub_code := gjson.Get(bodyStr, "data.code").Int()
			return code == 0 && sub_code == 0
		},
	}
	fmt.Println(createTobeSend)
	_, listOk, createErr := utils.NewHttpCall().PostJsonCall(createRequestOption)
	if createErr != nil {
		err = createErr
		return
	}
	if !listOk {
		err = errors.WithMessage(utils.GetNoEmptyError(err), "create secret error, assert return error:")
		fmt.Println(err)
		return
	}
	id = vaultName
	return
}
func (dwls *DexWalletLogicService) GetSecretFromSecretVault(vaultName string, accessToken string) (id string, err error) {
	id = ""
	createSecretUrl := fmt.Sprintf("http://%s/system-server/v1alpha1/secret/secret.infisical/v1/RetrieveSecret?workspace=otmoic", os.Getenv("OS_SYSTEM_SERVER"))
	createTobeSend := "{}"
	createTobeSend, _ = sjson.Set(createTobeSend, "name", vaultName)
	createTobeSend, _ = sjson.Set(createTobeSend, "env", "prod")
	createRequestOption := &utils.HttpCallRequestOption{
		Header:  map[string]string{"X-Access-Token": accessToken},
		Url:     createSecretUrl,
		Timeout: 1000 * 10,
		JsonStr: createTobeSend,
		TestOKFun: func(bodyStr string) bool {
			log.Println("bodyis:", bodyStr)
			code := gjson.Get(bodyStr, "code").Int()
			sub_code := gjson.Get(bodyStr, "data.code").Int()
			return code == 0 && sub_code == 0
		},
	}
	fmt.Println(createTobeSend)
	_, listOk, createErr := utils.NewHttpCall().PostJsonCall(createRequestOption)
	if createErr != nil {
		err = createErr
		return
	}
	if !listOk {
		err = errors.WithMessage(utils.GetNoEmptyError(err), "query secret, assert return error:")
		fmt.Println(err)
		return
	}
	id = vaultName
	return
}

func (dwls *DexWalletLogicService) StoreToSecretVault(vaultName string, privateKey string) (res string, err error) {
	accessToken, err := dwls.getSecretVaultAccessToken()
	fmt.Println(accessToken, err)
	fmt.Println("privateKey is ", privateKey)
	fmt.Println("walletName is ", vaultName)
	id, storeErr := dwls.SaveToSecretVault(vaultName, privateKey, accessToken)
	if storeErr != nil {
		err = errors.WithMessage(storeErr, "save to vault error:")
		return
	}
	res = id
	return
}

func (dwls *DexWalletLogicService) GetFromSecretVault(vaultName string) (res string, err error) {
	accessToken, err := dwls.getSecretVaultAccessToken()
	if err != nil {
		err = errors.WithMessage(err, "getSecretVaultAccessToken error:")
		return
	}
	fmt.Println("walletName is ", vaultName)
	secretId, err := dwls.GetSecretFromSecretVault(vaultName, accessToken)
	if err != nil {
		err = errors.WithMessage(err, "get secret from vault error:")
		return
	}
	res = secretId
	return
}
