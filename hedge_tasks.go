package adminapiservice

import (
	apps_hedge "admin-panel/apps/Hedge"
	hedgetasks "admin-panel/gen/hedge_tasks"
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/aws/smithy-go/ptr"
)

// hedgeTasks service example implementation.
// The example methods log the requests and return zero values.
type hedgeTaskssrvc struct {
	logger *log.Logger
}

// NewHedgeTasks returns the hedgeTasks service implementation.
func NewHedgeTasks(logger *log.Logger) hedgetasks.Service {
	return &hedgeTaskssrvc{logger}
}

// List all hedge tasks
func (s *hedgeTaskssrvc) ListTasks(ctx context.Context) (res *hedgetasks.ListTasksResult, err error) {
	s.logger.Print("hedgeTasks.listTasks")

	// Initialize result
	res = &hedgetasks.ListTasksResult{
		Code:    0,
		Message: "Success",
		Result:  []*hedgetasks.HTHedgeTask{},
	}
	hs := apps_hedge.NewHedgeService(s.logger)
	// Call HedgeService's ListHedges method to get hedge tasks
	hedges, err := hs.ListHedges(ctx)
	if err != nil {
		s.logger.Printf("Error listing hedges: %v", err)
		res.Code = 500
		res.Message = "Failed to retrieve hedge tasks: " + err.Error()
		return res, nil
	}

	// Convert Hedge objects to HTHedgeTask objects
	for _, hedge := range hedges {
		task := &hedgetasks.HTHedgeTask{
			ID:        hedge.ID,
			Name:      hedge.Name,
			Status:    hedge.Status,
			CreatedAt: hedge.CreatedAt,
			UpdatedAt: hedge.UpdatedAt,
			CreatedBy: ptr.String(hedge.CreatedBy),
		}

		// Set associated account ID (using CexAccountID)
		if hedge.CexAccountID != "" {
			task.AccountID = hedge.CexAccountID
		}

		// Set associated bridge ID
		if hedge.BridgeID != "" {
			task.BridgeID = hedge.BridgeID
		}

		res.Result = append(res.Result, task)
	}

	return res, nil
}

// Delete specified hedge task
func (s *hedgeTaskssrvc) DeleteTask(ctx context.Context, p *hedgetasks.DeleteTaskPayload) (res *hedgetasks.DeleteTaskResult, err error) {
	res = &hedgetasks.DeleteTaskResult{}
	s.logger.Print("hedgeTasks.deleteTask")
	return
}

// CreateTask handles requests to create or update hedge tasks
func (s *hedgeTaskssrvc) CreateTask(ctx context.Context, p *hedgetasks.HTCreateTaskPayload) (res *hedgetasks.CreateTaskResult, err error) {
	s.logger.Print("hedgeTasks.createTask")

	// Initialize result
	res = &hedgetasks.CreateTaskResult{
		Code:    0,
		Message: "Success",
	}

	// Create Hedge service instance
	hs := apps_hedge.NewHedgeService(s.logger)

	// Check if task ID is provided
	if p.ID != nil && *p.ID != "" {
		// If ID is provided, try to update existing task
		return s.updateExistingTask(ctx, hs, p)
	}

	// Create new task
	return s.createNewTask(ctx, hs, p)
}

// updateExistingTask updates an existing hedge task
func (s *hedgeTaskssrvc) updateExistingTask(ctx context.Context, hs *apps_hedge.HedgeService, p *hedgetasks.HTCreateTaskPayload) (*hedgetasks.CreateTaskResult, error) {
	s.logger.Printf("Updating existing hedge task with ID: %s", *p.ID)

	// Prepare parameters for updating hedge task
	updateParams := apps_hedge.UpdateHedgeParams{
		ID:           *p.ID,
		Name:         p.Name,
		CexAccountID: p.AccountID,
		BridgeID:     p.BridgeID,
		AmmName:      p.AmmName,
		Status:       "pending",
		// Other parameters can be added as needed
	}

	// Call HedgeService's UpdateHedge method to update the hedge task
	hedge, err := hs.UpdateHedge(ctx, updateParams)
	if err != nil {
		s.logger.Printf("Error updating hedge task: %v", err)
		return &hedgetasks.CreateTaskResult{
			Code:    1001,
			Message: "Failed to update hedge task: " + err.Error(),
		}, nil
	}

	// Return the updated task
	return s.convertHedgeToResult(hedge), nil
}

