package service

import (
	database "admin-panel/mongo_database"
	"admin-panel/types"
	"context"
	"sort"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

// MultiChainWalletAssetsService handles services related to cross-chain wallet assets and DEX
type MultiChainWalletAssetsService struct {
}

// GetWalletAssets gets assets for specified wallet addresses
func (s *MultiChainWalletAssetsService) GetWalletAssets(ctx context.Context, addresses []string) (*types.Type_WalletAssetResponse, error) {
	// 1. Query balance records
	balances, err := s.fetchWalletBalances(ctx, addresses)
	if err != nil {
		return nil, err
	}

	// 2. Get chain information
	chainInfoMap, err := s.fetchChainInfos(ctx, balances)
	if err != nil {
		return nil, err
	}

	// 3. Organize data structure
	response := s.buildResponse(addresses, balances, chainInfoMap)
	return response, nil
}

// fetchWalletBalances gets wallet balance records
func (s *MultiChainWalletAssetsService) fetchWalletBalances(ctx context.Context, addresses []string) ([]types.Schema_WalletBalance, error) {
	filter := bson.M{}
	if len(addresses) > 0 {
		filter["walletAddress"] = bson.M{"$in": addresses}
	}

	err, cursor := database.FindAll("main", "wallet_dex_all_balances", filter)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to query balance records")
	}
	defer cursor.Close(ctx)

	var balances []types.Schema_WalletBalance
	if err := cursor.All(ctx, &balances); err != nil {
		return nil, errors.WithMessage(err, "Failed to parse balance records")
	}

	return balances, nil
}

// fetchChainInfos gets chain information
func (s *MultiChainWalletAssetsService) fetchChainInfos(ctx context.Context, balances []types.Schema_WalletBalance) (map[int]types.Schema_ChainInfo, error) {
	// Extract all unique chain IDs
	chainIDs := make(map[int]bool)
	for _, balance := range balances {
		chainIDs[balance.ChainID] = true
	}

	chainIDsList := make([]int, 0, len(chainIDs))
	for id := range chainIDs {
		chainIDsList = append(chainIDsList, id)
	}

	// Query chain information
	chainFilter := bson.M{
		"chainId": bson.M{"$in": chainIDsList},
	}
	err, chainCursor := database.FindAll("main", "chainList", chainFilter)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to query chain information")
	}
	defer chainCursor.Close(ctx)

	var chainInfos []types.Schema_ChainInfo
	if err := chainCursor.All(ctx, &chainInfos); err != nil {
		return nil, errors.WithMessage(err, "Failed to parse chain information")
	}

	// Create mapping from chain ID to chain information
	chainInfoMap := make(map[int]types.Schema_ChainInfo)
	for _, info := range chainInfos {
		chainInfoMap[info.ChainID] = info
	}

	return chainInfoMap, nil
}

// buildResponse builds response data
func (s *MultiChainWalletAssetsService) buildResponse(addresses []string, balances []types.Schema_WalletBalance, chainInfoMap map[int]types.Schema_ChainInfo) *types.Type_WalletAssetResponse {
	// Group by chain ID
	chainGroups := s.groupBalancesByChain(balances)

	// Build response
	response := &types.Type_WalletAssetResponse{
		TotalAddresses: len(addresses),
		LastUpdated:    time.Now(),
		ChainAssets:    make([]types.Type_ChainAssetGroup, 0),
	}

	// Process assets for each chain
	for chainID, balances := range chainGroups {
		chainInfo, exists := chainInfoMap[chainID]
		if !exists {
			continue
		}

		chainAsset := s.buildChainAsset(chainID, chainInfo, balances)
		if len(chainAsset.AddressAssets) > 0 {
			response.ChainAssets = append(response.ChainAssets, chainAsset)
		}
	}

	// Sort chains, chains with assets come first
	s.sortChainAssets(response.ChainAssets)

	return response
}

// groupBalancesByChain groups balance records by chain ID
func (s *MultiChainWalletAssetsService) groupBalancesByChain(balances []types.Schema_WalletBalance) map[int][]types.Schema_WalletBalance {
	chainGroups := make(map[int][]types.Schema_WalletBalance)
	for _, balance := range balances {
		chainGroups[balance.ChainID] = append(chainGroups[balance.ChainID], balance)
	}
	return chainGroups
}

