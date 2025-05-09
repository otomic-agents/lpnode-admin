package main

import (
	"admin-panel/redis_database"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	// Redis key for storing market prices
	MarketPriceRedisKey = "market:prices:all"

	// CoinGecko API URL for getting prices (top 250 by market cap)
	CoinGeckoURL = "https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=800&page=1&sparkline=false"

	// Refresh interval (5 minutes)
	PriceRefreshInterval = 5 * time.Minute

	// Cache expiration time (10 minutes - longer than refresh to avoid cache miss)
	PriceCacheExpiration = 10 * 60 // 10 minutes in seconds
)

// Simple price structure to store only the needed data
type TokenPrice struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
}

// CoinGecko response structure (simplified)
type CoinGeckoResponse struct {
	ID           string  `json:"id"`
	Symbol       string  `json:"symbol"`
	CurrentPrice float64 `json:"current_price"`
}

// InitMarketPriceLoop starts a goroutine to periodically fetch and update market prices
func InitMarketPriceLoop() {
	log.Println("[MarketPrice] Initializing market price update service")

	// Start a goroutine for price updates
	go func() {
		log.Println("[MarketPrice] Price update goroutine started")

		// Initial iteration count for debugging
		iteration := 1

		for {
			log.Printf("[MarketPrice] Starting update iteration #%d", iteration)
			startTime := time.Now()

			// Fetch and update prices
			err := fetchAndUpdatePrices()
			if err != nil {
				log.Printf("[MarketPrice] ERROR: Update failed: %v", err)
			} else {
				duration := time.Since(startTime)
				log.Printf("[MarketPrice] Update iteration #%d completed in %v", iteration, duration)
			}

			// Log when we'll attempt next update
			nextUpdate := time.Now().Add(PriceRefreshInterval)
			log.Printf("[MarketPrice] Next update scheduled at: %s (in %v)",
				nextUpdate.Format("2006-01-02 15:04:05"),
				PriceRefreshInterval)

			// Sleep for the specified interval
			time.Sleep(PriceRefreshInterval)

			iteration++
		}
	}()

	log.Println("[MarketPrice] Market price update service initialized successfully")
}

// fetchAndUpdatePrices fetches prices from CoinGecko and stores them in Redis
func fetchAndUpdatePrices() error {
	log.Println("[MarketPrice] Fetching market prices from CoinGecko API")

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Make request to CoinGecko
	log.Printf("[MarketPrice] Making HTTP request to: %s", CoinGeckoURL)
	requestStart := time.Now()
	resp, err := client.Get(CoinGeckoURL)
	requestDuration := time.Since(requestStart)

	if err != nil {
		return fmt.Errorf("failed to fetch prices from CoinGecko (after %v): %w", requestDuration, err)
	}
	defer resp.Body.Close()

	log.Printf("[MarketPrice] Received response from CoinGecko in %v with status: %d %s",
		requestDuration, resp.StatusCode, resp.Status)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("CoinGecko API returned error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	log.Println("[MarketPrice] Parsing CoinGecko response")
	parseStart := time.Now()
	var coinGeckoResp []CoinGeckoResponse
	if err := json.NewDecoder(resp.Body).Decode(&coinGeckoResp); err != nil {
		return fmt.Errorf("failed to decode CoinGecko response: %w", err)
	}
	log.Printf("[MarketPrice] Parsed response in %v, found %d tokens",
		time.Since(parseStart), len(coinGeckoResp))

	// Log a few samples for verification
	if len(coinGeckoResp) > 0 {
		sampleSize := 5
		if len(coinGeckoResp) < sampleSize {
			sampleSize = len(coinGeckoResp)
		}

		log.Println("[MarketPrice] Sample token prices from response:")
		for i := 0; i < sampleSize; i++ {
			token := coinGeckoResp[i]
			log.Printf("[MarketPrice]   - %s: $%.6f", strings.ToUpper(token.Symbol), token.CurrentPrice)
		}
	}

	// Convert to our simple format
	prices := make([]TokenPrice, 0, len(coinGeckoResp))
	for _, coin := range coinGeckoResp {
		prices = append(prices, TokenPrice{
			Symbol: strings.ToUpper(coin.Symbol),
			Price:  coin.CurrentPrice,
		})
	}

	// Marshal to JSON
	log.Println("[MarketPrice] Converting prices to JSON")
	pricesJSON, err := json.Marshal(prices)
	if err != nil {
		return fmt.Errorf("failed to marshal prices to JSON: %w", err)
	}

	// Store in Redis
	log.Printf("[MarketPrice] Storing %d token prices in Redis with key: %s",
		len(prices), MarketPriceRedisKey)
	redisStart := time.Now()

	// Get Redis instance
	redisDb := redis_database.GetRedisDbIns("main")
	if redisDb == nil {
		return fmt.Errorf("failed to get Redis instance")
	}

	_, err = redisDb.SetEx(MarketPriceRedisKey, string(pricesJSON), PriceCacheExpiration)
	if err != nil {
		return fmt.Errorf("failed to store prices in Redis: %w", err)
	}

	log.Printf("[MarketPrice] Successfully stored prices in Redis in %v with expiration of %d seconds",
		time.Since(redisStart), PriceCacheExpiration)

	// Verify data was stored properly by reading it back
	log.Println("[MarketPrice] Verifying data was properly stored by reading it back")
	verifyData, verifyErr := redisDb.GetString(MarketPriceRedisKey)
	if verifyErr != nil {
		log.Printf("[MarketPrice] WARNING: Could not verify stored data: %v", verifyErr)
	} else {
		var storedPrices []TokenPrice
		if jsonErr := json.Unmarshal([]byte(verifyData), &storedPrices); jsonErr != nil {
			log.Printf("[MarketPrice] WARNING: Could not parse stored data: %v", jsonErr)
		} else {
			log.Printf("[MarketPrice] Verification successful, found %d tokens in Redis", len(storedPrices))
		}
	}

	log.Printf("[MarketPrice] Market price update completed successfully, stored %d tokens", len(prices))
	return nil
}
