package adminapiservice

import (
	chainconfig "admin-panel/gen/chain_config"
	database "admin-panel/mongo_database"
	"admin-panel/redis_database"
	"context"
	"log"

	"github.com/aws/smithy-go/ptr"
	"github.com/tidwall/sjson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// chainConfig service example implementation.
// The example methods log the requests and return zero values.
type chainConfigsrvc struct {
	logger *log.Logger
}

// NewChainConfig returns the chainConfig service implementation.
func NewChainConfig(logger *log.Logger) chainconfig.Service {
	return &chainConfigsrvc{logger}
}

func (s *chainConfigsrvc) SetChainList(ctx context.Context, p *chainconfig.SetChainListPayload) (res *chainconfig.SetChainListResult, err error) {

	res = &chainconfig.SetChainListResult{
		Code: ptr.Int64(0),
	}
	for _, item := range p.ChainList {
		filter := bson.M{
			"chainId": item.ChainID,
		}
		update := bson.M{
			"$set": map[string]interface{}{
				"chainName": item.ChainName,
				"name":      item.Name,
				"tokenName": item.TokenName,
			},
		}
		_, err := database.FindOneAndUpdate("main", "chainList", filter, update)
		if err != nil {
			log.Println(err)
			res.Code = ptr.Int64(1)
			res.Message = ptr.String(err.Error())
			return nil, err
		}
	}
	res.Message = ptr.String("ok")
	s.logger.Print("chainConfig.setChainList")
	return
}

func (s *chainConfigsrvc) DelChainList(ctx context.Context, p *chainconfig.DelChainListPayload) (res *chainconfig.DelChainListResult, err error) {
	res = &chainconfig.DelChainListResult{}
	s.logger.Print("chainConfig.delChainList")
	return
}

func (s *chainConfigsrvc) ChainList(ctx context.Context) (res *chainconfig.ChainListResult, err error) {
	res = &chainconfig.ChainListResult{}
	s.logger.Print("chainConfig.chainList")
	return
}

func (s *chainConfigsrvc) SetChainGasUsd(ctx context.Context, p *chainconfig.SetChainGasUsdPayload) (res *chainconfig.SetChainGasUsdResult, err error) {
	res = &chainconfig.SetChainGasUsdResult{}
	objectId, oErr := primitive.ObjectIDFromHex(p.ID)
	if oErr != nil {
		err = oErr
		return
	}
	find := bson.M{
		"_id":     objectId,
		"chainId": p.ChainID,
	}
	update := bson.M{
		"$set": map[string]int64{
			"tokenUsd": p.Usd,
		},
	}
	dbError := database.Update("main", "chainList", find, update)
	jsonStr := "{}"
	jsonStr, _ = sjson.Set(jsonStr, "action", "reload_chain_list")
	jsonStr, err = sjson.SetRaw(jsonStr, "payload", "{}")
	redis_database.GetDataRedis().Publish("system_event_bus", jsonStr)
	log.Println(dbError)
	s.logger.Print("chainConfig.setChainGasUsd")
	return
}

// SetChainClientConfig implements setChainClientConfig.
func (s *chainConfigsrvc) SetChainClientConfig(ctx context.Context, p *chainconfig.SetChainClientConfigPayload) (res *chainconfig.SetChainClientConfigResult, err error) {
	res = &chainconfig.SetChainClientConfigResult{}
	s.logger.Print("chainConfig.setChainClientConfig")
	return
}
