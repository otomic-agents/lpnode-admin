// apps/Hedge/Service.go
package apps_hedge

import (
	database "admin-panel/mongo_database"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HedgeService struct {
	logger *log.Logger
}

// NewHedgeService creates a new instance of HedgeService
func NewHedgeService(logger *log.Logger) *HedgeService {
	return &HedgeService{
		logger: logger,
	}
}

// ListHedges retrieves all hedge configurations
func (s *HedgeService) ListHedges(ctx context.Context) ([]*Hedge, error) {
	err, cursor := database.FindAll("main", "hedge_tasks", bson.M{})
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer cursor.Close(ctx)

	var dbHedges []DBHedge
	if err := cursor.All(ctx, &dbHedges); err != nil {
		return nil, fmt.Errorf("error reading hedges: %w", err)
	}

	// Convert to API response models
	hedges := make([]*Hedge, len(dbHedges))
	for i, dbHedge := range dbHedges {
		hedges[i] = s.toAPIModel(&dbHedge)
	}

	return hedges, nil
}

// CreateHedge creates a new hedge task
func (s *HedgeService) CreateHedge(ctx context.Context, params CreateHedgeParams) (*Hedge, error) {
	// Create new database model
	now := time.Now()

	// Convert string ID to ObjectID
	cexAccountObjID, err := primitive.ObjectIDFromHex(params.CexAccountID)
	if err != nil {
		return nil, fmt.Errorf("invalid account ID format: %w", err)
	}

	bridgeObjID, err := primitive.ObjectIDFromHex(params.BridgeID)
	if err != nil {
		return nil, fmt.Errorf("invalid bridge ID format: %w", err)
	}

	// Create unique HedgeID
	dbHedge := DBHedge{
		ID:           primitive.NewObjectID(),
		Name:         params.Name,
		CexAccountID: cexAccountObjID,
		BridgeID:     bridgeObjID,
		AmmName:      params.AmmName,
		Status:       "pending", // Initial status set to pending
		StatusDesc:   "Task created, waiting to start",
		CreatedAt:    now,
		UpdatedAt:    now,
		Version:      1,
	}

	spew.Dump(dbHedge)

	// Save hedge task to the database
	err = database.Insert("main", "hedge_tasks", dbHedge)
	if err != nil {
		return nil, fmt.Errorf("failed to create hedge task: %w", err)
	}

	// Convert to API model and return
	return s.toAPIModel(&dbHedge), nil
}

// CloseHedge closes a hedge task by setting its status to "closed"
func (s *HedgeService) CloseHedge(ctx context.Context, id string) error {
	// Convert string ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}

	// First get the existing hedge task to verify it exists
	var dbHedge DBHedge
	err = database.FindOne("main", "hedge_tasks", bson.M{"_id": objID}, &dbHedge)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return ErrHedgeNotFound
		}
		return fmt.Errorf("database error: %w", err)
	}

	// Check if the task is already closed
	if dbHedge.Status == "closed" {
		return fmt.Errorf("hedge task is already closed")
	}

	// Prepare update fields
	now := time.Now()
	updateFields := bson.M{
		"status":      "completed",
		"status_desc": "Hedge task closed by user",
		"updated_at":  now,
		"version":     dbHedge.Version + 1,
	}

	// If the task was active, record the end time in the active period
	if dbHedge.Status == "active" {
		if dbHedge.ActivePeriod == nil {
			// If active period wasn't set, create it with start time as the task's created time
			updateFields["active_period"] = bson.M{
				"start": dbHedge.CreatedAt,
				"end":   now,
			}
		} else {
			// Update just the end time
			updateFields["active_period.end"] = now
		}
	}

	// Log operation
	s.logger.Printf("Closing hedge task with ID: %s", id)

	// Execute update
	err = database.Update("main", "hedge_tasks", bson.M{"_id": objID}, bson.M{"$set": updateFields})
	if err != nil {
		return fmt.Errorf("failed to close hedge task: %w", err)
	}

	return nil
}

// Helper methods

