// apps/Hedge/Types.go
package apps_hedge

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Hedge is the API response model
type Hedge struct {
	ID            string   `json:"id"`
	HedgeID       string   `json:"hedge_id,omitempty"`
	Name          string   `json:"name"`
	Exchange      string   `json:"exchange"`
	APIKey        string   `json:"api_key"`
	BridgeID      string   `json:"bridge_id,omitempty"`
	Pair          string   `json:"pair,omitempty"`
	ChainPair     []string `json:"chain_pair,omitempty"`
	AmmName       string   `json:"amm_name,omitempty"`
	CexAccountID  string   `json:"cex_account_id,omitempty"`
	Status        string   `json:"status"`
	StatusDesc    string   `json:"status_desc,omitempty"`
	LockedByHedge bool     `json:"locked_by_hedge,omitempty"`
	ActivePeriod  *struct {
		Start string `json:"start"`
		End   string `json:"end"`
	} `json:"active_period,omitempty"`
	Exposure  map[string]float64 `json:"exposure,omitempty"`
	CreatedBy string             `json:"created_by,omitempty"`
	CreatedAt string             `json:"created_at"`
	UpdatedAt string             `json:"updated_at"`
	Version   int                `json:"version,omitempty"`
}

type DBHedge struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Name          string             `bson:"name"`
	Exchange      string             `bson:"exchange"` // Added based on context and Goa design
	APIKey        string             `bson:"api_key"`
	APISecret     string             `bson:"api_secret"`               // Sensitive, for internal use
	Passphrase    string             `bson:"passphrase,omitempty"`     // Optional, sensitive
	BridgeID      primitive.ObjectID `bson:"bridge_id,omitempty"`      // Optional, link to bridge
	Pair          string             `bson:"pair,omitempty"`           // Optional
	ChainPair     []string           `bson:"chain_pair,omitempty"`     // Optional
	AmmName       string             `bson:"amm_name,omitempty"`       // Optional
	CexAccountID  primitive.ObjectID `bson:"cex_account_id,omitempty"` // Optional, potentially redundant with _id
	Status        string             `bson:"status"`                   // e.g., "validating", "active", "inactive", "invalid_keys"
	StatusDesc    string             `bson:"status_desc,omitempty"`
	LockedByHedge bool               `bson:"locked_by_hedge,omitempty"`
	ActivePeriod  *struct {
		Start time.Time `bson:"start"`
		End   time.Time `bson:"end"`
	} `bson:"active_period,omitempty"`
	InitialSnapshot *struct {
		Cex       map[string]float64 `bson:"cex"`
		Dex       map[string]float64 `bson:"dex"`
		Timestamp time.Time          `bson:"timestamp"`
	} `bson:"initial_snapshot,omitempty"`
	CurrentSnapshot *struct {
		Cex       map[string]float64 `bson:"cex"`
		Dex       map[string]float64 `bson:"dex"`
		Timestamp time.Time          `bson:"timestamp"`
	} `bson:"current_snapshot,omitempty"`
	Exposure   map[string]float64 `bson:"exposure,omitempty"`
	RiskConfig *struct {
		MaxAssetExposure map[string]float64 `bson:"max_asset_exposure"`
		MinHedgeAmount   map[string]float64 `bson:"min_hedge_amount"`
		HedgeMode        string             `bson:"hedge_mode"`
	} `bson:"risk_config,omitempty"`
	CreatedBy string    `bson:"created_by,omitempty"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	Version   int       `bson:"version,omitempty"`
}

type CreateHedgeParams struct {
	Name         string
	CexAccountID string
	BridgeID     string
	AmmName      string
}

// UpdateHedgeParams defines parameters for updating a hedge task
type UpdateHedgeParams struct {
	ID           string
	Name         string
	CexAccountID string
	BridgeID     string
	AmmName      string
	Status       string
}
type WalletInfo struct {
	Address      string
	Name         string
	Token        string
	TokenAddress string
	ChainId      int64
}
type UpdateHedgeDataParams struct {
	ID                string
	Name              string
	BridgeID          string
	CexAccountID      string
	AmmName           string
	InitialSnapshot   map[string]map[string]float64
	RiskConfig        RiskConfig
	SourceWallet      WalletInfo
	DestinationWallet WalletInfo
}
type RiskConfig struct {
	MaxAssetExposure map[string]float64
	MinHedgeAmount   map[string]float64
	HedgeMode        string
}

// ErrHedgeNotFound is the error returned when a hedge task is not found
var ErrHedgeNotFound = fmt.Errorf("hedge task not found")
