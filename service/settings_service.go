package service

import (
	"admin-panel/logger"
	"admin-panel/types"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/pkg/errors"
)

type SettingsService struct {
}

func NewSettingsService() *SettingsService {
	return &SettingsService{}
}

func (settingsService *SettingsService) SaveSetting(settingsConfig *types.SettingsConfig) (ret bool, err error) {
	nameSpace := os.Getenv("NAMESPACE")
	val := reflect.ValueOf(*settingsConfig)

	if val.NumField() <= 0 {
		err = errors.New("there must be something modified")
		return
	}
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		value := val.Field(i)
		envTag := field.Tag.Get("env")
		if value.Kind() == reflect.String {
			log.Println(envTag, value.String())
			if value.String() == "" {
				err = errors.New(fmt.Sprintf("%s cannot be empty", envTag))
				return
			}
			if envTag == "RELAY_URI" {
				ok, testErr := settingsService.TestRpc(value.String())
				if testErr != nil {
					err = testErr
					return
				}
				if !ok {
					err = errors.New(fmt.Sprintf("RPC is not available:%s", value.String()))
					return
				}
			}
			pathed, pathErr := NewLpCluster().PathEnv(nameSpace, "otmoiclp", "otmoiclp", envTag, value.String())
			if pathErr != nil {
				err = pathErr
				return
			}
			if !pathed {
				err = errors.New("PathEnv err")
				return
			}
		}
	}
	ret = true
	return
}

func (settingsService *SettingsService) TestRpc(rpc string) (ok bool, err error) {
	rpcURL := rpc
	ok = false
	client := &http.Client{
		Timeout: 8 * time.Second,
	}
	log.Println("request to :", rpcURL)
	response, err := client.Get(rpcURL)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error connecting to RPC node:%s", err))
		return
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusOK {
		logger.System.Info("status code:", response.StatusCode)
		ok = true
		return
	}
	err = errors.New("RPC node is not available")
	return
}
