package types

import "time"

// Schema_WalletBalance represents the structure of wallet balance records in the database
type Schema_WalletBalance struct {
	ID               string    `bson:"_id,omitempty" json:"id,omitempty"`
	ChainID          int       `bson:"chainId" json:"chainId"`
	TokenAddress     string    `bson:"tokenAddress" json:"tokenAddress"`
	WalletAddress    string    `bson:"walletAddress" json:"walletAddress"`
	Balance          string    `bson:"balance" json:"balance"`
	ChainType        string    `bson:"chainType" json:"chainType"`
	Decimals         int       `bson:"decimals" json:"decimals"`
	FormattedBalance string    `bson:"formattedBalance" json:"formattedBalance"`
	Symbol           string    `bson:"symbol" json:"symbol"`
	UpdatedAt        time.Time `bson:"updatedAt" json:"updatedAt"`
	WalletNames      []string  `bson:"walletNames" json:"walletNames"`
}

// Schema_ChainInfo represents the structure of chain information in the database
type Schema_ChainInfo struct {
	ID          string `bson:"_id,omitempty" json:"id,omitempty"`
	ChainID     int    `bson:"chainId" json:"chainId"`
	ChainName   string `bson:"chainName" json:"chainName"`
	ChainType   string `bson:"chainType" json:"chainType"`
	TokenName   string `bson:"tokenName" json:"tokenName"`
	ServiceName string `bson:"serviceName" json:"serviceName"`
}

// Schema_AddressDoc used for aggregation query results
type Schema_AddressDoc struct {
	Address string `bson:"address"`
}

// Schema_UpdateTime used to query the last update time
type Schema_UpdateTime struct {
	UpdatedAt time.Time `bson:"updatedAt"`
}

// Type_WalletAssetResponse represents wallet asset response
type Type_WalletAssetResponse struct {
	TotalAddresses int                    `json:"totalAddresses"`
	LastUpdated    time.Time              `json:"lastUpdated"`
	ChainAssets    []Type_ChainAssetGroup `json:"chainAssets"`
}

// Type_ChainAssetGroup represents chain asset grouping
type Type_ChainAssetGroup struct {
	ChainID       int                      `json:"chainId"`
	ChainName     string                   `json:"chainName"`
	ChainType     string                   `json:"chainType"`
	NativeToken   string                   `json:"nativeToken"`
	AddressAssets []Type_AddressAssetGroup `json:"addressAssets"`
}

// Type_AddressAssetGroup represents address asset grouping
type Type_AddressAssetGroup struct {
	WalletAddress string              `json:"walletAddress"`
	WalletNames   []string            `json:"walletNames"`
	Tokens        []Type_TokenBalance `json:"tokens"`
}

// Type_TokenBalance represents token balance
type Type_TokenBalance struct {
	TokenAddress     string    `json:"tokenAddress"`
	Symbol           string    `json:"symbol"`
	FormattedBalance string    `json:"formattedBalance"`
	Decimals         int       `json:"decimals"`
	IsNative         bool      `json:"isNative"`
	UpdatedAt        time.Time `json:"updatedAt"`
}
