package service

import (
	database "admin-panel/mongo_database"
	"admin-panel/redis_database"
	"admin-panel/types"
	"context"
	"encoding/json"
	"log"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderCenterLogicService struct {
}

func (ocls *OrderCenterLogicService) AllRedis() (ret []types.CenterOrder, err error) {
	ret = make([]types.CenterOrder, 0)
	strList, err := redis_database.GetDataRedis().Smembers("KEY_BUSINESS_STATUS_INBUSINESS")
	if err != nil {
		return
	}
	for _, v := range strList {
		row := types.CenterOrder{}
		unmarshalErr := json.Unmarshal([]byte(v), &row)
		if unmarshalErr != nil {
			log.Println(unmarshalErr)
			continue
		}
		ret = append(ret, row)
	}
	return
}
func (ocls *OrderCenterLogicService) All(queryOption struct {
	Page     int64
	PageSize int64
	Status   int64
}, finder bson.M) (ret []types.CenterOrder, pageCount int64, err error) {
	var results []types.CenterOrder
	skip := queryOption.Page*queryOption.PageSize - queryOption.PageSize
	opts := options.Find().SetSort(bson.D{{"_id", -1}}).SetSkip(skip).SetLimit(queryOption.PageSize)
	pageCount = 0

	count, err := database.Count("businessHistory", "business", finder)
	if err != nil {
		return
	}
	err, cursor := database.FindAllOpt("businessHistory", "business", finder, opts)
	pageCount = count / queryOption.PageSize
	if count%queryOption.PageSize != 0 {
		pageCount++
	}
	if pageCount == 0 {
		pageCount = 1
	}
	if err != nil {
		return
	}

	if err = cursor.All(context.TODO(), &results); err != nil {
		err = errors.WithMessage(err, "cursor all error")
		return
	}

	for _, result := range results {
		log.Println(cursor.Current.String())
		cursor.Decode(&result)
	}

	ret = results
	return
}
