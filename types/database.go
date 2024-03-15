package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type MongoChainListRow struct {
	ID        primitive.ObjectID `bson:"_id"`
	ChainId   int64              `bson:"chainId" json:"chainId"`
	ChainName string             `bson:"chainName" json:"chainName"`
	Name      string             `bson:"name"`
	TokenName string             `bson:"to kenName" json:"tokenName"`
	MinUsd    int64              `bson:"tokenUsd" json:"minUsd"`
	ChainType string             `bson:"chainType" json:"chainType"`
}

type MongoTokenToSymbolRow struct {
	Symbol    string `bson:"symbol"`
	CoinType  string `bson:"coinType"`
	ChainId   int64  `bson:"chainId"`
	Token     string `bson:"token"`
	TokenName string `bson:"tokenName"`
	TokenId   string `bson:"tokenId"`
}
