package apps_cex_account

import "go.mongodb.org/mongo-driver/bson/primitive"

// CreateCexAccountPayload represents the data needed to create a CEX account
type CreateCexAccountPayload struct {
	Name       string `json:"name"`
	Exchange   string `json:"exchange"`
	APIKey     string `json:"api_key"`
	APISecret  string `json:"api_secret"`
	Passphrase string `json:"passphrase,omitempty"` // Optional
}

// CexAccount represents a CEX account in the system
type CexAccount struct {
	ID        string `json:"id" bson:"_id,omitempty"`
	Name      string `json:"name" bson:"name,omitempty"`
	Exchange  string `json:"exchange" bson:"exchange,omitempty"`
	APIKey    string `json:"api_key" bson:"api_key,omitempty"` // Masked in API response
	Status    string `json:"status" bson:"status,omitempty"`
	CreatedAt string `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at" bson:"updated_at,omitempty"`
}

// DBCexAccount represents the database model for CEX accounts
type DBCexAccount struct {
	ID         primitive.ObjectID  `bson:"_id,omitempty"`
	Name       string              `bson:"name"`
	Exchange   string              `bson:"exchange"`
	APIKey     string              `bson:"api_key"`
	APISecret  string              `bson:"api_secret"`
	Passphrase string              `bson:"passphrase,omitempty"`
	Status     string              `bson:"status"`
	CreatedAt  primitive.DateTime  `bson:"created_at"`
	UpdatedAt  primitive.DateTime  `bson:"updated_at"`
	IsDeleted  bool                `bson:"is_deleted,omitempty"`
	DeletedAt  *primitive.DateTime `bson:"deleted_at,omitempty"`
}

// TokenBalance represents the balance information for a single token
type TokenBalance struct {
	Asset  string `json:"asset" bson:"asset"`
	Free   string `json:"free" bson:"free"`
	Locked string `json:"locked" bson:"locked"`
	Total  string `json:"total" bson:"total"`
}

// ExchangeBalanceResult represents the balance result obtained from the exchange
type ExchangeBalanceResult struct {
	Info struct {
		Timestamp string `json:"timestamp" bson:"timestamp"`
		Exchange  string `json:"exchange" bson:"exchange"`
		Account   string `json:"account" bson:"account"`
	} `json:"info" bson:"info"`
	Balances []TokenBalance `json:"balances" bson:"balances"`
}

// CexWalletBalance represents the CEX wallet balance information stored in the database
type CexWalletBalance struct {
	ID             primitive.ObjectID    `json:"_id" bson:"_id,omitempty"`
	AccountId      string                `json:"accountId" bson:"accountId"`
	Exchange       string                `json:"exchange" bson:"exchange"`
	ExchangeResult ExchangeBalanceResult `json:"exchangeResult" bson:"exchangeResult"`
	LastUpdateTime primitive.DateTime    `json:"lastUpdateTime" bson:"lastUpdateTime"`
	Name           string                `json:"name" bson:"name"`
	UpdatedAt      primitive.DateTime    `json:"updatedAt" bson:"updatedAt"`
}
