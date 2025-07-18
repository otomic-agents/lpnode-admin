package main

import (
	"admin-panel/logger"
	database "admin-panel/mongo_database"
	"admin-panel/types"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/tidwall/sjson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChainListRow struct {
	Id          primitive.ObjectID  `bson:"_id"`
	Name        string              `bson:"name"`
	ChainName   string              `bson:"chainName"`
	Image       string              `bson:"image"`
	ServiceName string              `bson:"serviceName"`
	DeployName  string              `bson:"deployName"`
	ChainType   string              `bson:"chainType"`
	ChainId     int64               `bson:"chainId"`
	EnvList     []map[string]string `bson:"envList"`
}

func InitInstall() (err error) {
	log.Println("Init Install")
	err = install_init_chain_client()
	if err != nil {
		logger.System.Error(err)
	}
	err = install_init_market_adapter()
	if err != nil {
		logger.System.Warn(err)
	}
	err = install_init_amm_client()
	if err != nil {
		logger.System.Warn(err)
	}
	return
}
func install_init_chain_client() (err error) {
	var results []ChainListRow
	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}})
	err, cursor := database.FindAllOpt("main", "chainList", bson.M{}, opts)
	if err != nil {
		return
	}
	if err = cursor.All(context.TODO(), &results); err != nil {
		err = errors.WithMessage(err, "cursor all error")
		return
	}
	for _, result := range results {
		// log.Println(cursor.Current.String())
		cursor.Decode(&result)
		install_init_chain_client_item(result)
	}
	return
}
func install_init_chain_client_item(result ChainListRow) (err error) {
	logger.System.Info("setup client row", result.ChainName, result.ChainId)
	var v struct {
		Id primitive.ObjectID `bson:"_id"`
	} = struct {
		Id primitive.ObjectID `bson:"_id"`
	}{}
	chainName := strings.ToLower(result.ChainName)
	serviceName := result.ServiceName
	// log.Println(bson.M{"installType": "ammClient", "name": result.ChainName})
	database.FindOne("main", "install", bson.M{"installType": "ammClient", "name": chainName}, &v)
	if v.Id.Hex() != types.MongoEmptyIdHex {
		// if exists, continue to update instead of returning error
		// err = fmt.Errorf("the ammClient program already exists,chainName:%s", chainName)
		// logger.System.Warn(err)
		// return
	}
	installContextJson := `{}`
	name := chainName
	installType := "ammClient"
	installContextJson, _ = sjson.Set(installContextJson, "deployment.namespace", os.Getenv("POD_NAMESPACE"))
	installContextJson, _ = sjson.Set(installContextJson, "deployment.name", name)
	installContextJson, _ = sjson.Set(installContextJson, "deployment.image", result.Image)
	envArraySet := bson.A{}
	for _, v := range result.EnvList {
		for key, iv := range v {
			envArraySet = append(envArraySet, bson.M{
				"name":  key,
				"value": iv,
			})
		}
		// envArraySet = append(envArraySet, v)
	}

	// change Insert to FindOneAndUpdate
	_, err = database.FindOneAndUpdate(
		"main",
		"install",
		bson.M{"installType": "ammClient", "name": chainName},
		bson.M{"$set": bson.M{
			"installType":    installType,
			"chainType":      result.ChainType,
			"chainId":        result.ChainId,
			"name":           name,
			"deployName":     result.DeployName,
			"configStatus":   0,
			"lastinstall":    time.Now().UnixNano() / 1e6,
			"status":         1,
			"stderr":         "",
			"stdout":         "init install",
			"yaml":           "",
			"envList":        envArraySet,
			"installContext": installContextJson,
			"namespace":      os.Getenv("POD_NAMESPACE"),
			"un_stderr":      "",
			"un_stdout":      "",
			"serviceName":    serviceName,
		}},
	)
	return
}
func install_init_amm_client() (err error) {
	var v struct {
		Id primitive.ObjectID `bson:"_id"`
	} = struct {
		Id primitive.ObjectID `bson:"_id"`
	}{}
	database.FindOne("main", "install", bson.M{"installType": "amm"}, &v)
	if v.Id.Hex() != types.MongoEmptyIdHex {
		// if exists, continue to update instead of returning error
		// err = errors.New("the amm program already exists")
		// return
	}
	installContextJson := `{}`
	name := "amm-01"
	deployName := "amm-amm-01"
	installType := "amm"
	image := os.Getenv("AMM_APP_DISPLAY_IMAGE")
	installContextJson, _ = sjson.Set(installContextJson, "deployment.namespace", os.Getenv("POD_NAMESPACE"))
	installContextJson, _ = sjson.Set(installContextJson, "deployment.name", name)
	installContextJson, _ = sjson.Set(installContextJson, "deployment.image", image)

	// change Insert to FindOneAndUpdate
	_, err = database.FindOneAndUpdate(
		"main",
		"install",
		bson.M{"installType": "amm"},
		bson.M{"$set": bson.M{
			"installType":  installType,
			"name":         name,
			"deployName":   deployName,
			"configStatus": 0,
			"lastinstall":  time.Now().UnixNano() / 1e6,
			"status":       1,
			"stderr":       "",
			"stdout":       "init install",
			"yaml":         "",
			"envList": bson.A{
				bson.M{"name": "STATUS_KEY", "value": fmt.Sprintf("amm-status-report-%s", name)},
			},
			"installContext": installContextJson,
			"namespace":      os.Getenv("POD_NAMESPACE"),
			"un_stderr":      "",
			"un_stdout":      "",
		}},
	)
	return
}
func install_init_market_adapter() (err error) {
	var v struct {
		Id primitive.ObjectID `bson:"_id"`
	} = struct {
		Id primitive.ObjectID `bson:"_id"`
	}{}
	database.FindOne("main", "install", bson.M{"installType": "market"}, &v)
	if v.Id.Hex() != types.MongoEmptyIdHex {
		// if exists, continue to update instead of returning error
		// err = errors.New("the market program already exists")
		// return
	}
	installContextJson := `{}`
	name := "price"
	deployName := "amm-market-price"
	installType := "market"
	image := os.Getenv("MARKET_APP_DISPLAY_IMAGE")
	installContextJson, _ = sjson.Set(installContextJson, "deployment.namespace", os.Getenv("POD_NAMESPACE"))
	installContextJson, _ = sjson.Set(installContextJson, "deployment.name", name)
	installContextJson, _ = sjson.Set(installContextJson, "deployment.image", image)

	// change Insert to FindOneAndUpdate
	_, err = database.FindOneAndUpdate(
		"main",
		"install",
		bson.M{"installType": "market"},
		bson.M{"$set": bson.M{
			"installType":  installType,
			"name":         "price",
			"deployName":   deployName,
			"configStatus": 0,
			"lastinstall":  time.Now().UnixNano() / 1e6,
			"status":       1,
			"stderr":       "",
			"stdout":       "init install",
			"yaml":         "",
			"envList": bson.A{
				bson.M{"name": "STATUS_KEY", "value": fmt.Sprintf("amm-market-status-report-%s", "price")},
			},
			"installContext": installContextJson,
			"namespace":      os.Getenv("POD_NAMESPACE"),
			"un_stderr":      "",
			"un_stdout":      "",
		}},
	)
	return
}
