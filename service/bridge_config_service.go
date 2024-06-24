package service

import (
	bridgeconfig "admin-panel/gen/bridge_config"
	globalval "admin-panel/global_var"
	"admin-panel/logger"
	database "admin-panel/mongo_database"
	"admin-panel/types"
	"admin-panel/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BridgeConfigLogicService struct {
}

func NewBridgeConfigLogicService() *BridgeConfigLogicService {
	return &BridgeConfigLogicService{}
}

func (bcls *BridgeConfigLogicService) GetChainRowById(_id primitive.ObjectID) (res types.MongoChainListRow, err error) {
	res = types.MongoChainListRow{}
	filter := bson.M{
		"_id": _id,
	}
	err = database.FindOne("main", "chainList", filter, &res)
	return
}
func (bcls *BridgeConfigLogicService) GetAmmInstallRowByName(name string) (res types.InstallRow, err error) {
	res = types.InstallRow{}
	filter := bson.M{
		"installType": "amm",
		"name":        name,
	}
	err = database.FindOne("main", "install", filter, &res)
	return
}

func (bcls *BridgeConfigLogicService) GetTokenRowById(_id primitive.ObjectID) (res types.DBTokenRow, err error) {
	res = types.DBTokenRow{}
	filter := bson.M{
		"_id": _id,
	}
	err = database.FindOne("main", "tokens", filter, &res)
	return
}
func (bcls *BridgeConfigLogicService) GetWalletRowById(_id primitive.ObjectID) (res types.DBWalletRow, err error) {
	res = types.DBWalletRow{}
	filter := bson.M{
		"_id": _id,
	}
	err = database.FindOne("main", "wallets", filter, &res)
	return
}
func (bcls *BridgeConfigLogicService) CreateBridge(p *bridgeconfig.BridgeItem, idMap map[string]primitive.ObjectID) (commit func(commit bool) error, id string, err error) {
	id = ""
	commit = func(commit bool) error {
		if !commit {
			tobeDelId, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				return err
			}
			log.Println("delete:", tobeDelId)
			database.DeleteOne("main", "bridges", bson.M{
				"_id": tobeDelId,
			})
		}
		return nil
	}
	ammName := p.AmmName
	ammInfo, err := bcls.GetAmmInstallRowByName(ammName)
	if err != nil || ammInfo.ID.Hex() == types.MongoEmptyIdHex {
		err = utils.GetNoEmptyError(err)
		err = errors.WithMessage(err, "get amm install info error")
		return
	}
	// log.Println(idMap)
	srcChainInfo, err := bcls.GetChainRowById(idMap["srcChainId"])
	if err != nil || srcChainInfo.ID.Hex() == types.MongoEmptyIdHex {
		err = utils.GetNoEmptyError(err)
		err = errors.WithMessage(err, "get src chain info error")
		return
	}
	dstChainInfo, err := bcls.GetChainRowById(idMap["dstChainId"])
	if err != nil || dstChainInfo.ID.Hex() == types.MongoEmptyIdHex {
		err = utils.GetNoEmptyError(err)
		err = errors.WithMessage(err, "get dst chain info error")
		return
	}
	srcTokenInfo, err := bcls.GetTokenRowById(idMap["srcTokenId"])
	if err != nil || srcTokenInfo.ID.Hex() == types.MongoEmptyIdHex {
		err = utils.GetNoEmptyError(err)
		err = errors.WithMessage(err, "srcTokenInfo get Error")
		return
	}
	if srcTokenInfo.ChainId != srcChainInfo.ChainId {
		err = errors.WithMessage(utils.GetNoEmptyError(err), fmt.Sprintf("source chain token and id mismatch, chainId:%d, token:%s", srcChainInfo.ChainId, srcTokenInfo.Address))
		return
	}
	dstTokenInfo, err := bcls.GetTokenRowById(idMap["dstTokenId"])
	if err != nil || dstTokenInfo.ID.Hex() == types.MongoEmptyIdHex {
		err = utils.GetNoEmptyError(err)
		err = errors.WithMessage(err, "dstTokenInfo get Error")
		return
	}
	if dstTokenInfo.ChainId != dstChainInfo.ChainId {
		err = errors.WithMessage(utils.GetNoEmptyError(err), fmt.Sprintf("target token not belong to target chain, chainId:%d, token:%s", dstChainInfo.ChainId, dstTokenInfo.Address))
		return
	}
	walletInfo, err := bcls.GetWalletRowById(idMap["walletId"])
	if err != nil || walletInfo.ID.Hex() == types.MongoEmptyIdHex {
		err = utils.GetNoEmptyError(err)
		err = errors.WithMessage(err, "walletInfo get Error")
		return
	}
	if walletInfo.ChainId != dstTokenInfo.ChainId {
		err = errors.WithMessage(utils.GetNoEmptyError(err), "wallet not belong to target chain, please check config")
		return
	}
	srcWalletInfo, err := bcls.GetWalletRowById(idMap["srcWalletId"])
	if err != nil || srcWalletInfo.ID.Hex() == types.MongoEmptyIdHex {
		err = utils.GetNoEmptyError(err)
		err = errors.WithMessage(err, "src walletInfo get Error")
		return
	}
	if srcWalletInfo.ChainId != srcTokenInfo.ChainId {
		err = errors.WithMessage(utils.GetNoEmptyError(err), "wallet not belong to source chain, please check config")
		return
	}
	bridgeExist := bcls.HasBridge(srcTokenInfo.Address, dstTokenInfo.Address, srcChainInfo.ChainId, dstChainInfo.ChainId, ammInfo.Name)
	if bridgeExist {
		err = utils.GetNoEmptyError(err)
		err = errors.WithMessage(err, "bridge already exist")
		return
	}
	srcClientUrl, err := bcls.GetClientUrl(srcChainInfo.ChainId)
	if err != nil {
		err = errors.WithMessage(err, "get src chain url error occur:")
		return
	}
	if srcClientUrl == "" {
		err = errors.New("cannot get correct SrcClientUrl, cannot create bridge")
		return
	}
	clientUrl, err := bcls.GetClientUrl(dstChainInfo.ChainId)
	if err != nil {
		err = errors.WithMessage(err, "get dst chain Url error occur:")
		return
	}
	if clientUrl == "" {
		err = errors.New("cannot get correct DstClientUrl, cannot create bridge")
		return
	}
	srcToken, err := utils.GetHexAddress(srcTokenInfo.Address, srcChainInfo.ChainType)
	if err != nil {
		err = errors.WithMessage(err, "srcToken process error")
		return
	}
	dstToken, err := utils.GetHexAddress(dstTokenInfo.Address, dstChainInfo.ChainType)
	if err != nil {
		err = errors.WithMessage(err, "dstToken process error")
		return
	}
	mongoUpsert, err := database.FindOneAndUpdate("main", "bridges", bson.M{
		"srcChain_id": srcChainInfo.ID,
		"dstChain_id": dstChainInfo.ID,
		"srcToken_id": srcTokenInfo.ID,
		"dstToken_id": dstTokenInfo.ID,
	}, bson.M{
		"$set": bson.M{
			"enableLimiter":     p.EnableLimiter,
			"enableHedge":       p.EnableHedge,
			"bridgeName":        p.BridgeName,
			"srcChainId":        srcChainInfo.ChainId,
			"dstChainId":        dstChainInfo.ChainId,
			"srcToken":          srcToken,
			"dstToken":          dstToken,
			"walletName":        walletInfo.WalletName,
			"wallet_id":         walletInfo.ID,
			"src_wallet_id":     srcWalletInfo.ID,
			"lpReceiverAddress": srcWalletInfo.Address,
			"msmqName":          bcls.GetMsmqName(srcToken, dstToken, srcChainInfo.ChainId, dstChainInfo.ChainId),
			"srcClientUri":      srcClientUrl,
			"dstClientUri":      clientUrl,
			"createdAt":         time.Now().UnixNano() / 1e6,
			"ammName":           ammInfo.Name,
		},
	})
	if err != nil {
		return
	}
	if mongoUpsert.UpsertedID == nil {
		err = errors.New("did not successfully create bridge")
		return
	}
	log.Println(mongoUpsert.UpsertedID, "🚵‍♀️🚵‍♀️🚵‍♀️🚵‍♀️🚵‍♀️🚵‍♀️🚵‍♀️🚵‍♀️")
	id = mongoUpsert.UpsertedID.(primitive.ObjectID).Hex()
	return
}
func (bcls *BridgeConfigLogicService) GetMsmqName(token0 string, token1 string, chain0 int64, chain1 int64) string {
	return fmt.Sprintf("%s/%s_%d_%d", token0, token1, chain0, chain1)
}
func (bcls *BridgeConfigLogicService) HasBridge(srcToken string, dstToken string, srcChainId int64, dstChainId int64, ammName string) bool {
	filter := bson.M{
		"srcChainId": srcChainId,
		"dstChainId": dstChainId,
		"srcToken":   srcToken,
		"dstToken":   dstToken,
		"ammName":    ammName,
	}
	log.Println(filter)
	hasBridge, err := database.MatchOne("main", "bridges", filter)
	log.Println(hasBridge, err)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		return true
	}
	return hasBridge
}
func (bcls *BridgeConfigLogicService) GetBridgeListByFilter(filter bson.M) (ret []types.DBBridgeRow, err error) {
	emptyList := []types.DBBridgeRow{}
	ret = emptyList
	err, cursor := database.FindAll("main", "bridges", filter)
	if err != nil {
		return
	}
	var results []types.DBBridgeRow
	if err = cursor.All(context.TODO(), &results); err != nil {
		return
	}
	for _, result := range results {
		cursor.Decode(&result)
	}
	ret = results
	return
}

