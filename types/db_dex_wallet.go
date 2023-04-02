package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type DBWalletRow struct {
	ID              primitive.ObjectID `bson:"_id"`
	WalletName      string             `bson:"walletName"`
	PrivateKey      string             `bson:"privateKey"`
	Address         string             `bson:"address"`
	AddressLower    string             `bson:"addressLower"`
	ChainType       string             `bson:"chainType"`
	ChainId         int64              `bson:"chainId"`
	AccountId       string             `bson:"accountId"`
	StoreId         string             `bson:"storeId"`
	WalletType      string             `bson:"walletType"`
	VaultHostType   string             `bson:"vaultHostType"`
	VaultName       string             `bson:"vaultName"`
	VaultSecertType string             `bson:"vaultSecertType"`
}
