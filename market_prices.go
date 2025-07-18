package adminapiservice

import (
	marketprices "admin-panel/gen/market_prices"
	"admin-panel/redis_database"
	"context"
	"encoding/json"
	"log"
)

const (
	// Redis key for storing market prices
	MarketPriceRedisKey = "market:prices:all"
)

// marketPrices service example implementation.
// The example methods log the requests and return zero values.
type marketPricessrvc struct {
	logger *log.Logger
}

// NewMarketPrices returns the marketPrices service implementation.
func NewMarketPrices(logger *log.Logger) marketprices.Service {
	return &marketPricessrvc{logger}
}

// Get current USDT prices for all available tokens
func (s *marketPricessrvc) GetAllPrices(ctx context.Context) (res *marketprices.PriceResponse, err error) {
	s.logger.Print("marketPrices.getAllPrices")

	// Initialize response structure
	res = &marketprices.PriceResponse{
		Code:    0,
		Message: "Success",
		Result:  []*marketprices.TokenPrice{},
	}

	// Get Redis instance
	redisDb := redis_database.GetRedisDbIns("main")
	if redisDb == nil {
		res.Code = 1001
		res.Message = "Failed to get Redis instance"
		return res, nil
	}

	// Get price data from Redis
	jsonData, err := redisDb.GetString(MarketPriceRedisKey)
	if err != nil {
		s.logger.Printf("Error retrieving market prices from Redis: %v", err)
		res.Code = 1002
		res.Message = "Failed to retrieve market prices"
		return res, nil
	}

	// If no data in Redis
	if jsonData == "" {
		s.logger.Print("No market prices found in Redis")
		res.Code = 1003
		res.Message = "No market price data available"
		return res, nil
	}

	// Parse the JSON data
	var tokenPrices []struct {
		Symbol string  `json:"symbol"`
		Price  float64 `json:"price"`
	}

	if err := json.Unmarshal([]byte(jsonData), &tokenPrices); err != nil {
		s.logger.Printf("Error parsing market price data: %v", err)
		res.Code = 1004
		res.Message = "Failed to parse market price data"
		return res, nil
	}

	// Convert to response format
	for _, tp := range tokenPrices {
		res.Result = append(res.Result, &marketprices.TokenPrice{
			Symbol: tp.Symbol,
			Price:  tp.Price,
		})
	}

	s.logger.Printf("Returning %d token prices", len(res.Result))
	return res, nil
}
