package adminapiservice

import (
	bridgeconfig "admin-panel/gen/bridge_config"
	database "admin-panel/mongo_database"
	redisbus "admin-panel/redis_bus"
	"admin-panel/service"
	"admin-panel/types"
	"admin-panel/utils"

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

// 用于创建跨链配置
func (s *bridgeConfigsrvc) BridgeCreate(ctx context.Context, p *bridgeconfig.BridgeItem) (res *bridgeconfig.BridgeCreateResult, err error) {
	res = &bridgeconfig.BridgeCreateResult{}

	idMap, err := s.GetIdList(p)
	if err != nil {
		err = errors.WithMessage(err, "获取Id列表发生了错误:")
		return
	}
	bcls := service.NewBridgeConfigLogicService()
	commit, _id, err := bcls.CreateBridge(p, idMap) //创建bridge
	if err != nil {
		return
	}
onReturn:
	if err != nil {
		commit(false)
		return
	}
	configResult, err := bcls.ConfigLp() // 根据最新的List ConfigLp
	if err != nil {
		goto onReturn
	}
	if !configResult {
		err = errors.Errorf("configLp失败")
		goto onReturn
	}
	configClient, err := bcls.ConfigClient() // 根据bridge信息 config all client
	if err != nil {
		err = errors.WithMessage(err, "配置Client发生了错误")
		goto onReturn
	}
	if !configClient {
		err = errors.New("配置client发生了错误")
		goto onReturn
	}
	res.Code = ptr.Int64(0)
	res.Result = ptr.Int64(0)
	log.Println("创建bridge成功🧂🧂🧂🧂🧂🧂", _id)
	res.Message = ptr.String(_id)
	redisbus.GetRedisBus().PublishEvent("LP_SYSTEM_Notice", `{"type":"bridgeUpdate","payload":"{}"}`)
	s.logger.Print("bridgeConfig.bridgeCreate", _id)
	return
}
func (s *bridgeConfigsrvc) GetIdList(bridgeItem *bridgeconfig.BridgeItem) (res map[string]primitive.ObjectID, err error) {
	res = make(map[string]primitive.ObjectID, 0)
	srcChainId, err := primitive.ObjectIDFromHex(bridgeItem.SrcChainID)
	if err != nil {
		err = errors.WithMessage(err, "srcChainId不正确:")
		return
	}
	dstChainId, err := primitive.ObjectIDFromHex(bridgeItem.DstChainID)
	if err != nil {
		err = errors.WithMessage(err, "dstChainId不正确:")
		return
	}
	srcTokenId, err := primitive.ObjectIDFromHex(bridgeItem.SrcTokenID)
	if err != nil {
		err = errors.WithMessage(err, "srcTokenId不正确:")
		return
	}
	dstTokenId, err := primitive.ObjectIDFromHex(bridgeItem.DstTokenID)
	if err != nil {
		err = errors.WithMessage(err, "dstTokenId不正确:")
		return
	}
	walletId, err := primitive.ObjectIDFromHex(bridgeItem.WalletID)
	if err != nil {
		err = errors.WithMessage(err, "walletId不正确:")
	}
	srcWalletId, err := primitive.ObjectIDFromHex(bridgeItem.SrcWalletID)
	if err != nil {
		err = errors.WithMessage(err, "收款钱包id不正确:")
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
		err = errors.WithMessage(err, "id格式不正确无法转为Mongoid")
		return
	}
	var v struct {
		Id primitive.ObjectID `bson:"_id"`
	} = struct {
		Id primitive.ObjectID `bson:"_id"`
	}{}
	err = database.FindOne("main", "bridges", bson.M{"_id": objId}, &v)
	if err != nil {
		err = errors.WithMessage(err, "查询数据库错误")
		return
	}
	if v.Id.Hex() == types.MongoEmptyIdHex {
		err = errors.New("没有找到对应的Bridge")
		return
	}
	delCount, err := database.DeleteOne("main", "bridges", bson.M{"_id": objId})
	if err != nil {
		return
	}
	bcls := service.NewBridgeConfigLogicService()
	configLpResult, err := bcls.ConfigLp() // 通过数据库记录重新config Lp
	if err != nil || !configLpResult {
		err = errors.WithMessage(utils.GetNoEmptyError(err), "重新配置Lp发生了错误")
		return
	}
	//configResult, err := bcls.ConfigAllClient() // 根据 amm install 记录重新config Client
	configResult, err := bcls.ConfigClient() // 根据数据库刷新client 配置，至于删掉多出来的Client配置，暂时无需处理，如最后一项bsc 设置
	if err != nil || !configResult {
		err = errors.WithMessage(utils.GetNoEmptyError(err), "配置Client 发生了错误")
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
	// flag, err := bcls.ConfigLp()
	bcls.ConfigClient()
	//bcls.ConfigClient()

	// log.Println(flag)
	// _, err = bcls.ConfigAllClient()
	log.Print(err)
	if err != nil {
		return
	}
	// log.Println(val)
	return
}
