package database_config

import (
	"fmt"
	"log"
	"os"
)

type MongoDbConnectInfoItem struct {
	Url      string `bson:"mongoUrl"`
	Database string `bson:"mongoDatabase"`
	UserName string ``
	Password string
}

var MongoDataBaseConfigIns = make(map[string]MongoDbConnectInfoItem)

func InitMongoConfig() {
	prodMongoHost := os.Getenv("MONGODB_HOST")
	if prodMongoHost != "" {
		log.Println("使用环境变量中的Mongodb配置")
		prodMongoName := os.Getenv("MONGODB_USER")
		prodMongoPass := os.Getenv("MONGODBPASS")
		prodMongoHost := os.Getenv("MONGODB_HOST")
		prodMongoPort := os.Getenv("MONGODB_PORT")
		prodMongoDBNameStore := os.Getenv("MONGODB_DBNAME_LP_STORE")
		prodMongoDBNameHistory := os.Getenv("MONGODB_DBNAME_BUSINESS_HISTORY")

		url := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=%s",prodMongoName, prodMongoPass, prodMongoHost, prodMongoPort, prodMongoDBNameStore, prodMongoDBNameStore)
		item := MongoDbConnectInfoItem{Url: url, Database: prodMongoDBNameStore}

		MongoDataBaseConfigIns["main"] = item
		businessUrl := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=%s", prodMongoName, prodMongoPass, prodMongoHost, prodMongoPort, prodMongoDBNameHistory, prodMongoDBNameHistory)
		itemBusiness := MongoDbConnectInfoItem{Url: businessUrl, Database: prodMongoDBNameHistory}
		MongoDataBaseConfigIns["businessHistory"] = itemBusiness
		return
	}
	item := MongoDbConnectInfoItem{Url: "mongodb://admin:123456@127.0.0.1:27017/lp_store?authSource=lp_store", Database: "lp_store"}
	MongoDataBaseConfigIns["main"] = item
	itemBusiness := MongoDbConnectInfoItem{Url: "mongodb://admin:123456@127.0.0.1:27017/businessHistory?authSource=businessHistory", Database: "businessHistory"}
	MongoDataBaseConfigIns["businessHistory"] = itemBusiness
}
