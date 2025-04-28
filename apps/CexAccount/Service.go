package apps_cex_account

import (
	database "admin-panel/mongo_database"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CexAccountService struct {
	logger *log.Logger
}

// NewCexAccountService creates a new instance of CexAccountService
func NewCexAccountService(logger *log.Logger) *CexAccountService {
	return &CexAccountService{
		logger: logger,
	}
}

// CreateAccount creates a new CEX account
func (s *CexAccountService) CreateAccount(ctx context.Context, payload *CreateCexAccountPayload) (*CexAccount, error) {
	// Validate input
	if payload.Name == "" {
		return nil, fmt.Errorf("account name is required")
	}
	if payload.Exchange == "" {
		return nil, fmt.Errorf("exchange name is required")
	}
	if payload.APIKey == "" {
		return nil, fmt.Errorf("API key is required")
	}
	if payload.APISecret == "" {
		return nil, fmt.Errorf("API secret is required")
	}

	// Check if account with same name already exists
	var existingAccount DBCexAccount
	err := database.FindOne("main", "cex_accounts", bson.M{"name": payload.Name}, &existingAccount)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if existingAccount.ID != primitive.NilObjectID {
		return nil, fmt.Errorf("an account with name '%s' already exists", payload.Name)
	}

	// Create new account document
	now := primitive.NewDateTimeFromTime(time.Now())
	dbAccount := DBCexAccount{
		Name:       payload.Name,
		Exchange:   payload.Exchange,
		APIKey:     payload.APIKey,
		APISecret:  payload.APISecret,
		Passphrase: payload.Passphrase,
		Status:     "active", // Default status
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// Insert into database
	err = database.Insert("main", "cex_accounts", dbAccount)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	// Retrieve the inserted document to get the ID
	var insertedAccount DBCexAccount
	err = database.FindOne("main", "cex_accounts", bson.M{
		"name":      payload.Name,
		"exchange":  payload.Exchange,
		"api_key":   payload.APIKey,
		"createdAt": now,
	}, &insertedAccount)
	if err != nil {
		s.logger.Printf("Warning: Account created but failed to retrieve: %v", err)
		// Return a partial response without ID
		return &CexAccount{
			Name:      payload.Name,
			Exchange:  payload.Exchange,
			APIKey:    maskAPIKey(payload.APIKey),
			Status:    "active",
			CreatedAt: now.Time().Format(time.RFC3339),
			UpdatedAt: now.Time().Format(time.RFC3339),
		}, nil
	}

	// Convert to API response model
	return s.toAPIModel(&insertedAccount), nil
}

// GetAccount retrieves a CEX account by ID
func (s *CexAccountService) GetAccount(ctx context.Context, id string) (*CexAccount, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid account ID format: %w", err)
	}

	var dbAccount DBCexAccount
	err = database.FindOne("main", "cex_accounts", bson.M{"_id": objectID}, &dbAccount)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	if dbAccount.ID == primitive.NilObjectID {
		return nil, fmt.Errorf("account not found")
	}

	return s.toAPIModel(&dbAccount), nil
}

// ListAccounts retrieves all CEX accounts
func (s *CexAccountService) ListAccounts(ctx context.Context) ([]*CexAccount, error) {
	err, cursor := database.FindAll("main", "cex_accounts", bson.M{})
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer cursor.Close(ctx)

	var dbAccounts []DBCexAccount
	if err := cursor.All(ctx, &dbAccounts); err != nil {
		return nil, fmt.Errorf("error reading accounts: %w", err)
	}

	// Convert to API response models
	accounts := make([]*CexAccount, len(dbAccounts))
	for i, dbAccount := range dbAccounts {
		accounts[i] = s.toAPIModel(&dbAccount)
	}

	return accounts, nil
}

// UpdateAccount updates an existing CEX account
func (s *CexAccountService) UpdateAccount(ctx context.Context, id string, payload *CreateCexAccountPayload) (*CexAccount, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid account ID format: %w", err)
	}

	// Prepare update document
	update := bson.M{
		"$set": bson.M{
			"updated_at": primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	// Only update fields that are provided
	if payload.Name != "" {
		update["$set"].(bson.M)["name"] = payload.Name
	}
	if payload.Exchange != "" {
		update["$set"].(bson.M)["exchange"] = payload.Exchange
	}
	if payload.APIKey != "" {
		update["$set"].(bson.M)["api_key"] = payload.APIKey
	}
	if payload.APISecret != "" {
		update["$set"].(bson.M)["api_secret"] = payload.APISecret
	}
	if payload.Passphrase != "" {
		update["$set"].(bson.M)["passphrase"] = payload.Passphrase
	}

	// Update the document
	err = database.Update("main", "cex_accounts", bson.M{"_id": objectID}, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update account: %w", err)
	}

	// Retrieve the updated document
	var updatedAccount DBCexAccount
	err = database.FindOne("main", "cex_accounts", bson.M{"_id": objectID}, &updatedAccount)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	if updatedAccount.ID == primitive.NilObjectID {
		return nil, fmt.Errorf("account not found after update")
	}

	return s.toAPIModel(&updatedAccount), nil
}

// DeleteAccount deletes a CEX account
func (s *CexAccountService) DeleteAccount(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid account ID format: %w", err)
	}

	deletedCount, err := database.DeleteOne("main", "cex_accounts", bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	if deletedCount == 0 {
		return fmt.Errorf("account not found")
	}

	return nil
}

// GetTokenBalance gets the specific token balance for the specified account
func (s *CexAccountService) GetTokenBalance(ctx context.Context, accountID string, symbol string) (*TokenBalance, error) {
	s.logger.Printf("Querying account %s for %s token balance", accountID, symbol)

	// Query the database for balance information
	var walletBalance CexWalletBalance
	err := database.FindOne("main", "cex_wallet_balances", bson.M{"accountId": accountID}, &walletBalance)
	if err != nil {
		return nil, fmt.Errorf("database query error: %w", err)
	}

	// Check if account was found
	if walletBalance.ID == primitive.NilObjectID {
		// If no balance information is found, return default values
		s.logger.Printf("No balance information found for account %s, returning default values", accountID)
		return &TokenBalance{
			Asset:  symbol,
			Free:   "0",
			Locked: "0",
			Total:  "0",
		}, nil
	}

	// Search for the specified token in the balance array
	for _, balance := range walletBalance.ExchangeResult.Balances {
		if balance.Asset == symbol {
			// Found matching token, return balance information
			return &TokenBalance{
				Asset:  balance.Asset,
				Free:   balance.Free,
				Locked: balance.Locked,
				Total:  balance.Total,
			}, nil
		}
	}

	// Specified token not found, return zero balance
	s.logger.Printf("No balance information found for token %s in account %s, returning zero balance", symbol, accountID)
	return &TokenBalance{
		Asset:  symbol,
		Free:   "0",
		Locked: "0",
		Total:  "0",
	}, nil
}

// GetAllTokenBalances gets all token balances for the specified account
func (s *CexAccountService) GetAllTokenBalances(ctx context.Context, accountID string) ([]*TokenBalance, error) {
	s.logger.Printf("Querying all token balances for account %s", accountID)

	// Query the database for balance information
	var walletBalance struct {
		ID             primitive.ObjectID `bson:"_id"`
		AccountID      string             `bson:"accountId"`
		Exchange       string             `bson:"exchange"`
		ExchangeResult struct {
			Info struct {
				Timestamp string `bson:"timestamp"`
				Exchange  string `bson:"exchange"`
				Account   string `bson:"account"`
			} `bson:"info"`
			Balances []struct {
				Asset  string `bson:"asset"`
				Free   string `bson:"free"`
				Locked string `bson:"locked"`
				Total  string `bson:"total"`
			} `bson:"balances"`
		} `bson:"exchangeResult"`
	}

	err := database.FindOne("main", "cex_wallet_balances", bson.M{"accountId": accountID}, &walletBalance)
	if err != nil {
		return nil, fmt.Errorf("database query error: %w", err)
	}

	// Check if account was found
	if walletBalance.ID == primitive.NilObjectID {
		// If no balance information is found, return empty list
		s.logger.Printf("No balance information found for account %s, returning empty list", accountID)
		return []*TokenBalance{}, nil
	}

	// Convert balance array to target format
	tokenBalances := make([]*TokenBalance, len(walletBalance.ExchangeResult.Balances))
	for i, balance := range walletBalance.ExchangeResult.Balances {
		tokenBalances[i] = &TokenBalance{
			Asset:  balance.Asset,
			Free:   balance.Free,
			Locked: balance.Locked,
			Total:  balance.Total,
		}
	}

	return tokenBalances, nil
}

// TestConnection tests the connection to the exchange using the account credentials
func (s *CexAccountService) TestConnection(ctx context.Context, id string) (bool, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, fmt.Errorf("invalid account ID format: %w", err)
	}

	// Get the account with secrets from DB
	var dbAccount DBCexAccount
	err = database.FindOne("main", "cex_accounts", bson.M{"_id": objectID}, &dbAccount)
	if err != nil {
		return false, fmt.Errorf("database error: %w", err)
	}

	if dbAccount.ID == primitive.NilObjectID {
		return false, fmt.Errorf("account not found")
	}

	// Here you would implement the actual exchange connection test
	// This is a placeholder - you'll need to implement exchange-specific logic
	var connectionSuccess bool
	switch dbAccount.Exchange {
	case "binance":
		connectionSuccess, err = testBinanceConnection(dbAccount.APIKey, dbAccount.APISecret, dbAccount.Passphrase)
	case "coinbase":
		connectionSuccess, err = testCoinbaseConnection(dbAccount.APIKey, dbAccount.APISecret, dbAccount.Passphrase)
	// Add more exchanges as needed
	default:
		return false, fmt.Errorf("unsupported exchange: %s", dbAccount.Exchange)
	}

	if err != nil {
		// Update account status to error
		updateErr := database.Update(
			"main",
			"cex_accounts",
			bson.M{"_id": objectID},
			bson.M{
				"$set": bson.M{
					"status":     "error",
					"updated_at": primitive.NewDateTimeFromTime(time.Now()),
				},
			},
		)
		if updateErr != nil {
			// Log the error but don't return it
			s.logger.Printf("Failed to update account status: %v", updateErr)
		}
		return false, fmt.Errorf("connection test failed: %w", err)
	}

	// Update account status to active
	updateErr := database.Update(
		"main",
		"cex_accounts",
		bson.M{"_id": objectID},
		bson.M{
			"$set": bson.M{
				"status":     "active",
				"updated_at": primitive.NewDateTimeFromTime(time.Now()),
			},
		},
	)
	if updateErr != nil {
		// Log the error but don't return it
		s.logger.Printf("Failed to update account status: %v", updateErr)
	}

	return connectionSuccess, nil
}

// Helper methods

// toAPIModel converts a database model to an API response model
func (s *CexAccountService) toAPIModel(dbAccount *DBCexAccount) *CexAccount {
	if dbAccount == nil {
		return nil
	}

	return &CexAccount{
		ID:        dbAccount.ID.Hex(),
		Name:      dbAccount.Name,
		Exchange:  dbAccount.Exchange,
		APIKey:    maskAPIKey(dbAccount.APIKey),
		Status:    dbAccount.Status,
		CreatedAt: dbAccount.CreatedAt.Time().Format(time.RFC3339),
		UpdatedAt: dbAccount.UpdatedAt.Time().Format(time.RFC3339),
	}
}

// maskAPIKey masks the API key for security
func maskAPIKey(apiKey string) string {
	if len(apiKey) < 8 {
		return "****"
	}
	return apiKey[:4] + "..." + apiKey[len(apiKey)-4:]
}

// Exchange-specific connection test implementations
// These would need to be implemented based on the exchange APIs you're using

func testBinanceConnection(apiKey, apiSecret, passphrase string) (bool, error) {
	// Implement Binance API connection test
	// This is a placeholder - you'll need to implement actual API calls
	return true, nil
}

func testCoinbaseConnection(apiKey, apiSecret, passphrase string) (bool, error) {
	// Implement Coinbase API connection test
	// This is a placeholder - you'll need to implement actual API calls
	return true, nil
}
