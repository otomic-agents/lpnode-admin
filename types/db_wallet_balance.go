package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type DBWalletBalanceValue struct {
	Type string `bson:"type"`
	Hex  string `bson:"hex"`
}

type DBWalletBalance struct {
	ID             primitive.ObjectID   `bson:"_id"`
	ChainId        int64                `bson:"chainId"`
	Token          string               `bson:"token"`
	WalletAddress  string               `bson:"wallet_address"`
	BalanceValue   DBWalletBalanceValue `bson:"balance_value"`
	ChainName      string               `bson:"chainName"`
	Decimals       int64                `bson:"decimals"`
	LastUpdateTime primitive.DateTime   `bson:"lastUpdateTime"`
	UpdatedAt      primitive.DateTime   `bson:"updatedAt"`
	WalletName     string               `bson:"wallet_name"`
}