// checkExistingActiveTask checks if there's an active hedge task, returns error result or nil
func (s *hedgeTaskssrvc) checkExistingActiveTask(ctx context.Context, hs *apps_hedge.HedgeService, bridgeID string) (*hedgetasks.CreateTaskResult, bool) {
	// Get all hedge tasks associated with the bridgeID
	hedgeTasks, err := hs.GetHedgesByBridgeID(ctx, bridgeID)
	if err != nil {
		s.logger.Printf("Error fetching hedge tasks for bridge ID %s: %v", bridgeID, err)
		return &hedgetasks.CreateTaskResult{
			Code:    1002,
			Message: "Failed to check existing hedge tasks: " + err.Error(),
		}, false
	}

	// Check if there are any tasks in pending or active status
	for _, task := range hedgeTasks {
		if task.Status == "pending" || task.Status == "active" {
			s.logger.Printf("Found existing %s hedge task (ID: %s) for bridge ID: %s",
				task.Status, task.ID, bridgeID)
			return &hedgetasks.CreateTaskResult{
				Code:    1003,
				Message: fmt.Sprintf("A %s hedge task already exists for this bridge. Please use a different bridge or update the existing task.", task.Status),
			}, false
		}
	}

	// No active task found, can proceed with creation
	return nil, true
}

// createNewTask creates a new hedge task
func (s *hedgeTaskssrvc) createNewTask(ctx context.Context, hs *apps_hedge.HedgeService, p *hedgetasks.HTCreateTaskPayload) (*hedgetasks.CreateTaskResult, error) {
	// First check if there's an active task
	if errorResult, canProceed := s.checkExistingActiveTask(ctx, hs, p.BridgeID); !canProceed {
		return errorResult, nil
	}

	// Prepare parameters for creating hedge task
	createParams := apps_hedge.CreateHedgeParams{
		Name:         p.Name,
		CexAccountID: p.AccountID,
		BridgeID:     p.BridgeID,
		AmmName:      p.AmmName,
		// Other parameters can be added as needed
	}

	// Call HedgeService's CreateHedge method to create the hedge task
	hedge, err := hs.CreateHedge(ctx, createParams)
	if err != nil {
		s.logger.Printf("Error creating hedge task: %v", err)
		return &hedgetasks.CreateTaskResult{
			Code:    1004,
			Message: "Failed to create hedge task: " + err.Error(),
		}, nil
	}

	// Return the created task
	return s.convertHedgeToResult(hedge), nil
}

// Get a single hedge task
func (s *hedgeTaskssrvc) GetTask(ctx context.Context, p *hedgetasks.GetTaskPayload) (res *hedgetasks.GetTaskResult, err error) {
	s.logger.Printf("hedgeTasks.getTask with ID: %s", p.TaskID)

	// Initialize result
	res = &hedgetasks.GetTaskResult{
		Code:    0,
		Message: "Success",
	}

	// Create Hedge service instance
	hs := apps_hedge.NewHedgeService(s.logger)

	// Call HedgeService's GetHedge method to get the hedge task
	hedge, err := hs.GetHedge(ctx, p.TaskID)
	if err != nil {
		s.logger.Printf("Error getting hedge task: %v", err)
		res.Code = 404
		res.Message = "Failed to retrieve hedge task: " + err.Error()
		return res, nil
	}

	// If task not found
	if hedge == nil {
		res.Code = 404
		res.Message = "Hedge task not found"
		return res, nil
	}

	// Convert Hedge object to HTHedgeTask object
	res.Result = &hedgetasks.HTHedgeTask{
		AmmName:   ptr.String(hedge.AmmName),
		ID:        hedge.ID,
		Name:      hedge.Name,
		Status:    hedge.Status,
		AccountID: hedge.CexAccountID,
		BridgeID:  hedge.BridgeID,
		CreatedAt: hedge.CreatedAt,
		UpdatedAt: hedge.UpdatedAt,
		CreatedBy: ptr.String(hedge.CreatedBy),
	}

	return res, nil
}

