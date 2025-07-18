package apps_exchange_rate

import (
	database "admin-panel/mongo_database"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DBExchangeRate represents the exchange rate information model in the database
type DBExchangeRate struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Symbol            string             `bson:"symbol" json:"symbol"`
	Asks              [][]string         `bson:"asks" json:"asks"`
	Bids              [][]string         `bson:"bids" json:"bids"`
	IncomingTimestamp int64              `bson:"incomingTimestamp" json:"incomingTimestamp"`
	LastUpdateId      int64              `bson:"lastUpdateId" json:"lastUpdateId"`
	StdSymbol         string             `bson:"stdSymbol" json:"stdSymbol"`
	Stream            string             `bson:"stream" json:"stream"`
	Timestamp         int64              `bson:"timestamp" json:"timestamp"`
}

// ExchangeRateService handles services related to exchange rates
type ExchangeRateService struct {
	logger *log.Logger
}

// NewExchangeRateService creates a new instance of ExchangeRateService
func NewExchangeRateService(logger *log.Logger) *ExchangeRateService {
	return &ExchangeRateService{
		logger: logger,
	}
}

// GetAllRates retrieves exchange rate information for all trading pairs
func (s *ExchangeRateService) GetAllRates(ctx context.Context) ([]DBExchangeRate, error) {
	s.logger.Printf("Retrieving exchange rate information for all trading pairs")

	// Query all exchange rate information from the database
	err, cursor := database.FindAll("main", "market_orderbooks", bson.M{})
	if err != nil {
		s.logger.Printf("Database query error: %v", err)
		return nil, fmt.Errorf("database query error: %w", err)
	}
	defer cursor.Close(ctx)

	// Parse query results
	var exchangeRates []DBExchangeRate
	if err := cursor.All(ctx, &exchangeRates); err != nil {
		s.logger.Printf("Error parsing exchange rate data: %v", err)
		return nil, fmt.Errorf("error parsing exchange rate data: %w", err)
	}

	return exchangeRates, nil
}