func (bcls *BridgeConfigLogicService) GetConfigLpStruct() (res []types.BridgeConfigLpConfigItem, err error) {
	res = make([]types.BridgeConfigLpConfigItem, 0)

	err, cursor := database.FindAll("main", "bridges", bson.M{})
	if err != nil {
		errors.WithMessage(err, "query error")
		return
	}
	var results []types.DBBridgeRow
	if err = cursor.All(context.TODO(), &results); err != nil {
		errors.WithMessage(err, "cursor handle error")
		return
	}
	for _, result := range results {
		cursor.Decode(&result)
		res = append(res, types.BridgeConfigLpConfigItem{
			Bridge: types.BridgeConfigLpConfigItemBridge{
				SrcChainId: result.SrcChainId,
				DstChainId: result.DstChainId,
				SrcToken:   result.SrcToken,
				DstToken:   result.DstToken,
			},
			Wallet: struct {
				Name string "json:\"name\""
			}{Name: result.WalletName},
			LpReceiverAddress: result.LpReceiverAddress,
			MsmqName:          result.MsmqName,
			DstClientUri:      result.DstClientUri,
			SrcClientUri:      result.SrcClientUri,
			EnableLimiter:     result.EnableLimiter, // whether enable permission limit
		})
	}
	return
}
func (bcls *BridgeConfigLogicService) GetClientUrl(chainId int64) (ret string, err error) {
	result := struct {
		ID          primitive.ObjectID `bson:"_id"`
		ServiceName string             `bson:"serviceName"`
		ChainType   string             `bson:"chainType"`
		ChainId     int64              `bson:"chainId"`
	}{}
	err = database.FindOne("main", "install", bson.M{"chainId": chainId}, &result)
	if err != nil {
		logger.System.Warn("cannot find corresponding install record", chainId, "service")
		return
	}

	// http://chain-client-evm-avax-server-9000:9100/evm-client-9000
	ret = fmt.Sprintf("http://%s:9100/%s-client-%d", result.ServiceName, result.ChainType, result.ChainId)
	return
}
func (bcls *BridgeConfigLogicService) GetClientSetWalletUrl(chainId int64) (ret string, err error) {
	result := struct {
		ID          primitive.ObjectID `bson:"_id"`
		ServiceName string             `bson:"serviceName"`
		ChainType   string             `bson:"chainType"`
		ChainId     int64              `bson:"chainId"`
		Status      int64              `bson:"status"`
	}{}
	err = database.FindOne("main", "install", bson.M{
		"chainId": chainId,
		"status": bson.M{
			"$gte": 1,
		}},
		&result)
	if err != nil {
		return
	}
	if result.ID.Hex() == types.MongoEmptyIdHex {
		err = errors.WithMessage(utils.GetNoEmptyError(err), "cannot find install record, cannot return url")
		return
	}
	// https://api-lpnode-3.edge-dev.xyz/evm-client-9000/lpnode_admin_panel/set_wallet
	ret = fmt.Sprintf("http://%s:9100/%s-client-%d/lpnode_admin_panel/set_wallet", result.ServiceName, result.ChainType, result.ChainId)
	return
}

