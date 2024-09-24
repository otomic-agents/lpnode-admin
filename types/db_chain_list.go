package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DbChainListRow struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	ChainId    int                `bson:"chainId"`
	ChainName  string             `bson:"chainName"`
	ChainType  string             `bson:"chainType"`
	DeployName string             `bson:"deployName"`
	RpcTx      string             `bson:"rpcTx"`
	EnvList    []struct {
		STATUS_KEY string `bson:"STATUS_KEY"`
	} `bson:"envList"`
	Image       string  `bson:"image"`
	Name        string  `bson:"name"`
	ServiceName string  `bson:"serviceName"`
	TokenName   string  `bson:"tokenName"`
	TokenUsd    float64 `bson:"tokenUsd"`
	Precision   int     `bson:"precision"`
}