// buildChainAsset builds chain asset group
func (s *MultiChainWalletAssetsService) buildChainAsset(chainID int, chainInfo types.Schema_ChainInfo, balances []types.Schema_WalletBalance) types.Type_ChainAssetGroup {
	// Group by address
	addressGroups := make(map[string][]types.Schema_WalletBalance)
	for _, balance := range balances {
		addressGroups[balance.WalletAddress] = append(addressGroups[balance.WalletAddress], balance)
	}

	// Create chain asset group
	chainAsset := types.Type_ChainAssetGroup{
		ChainID:       chainID,
		ChainName:     chainInfo.ChainName,
		ChainType:     chainInfo.ChainType,
		NativeToken:   chainInfo.TokenName,
		AddressAssets: make([]types.Type_AddressAssetGroup, 0),
	}

	// Process assets for each address
	for address, balances := range addressGroups {
		addressAsset := s.buildAddressAsset(address, balances)
		chainAsset.AddressAssets = append(chainAsset.AddressAssets, addressAsset)
	}

	return chainAsset
}

// buildAddressAsset builds address asset group
func (s *MultiChainWalletAssetsService) buildAddressAsset(address string, balances []types.Schema_WalletBalance) types.Type_AddressAssetGroup {
	addressAsset := types.Type_AddressAssetGroup{
		WalletAddress: address,
		WalletNames:   balances[0].WalletNames, // Assume wallet names are the same for the same address
		Tokens:        make([]types.Type_TokenBalance, 0),
	}

	// Process each token balance
	for _, balance := range balances {
		tokenBalance := s.buildTokenBalance(balance)
		addressAsset.Tokens = append(addressAsset.Tokens, tokenBalance)
	}

	// Sort tokens, native tokens come first
	s.sortTokenBalances(addressAsset.Tokens)

	return addressAsset
}

// buildTokenBalance builds token balance
func (s *MultiChainWalletAssetsService) buildTokenBalance(balance types.Schema_WalletBalance) types.Type_TokenBalance {
	// Determine if it's a native token
	isNative := balance.TokenAddress == "0x0000000000000000000000000000000000000000" ||
		balance.TokenAddress == "11111111111111111111111111111111" // Solana system program ID

	return types.Type_TokenBalance{
		TokenAddress:     balance.TokenAddress,
		Symbol:           balance.Symbol,
		FormattedBalance: balance.FormattedBalance,
		Decimals:         balance.Decimals,
		IsNative:         isNative,
		UpdatedAt:        balance.UpdatedAt,
	}
}

// sortTokenBalances sorts token balances
func (s *MultiChainWalletAssetsService) sortTokenBalances(tokens []types.Type_TokenBalance) {
	sort.Slice(tokens, func(i, j int) bool {
		if tokens[i].IsNative && !tokens[j].IsNative {
			return true
		}
		if !tokens[i].IsNative && tokens[j].IsNative {
			return false
		}
		return tokens[i].Symbol < tokens[j].Symbol
	})
}

// sortChainAssets sorts chain assets
func (s *MultiChainWalletAssetsService) sortChainAssets(chainAssets []types.Type_ChainAssetGroup) {
	sort.Slice(chainAssets, func(i, j int) bool {
		return len(chainAssets[i].AddressAssets) > len(chainAssets[j].AddressAssets)
	})
}

// GetMonitoredAddresses gets all monitored addresses
func (s *MultiChainWalletAssetsService) GetMonitoredAddresses(ctx context.Context) ([]string, error) {
	// Aggregation query to get all unique wallet addresses
	pipeline := []bson.M{
		{
			"$group": bson.M{
				"_id": "$walletAddress",
			},
		},
		{
			"$project": bson.M{
				"_id":     0,
				"address": "$_id",
			},
		},
	}

	cursor, err := database.Aggregate("main", "wallet_dex_all_balances", pipeline)
	if err != nil {
		return nil, errors.WithMessage(err, "Aggregation query failed")
	}
	defer cursor.Close(ctx)

	var results []types.Schema_AddressDoc
	if err := cursor.All(ctx, &results); err != nil {
		return nil, errors.WithMessage(err, "Failed to parse aggregation results")
	}

	addresses := make([]string, len(results))
	for i, doc := range results {
		addresses[i] = doc.Address
	}

	return addresses, nil
}

// GetLastUpdateTime gets the last update time
func (s *MultiChainWalletAssetsService) GetLastUpdateTime(ctx context.Context) (time.Time, error) {
	pipeline := []bson.M{
		{
			"$sort": bson.M{"updatedAt": -1},
		},
		{
			"$limit": 1,
		},
		{
			"$project": bson.M{
				"_id":       0,
				"updatedAt": 1,
			},
		},
	}

	cursor, err := database.Aggregate("main", "wallet_dex_all_balances", pipeline)
	if err != nil {
		return time.Time{}, errors.WithMessage(err, "Failed to query last update time")
	}
	defer cursor.Close(ctx)

	var result types.Schema_UpdateTime
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return time.Time{}, errors.WithMessage(err, "Failed to parse last update time")
		}
		return result.UpdatedAt, nil
	}

	return time.Time{}, nil // No records
}
