package service

import (
	database "admin-panel/mongo_database"
	"admin-panel/types"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RelayAccountLogicService struct {
}

func NewRelayAccountLogicService() *RelayAccountLogicService {
	return &RelayAccountLogicService{}
}

func (*RelayAccountLogicService) GetRelayApiKey() (apiKey string, err error) {
	val := struct {
		ID          primitive.ObjectID `bson:"_id"`
		RelayApiKey string             `bson:"relayApiKey"`
	}{}
	err = database.FindOne("main", "relayAccounts", bson.M{}, &val)
	if err != nil {
		err = errors.WithMessage(err, "查询数据库发生了错误")
		return
	}
	if val.ID.Hex() == types.MongoEmptyIdHex {
		err = errors.New("没有找到任何一个relay账号")
		return
	}
	apiKey = val.RelayApiKey
	return
}
