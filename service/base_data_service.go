package service

import (
	"admin-panel/logger"
	database "admin-panel/mongo_database"
	"admin-panel/types"
	"admin-panel/utils"
	"strings"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

type BaseDataLogicService struct {
}

func NewBaseDataLogicService() *BaseDataLogicService {
	return &BaseDataLogicService{}
}
func (bds *BaseDataLogicService) GetChainRowByName(chainName string) (retRow types.MongoChainListRow, err error) {
	filter := bson.M{
		"chainName": strings.ToUpper(chainName),
	}
	logger.System.Debug(filter)
	retRow = types.MongoChainListRow{}
	err = database.FindOne("main", "chainList", filter, &retRow)
	if err != nil {
		return
	}
	if retRow.ID.Hex() == types.MongoEmptyIdHex {
		err = errors.WithMessage(utils.GetNoEmptyError(err), "chain not found")
		return
	}
	return
}
func (bds *BaseDataLogicService) GetChainRowById(chainId int64) (retRow types.MongoChainListRow, err error) {
	filter := bson.M{
		"chainId": chainId,
	}
	logger.System.Debug(filter)
	retRow = types.MongoChainListRow{}
	err = database.FindOne("main", "chainList", filter, &retRow)
	if err != nil {
		return
	}
	return
}