// toAPIModel converts a database model to an API response model
func (s *HedgeService) toAPIModel(dbHedge *DBHedge) *Hedge {
	if dbHedge == nil {
		return nil
	}

	hedge := &Hedge{
		ID:         dbHedge.ID.Hex(),
		Name:       dbHedge.Name,
		Exchange:   dbHedge.Exchange,
		APIKey:     maskAPIKey(dbHedge.APIKey),
		Status:     dbHedge.Status,
		StatusDesc: dbHedge.StatusDesc,
		Pair:       dbHedge.Pair,
		ChainPair:  dbHedge.ChainPair,
		AmmName:    dbHedge.AmmName,
		CreatedBy:  dbHedge.CreatedBy,
		CreatedAt:  dbHedge.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  dbHedge.UpdatedAt.Format(time.RFC3339),
		Version:    dbHedge.Version,
	}

	if dbHedge.BridgeID != primitive.NilObjectID {
		hedge.BridgeID = dbHedge.BridgeID.Hex()
	}

	if dbHedge.CexAccountID != primitive.NilObjectID {
		hedge.CexAccountID = dbHedge.CexAccountID.Hex()
	}

	if dbHedge.ActivePeriod != nil {
		hedge.ActivePeriod = &struct {
			Start string `json:"start"`
			End   string `json:"end"`
		}{
			Start: dbHedge.ActivePeriod.Start.Format(time.RFC3339),
			End:   dbHedge.ActivePeriod.End.Format(time.RFC3339),
		}
	}

	// Add exposure data if available
	if dbHedge.Exposure != nil {
		hedge.Exposure = dbHedge.Exposure
	}

	return hedge
}

// GetHedgesByBridgeID retrieves all hedge tasks related to a bridge ID
func (s *HedgeService) GetHedgesByBridgeID(ctx context.Context, bridgeID string) ([]*Hedge, error) {
	// Convert string ID to ObjectID
	bridgeObjID, err := primitive.ObjectIDFromHex(bridgeID)
	if err != nil {
		return nil, fmt.Errorf("invalid bridge ID format: %w", err)
	}

	// Query the database
	err, cursor := database.FindAll("main", "hedge_tasks", bson.M{"bridge_id": bridgeObjID})
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer cursor.Close(ctx)

	var dbHedges []DBHedge
	if err := cursor.All(ctx, &dbHedges); err != nil {
		return nil, fmt.Errorf("error reading hedges: %w", err)
	}

	// If no records found
	if len(dbHedges) == 0 {
		return []*Hedge{}, nil
	}

	// Convert to API response models
	hedges := make([]*Hedge, len(dbHedges))
	for i, dbHedge := range dbHedges {
		hedges[i] = s.toAPIModel(&dbHedge)
	}

	return hedges, nil
}

// UpdateHedge updates an existing hedge task
func (s *HedgeService) UpdateHedge(ctx context.Context, params UpdateHedgeParams) (*Hedge, error) {
	// Convert string ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(params.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	// First get the existing hedge task
	var dbHedge DBHedge
	err = database.FindOne("main", "hedge_tasks", bson.M{"_id": objID}, &dbHedge)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, ErrHedgeNotFound
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Prepare update fields
	now := time.Now()
	updateFields := bson.M{
		"status":     params.Status,
		"updated_at": now,
		"version":    dbHedge.Version + 1,
	}

	// Only update non-empty provided fields
	if params.Name != "" {
		updateFields["name"] = params.Name
	}

	if params.AmmName != "" {
		updateFields["amm_name"] = params.AmmName
	}

	// If CexAccountID is provided, convert and update
	if params.CexAccountID != "" {
		cexAccountObjID, err := primitive.ObjectIDFromHex(params.CexAccountID)
		if err != nil {
			return nil, fmt.Errorf("invalid account ID format: %w", err)
		}
		updateFields["cex_account_id"] = cexAccountObjID
	}

	// If BridgeID is provided, convert and update
	if params.BridgeID != "" {
		bridgeObjID, err := primitive.ObjectIDFromHex(params.BridgeID)
		if err != nil {
			return nil, fmt.Errorf("invalid bridge ID format: %w", err)
		}
		updateFields["bridge_id"] = bridgeObjID
	}

	// Execute update
	err = database.Update("main", "hedge_tasks", bson.M{"_id": objID}, bson.M{"$set": updateFields})
	if err != nil {
		return nil, fmt.Errorf("failed to update hedge task: %w", err)
	}

	// Get the updated hedge task
	err = database.FindOne("main", "hedge_tasks", bson.M{"_id": objID}, &dbHedge)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve updated hedge task: %w", err)
	}

	// Convert to API response model and return
	return s.toAPIModel(&dbHedge), nil
}

