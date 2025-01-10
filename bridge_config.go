package adminapiservice

import (
	bridgeconfig "admin-panel/gen/bridge_config"
	database "admin-panel/mongo_database"
	redisbus "admin-panel/redis_bus"
	"admin-panel/service"
	"admin-panel/types"
	"admin-panel/utils"
	"fmt"

	"context"
	"log"

	"github.com/aws/smithy-go/ptr"
	"github.com/pkg/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// bridgeConfig service example implementation.
// The example methods log the requests and return zero values.
type bridgeConfigsrvc struct {
	logger *log.Logger
}

// NewBridgeConfig returns the bridgeConfig service implementation.
func NewBridgeConfig(logger *log.Logger) bridgeconfig.Service {
	return &bridgeConfigsrvc{logger}
}

// create bridge config
func (s *bridgeConfigsrvc) BridgeCreate(ctx context.Context, p *bridgeconfig.BridgeItem) (res *bridgeconfig.BridgeCreateResult, err error) {
	res = &bridgeconfig.BridgeCreateResult{}

	idMap, err := s.GetIdList(p)
	if err != nil {
		err = errors.WithMessage(err, "getting id list occur error:")
		return
	}
	bcls := service.NewBridgeConfigLogicService()
	commit, _id, err := bcls.CreateBridge(p, idMap) //create bridge
	if err != nil {
		return
	}
onReturn:
	if err != nil {
		commit(false)
		return
	}
	configResult, err := bcls.ConfigLp() // according to the latest list configlp
	if err != nil {
		goto onReturn
	}
	if !configResult {
		err = errors.Errorf("configLp failed")
		goto onReturn
	}
	configClient, err := bcls.ConfigClient() // config all client according to bridge info
	if err != nil {
		err = errors.WithMessage(err, "configuring client occur error")
		goto onReturn
	}
	if !configClient {
		err = errors.New("configuring client occur error")
		goto onReturn
	}
	res.Code = ptr.Int64(0)
	res.Result = ptr.Int64(0)
	log.Println("creating bridge successfully", _id)
	res.Message = ptr.String(_id)
	redisbus.GetRedisBus().PublishEvent("LP_SYSTEM_Notice", `{"type":"bridgeUpdate","payload":"{}"}`)
	s.logger.Print("bridgeConfig.bridgeCreate", _id)
	return
}
func (s *bridgeConfigsrvc) GetIdList(bridgeItem *bridgeconfig.BridgeItem) (res map[string]primitive.ObjectID, err error) {
	res = make(map[string]primitive.ObjectID, 0)
	srcChainId, err := primitive.ObjectIDFromHex(bridgeItem.SrcChainID)
	if err != nil {
		err = errors.WithMessage(err, "srcChainId incorrect:")
		return
	}
	dstChainId, err := primitive.ObjectIDFromHex(bridgeItem.DstChainID)
	if err != nil {
		err = errors.WithMessage(err, "dstChainId incorrect:")
		return
	}
	srcTokenId, err := primitive.ObjectIDFromHex(bridgeItem.SrcTokenID)
	if err != nil {
		err = errors.WithMessage(err, "srcTokenId incorrect:")
		return
	}
	dstTokenId, err := primitive.ObjectIDFromHex(bridgeItem.DstTokenID)
	if err != nil {
		err = errors.WithMessage(err, "dstTokenId incorrect:")
		return
	}
	walletId, err := primitive.ObjectIDFromHex(bridgeItem.WalletID)
	if err != nil {
		err = errors.WithMessage(err, "walletId incorrect:")
		return
	}
	srcWalletId, err := primitive.ObjectIDFromHex(bridgeItem.SrcWalletID)
	if err != nil {
		err = errors.WithMessage(err, "receiving wallet id incorrect:")
	}
	res["srcChainId"] = srcChainId
	res["dstChainId"] = dstChainId
	res["srcTokenId"] = srcTokenId
	res["dstTokenId"] = dstTokenId
	res["walletId"] = walletId
	res["srcWalletId"] = srcWalletId
	return
}

// BridgeList implements bridgeList.
func (s *bridgeConfigsrvc) BridgeList(ctx context.Context) (res *bridgeconfig.BridgeListResult, err error) {
	res = &bridgeconfig.BridgeListResult{}
	bcls := service.NewBridgeConfigLogicService()
	list, err := bcls.GetBridgeListByFilter(bson.M{})
	if err != nil {
		return
	}
	log.Println(list)
	log.Println(len(list))
	retList := make([]*bridgeconfig.ListBridgeItem, 0)
	for _, v := range list {

		var srcTokenBalance, dstTokenBalance string
		if v.SrcTokenBalance != "" {
			srcTokenBalance = v.SrcTokenBalance
		}
		if v.DstTokenBalance != "" {
			dstTokenBalance = v.DstTokenBalance
		}
		retList = append(retList, &bridgeconfig.ListBridgeItem{
			ID:                ptr.String(v.ID.Hex()),
			DstChainID:        ptr.String(v.DstChain_ID.Hex()),
			DstTokenID:        ptr.String(v.DstToken_ID.Hex()),
			SrcChainID:        ptr.String(v.SrcChain_ID.Hex()),
			SrcTokenID:        ptr.String(v.SrcToken_ID.Hex()),
			AmmName:           ptr.String(v.AmmName),
			BridgeName:        ptr.String(v.BridgeName),
			DstChainRawID:     ptr.Int64(v.DstChainId),
			DstClientURI:      ptr.String(v.DstClientUri),
			DstToken:          ptr.String(v.DstToken),
			LpReceiverAddress: ptr.String(v.LpReceiverAddress),
			MsmqName:          ptr.String(v.MsmqName),
			SrcChainRawID:     ptr.Int64(v.SrcChainId),
			SrcToken:          ptr.String(v.SrcToken),
			WalletName:        ptr.String(v.WalletName),
			WalletID:          ptr.String(v.Wallet_ID.Hex()),
			EnableHedge:       ptr.Bool(v.EnableHedge),
			SrcTokenBalance:   srcTokenBalance,
			DstTokenBalance:   dstTokenBalance,
			SrcTokenDecimals:  v.SrcTokenDecimals,
			DstTokenDecimals:  v.DstTokenDecimals,
		})
	}
	res.Code = ptr.Int64(0)
	res.Result = retList
	s.logger.Print("bridgeConfig.bridgeList")
	return
}

// BridgeDelete implements bridgeDelete.
func (s *bridgeConfigsrvc) BridgeDelete(ctx context.Context, p *bridgeconfig.DeleteBridgeFilter) (res *bridgeconfig.BridgeDeleteResult, err error) {
	res = &bridgeconfig.BridgeDeleteResult{}
	objId, err := primitive.ObjectIDFromHex(p.ID)
	if err != nil {
		err = errors.WithMessage(err, "id format incorrect, unable to convert to mongoid")
		return
	}
	var v struct {
		Id primitive.ObjectID `bson:"_id"`
	} = struct {
		Id primitive.ObjectID `bson:"_id"`
	}{}
	err = database.FindOne("main", "bridges", bson.M{"_id": objId}, &v)
	if err != nil {
		err = errors.WithMessage(err, "query database error")
		return
	}
	if v.Id.Hex() == types.MongoEmptyIdHex {
		err = errors.New("did not find corresponding bridge")
		return
	}
	delCount, err := database.DeleteOne("main", "bridges", bson.M{"_id": objId})
	if err != nil {
		return
	}
	bcls := service.NewBridgeConfigLogicService()
	configLpResult, err := bcls.ConfigLp() // reconfig lp through database record
	if err != nil || !configLpResult {
		err = errors.WithMessage(utils.GetNoEmptyError(err), "reconfiguring lp occur error")
		return
	}
	configResult, err := bcls.ConfigClient() // refresh client config according to database, as for deleting redundant client config, no need to handle for now, like the last bsc item
	if err != nil || !configResult {
		err = errors.WithMessage(utils.GetNoEmptyError(err), "configuring client occur error")
		return
	}

	res.Code = ptr.Int64(0)
	res.Result = ptr.Int64(delCount)
	res.Message = ptr.String("")
	redisbus.GetRedisBus().PublishEvent("LP_SYSTEM_Notice", `{"type":"bridgeUpdate","payload":"{}"}`)
	s.logger.Print("bridgeConfig.bridgeDelete")
	return
}
func (s *bridgeConfigsrvc) BridgeTest(ctx context.Context, p *bridgeconfig.BridgeTestPayload) (res *bridgeconfig.BridgeTestResult, err error) {
	res = &bridgeconfig.BridgeTestResult{}
	redisbus.GetRedisBus().PublishEvent("LP_SYSTEM_Notice", `{"type":"bridgeUpdate","payload":"{}"}`)
	// s.logger.Print("bridgeConfig.bridgeTest")
	bcls := service.NewBridgeConfigLogicService()
	//bcls.GetConfigClientStruct()
	//bcls.GetConfigJsonData()
	// return
	flag, err := bcls.ConfigLp()
	fmt.Println(flag)
	fmt.Println(err)
	// bcls.ConfigClient()
	//bcls.ConfigClient()

	// log.Println(flag)
	// _, err = bcls.ConfigAllClient()

	if err != nil {
		return
	}
	// log.Println(val)
	return
}
