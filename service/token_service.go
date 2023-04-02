package service

import (
	database "admin-panel/mongo_database"
	"admin-panel/types"
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

type TokenManagerLogicService struct {
}

func NewTokenManagerLogicService() *TokenManagerLogicService {
	return &TokenManagerLogicService{}
}

func (*TokenManagerLogicService) InsertRow(interface{}) (insertId string, err error) {
	return
}
func (*TokenManagerLogicService) FindOneByFilter(filter bson.M) (ret types.DBTokenRow, err error) {
	ret = types.DBTokenRow{}
	err = database.FindOne("main", "tokens", filter, &ret)
	if err != nil && strings.Contains(err.Error(), "no documents in result") {
		err = nil
	}
	return
}
func (*TokenManagerLogicService) ListAll(data bson.M) (ret []types.DBTokenRow, err error) {
	emptyList := []types.DBTokenRow{}
	ret = emptyList
	err, cursor := database.FindAll("main", "tokens", data)
	if err != nil {
		return
	}
	var results []types.DBTokenRow
	if err = cursor.All(context.TODO(), &results); err != nil {
		return
	}
	for _, result := range results {
		cursor.Decode(&result)
	}
	ret = results
	return
}
