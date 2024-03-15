package main

import (
	"admin-panel/database_config"
	"admin-panel/logger"
	database "admin-panel/mongo_database"
	"admin-panel/service"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/robertkrimen/otto"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	service.NewLpCluster()
	database_config.Init()
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(2)
	startTime := time.Now().UnixNano() / 1e6
	go func() {
		for {
			logger.System.Info("Link")
			nowTime := time.Now().UnixNano() / 1e6
			if nowTime-startTime > 1000*120 {
				logger.System.Error("exit if db not connected after timeout")
				os.Exit(5)
			}
			logger.System.Debug("preparing db connection [lp_store]...")
			err := database.InitConnect("main", nil)
			if err == nil {
				waitGroup.Done()
				return
			}
			logger.System.Error(err)
			time.Sleep(time.Second * 3)
		}

	}()
	go func() {
		for {
			logger.System.Info("Link")
			nowTime := time.Now().UnixNano() / 1e6
			if nowTime-startTime > 1000*120 {
				logger.System.Error("exit if db not connected after timeout")
				os.Exit(5)
			}
			logger.System.Debug("preparing db connection [businessHistory]...")
			err := database.InitConnect("businessHistory", nil)
			if err == nil {
				waitGroup.Done()
				return
			}
			logger.System.Error(err)
			time.Sleep(time.Second * 3)
		}

	}()
	waitGroup.Wait()
	logger.System.Debug("database connection completed...")

	initDbData()
	err := initIndex()
	if err != nil {
		log.Println("index creation error", err)
	}
	err = initTokenIndex()
	if err != nil {
		log.Println("index creation error", err)
	}
	err = indexWallet()
	if err != nil {
		log.Println("index creation error", err)
	}
	if err != nil {
		log.Println("base data error", err)
	}
	err = InitMonitor()
	if err != nil {
		log.Println("init Monitor failed", err)
	}
	InitInstall()

}
func initDbData() {
	initData, err := ioutil.ReadFile("./init_data/init_data.js")
	if err != nil {
		log.Fatalln("failed reading init data")
	}
	vm := otto.New()
	vm.Run(string(initData))
	value, err := vm.Get("json_init_data")
	if err != nil {
		log.Fatalln("get var from init data failed")
	}
	for _, v := range gjson.Get(value.String(), "data").Array() {
		collectionName := v.Get("collectionName").String()
		listData := v.Get("data").Array()
		if len(listData) > 0 {
			filters := v.Get("filter").Array()
			sets := v.Get("set").Array()
			for _, rowData := range listData { // loop through collection data
				filter := bson.M{}
				set := bson.M{}
				for _, filterItem := range filters {
					if filterItem.Get("type").String() == "string" {
						filter[filterItem.Get("name").String()] = rowData.Get(filterItem.Get("name").String()).String()
					}
					if filterItem.Get("type").String() == "int" {
						filter[filterItem.Get("name").String()] = rowData.Get(filterItem.Get("name").String()).Int()

					}
					for _, setItem := range sets {
						if setItem.Get("type").String() == "string" {
							set[setItem.Get("name").String()] = rowData.Get(setItem.Get("name").String()).String()
						}
						if setItem.Get("type").String() == "int" {
							set[setItem.Get("name").String()] = rowData.Get(setItem.Get("name").String()).Int()
						}
						if setItem.Get("type").String() == "array" {
							tobeSet := bson.A{}
							for _, v := range rowData.Get(setItem.Get("name").String()).Array() {
								tobeSetItem := bson.M{}
								for k, sv := range v.Map() {
									tobeSetItem[k] = sv.String()
								}
								tobeSet = append(tobeSet, tobeSetItem)
							}
							set[setItem.Get("name").String()] = tobeSet
						}
					}
				}
				fmt.Println("fish.......", set)
				database.FindOneAndUpdate("main", collectionName, filter, bson.M{"$set": set})
			}
		}
	}
	time.Sleep(time.Microsecond * 50)
}

func initIndex() (err error) {
	session, err := database.GetSession("main")
	if err != nil {
		log.Println("no valid db connection obtained")
		return
	}

	//chainList
	indexList := session.Collection("chainList").Indexes()

	indexList.DropAll(context.TODO())

	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{"_id", 1},
		},
	}
	chainIdModel := mongo.IndexModel{
		Keys: bson.D{
			{"chainId", -1},
		},
		Options: options.Index().SetUnique(true),
	}
	messages, err := indexList.CreateMany(context.TODO(), []mongo.IndexModel{indexModel, chainIdModel})
	if err != nil {
		return
	}
	log.Println("create result is:")
	strings.Join(messages, "\r\n")

	//bridges
	bridgesIndexList := session.Collection("bridges").Indexes()

	bridgesIndexList.DropAll(context.TODO())

	bridgesindexModel := mongo.IndexModel{
		Keys: bson.D{
			{"_id", 1},
		},
	}
	bridgesNameUniqModel := mongo.IndexModel{
		Keys: bson.D{
			{"bridgeName", -1},
		},
		Options: options.Index().SetUnique(true),
	}
	bridgesUniqlModel := mongo.IndexModel{
		Keys: bson.D{
			{"srcChain_id", -1},
			{"dstChain_id", -1},
			{"srcToken_id", -1},
			{"dstToken_id", -1},
			{"ammName", -1},
		},
		Options: options.Index().SetUnique(true),
	}
	bridgesMessages, err := bridgesIndexList.CreateMany(context.TODO(), []mongo.IndexModel{bridgesindexModel, bridgesNameUniqModel, bridgesUniqlModel})
	if err != nil {
		return
	}
	log.Println("create result is:")
	strings.Join(bridgesMessages, "\r\n")

	return
}

func initTokenIndex() (err error) {
	session, err := database.GetSession("main")
	if err != nil {
		log.Println("no valid db connection obtained")
		return
	}

	indexList := session.Collection("tokens").Indexes()

	indexList.DropAll(context.TODO())

	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{"_id", 1},
		},
	}
	addressLowerModel := mongo.IndexModel{
		Keys: bson.D{
			{"addressLower", -1},
			{"chainId", -1},
		},
		Options: options.Index().SetUnique(true),
	}
	messages, err := indexList.CreateMany(context.TODO(), []mongo.IndexModel{indexModel, addressLowerModel})
	if err != nil {
		return
	}
	log.Println("creation result is [Token]:")
	strings.Join(messages, "\r\n")
	return
}

func indexWallet() (err error) {
	session, err := database.GetSession("main")
	if err != nil {
		log.Println("no valid db connection obtained")
		return
	}

	indexList := session.Collection("wallets").Indexes()

	indexList.DropAll(context.TODO())

	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{"_id", 1},
		},
	}
	nameUniqModel := mongo.IndexModel{
		Keys: bson.D{
			{"walletName", -1},
		},
		Options: options.Index().SetUnique(true)}
	uniqModel := mongo.IndexModel{
		Keys: bson.D{
			{"privateKey", -1},
			{"chainId", -1},
		},
		Options: options.Index().SetUnique(true),
	}
	messages, err := indexList.CreateMany(context.TODO(), []mongo.IndexModel{indexModel, nameUniqModel, uniqModel})
	if err != nil {
		return
	}
	log.Println("creation result is [wallet]:")
	strings.Join(messages, "\r\n")
	return
}
