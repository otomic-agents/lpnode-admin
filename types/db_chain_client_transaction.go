package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type View_ChainTransaction struct {
	BusinessID      string                 `json:"businessId" bson:"businessId"`
	EventName       string                 `json:"eventName" bson:"eventName"`
	SystemChainId   string                 `json:"systemChainId" bson:"systemChainId"`
	SrcChainName    string                 `json:"srcChainName" bson:"srcChainName"`
	DstChainName    string                 `json:"dstChainName" bson:"dstChainName"`
	SrcToken        string                 `json:"srcToken" bson:"srcToken"`
	DstToken        string                 `json:"dstToken" bson:"dstToken"`
	Status          string                 `json:"status" bson:"status"`
	LatestGasPrice  string                 `json:"latestGasPrice" bson:"latestGasPrice"`
	IncreasedNumber int                    `json:"increasedNumber" bson:"increasedNumber"`
	LatestSend      int64                  `json:"latestSend" bson:"latestSend"`
	SendList        []Db_SendRecord        `json:"sendList" bson:"sendList"`
	IncGasList      []Db_GasIncreaseRecord `json:"incGasList" bson:"incGasList"`
}

// Db_ChainTransaction Blockchain transaction record
type Db_ChainTransaction struct {
	ID              primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	BusinessID      string                 `json:"businessId" bson:"businessId"`
	EventName       string                 `json:"eventName" bson:"eventName"`
	ChainID         int                    `json:"chainId" bson:"chainId"`
	SendTime        int64                  `json:"sendTime" bson:"sendTime"`
	GasPrice        string                 `json:"gasPrice" bson:"gasPrice"`
	SystemChainID   string                 `json:"systemChainId" bson:"systemChainId"`
	ErrorMessage    string                 `json:"errorMessage" bson:"errorMessage,omitempty"`
	SendID          string                 `json:"sendId" bson:"sendId"`
	TxHash          string                 `json:"txHash" bson:"txHash"`
	Status          string                 `json:"status" bson:"status"`
	NeedNewGasPrice bool                   `json:"needNewGasPrice" bson:"needNewGasPrice"`
	IncreasedNumber int                    `json:"increasedNumber" bson:"increasedNumber"`
	CreatedAt       int64                  `json:"createdAt" bson:"createdAt"`
	UpdatedAt       int64                  `json:"updatedAt" bson:"updatedAt"`
	SendList        []Db_SendRecord        `json:"sendList" bson:"sendList"`
	IncGasList      []Db_GasIncreaseRecord `json:"incGasList" bson:"incGasList"`
}

// Db_SendRecord Send record
type Db_SendRecord struct {
	Timestamp       int64  `json:"timestamp" bson:"timestamp"`
	GasPrice        string `json:"gasPrice" bson:"gasPrice"`
	SendID          string `json:"sendId" bson:"sendId"`
	TxHash          string `json:"txHash" bson:"txHash"`
	Status          string `json:"status" bson:"status"`
	NeedNewGasPrice bool   `json:"needNewGasPrice" bson:"needNewGasPrice"`
	Error           string `json:"error,omitempty" bson:"error,omitempty"`
}

// Db_GasIncreaseRecord Gas increase record
type Db_GasIncreaseRecord struct {
	Timestamp        int64  `json:"timestamp" bson:"timestamp"`
	OriginalGasPrice string `json:"originalGasPrice" bson:"originalGasPrice"`
	NewGasPrice      string `json:"newGasPrice" bson:"newGasPrice"`
}

// CollectionName Returns collection name
func (Db_ChainTransaction) CollectionName() string {
	return "chain_transactions"
}
