package service

import (
	database "admin-panel/mongo_database"
	"admin-panel/types"
	"context"
	"fmt"
	"log"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AmmOrderCenterLogicService struct {
}

func (ammocls *AmmOrderCenterLogicService) All(queryOption struct {
	Page     int64
	PageSize int64
	Status   int64
}, ammName string, finder bson.M) (ret []types.AmmContext, pageCount int64, err error) {
	var results []types.AmmContext
	skip := queryOption.Page*queryOption.PageSize - queryOption.PageSize
	collectionName := fmt.Sprintf("ammContext_%s", ammName)
	count, err := database.Count("main", collectionName, finder)
	if err != nil {
		return
	}

	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}}).SetSkip(skip).SetLimit(queryOption.PageSize)
	pageCount = count / queryOption.PageSize
	if count%queryOption.PageSize != 0 {
		pageCount++
	}
	if pageCount == 0 {
		pageCount = 1
	}
	err, cursor := database.FindAllOpt("main", collectionName, bson.M{}, opts)
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
