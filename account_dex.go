package adminapiservice

import (
	accountdex "admin-panel/gen/account_dex"
	"admin-panel/service"
	"context"
	"log"
	"time"

	"github.com/aws/smithy-go/ptr"
	"github.com/davecgh/go-spew/spew"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// accountDex service example implementation.
// The example methods log the requests and return zero values.
type accountDexsrvc struct {
	logger *log.Logger
}

// NewAccountDex returns the accountDex service implementation.
func NewAccountDex(logger *log.Logger) accountdex.Service {
	return &accountDexsrvc{logger}
}

// WalletInfo implements walletInfo.
func (s *accountDexsrvc) WalletInfo(ctx context.Context, p *accountdex.WalletInfoPayload) (res *accountdex.WalletInfoResult, err error) {
	res = &accountdex.WalletInfoResult{}
	objID, err := primitive.ObjectIDFromHex(p.ID)
	if err != nil {
		return nil, err
	}
	spew.Dump(p.ID)
	wallet, findErr := service.NewDexWalletLogicService().FindOneByFilter(bson.M{
		"_id": objID,
	})
	if findErr != nil {
		err = findErr
		return
	}
	spew.Dump(wallet)
	res.Code = ptr.Int64(0)
	res.Result = &accountdex.ADBWalletInfo{
		ID:           ptr.String(wallet.ID.Hex()),
		Address:      ptr.String(wallet.Address),
		AddressLower: ptr.String(wallet.AddressLower),
		ChainID:      ptr.Int(int(wallet.ChainId)),
		ChainType:    ptr.String(wallet.ChainType),
		WalletName:   ptr.String(wallet.WalletName),
	}
	s.logger.Print("accountDex.walletInfo")
	return
}

// GetWalletAssets implements getWalletAssets.
func (s *accountDexsrvc) GetWalletAssets(ctx context.Context, p *accountdex.GetWalletAssetsPayload) (res *accountdex.GetWalletAssetsResult, err error) {
	s.logger.Print("accountDex.getWalletAssets")

	// Initialize result
	res = &accountdex.GetWalletAssetsResult{
		Code:    ptr.Int(0),
		Message: "success",
	}

	// Create service instance
	walletService := &service.MultiChainWalletAssetsService{}

	// Get monitored addresses
	monitoredAddresses, err := walletService.GetMonitoredAddresses(ctx)
	if err != nil {
		return &accountdex.GetWalletAssetsResult{
			Code:    ptr.Int(30001),
			Message: "Failed to get monitored addresses: " + err.Error(),
		}, nil
	}

	// Get last update time
	lastUpdateTime, err := walletService.GetLastUpdateTime(ctx)
	if err != nil {
		return &accountdex.GetWalletAssetsResult{
			Code:    ptr.Int(30001),
			Message: "Failed to get last update time: " + err.Error(),
		}, nil
	}

	// Determine query addresses
	addresses := monitoredAddresses
	if p.Addresses != nil && len(p.Addresses) > 0 {
		addresses = p.Addresses
	}

	// Get wallet assets
	walletAssets, err := walletService.GetWalletAssets(ctx, addresses)
	if err != nil {
		return &accountdex.GetWalletAssetsResult{
			Code:    ptr.Int(30001),
			Message: "Failed to get wallet assets: " + err.Error(),
		}, nil
	}

	// Build response
	responseData := &accountdex.ADBWalletAssetResponse{
		TotalAddresses: ptr.Int(walletAssets.TotalAddresses),
		LastUpdated:    ptr.String(lastUpdateTime.Format(time.RFC3339)),
		TotalValue:     ptr.String("0"),
		ChainAssets:    make([]*accountdex.ADBChainAssetGroup, 0),
	}

	// Convert chain asset data
	for _, chainAsset := range walletAssets.ChainAssets {
		chainGroup := &accountdex.ADBChainAssetGroup{
			ChainID:       ptr.Int(chainAsset.ChainID),
			ChainName:     ptr.String(chainAsset.ChainName),
			ChainType:     ptr.String(chainAsset.ChainType),
			NativeToken:   ptr.String(chainAsset.NativeToken),
			AddressAssets: make([]*accountdex.ADBAddressAssetGroup, 0),
			TotalValue:    ptr.String("0"),
		}

		// Convert address assets
		for _, addressAsset := range chainAsset.AddressAssets {
			addressGroup := &accountdex.ADBAddressAssetGroup{
				WalletAddress: ptr.String(addressAsset.WalletAddress),
				WalletNames:   addressAsset.WalletNames,
				Tokens:        make([]*accountdex.ADBTokenBalance, 0),
				TotalValue:    ptr.String("0"),
			}

			// Convert tokens
			for _, token := range addressAsset.Tokens {
				tokenBalance := &accountdex.ADBTokenBalance{
					TokenAddress:     ptr.String(token.TokenAddress),
					Symbol:           ptr.String(token.Symbol),
					FormattedBalance: ptr.String(token.FormattedBalance),
					Decimals:         ptr.Int(token.Decimals),
					IsNative:         ptr.Bool(token.IsNative),
					Price:            ptr.String("0"),
					Value:            ptr.String("0"),
					UpdatedAt:        ptr.String(token.UpdatedAt.Format(time.RFC3339)),
				}

				addressGroup.Tokens = append(addressGroup.Tokens, tokenBalance)
			}

			chainGroup.AddressAssets = append(chainGroup.AddressAssets, addressGroup)
		}

		responseData.ChainAssets = append(responseData.ChainAssets, chainGroup)
	}
	res.Code = ptr.Int(0)
	res.Message = "sucess"
	res.Result = responseData
	return res, nil
}