func (bcls *BridgeConfigLogicService) GetConfigJsonData() (res string, err error) {

	mdb, err := database.GetSession("main")
	if err != nil {
		log.Println("get database instance error occur")
		return
	}
	var results []types.DBBridgeRow

	cursor, err := mdb.Collection("bridges").Find(context.TODO(), bson.M{})
	if err = cursor.All(context.TODO(), &results); err != nil {
		err = errors.WithMessage(err, "cursor handle error")
		return
	}
	baseJson := "{}"
	for _, result := range results {
		cursor.Decode(&result)
		srcChainInfo, getChainInfoErr := bcls.GetChainRowById(result.SrcChain_ID)
		if getChainInfoErr != nil {
			err = errors.WithMessage(getChainInfoErr, "getChainInfoErr error:")
			return
		}
		srcChainInfoStr, _ := json.Marshal(srcChainInfo)
		dstChainInfo, getDstChainInfoErr := bcls.GetChainRowById(result.DstChain_ID)
		if getDstChainInfoErr != nil {
			err = errors.WithMessage(getDstChainInfoErr, "getDstChainInfoErr error:")
			return
		}
		dstChainInfoStr, _ := json.Marshal(dstChainInfo)
		srcWalletInfo, getSrcWalletErr := bcls.GetWalletRowById(result.Src_Wallet_Id)
		if getSrcWalletErr != nil {
			err = errors.WithMessage(getSrcWalletErr, " getSrcWalletErr error:")
			return
		}
		dstWalletInfo, getDstWalletErr := bcls.GetWalletRowById(result.Wallet_ID)
		if getDstWalletErr != nil {
			err = errors.WithMessage(getDstWalletErr, " getDstWalletErr error:")
			return
		}
		srcTokenInfo, getSrcTokenErr := bcls.GetTokenRowById(result.SrcToken_ID)
		if getSrcTokenErr != nil {
			err = errors.WithMessage(getSrcTokenErr, " getSrcTokenErr error:")
			return
		}
		dstTokenInfo, getDstTokenErr := bcls.GetTokenRowById(result.DstToken_ID)
		if getDstTokenErr != nil {
			err = errors.WithMessage(getDstTokenErr, " getDstTokenErr error:")
			return
		}
		srcTokenBase := "{}"
		srcTokenBase, _ = sjson.Set(srcTokenBase, "address", srcTokenInfo.Address)
		srcTokenBase, _ = sjson.Set(srcTokenBase, "tokenId", srcTokenInfo.TokenId)

		nativeTokenBase := "{}"
		nativeTokenBase, _ = sjson.Set(nativeTokenBase, "address", "0x0000000000000000000000000000000000000000")
		nativeTokenBase, _ = sjson.Set(nativeTokenBase, "tokenId", srcTokenInfo.TokenId)

		dstTokenBase := "{}"
		dstTokenBase, _ = sjson.Set(dstTokenBase, "address", dstTokenInfo.Address)
		dstTokenBase, _ = sjson.Set(dstTokenBase, "tokenId", dstTokenInfo.TokenId)
		// chain Info set
		baseJson, _ = sjson.SetRaw(baseJson, fmt.Sprintf("%d.chainInfo", result.SrcChainId), string(srcChainInfoStr))
		baseJson, _ = sjson.SetRaw(baseJson, fmt.Sprintf("%d.chainInfo", result.DstChainId), string(dstChainInfoStr))

		//## wallet set

		baseJson, _ = sjson.Set(baseJson, fmt.Sprintf("%d.walletInfo.%s.walletName", result.SrcChainId, srcWalletInfo.ID.Hex()), srcWalletInfo.WalletName)
		baseJson, _ = sjson.Set(baseJson, fmt.Sprintf("%d.walletInfo.%s.accountId", result.SrcChainId, srcWalletInfo.ID.Hex()), srcWalletInfo.AccountId)
		baseJson, _ = sjson.Set(baseJson, fmt.Sprintf("%d.walletInfo.%s.privateKey", result.SrcChainId, srcWalletInfo.ID.Hex()), srcWalletInfo.PrivateKey)
		baseJson, _ = sjson.Set(baseJson, fmt.Sprintf("%d.walletInfo.%s.walletType", result.SrcChainId, srcWalletInfo.ID.Hex()), srcWalletInfo.WalletType)
		baseJson, _ = sjson.Set(baseJson, fmt.Sprintf("%d.walletInfo.%s.storeId", result.SrcChainId, srcWalletInfo.ID.Hex()), srcWalletInfo.StoreId)
		baseJson, _ = sjson.Set(baseJson, fmt.Sprintf("%d.walletInfo.%s.address", result.SrcChainId, srcWalletInfo.ID.Hex()), srcWalletInfo.Address)
		baseJson, _ = sjson.Set(baseJson, fmt.Sprintf("%d.walletInfo.%s.vaultHostType", result.SrcChainId, srcWalletInfo.ID.Hex()), srcWalletInfo.VaultHostType)
		baseJson, _ = sjson.Set(baseJson, fmt.Sprintf("%d.walletInfo.%s.vaultName", result.SrcChainId, srcWalletInfo.ID.Hex()), srcWalletInfo.VaultName)
		baseJson, _ = sjson.Set(baseJson, fmt.Sprintf("%d.walletInfo.%s.vaultSecertType", result.SrcChainId, srcWalletInfo.ID.Hex()), srcWalletInfo.VaultSecertType)

		baseJson, _ = sjson.Set(baseJson, fmt.Sprintf("%d.walletInfo.%s.walletName", result.DstChainId, dstWalletInfo.ID.Hex()), dstWalletInfo.WalletName)
		baseJson, _ = sjson.Set(baseJson, fmt.Sprintf("%d.walletInfo.%s.accountId", result.DstChainId, dstWalletInfo.ID.Hex()), dstWalletInfo.AccountId)
		baseJson, _ = sjson.Set(baseJson, fmt.Sprintf("%d.walletInfo.%s.privateKey", result.DstChainId, dstWalletInfo.ID.Hex()), dstWalletInfo.PrivateKey)
		baseJson, _ = sjson.Set(baseJson, fmt.Sprintf("%d.walletInfo.%s.walletType", result.DstChainId, dstWalletInfo.ID.Hex()), dstWalletInfo.WalletType)
		baseJson, _ = sjson.Set(baseJson, fmt.Sprintf("%d.walletInfo.%s.storeId", result.DstChainId, dstWalletInfo.ID.Hex()), dstWalletInfo.StoreId)
		baseJson, _ = sjson.Set(baseJson, fmt.Sprintf("%d.walletInfo.%s.address", result.DstChainId, dstWalletInfo.ID.Hex()), dstWalletInfo.Address)
		baseJson, _ = sjson.Set(baseJson, fmt.Sprintf("%d.walletInfo.%s.vaultHostType", result.DstChainId, dstWalletInfo.ID.Hex()), dstWalletInfo.VaultHostType)
		baseJson, _ = sjson.Set(baseJson, fmt.Sprintf("%d.walletInfo.%s.vaultName", result.DstChainId, dstWalletInfo.ID.Hex()), dstWalletInfo.VaultName)
		baseJson, _ = sjson.Set(baseJson, fmt.Sprintf("%d.walletInfo.%s.vaultSecertType", result.DstChainId, dstWalletInfo.ID.Hex()), dstWalletInfo.VaultSecertType)

		// # token set
		baseJson, _ = sjson.SetRaw(baseJson, fmt.Sprintf("%d.walletInfo.%s.tokenInfo.%s", result.SrcChainId, srcWalletInfo.ID.Hex(), result.SrcToken_ID.Hex()), srcTokenBase)
		baseJson, _ = sjson.SetRaw(baseJson, fmt.Sprintf("%d.walletInfo.%s.tokenInfo.%s", result.DstChainId, dstWalletInfo.ID.Hex(), result.DstToken_ID.Hex()), dstTokenBase)
		baseJson, _ = sjson.SetRaw(baseJson, fmt.Sprintf("%d.walletInfo.%s.tokenInfo.%s", result.SrcChainId, srcWalletInfo.ID.Hex(), "0x0000000000000000000000000000000000000000"), nativeTokenBase)
		baseJson, _ = sjson.SetRaw(baseJson, fmt.Sprintf("%d.walletInfo.%s.tokenInfo.%s", result.DstChainId, dstWalletInfo.ID.Hex(), "0x0000000000000000000000000000000000000000"), nativeTokenBase)
	}
	res = baseJson
	logger.System.Debug("got configJson", "\r\n", gjson.Get(baseJson, "@pretty").String())
	return
}
func (bcls *BridgeConfigLogicService) GetUniqDstToken(dstChainId int64, walletName string) (res []types.TDBBridgeUniqDstToken, err error) {
	err, cursor := database.FindAll("main", "bridges", bson.M{"dstChainId": dstChainId, "walletName": walletName})
	if err != nil {
		err = errors.WithMessage(err, "query unique dstToken error occur")
		return
	}
	var results []types.DBBridgeUniqDstToken // DstChain group
	res = make([]types.TDBBridgeUniqDstToken, 0)

	if err = cursor.All(context.TODO(), &results); err != nil {
		err = errors.WithMessage(err, "cursor handle error")
		return
	}
	for _, result := range results {
		cursor.Decode(&result)
		tokenId := result.DstTokenId

		tokenRow := types.DBBridgeUniqDstTokenInfo{}
		findTokenErr := database.FindOne("main", "tokens", bson.M{"_id": tokenId}, &tokenRow)
		if findTokenErr != nil {
			return
		}
		if tokenRow.Id.Hex() == types.MongoEmptyIdHex {
			err = errors.WithMessage(utils.GetNoEmptyError(err), "query token actual value error occur")
			return
		}
		res = append(res, types.TDBBridgeUniqDstToken{
			DstTokenId: result.DstTokenId.Hex(),
			DstToken:   result.DstToken,
			WalletName: result.WalletName,
			Info: struct {
				TokenId      string
				Id           primitive.ObjectID "bson:\"_id\""
				Address      string
				ReceiptId    string
				AddressLower string
			}{
				TokenId:      tokenRow.TokenId,
				Id:           tokenRow.Id,
				Address:      tokenRow.Address,
				ReceiptId:    tokenRow.Address,
				AddressLower: tokenRow.AddressLower,
			},
		})
	}

	return
}

