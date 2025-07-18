package database

import (
	"admin-panel/database_config"
	sysLogger "admin-panel/logger"
	"admin-panel/types"
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/aws/smithy-go/ptr"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoSet struct {
	Client  *mongo.Client
	Session *mongo.Database
	DbName  string
}

var MongoSetLock sync.Mutex

var DbList map[string]MongoSet = make(map[string]MongoSet)

type UrlOption struct {
	Url      string
	DataBase string
}

// IsInit check if db already connected
func IsInit(dbKey string) bool {
	_, ok := DbList[dbKey]
	if !ok {
		return false
	}
	return true
}
func InitConnect(dbKey string, option *UrlOption) (err error) {
	if IsInit(dbKey) {
		log.Printf("key for db already connected [%s]", dbKey)
		return
	}
	var url string
	var dbName string
	if option == nil {
		mongoConfig, ok := database_config.MongoDataBaseConfigIns[dbKey]
		url = mongoConfig.Url
		dbName = mongoConfig.Database
		if !ok {
			sysLogger.Config.Errorf("mongodb config file does not exist ,Key[%s]", dbKey)
			os.Exit(0)
		}
	} else {
		url = option.Url
		dbName = option.DataBase
	}

	log.Println("start connecting database", url)
	ctx, _ := context.WithTimeout(context.TODO(), time.Second*2)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		log.Println("error occurred connecting", err)
		return
	}
	dbSession := client.Database(dbName)
	log.Println("select database", dbName)

	selectCtx, _ := context.WithDeadline(context.Background(), time.Now().Add(time.Second*2))
	selectCtx.Deadline()
	collections, err := dbSession.ListCollections(selectCtx, bson.M{})
	if err != nil {
		return
	}
	var _ = collections
	MongoSetLock.Lock()
	DbList[dbKey] = MongoSet{
		Client:  client,
		Session: dbSession,
		DbName:  dbName,
	}
	MongoSetLock.Unlock()
	log.Println("finish selecting database", dbName)
	return
}

func FindOne(dbKey string, collection string, filter bson.M, v interface{}) error {
	if !IsInit(dbKey) {
		return errors.New(fmt.Sprintf("database not initialized%s", dbKey))
	}
	session := DbList[dbKey].Session
	// database := DbList[dbKey].DbName
	err := session.Collection(collection).FindOne(context.Background(), filter).Decode(v)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		return err
	}
	return nil
}
func MatchOne(dbKey string, collection string, filter bson.M) (res bool, err error) {
	res = false
	err = nil
	if !IsInit(dbKey) {
		err = errors.WithMessage(err, fmt.Sprintf("database not initialized%s", dbKey))
		return
	}
	session := DbList[dbKey].Session
	docStruct := &struct {
		ID primitive.ObjectID `bson:"_id"`
	}{}
	err = session.Collection(collection).FindOne(context.Background(), filter).Decode(docStruct)
	if err != nil {
		return
	}
	if docStruct.ID.Hex() == types.MongoEmptyIdHex {
		res = true
		return
	}
	res = true
	return
}
func FindAll(dbKey string, collection string, filter bson.M) (error, *mongo.Cursor) {
	if !IsInit(dbKey) {
		return errors.New(fmt.Sprintf("database not initialized%s", dbKey)), nil
	}
	session := DbList[dbKey].Session
	cursor, err := session.Collection(collection).Find(context.Background(), filter)
	if err != nil {
		return err, nil
	}
	return nil, cursor
}
func Aggregate(dbKey string, collection string, pipeline []bson.M) (*mongo.Cursor, error) {
	if !IsInit(dbKey) {
		return nil, errors.New(fmt.Sprintf("database not initialized%s", dbKey))
	}
	session := DbList[dbKey].Session
	cursor, err := session.Collection(collection).Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	return cursor, nil
}
func FindAllOpt(dbKey string, collection string, filter bson.M, opts *options.FindOptions) (error, *mongo.Cursor) {
	if !IsInit(dbKey) {
		return errors.New(fmt.Sprintf("database not initialized%s", dbKey)), nil
	}
	session := DbList[dbKey].Session
	cursor, err := session.Collection(collection).Find(context.Background(), filter, opts)
	if err != nil {
		return err, nil
	}
	return nil, cursor
}
func Count(dbKey string, collection string, filter bson.M) (int64, error) {
	if !IsInit(dbKey) {
		return 0, errors.New(fmt.Sprintf("database not initialized%s", dbKey))
	}
	session := DbList[dbKey].Session
	count, err := session.Collection(collection).CountDocuments(context.Background(), filter)

	if err != nil {
		return 0, err
	}
	return count, nil
}
func FindOneAndUpdate(dbKey string, collection string, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	if !IsInit(dbKey) {
		return nil, errors.New(fmt.Sprintf("database not initialized%s", dbKey))
	}
	session := DbList[dbKey].Session
	result, err := session.Collection(collection).UpdateOne(context.Background(), filter, update, &options.UpdateOptions{Upsert: ptr.Bool(true)})
	if err != nil {
		return nil, err
	}
	return result, nil
}
func Update(dbKey string, collection string, filter interface{}, update interface{}) error {
	if !IsInit(dbKey) {
		return errors.New(fmt.Sprintf("database not initialized%s", dbKey))
	}
	session := DbList[dbKey].Session
	_, err := session.Collection(collection).UpdateOne(context.Background(), filter, update, &options.UpdateOptions{Upsert: ptr.Bool(false)})
	if err != nil {
		return err
	}
	return nil
}
func Insert(dbKey string, collection string, set interface{}) error {
	if !IsInit(dbKey) {
		return errors.New(fmt.Sprintf("database not initialized%s", dbKey))
	}
	session := DbList[dbKey].Session
	_, err := session.Collection(collection).InsertOne(context.Background(), set)
	if err != nil {
		return err
	}
	return nil
}
func DeleteOne(dbKey string, collection string, filter interface{}) (int64, error) {
	if !IsInit(dbKey) {
		return 0, errors.New(fmt.Sprintf("database not initialized%s", dbKey))
	}
	session := DbList[dbKey].Session
	deleteCount, err := session.Collection(collection).DeleteOne(context.TODO(), filter)
	if err != nil {
		return 0, err
	}
	return deleteCount.DeletedCount, nil
}

func GetSession(dbKey string) (*mongo.Database, error) {
	if !IsInit(dbKey) {
		return nil, errors.New(fmt.Sprintf("database not initialized%s", dbKey))
	}
	session := DbList[dbKey].Session
	return session, nil
}
