package service

import (
	database "admin-panel/mongo_database"
	"admin-panel/redis_database"
	"admin-panel/types"
	"admin-panel/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/big"
	"strings"

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
	ret, renderErr := ocls.OrderRender(ret)
	if renderErr != nil {
		err = errors.WithMessage(renderErr, "render order error:")
		return
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
	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}}).SetSkip(skip).SetLimit(queryOption.PageSize)
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
		// log.Println(cursor.Current.String())
		cursor.Decode(&result)
	}

	ret, renderErr := ocls.OrderRender(results)
	if renderErr != nil {
		err = errors.WithMessage(renderErr, "render order error:")
		return
	}
	return
}
func (ocls *OrderCenterLogicService) OrderRender(source []types.CenterOrder) (ret []types.CenterOrder, err error) {
	log.Println("OrderRender")
	var bridgeNameList []string = []string{}
	seenBridge := make(map[string]bool)
	ret = source
	var tokenList []string
	seenToken := make(map[string]bool)

	for _, v := range source {
		srcToken := strings.ToLower(v.PreBusiness.SwapAssetInformation.Quote.QuoteBase.Bridge.SrcToken)
		dstToken := strings.ToLower(v.PreBusiness.SwapAssetInformation.Quote.QuoteBase.Bridge.DstToken)
		srcChainId := v.PreBusiness.SwapAssetInformation.Quote.QuoteBase.Bridge.SrcChainId
		dstChainId := v.PreBusiness.SwapAssetInformation.Quote.QuoteBase.Bridge.DstChainId
		msmqName := fmt.Sprintf("%s/%s_%d_%d", srcToken, dstToken, srcChainId, dstChainId)
		if !seenBridge[msmqName] {
			bridgeNameList = append(bridgeNameList, msmqName)
			seenBridge[msmqName] = true
		}

		if !seenToken[srcToken] {
			tokenList = append(tokenList, srcToken)
			seenToken[srcToken] = true
		}
		if !seenToken[dstToken] {
			tokenList = append(tokenList, dstToken)
			seenToken[dstToken] = true
		}
	}
	log.Println(bridgeNameList)
	log.Println(tokenList)
	ret, renderTokenErr := ocls.orderRenderTokens(ret, tokenList)
	if renderTokenErr != nil {
		err = errors.WithMessage(renderTokenErr, "render token err:")
		return
	}

	ret, renderChainInfoErr := ocls.orderRenderChainInfo(ret)

	if renderChainInfoErr != nil {
		err = errors.WithMessage(renderChainInfoErr, "render chainInfo err:")
		return
	}

	ret, renderBridgeErr := ocls.orderRenderBridges(ret)
	if renderBridgeErr != nil {
		err = errors.WithMessage(renderBridgeErr, "render bridges err:")
		return
	}
	ret, renderViewErr := ocls.orderRenderViews(ret)
	if renderViewErr != nil {
		err = errors.WithMessage(renderViewErr, "render view err:")
		return
	}
	// pp.Print(ret)
	return
}
func (ocls *OrderCenterLogicService) orderRenderTokens(source []types.CenterOrder, tokenAddress []string) (ret []types.CenterOrder, err error) {
	log.Println(tokenAddress, "🚛🚛🚛🚛")
	ret = source
	var results []types.DBTokenRow
	err, cursor := database.FindAll("main", "tokens", bson.M{})
	if err != nil {
		return
	}
	if err = cursor.All(context.TODO(), &results); err != nil {
		return
	}
	for _, result := range results {
		cursor.Decode(&result)
	}
	log.Println(results)
	tokenMap := make(map[string]types.DBTokenRow)
	for _, result := range results {
		if result.ChainId == 501 {
			tokenHexAddress, convertErr := utils.Base58ToHexString(result.Address)
			if convertErr != nil {
				log.Println("convert error")
				break
			}
			tokenMap[tokenHexAddress] = result
		} else {
			tokenMap[result.Address] = result
		}
	}

	for i, order := range source {
		if token, ok := tokenMap[order.PreBusiness.SwapAssetInformation.Quote.QuoteBase.Bridge.SrcToken]; ok {
			log.Println("set view", "🚁🚁🚁🚁🚁", token.TokenName, token.Precision)
			order.ViewInfo.SrcTokenName = token.TokenName
			order.ViewInfo.SrcTokenPrecision = token.Precision
			source[i] = order
		} else {
			log.Println("not Found", order.PreBusiness.SwapAssetInformation.Quote.QuoteBase.Bridge.SrcToken)
		}
		if token, ok := tokenMap[order.PreBusiness.SwapAssetInformation.Quote.QuoteBase.Bridge.DstToken]; ok {
			log.Println("set view", "🚁🚁🚁🚁🚁", token.TokenName, token.Precision)
			order.ViewInfo.DstTokenName = token.TokenName
			order.ViewInfo.DstTokenPrecision = token.Precision
			source[i] = order
		}
	}
	// log.Println(tokenMap)
	for _, view := range source {
		log.Println(view.PreBusiness.SwapAssetInformation.Quote.QuoteBase.Bridge.SrcToken)
		log.Println("🎰🎰🎰🎰🎰", view.ViewInfo.SrcTokenName, view.ViewInfo.DstTokenName)
	}
	return
}
func (ocls *OrderCenterLogicService) orderRenderBridges(source []types.CenterOrder) (ret []types.CenterOrder, err error) {
	bridgeMap := make(map[string]types.DBBridgeWalletDetailsAggregateItem)
	ret = source
	mdb, err := database.GetSession("main")
	if err != nil {
		return
	}
	var results []types.DBBridgeWalletDetailsAggregateItem

	cursor, err := mdb.Collection("bridges").Aggregate(context.TODO(), bson.A{
		bson.M{
			"$lookup": bson.M{
				"from":         "wallets",
				"localField":   "src_wallet_id",
				"foreignField": "_id",
				"as":           "wallet_details",
			},
		},
		bson.M{
			"$unwind": "$wallet_details",
		},
		bson.M{
			"$project": bson.M{
				"walletDetailsName":       "$wallet_details.walletName",
				"walletDetailsPrivateKey": "$wallet_details.privateKey",
				"_id":                     1,
				"wallet_id":               1,
				"msmqName":                1,
				"walletName":              1,
			},
		},
	})
	if err != nil {
		err = errors.WithMessage(err, "Aggregate err:")
		return
	}
	if err = cursor.All(context.TODO(), &results); err != nil {
		return
	}
	for _, result := range results {
		cursor.Decode(&result)
		bridgeMap[result.MsmqName] = result
	}
	for i, v := range source {
		srcToken := v.PreBusiness.SwapAssetInformation.Quote.QuoteBase.Bridge.SrcToken
		dstToken := v.PreBusiness.SwapAssetInformation.Quote.QuoteBase.Bridge.DstToken
		srcChain := v.PreBusiness.SwapAssetInformation.Quote.QuoteBase.Bridge.SrcChainId
		dstChain := v.PreBusiness.SwapAssetInformation.Quote.QuoteBase.Bridge.DstChainId
		msmqName := fmt.Sprintf("%s/%s_%d_%d", srcToken, dstToken, srcChain, dstChain)
		if bridge, ok := bridgeMap[msmqName]; ok {
			log.Println("set view", "🚁🚁🚁🚁🚁🚁🚁🚁🚁🚁🚁🚁🚁🚁🚁🚁🚁🚁🚁🚁🚁🚁🚁🚁🚁", bridge.WalletDetailsName)
			v.ViewInfo.ReceiverWallet = bridge.WalletDetailsName
			v.ViewInfo.PaymentWallet = bridge.WalletName
			source[i] = v
		}
	}

	return
}
func (ocls *OrderCenterLogicService) orderRenderChainInfo(source []types.CenterOrder) (ret []types.CenterOrder, err error) {
	ret = source
	var results []types.DbChainListRow
	err, cursor := database.FindAll("main", "chainList", bson.M{})
	if err != nil {
		return
	}
	if err = cursor.All(context.TODO(), &results); err != nil {
		return
	}
	for _, result := range results {
		cursor.Decode(&result)
	}
	chainMap := make(map[int]types.DbChainListRow)
	for _, result := range results {
		chainMap[result.ChainId] = result
	}
	for i, order := range source {
		if chain, ok := chainMap[int(order.PreBusiness.SwapAssetInformation.Quote.QuoteBase.Bridge.SrcChainId)]; ok {
			log.Println("set view", "🚁🚁🚁🚁🚁", order.EventTransferIn.SrcChainID)
			order.ViewInfo.SrcChainName = chain.ChainName
			order.ViewInfo.SrcChainPrecision = int64(chain.Precision)
			order.ViewInfo.SrcChainRpcTx = chain.RpcTx
		}
		if chain, ok := chainMap[int(order.PreBusiness.SwapAssetInformation.Quote.QuoteBase.Bridge.DstChainId)]; ok {
			log.Println("set view", "🚁🚁🚁🚁🚁", order.EventTransferIn.SrcChainID)
			order.ViewInfo.DstChainName = chain.ChainName
			order.ViewInfo.DstChainNativeTokenName = chain.TokenName
			order.ViewInfo.DstChainPrecision = int64(chain.Precision)
			order.ViewInfo.DstChainRpcTx = chain.RpcTx
		}
		source[i] = order
	}
	return
}
func (ocls *OrderCenterLogicService) orderRenderViews(source []types.CenterOrder) (ret []types.CenterOrder, err error) {
	for i, v := range source {
		inputAmount := new(big.Float)
		inputAmount.SetString(v.PreBusiness.SwapAssetInformation.Amount)
		srcPrecision := int(v.ViewInfo.SrcTokenPrecision)
		v.ViewInfo.ReceiverAmount = new(big.Float).Quo(inputAmount, big.NewFloat(math.Pow10(srcPrecision))).String()
		log.Println(v.PreBusiness.SwapAssetInformation.Amount, srcPrecision, v.ViewInfo.ReceiverAmount)

		sendAmount := new(big.Float)
		sendAmount.SetString(v.PreBusiness.SwapAssetInformation.DstAmount)
		dstPrecision := int(v.ViewInfo.DstTokenPrecision)
		v.ViewInfo.PaymentAmount = new(big.Float).Quo(sendAmount, big.NewFloat(math.Pow10(dstPrecision))).String()
		nativePrecision := int(v.ViewInfo.DstChainPrecision)
		sendNativeAmount := new(big.Float)
		sendNativeAmount.SetString(v.PreBusiness.SwapAssetInformation.DstNativeAmount)
		v.ViewInfo.PaymentNativeAmount = new(big.Float).Quo(sendNativeAmount, big.NewFloat(math.Pow10(nativePrecision))).String()

		// time info
		v.ViewInfo.QuoteTimestamp = v.PreBusiness.SwapAssetInformation.Quote.Timestamp

		// assign
		source[i] = v
	}
	ret = source
	return
}