// SaveHedgeData saves hedge data to update an existing hedge configuration
func (s *hedgeTaskssrvc) SaveHedgeData(ctx context.Context, p *hedgetasks.HTSaveHedgeDataPayload) (res *hedgetasks.SaveHedgeDataResult, err error) {
	s.logger.Print("hedgeTasks.saveHedgeData")

	// Initialize result
	res = &hedgetasks.SaveHedgeDataResult{
		Code:    0,
		Message: "Success",
	}

	// Create Hedge service instance
	hs := apps_hedge.NewHedgeService(s.logger)

	// Prepare initial snapshot data from the payload
	initialSnapshot := make(map[string]map[string]float64)

	// Process source token data
	initialSnapshot["cex"] = make(map[string]float64)
	initialSnapshot["dex"] = make(map[string]float64)

	// Convert string balances to float64
	sourceCexBalance, _ := strconv.ParseFloat(ptr.ToString(p.InitialBalances.Source.Cex), 64)
	sourceDexBalance, _ := strconv.ParseFloat(ptr.ToString(p.InitialBalances.Source.Dex), 64)
	destCexBalance, _ := strconv.ParseFloat(ptr.ToString(p.InitialBalances.Destination.Cex), 64)
	destDexBalance, _ := strconv.ParseFloat(ptr.ToString(p.InitialBalances.Destination.Dex), 64)

	// Add source token balances
	initialSnapshot["cex"][ptr.ToString(p.InitialBalances.Source.Token)] = sourceCexBalance
	initialSnapshot["dex"][ptr.ToString(p.InitialBalances.Source.Token)] = sourceDexBalance

	// Add destination token balances
	initialSnapshot["cex"][ptr.ToString(p.InitialBalances.Destination.Token)] = destCexBalance
	initialSnapshot["dex"][ptr.ToString(p.InitialBalances.Destination.Token)] = destDexBalance

	// Prepare risk config
	riskConfig := apps_hedge.RiskConfig{
		HedgeMode:        ptr.ToString(p.RiskConfig.HedgeMode),
		MaxAssetExposure: p.RiskConfig.MaxAssetExposure,
		MinHedgeAmount:   p.RiskConfig.MinHedgeAmount,
	}

	// Prepare update parameters
	updateParams := apps_hedge.UpdateHedgeDataParams{
		ID:              p.ID,
		Name:            p.Name,
		BridgeID:        p.BridgeID,
		CexAccountID:    p.CexAccountID,
		AmmName:         p.AmmName,
		InitialSnapshot: initialSnapshot,
		RiskConfig:      riskConfig,
		SourceWallet: apps_hedge.WalletInfo{
			Address: ptr.ToString(p.InitialBalances.Source.Wallet),
			Name:    ptr.ToString(p.InitialBalances.Source.WalletName),
			Token:   ptr.ToString(p.InitialBalances.Source.Token),
		},
		DestinationWallet: apps_hedge.WalletInfo{
			Address: ptr.ToString(p.InitialBalances.Destination.Wallet),
			Name:    ptr.ToString(p.InitialBalances.Destination.WalletName),
			Token:   ptr.ToString(p.InitialBalances.Destination.Token),
		},
	}

	// Call HedgeService's UpdateHedgeData method to update the hedge data
	hedge, err := hs.UpdateHedgeData(ctx, updateParams)
	if err != nil {
		s.logger.Printf("Error updating hedge data: %v", err)
		return &hedgetasks.SaveHedgeDataResult{
			Code:    1001,
			Message: "Failed to update hedge data: " + err.Error(),
		}, nil
	}

	// Return the updated hedge information
	res.Result = &hedgetasks.HTSaveHedgeDataResult{
		ID:     hedge.ID,
		Name:   hedge.Name,
		Status: hedge.Status,
	}

	return res, nil
}
// Close a hedge task
func (s *hedgeTaskssrvc) CloseTask(ctx context.Context, p *hedgetasks.CloseTaskPayload) (res *hedgetasks.CloseTaskResult, err error) {
	s.logger.Printf("hedgeTasks.closeTask with ID: %s", p.TaskID)

	// Initialize result
	res = &hedgetasks.CloseTaskResult{
		Code:    0,
		Message: "Success",
	}

	// Create Hedge service instance
	hs := apps_hedge.NewHedgeService(s.logger)

	// Call HedgeService's CloseHedge method to close the hedge task
	err = hs.CloseHedge(ctx, p.TaskID)
	if err != nil {
		s.logger.Printf("Error closing hedge task: %v", err)
		res.Code = 1001
		res.Message = "Failed to close hedge task: " + err.Error()
		return res, nil
	}

	res.Message = "Hedge task closed successfully"
	return res, nil
}

// convertHedgeToResult converts a Hedge object to API response result
func (s *hedgeTaskssrvc) convertHedgeToResult(hedge *apps_hedge.Hedge) *hedgetasks.CreateTaskResult {
	return &hedgetasks.CreateTaskResult{
		Code:    0,
		Message: "Success",
		Result: &hedgetasks.HTHedgeTask{
			ID:        hedge.ID,
			Name:      hedge.Name,
			Status:    hedge.Status,
			AccountID: hedge.CexAccountID,
			BridgeID:  hedge.BridgeID,
			CreatedAt: hedge.CreatedAt,
			UpdatedAt: hedge.UpdatedAt,
			CreatedBy: ptr.String(hedge.CreatedBy),
		},
	}
}
