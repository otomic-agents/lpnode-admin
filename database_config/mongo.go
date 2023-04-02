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
	prodMongoHost := os.Getenv("OBRIDGE_LPNODE_DB_MONGODB_SERVICE_HOST")
	if prodMongoHost != "" {
		log.Println("使用环境变量中的Mongodb配置")
		prodMongoPass := os.Getenv("MONGODBPASS")
		prodMongoHost := os.Getenv("OBRIDGE_LPNODE_DB_MONGODB_SERVICE_HOST")
		prodMongoPort := os.Getenv("OBRIDGE_LPNODE_DB_MONGODB_SERVICE_PORT")
		url := fmt.Sprintf("mongodb://root:%s@%s:%s/lp_store?authSource=admin", prodMongoPass, prodMongoHost, prodMongoPort)
		item := MongoDbConnectInfoItem{Url: url, Database: "lp_store"}

		MongoDataBaseConfigIns["main"] = item
		businessUrl := fmt.Sprintf("mongodb://root:%s@%s:%s/businessHistory?authSource=admin", prodMongoPass, prodMongoHost, prodMongoPort)
		itemBusiness := MongoDbConnectInfoItem{Url: businessUrl, Database: "businessHistory"}
		MongoDataBaseConfigIns["businessHistory"] = itemBusiness
		return
	}
	item := MongoDbConnectInfoItem{Url: "mongodb://root:123456@127.0.0.1:27017/lp_store?authSource=admin", Database: "lp_store"}
	MongoDataBaseConfigIns["main"] = item
	itemBusiness := MongoDbConnectInfoItem{Url: "mongodb://root:123456@127.0.0.1:27017/businessHistory?authSource=admin", Database: "businessHistory"}
	MongoDataBaseConfigIns["businessHistory"] = itemBusiness
}
