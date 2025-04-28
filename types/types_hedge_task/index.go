package types_hedge_task

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// HedgeTask hedge task model
type HedgeTask struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Status    string             `bson:"status" json:"status"` // pending, active, completed, failed, cancelled
	AccountID primitive.ObjectID `bson:"account_id" json:"account_id"`
	BridgeID  primitive.ObjectID `bson:"bridge_id" json:"bridge_id"`
	Balances  map[string]Balance `bson:"balances" json:"balances"` // key is token symbol, like "ETH", "BNB"
	Exposure  Exposure           `bson:"exposure" json:"exposure"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	ClosedAt  *time.Time         `bson:"closed_at,omitempty" json:"closed_at,omitempty"`
	CreatedBy string             `bson:"created_by" json:"created_by"`
	Notes     string             `bson:"notes" json:"notes"`
}

// Balance token balance information
type Balance struct {
	DEX DEXBalance `bson:"dex" json:"dex"`
	CEX CEXBalance `bson:"cex" json:"cex"`
}

// DEXBalance balance information on DEX
type DEXBalance struct {
	Chain    string `bson:"chain" json:"chain"`
	Amount   string `bson:"amount" json:"amount"`
	PriceUSD string `bson:"price_usd" json:"price_usd"`
}

// CEXBalance balance information on CEX
type CEXBalance struct {
	Amount   string `bson:"amount" json:"amount"`
	PriceUSD string `bson:"price_usd" json:"price_usd"`
}

// Exposure exposure information
type Exposure struct {
	TokenExposures map[string]TokenExposure `bson:"token_exposures" json:"token_exposures"` // key is token symbol, like "ETH", "BNB"
	TotalUSDValue  string                   `bson:"total_usd_value" json:"total_usd_value"` // total USD value
	Percentage     string                   `bson:"percentage" json:"percentage"`           // exposure percentage
}

// TokenExposure exposure information for a single token
type TokenExposure struct {
	NetAmount   string `bson:"net_amount" json:"net_amount"`       // DEX amount - CEX amount
	NetUSDValue string `bson:"net_usd_value" json:"net_usd_value"` // net USD value
}

// HedgeTaskCreateParams parameters for creating a hedge task
type HedgeTaskCreateParams struct {
	Name      string             `json:"name" validate:"required"`
	AccountID primitive.ObjectID `json:"account_id" validate:"required"`
	BridgeID  primitive.ObjectID `json:"bridge_id" validate:"required"`
	Notes     string             `json:"notes"`
}

// HedgeTaskStatusUpdateParams parameters for updating hedge task status
type HedgeTaskStatusUpdateParams struct {
	Status string `json:"status" validate:"required,oneof=active completed cancelled"`
}
