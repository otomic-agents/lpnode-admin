package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type DBTokenRow struct {
	ID           primitive.ObjectID `bson:"_id"`
	ReceiptId    string             `bson:"receiptId"` // Near only
	ChainId      int64              `bson:"chainId"`
	Address      string             `bson:"address"`
	TokenId      string             `bson:"tokenId"`
	AddressLower string             `bson:"addressLower"`
	AddressIndex string             `bson:"addressIndex"`
	TokenName    string             `bson:"tokenName"`
	MarketName   string             `bson:"marketName"`
	Precision    int64              `bson:"precision"`
	CoinType     string             `bson:"coinType"`
	ChainType    string             `bson:"chainType"`
}
