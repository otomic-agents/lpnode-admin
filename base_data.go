package adminapiservice

import (
	basedata "admin-panel/gen/base_data"
	database "admin-panel/mongo_database"
	"admin-panel/types"
	"context"
	"log"
	"os"

	"github.com/aws/smithy-go/ptr"
	"go.mongodb.org/mongo-driver/bson"
)

// baseData service example implementation.
// The example methods log the requests and return zero values.
type baseDatasrvc struct {
	logger *log.Logger
}

// NewBaseData returns the baseData service implementation.
func NewBaseData(logger *log.Logger) basedata.Service {
	return &baseDatasrvc{logger}
}

// ChainDataList implements chainDataList.
func (s *baseDatasrvc) ChainDataList(ctx context.Context) (res *basedata.ChainDataListResult, err error) {
	res = &basedata.ChainDataListResult{Result: make([]*basedata.ChainDataItem, 0)}
	var results []types.ChainInfoStoreItem
	err, cursor := database.FindAll("main", "chainList", bson.M{})
	if err != nil {
		return
	}
	if err = cursor.All(context.TODO(), &results); err != nil {
		return
	}
	for _, result := range results {
		cursor.Decode(&result)
		res.Result = append(res.Result, &basedata.ChainDataItem{
			ID:        ptr.String(result.Id.Hex()),
			ChainID:   ptr.Int64(result.ChainId),
			Name:      ptr.String(result.Name),
			ChainName: ptr.String(result.ChainName),
			TokenName: ptr.String(result.TokenName),
		})
	}
	res.Code = ptr.Int64(0)
	res.Message = ptr.String("")

	s.logger.Print("baseData.chainDataList")
	return
}
func (s *baseDatasrvc) RunTimeEnv(ctx context.Context) (res *basedata.RunTimeEnvResult, err error) {
	res = &basedata.RunTimeEnvResult{Result: ptr.String(""), Code: ptr.Int64(0), Message: ptr.String("")}
	env := os.Getenv("DEPLOY_ENV")
	if env != "" {
		res.Result = ptr.String(env)
	} else {
		res.Result = ptr.String("dev")
	}
	return
}
