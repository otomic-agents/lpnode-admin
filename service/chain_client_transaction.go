package service

import (
	"admin-panel/logger"
	database "admin-panel/mongo_database"
	"admin-panel/types"
	"context"
	"fmt"
	"log"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChainClientTransactionService struct {
}

func (ccts *ChainClientTransactionService) All(queryOption struct {
	Page     int64
	PageSize int64
	Status   int64
}, ammName string, finder bson.M) (res []types.View_ChainTransaction, pageCount int64, err error) {

	res = make([]types.View_ChainTransaction, 0)
	var results []types.Db_ChainTransaction
	skip := queryOption.Page*queryOption.PageSize - queryOption.PageSize
	collectionName := "chain_clients_sended_transactions"
	count, err := database.Count("main", collectionName, finder)
	if err != nil {
		return
	}
	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}}).SetSkip(int64(skip)).SetLimit(int64(queryOption.PageSize))
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
	// Extract business IDs for AMM context lookup
	businessHashList := make([]string, 0)
	for _, v := range results {
		businessHashList = append(businessHashList, v.BusinessID)
	}
	// Load AMM data for these business IDs
	ammContextList, ammLoadErr := ccts.loadAmmData(ammName, businessHashList)
	if ammLoadErr != nil {
		err = errors.WithMessage(ammLoadErr, "load amm error:")
	}
	// Create a map for quick lookup of AMM contexts by business hash
	ammContextMap := make(map[string]types.AmmContext)
	for _, ammContext := range ammContextList {
		ammContextMap[ammContext.BusinessHash] = ammContext
	}

	// Process each transaction and build the view objects
	for _, txn := range results {
		viewTxn := types.View_ChainTransaction{
			BusinessID:      txn.BusinessID,
			EventName:       txn.EventName,
			Status:          txn.Status,
			LatestGasPrice:  txn.GasPrice,
			IncreasedNumber: txn.IncreasedNumber,
			SystemChainId:   txn.SystemChainID,
			SendList:        txn.SendList,
			IncGasList:      txn.IncGasList,
			LatestSend:      txn.SendTime,
		}

		// Find corresponding AMM context and populate chain/token info
		if ammContext, exists := ammContextMap[txn.BusinessID]; exists {
			// Set source chain name
			viewTxn.SrcChainName = ammContext.SystemOrder.BaseInfo.SrcChain.Name

			// Set destination chain name
			viewTxn.DstChainName = ammContext.SystemOrder.BaseInfo.DstChain.Name

			// Set source token symbol
			viewTxn.SrcToken = ammContext.SystemOrder.BaseInfo.SrcToken.Symbol

			// Set destination token symbol
			viewTxn.DstToken = ammContext.SystemOrder.BaseInfo.DstToken.Symbol
		}

		res = append(res, viewTxn)
	}
	// spew.Dump(res)
	return
}
func (ccts *ChainClientTransactionService) loadAmmData(ammName string, businessIdList []string) (res []types.AmmContext, err error) {
	collectionName := fmt.Sprintf("ammContext_%s", ammName)
	var results []types.AmmContext
	opts := options.Find()
	filter := bson.M{}

	if len(businessIdList) > 0 {
		logger.System.Info("Querying with $in operator for business IDs:", businessIdList)
		filter = bson.M{"businessHash": bson.M{"$in": businessIdList}}
	}

	err, cursor := database.FindAllOpt("main", collectionName, filter, opts)
	if err != nil {
		logger.System.Error("Failed to execute database query", "error", err)
		return
	}

	if err = cursor.All(context.TODO(), &results); err != nil {
		err = errors.WithMessage(err, "failed to retrieve all documents from cursor")
		logger.System.Error("Cursor iteration failed", "error", err)
		return
	}

	// Note: The following loop seems unnecessary since cursor.All() already populated results
	// If you need to process each document individually, consider removing cursor.All() above
	// and use cursor.Next() instead
	for _, result := range results {
		logger.System.Debug("Processing document", "document", cursor.Current.String())
		cursor.Decode(&result)
	}

	res = results
	logger.System.Debug("Successfully loaded AMM data", "count", len(results))
	return
}