func (bcls *BridgeConfigLogicService) ConfigLp() (configResult bool, err error) {
	lprs := NewLpRegisterLogicService()
	als := NewAuthenticationLimiterService()
	relayApiKey, err := lprs.GetRelayApiKey()
	if err != nil {
		err = errors.WithMessage(err, "get relayApiKey error, lp may not register account yet")
		return
	}
	lpName, err := lprs.GetLpName()
	if err != nil {
		err = errors.WithMessage(err, "get lpname error, lp may not register account yet")
		return
	}
	configResult = false
	lpItemList, err := bcls.GetConfigLpStruct()
	if err != nil {
		return
	}
	limiterConf, err := als.Get()
	if err != nil {
		return
	}
	jsonStr := `{"data":[]}`

	for i, v := range lpItemList {
		fmt.Println(v)
		jsonStr, _ = sjson.Set(jsonStr, fmt.Sprintf("data.%d.bridge.src_chain_id", i), v.Bridge.SrcChainId)
		jsonStr, _ = sjson.Set(jsonStr, fmt.Sprintf("data.%d.bridge.dst_chain_id", i), v.Bridge.DstChainId)
		jsonStr, _ = sjson.Set(jsonStr, fmt.Sprintf("data.%d.bridge.src_token", i), v.Bridge.SrcToken)
		jsonStr, _ = sjson.Set(jsonStr, fmt.Sprintf("data.%d.bridge.dst_token", i), v.Bridge.DstToken)
		jsonStr, _ = sjson.Set(jsonStr, fmt.Sprintf("data.%d.wallet.name", i), v.Wallet.Name)
		jsonStr, _ = sjson.Set(jsonStr, fmt.Sprintf("data.%d.lp_receiver_address", i), v.LpReceiverAddress)
		jsonStr, _ = sjson.Set(jsonStr, fmt.Sprintf("data.%d.msmq_name", i), v.MsmqName)
		jsonStr, _ = sjson.Set(jsonStr, fmt.Sprintf("data.%d.src_client_uri", i), v.SrcClientUri)
		jsonStr, _ = sjson.Set(jsonStr, fmt.Sprintf("data.%d.dst_client_uri", i), v.DstClientUri)
		jsonStr, _ = sjson.Set(jsonStr, fmt.Sprintf("data.%d.relay_api_key", i), relayApiKey)
		jsonStr, _ = sjson.Set(jsonStr, fmt.Sprintf("data.%d.lp_id", i), lpName)
		if limiterConf.Data == "" {
			limiterConf.Data = "{}"
		}
		if v.EnableLimiter { // if record enable limiter, then add in configlp
			jsonStr, _ = sjson.SetRaw(jsonStr, fmt.Sprintf("data.%d.authentication_limiter", i), limiterConf.Data)
			jsonStr, _ = sjson.Set(jsonStr, fmt.Sprintf("data.%d.authentication_limiter.limiter_state", i), "on")
		} else {
			jsonStr, _ = sjson.SetRaw(jsonStr, fmt.Sprintf("data.%d.authentication_limiter", i), limiterConf.Data)
			jsonStr, _ = sjson.Set(jsonStr, fmt.Sprintf("data.%d.authentication_limiter.limiter_state", i), "off")
		}

	}
	// toSendArr := gjson.Get(jsonStr, "data").Array()
	// if len(toSendArr) <= 0 {
	// 	log.Println("[currently no data to send]")
	// 	return
	// }
	toSend := gjson.Get(jsonStr, "data").Raw
	postOption := &utils.HttpCallRequestOption{
		Url:     fmt.Sprintf("http://%s:%d/lpnode/lpnode_admin_panel/config_lp", globalval.LpNodeHost, globalval.LpNodePort),
		Timeout: 2000,
		JsonStr: toSend,
		Header:  map[string]string{"Content-Type": "application/json"},
		TestOKFun: func(bodyStr string) bool {
			log.Println("bodyis:", bodyStr)
			return gjson.Get(bodyStr, "code").Int() == 200
		},
	}
	_, configResult, err = utils.NewHttpCall().PostJsonCall(postOption)
	log.Println("___________________")
	log.Println(postOption.Url)
	log.Println("send Body json")
	log.Println(toSend)
	log.Println("err:", err)
	log.Println("___________________")
	return
}

