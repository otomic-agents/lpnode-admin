package adminapiservice

import (
	basedata "admin-panel/gen/base_data"
	database "admin-panel/mongo_database"
	"admin-panel/service"
	"admin-panel/types"
	"context"
	"log"
	"os"

	"github.com/aws/smithy-go/ptr"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
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
func (s *baseDatasrvc) GetLpInfo(ctx context.Context) (res *basedata.GetLpInfoResult, err error) {
	res = &basedata.GetLpInfoResult{Result: nil, Code: ptr.Int64(0), Message: ptr.String("")}
	ret := &basedata.LpInfo{}
	err = database.FindOne("main", "relayAccounts", bson.M{}, ret)
	if err != nil {
		return
	}
	res.Result = ret
	return
}

// Get wallet list with their associated tokens
func (s *baseDatasrvc) GetWalletAndTokens(ctx context.Context, p *basedata.GetWalletAndTokensPayload) (res *basedata.GetWalletAndTokensResult, err error) {
	res = &basedata.GetWalletAndTokensResult{}
	s.logger.Print("🔍 baseData.getWalletAndTokens - Starting to fetch wallet and tokens")

	// Fetch configuration data
	bcls := service.NewBridgeConfigLogicService()
	configStr, err := bcls.GetConfigData(p.ChainID)
	if err != nil {
		s.logger.Printf("❌ Failed to get config data for chain ID %s: %v", p.ChainID, err)
		err = errors.WithMessage(err, "get config error:")
		return
	}
	s.logger.Printf("✅ Successfully fetched config data for chain ID %d", p.ChainID)

	// Parse JSON data using gjson
	walletArray := gjson.Parse(configStr).Array()
	var walletItems []*basedata.WalletItem

	// Iterate through each wallet
	for _, wallet := range walletArray {
		walletItem := &basedata.WalletItem{
			WalletName:              ptr.String(wallet.Get("wallet_name").String()),
			Address:                 ptr.String(wallet.Get("address").String()),
			CanSign:                 ptr.Bool(wallet.Get("can_sign").Bool()),
			CanSign712:              ptr.Bool(wallet.Get("can_sign_712").Bool()),
			Type:                    ptr.String(wallet.Get("type").String()),
			SignatureServiceAddress: ptr.String(wallet.Get("signature_service_address").String()),
		}

		// Process token list
		tokenList := wallet.Get("token_list").Array()
		var tokens []*basedata.WalletTokenItem
		for _, token := range tokenList {
			tokens = append(tokens, &basedata.WalletTokenItem{
				Address:  ptr.String(token.String()),
				Symbol:   ptr.String(""), // Symbol is optional, leave empty if not provided
				Decimals: ptr.Int32(0),   // Decimals is optional, default to 0 if not provided
			})
		}
		walletItem.Tokens = tokens

		// Add wallet item to the result
		walletItems = append(walletItems, walletItem)
		s.logger.Printf("➕ Added wallet: %s with %d tokens", *walletItem.WalletName, len(tokens))
	}

	// Set response data
	res.Code = ptr.Int64(0)
	res.Result = walletItems
	res.Message = ptr.String("success")
	s.logger.Print("🎉 Successfully processed wallet and tokens data")
	return res, nil
}
