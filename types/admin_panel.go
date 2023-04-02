package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type ChainInfoStoreItem struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id"`
	ChainName string             `json:"chainName"`
	Name      string             `json:"name"`
	ChainId   int64              `json:"chainId"`
	TokenName string             `json:"tokenName"`
	TokenUsd  int64              `json:"tokenUsd"`
}