func (bcls *BridgeConfigLogicService) ConfigClient() (configResult bool, err error) {

	configResult = false
	chainListStr, err := bcls.GetConfigJsonData()
	if err != nil {
		err = errors.WithMessage(err, "cannot get correct config structure, please check datasource")
		return
	}

	//dwls := NewDexWalletLogicService()
	log.Printf("total need request %d chains", len(gjson.Get(chainListStr, "@this").Map()))
	for chainKey, chainItem := range gjson.Get(chainListStr, "@this").Map() { // chain id level
		log.Println(chainKey, "ChainId:🟥🟥🟥🟥🟥🟥")
		dataStr := `{"data":[]}`
		walletIndex := 0
		for _, wallet := range chainItem.Get("walletInfo").Map() {
			walletName := wallet.Get("walletName").String()
			address := wallet.Get("address").String()
			accountId := wallet.Get("accountId").String()
			privateKey := wallet.Get("privateKey").String()
			walletType := wallet.Get("walletType").String()
			storeId := wallet.Get("storeId").String()
			vaultHostType := wallet.Get("vaultHostType").String()
			vaultName := wallet.Get("vaultName").String()
			vaultSecertType := wallet.Get("vaultSecertType").String()
			dataStr, _ = sjson.Set(dataStr, fmt.Sprintf("data.%d.wallet_name", walletIndex), walletName)
			dataStr, _ = sjson.Set(dataStr, fmt.Sprintf("data.%d.can_sign_712", walletIndex), true)
			dataStr, _ = sjson.Set(dataStr, fmt.Sprintf("data.%d.account_id", walletIndex), accountId)
			dataStr, _ = sjson.Set(dataStr, fmt.Sprintf("data.%d.private_key", walletIndex), privateKey)
			dataStr, _ = sjson.Set(dataStr, fmt.Sprintf("data.%d.address", walletIndex), address)
			isTypeSet := false
			if walletType == "storeId" {
				dataStr, _ = sjson.Set(dataStr, fmt.Sprintf("data.%d.type", walletIndex), "vault")
				dataStr, _ = sjson.Set(dataStr, fmt.Sprintf("data.%d.secert_id", walletIndex), storeId)
				dataStr, _ = sjson.Set(dataStr, fmt.Sprintf("data.%d.private_key", walletIndex), "")
				dataStr, _ = sjson.Set(dataStr, fmt.Sprintf("data.%d.vault_host_type", walletIndex), vaultHostType)
				dataStr, _ = sjson.Set(dataStr, fmt.Sprintf("data.%d.vault_name", walletIndex), vaultName)
				dataStr, _ = sjson.Set(dataStr, fmt.Sprintf("data.%d.vault_secert_type", walletIndex), vaultSecertType)
				isTypeSet = true
			}
			if walletType == "secretVault" {
				dataStr, _ = sjson.Set(dataStr, fmt.Sprintf("data.%d.type", walletIndex), "secret_vault")
				dataStr, _ = sjson.Set(dataStr, fmt.Sprintf("data.%d.vault_host_type", walletIndex), vaultHostType)
				dataStr, _ = sjson.Set(dataStr, fmt.Sprintf("data.%d.vault_name", walletIndex), vaultName)
				dataStr, _ = sjson.Set(dataStr, fmt.Sprintf("data.%d.secert_id", walletIndex), vaultName)
				dataStr, _ = sjson.Set(dataStr, fmt.Sprintf("data.%d.vault_secert_type", walletIndex), vaultSecertType)
				isTypeSet = true
			}
			if !isTypeSet {
				dataStr, _ = sjson.Set(dataStr, fmt.Sprintf("data.%d.type", walletIndex), "key")
				dataStr, _ = sjson.Set(dataStr, fmt.Sprintf("data.%d.vault_host_type", walletIndex), vaultHostType)
				dataStr, _ = sjson.Set(dataStr, fmt.Sprintf("data.%d.vault_name", walletIndex), vaultName)
				dataStr, _ = sjson.Set(dataStr, fmt.Sprintf("data.%d.vault_secert_type", walletIndex), vaultSecertType)
			}

			tokenIndex := 0
			for _, tokens := range wallet.Get("tokenInfo").Map() {
				address := tokens.Get("address").String()
				tokenId := tokens.Get("tokenId").String()
				if chainItem.Get("chainInfo.chainType").String() == "near" {
					dataStr, _ = sjson.Set(dataStr, fmt.Sprintf("data.%d.token_list.%d.token_id", walletIndex, tokenIndex), tokenId)
					dataStr, _ = sjson.Set(dataStr, fmt.Sprintf("data.%d.token_list.%d.create_receipt_id", walletIndex, tokenIndex), address)
				} else {
					dataStr, _ = sjson.Set(dataStr, fmt.Sprintf("data.%d.token_list.%d", walletIndex, tokenIndex), address)
				}
				tokenIndex++
			}
			walletIndex++
			chainId, _ := strconv.ParseInt(chainKey, 10, 64)
			url, getUrlErr := bcls.GetClientSetWalletUrl(chainId)
			log.Println("address need request is", url)
			if getUrlErr != nil {
				if strings.Contains(getUrlErr.Error(), "install record not found") {
					logger.System.Warnf("cannot find install record and url for this chain, stop config temporary 🌏🌏🌏🌏🌏%s", chainKey)
					continue
				}
				err = getUrlErr
				return
			}
			tobeSend := gjson.Get(dataStr, "data").Raw
			requestOption := &utils.HttpCallRequestOption{
				Url:     url,
				Timeout: 10000,
				JsonStr: tobeSend,
				TestOKFun: func(bodyStr string) bool {
					log.Println("bodyis:", bodyStr)
					return true
					// return gjson.Get(bodyStr, "code").Int() == 0
				},
			}
			log.Println("___________________")
			log.Println(tobeSend)
			log.Println(url, chainKey)
			log.Println("___________________")
			_, ok, setWalletErr := utils.NewHttpCall().PostJsonCall(requestOption)
			if setWalletErr != nil {
				err = setWalletErr
				return
			}
			if !ok {
				err = errors.New(fmt.Sprintf("target service return parse result incorrect %s", requestOption.Url))
				log.Println("config error occur", "🟥🟥🟥🟥🟥🟥")
			}
		}
	}

	configResult = true
	return
}
