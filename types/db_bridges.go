package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type DBBridgeRow struct {
	ID                primitive.ObjectID `bson:"_id"`
	BridgeName        string             `bson:"bridgeName"`
	SrcChainId        int64              `bson:"srcChainId"`
	SrcChain_ID       primitive.ObjectID `bson:"srcChain_id"`
	DstChainId        int64              `bson:"dstChainId"`
	DstChain_ID       primitive.ObjectID `bson:"dstChain_id"`
	SrcToken          string             `bson:"srcToken"`
	SrcToken_ID       primitive.ObjectID `bson:"srcToken_id"`
	DstToken          string             `bson:"dstToken"`
	DstToken_ID       primitive.ObjectID `bson:"dstToken_id"`
	WalletName        string             `bson:"walletName"`
	Wallet_ID         primitive.ObjectID `bson:"wallet_id"`
	Src_Wallet_Id     primitive.ObjectID `bson:"src_wallet_id"`
	LpReceiverAddress string             `bson:"lpReceiverAddress"`
	MsmqName          string             `bson:"msmqName"`
	SrcClientUri      string             `bson:"srcClientUri"`
	DstClientUri      string             `bson:"dstClientUri"`
	RelayApiKey       string             `bson:"relayApiKey"`
	AmmName           string             `bson:"ammName"`
	EnableHedge       bool               `bson:"enableHedge"`
	EnableLimiter     bool               `bson:"enableLimiter"`
}

type DBBridgeDstChainAggregateItem struct {
	DstChainId int64
	WalletName string
	AccountId  string
	PrivateKey string
	TokenList  []struct {
		Address string
		TokenId string
	}
}

type DBBridgeUniqDstTokenInfo struct {
	Id           primitive.ObjectID `bson:"_id"`
	Address      string
	TokenId      string `bson:"tokenId"`
	AddressLower string
}
type DBBridgeUniqDstToken struct {
	DstTokenId primitive.ObjectID `bson:"dstToken_id"`
	DstToken   string             `bson:"dstToken"`
	WalletName string             `bson:"walletName"`
}

type TDBBridgeUniqDstToken struct {
	TokenId    string // near only
	DstTokenId string
	DstToken   string
	WalletName string
	Info       struct {
		TokenId      string             // near only
		Id           primitive.ObjectID `bson:"_id"`
		Address      string
		ReceiptId    string
		AddressLower string
	}
}