// GetHedge retrieves a single hedge task by ID
func (s *HedgeService) GetHedge(ctx context.Context, id string) (*Hedge, error) {
	// Convert string ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	// Query the database
	var dbHedge DBHedge
	err = database.FindOne("main", "hedge_tasks", bson.M{"_id": objID}, &dbHedge)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, ErrHedgeNotFound
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Convert to API response model and return
	return s.toAPIModel(&dbHedge), nil
}

// UpdateHedgeData updates hedge task data, including initial snapshot and risk configuration
func (s *HedgeService) UpdateHedgeData(ctx context.Context, params UpdateHedgeDataParams) (*Hedge, error) {
	// Convert string ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(params.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid hedge ID format: %w", err)
	}

	// Convert bridge ID to ObjectID
	bridgeObjID, err := primitive.ObjectIDFromHex(params.BridgeID)
	if err != nil {
		return nil, fmt.Errorf("invalid bridge ID format: %w", err)
	}

	// Convert CEX account ID to ObjectID
	cexAccountObjID, err := primitive.ObjectIDFromHex(params.CexAccountID)
	if err != nil {
		return nil, fmt.Errorf("invalid CEX account ID format: %w", err)
	}

	// First get the existing hedge task
	var dbHedge DBHedge
	err = database.FindOne("main", "hedge_tasks", bson.M{"_id": objID}, &dbHedge)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, ErrHedgeNotFound
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Prepare update data
	now := time.Now()

	// Create initial snapshot structure
	initialSnapshotData := &struct {
		Cex       map[string]float64 `bson:"cex"`
		Dex       map[string]float64 `bson:"dex"`
		Timestamp time.Time          `bson:"timestamp"`
	}{
		Cex:       params.InitialSnapshot["cex"],
		Dex:       params.InitialSnapshot["dex"],
		Timestamp: now,
	}

	// Create risk configuration structure
	riskConfigData := &struct {
		MaxAssetExposure map[string]float64 `bson:"max_asset_exposure"`
		MinHedgeAmount   map[string]float64 `bson:"min_hedge_amount"`
		HedgeMode        string             `bson:"hedge_mode"`
	}{
		MaxAssetExposure: params.RiskConfig.MaxAssetExposure,
		MinHedgeAmount:   params.RiskConfig.MinHedgeAmount,
		HedgeMode:        params.RiskConfig.HedgeMode,
	}

	// Prepare chain pair data
	chainPair := []string{params.SourceWallet.Token, params.DestinationWallet.Token}
	chainAddress := []string{params.SourceWallet.TokenAddress, params.DestinationWallet.TokenAddress}
	chainIdList := []int64{params.SourceWallet.ChainId, params.DestinationWallet.ChainId}
	walletList := []string{params.SourceWallet.Address, params.DestinationWallet.Address}
	// Prepare update fields
	updateFields := bson.M{
		"name":             params.Name,
		"bridge_id":        bridgeObjID,
		"cex_account_id":   cexAccountObjID,
		"amm_name":         params.AmmName,
		"initial_snapshot": initialSnapshotData,
		"risk_config":      riskConfigData,
		"chain_pair":       chainPair,
		"chain_address":    chainAddress,
		"chain_id_list":    chainIdList,
		"wallet_list":      walletList,
		"status":           "active", // Set status to active as data is saved
		"status_desc":      "Hedge task initialized with data",
		"updated_at":       now,
		"version":          dbHedge.Version + 1,
	}

	// Log operation
	s.logger.Printf("Updating hedge data for ID: %s with data: %+v", params.ID, updateFields)

	// Execute update
	err = database.Update("main", "hedge_tasks", bson.M{"_id": objID}, bson.M{"$set": updateFields})
	if err != nil {
		return nil, fmt.Errorf("failed to update hedge data: %w", err)
	}

	// Get the updated hedge task
	err = database.FindOne("main", "hedge_tasks", bson.M{"_id": objID}, &dbHedge)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve updated hedge task: %w", err)
	}

	// Convert to API response model and return
	return s.toAPIModel(&dbHedge), nil
}

// maskAPIKey masks the API key for security
func maskAPIKey(apiKey string) string {
	if len(apiKey) < 8 {
		return "****"
	}
	return apiKey[:4] + "..." + apiKey[len(apiKey)-4:]
}
