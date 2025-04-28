package adminapiservice

import (
	apps_cex_account "admin-panel/apps/CexAccount"
	accountcex "admin-panel/gen/account_cex"
	"context"
	"fmt"
	"log"

	"github.com/aws/smithy-go/ptr"
	"github.com/pkg/errors"
)

// accountCex service example implementation.
// The example methods log the requests and return zero values.
type accountCexsrvc struct {
	logger *log.Logger
}

// NewAccountCex returns the accountCex service implementation.
func NewAccountCex(logger *log.Logger) accountcex.Service {
	return &accountCexsrvc{logger}
}

// Get wallet balance information for the specified CEX account
func (s *accountCexsrvc) WalletInfo(ctx context.Context, p *accountcex.WalletInfoPayload) (res *accountcex.WalletInfoResult, err error) {
	res = &accountcex.WalletInfoResult{}
	s.logger.Print("accountCex.walletInfo")
	return
}

// Create a new CEX account configuration
func (s *accountCexsrvc) CreateAccount(ctx context.Context, p *accountcex.CACexAccountPayload) (res *accountcex.CreateAccountResult, err error) {
	s.logger.Print("accountCex.createAccount")

	// Create a new CEX account service
	cexAccountService := apps_cex_account.NewCexAccountService(s.logger)

	// Convert the payload to the format expected by the service
	payload := &apps_cex_account.CreateCexAccountPayload{
		Name:      p.Name,
		Exchange:  p.Exchange,
		APIKey:    p.APIKey,
		APISecret: p.APISecret,
	}

	// Add passphrase if provided
	if p.Passphrase != nil {
		payload.Passphrase = *p.Passphrase
	}

	// Call the service to create the account
	account, err := cexAccountService.CreateAccount(ctx, payload)
	if err != nil {
		return &accountcex.CreateAccountResult{
			Code:    500,
			Message: fmt.Sprintf("Failed to create CEX account: %v", err),
		}, nil
	}

	// Map the result to the expected response format
	return &accountcex.CreateAccountResult{
		Code: 200,
		Result: &accountcex.CACexAccount{
			Name:     account.Name,
			Exchange: account.Exchange,
			APIKey:   account.APIKey, // This will be masked by the service
			Status:   account.Status,
		},
		Message: "CEX account created successfully",
	}, nil
}

// Get the specific token balance for the specified CEX account
func (s *accountCexsrvc) TokenBalance(ctx context.Context, p *accountcex.TokenBalancePayload) (res *accountcex.TokenBalanceResult, err error) {
	s.logger.Printf("accountCex.tokenBalance - Account ID: %s, Token: %s", p.AccountID, p.Symbol)

	// Create CEX account service instance
	cexAccountService := apps_cex_account.NewCexAccountService(s.logger)

	// Call service to get token balance
	tokenBalance, err := cexAccountService.GetTokenBalance(ctx, p.AccountID, p.Symbol)
	if err != nil {
		// Return error response
		return &accountcex.TokenBalanceResult{
			Code:    30001,
			Message: err.Error(),
			Result:  nil,
		}, nil
	}

	// Return successful response
	return &accountcex.TokenBalanceResult{
		Code:    0,
		Message: "success",
		Result: &accountcex.CATokenBalance{
			Asset:  tokenBalance.Asset,
			Free:   tokenBalance.Free,
			Locked: tokenBalance.Locked,
			Total:  tokenBalance.Total,
		},
	}, nil
}

// List all configured CEX accounts
func (s *accountCexsrvc) ListAccounts(ctx context.Context) (res *accountcex.ListAccountsResult, err error) {
	s.logger.Print("accountCex.listAccounts")

	// Create CEX account service instance
	cexAccountService := apps_cex_account.NewCexAccountService(s.logger)

	// Call service to get all accounts
	accounts, err := cexAccountService.ListAccounts(ctx)
	if err != nil {
		return &accountcex.ListAccountsResult{
			Code:    30001,
			Message: fmt.Sprintf("Failed to list CEX accounts: %v", err),
		}, nil
	}

	// Convert service returned account list to API response format
	apiAccounts := make([]*accountcex.CACexAccount, len(accounts))
	for i, account := range accounts {
		apiAccounts[i] = &accountcex.CACexAccount{
			ID:       ptr.String(account.ID),
			Name:     account.Name,
			Exchange: account.Exchange,
			APIKey:   account.APIKey, // Sensitive information has been masked by the service layer
			Status:   account.Status,
		}
	}

	// Return successful response
	return &accountcex.ListAccountsResult{
		Code:    0,
		Result:  apiAccounts,
		Message: "CEX accounts retrieved successfully",
	}, nil
}

// Delete the specified CEX account
func (s *accountCexsrvc) DeleteAccount(ctx context.Context, p *accountcex.DeleteAccountPayload) (res *accountcex.DeleteAccountResult, err error) {
	res = &accountcex.DeleteAccountResult{}
	s.logger.Print("accountCex.deleteAccount")
	accountDbId := p.AccountID
	cexAccountService := apps_cex_account.NewCexAccountService(s.logger)
	err = cexAccountService.DeleteAccount(ctx, accountDbId)
	if err != nil {
		err = errors.WithMessage(err, "delete account error")
		return
	}
	res.Code = 0
	res.Message = "delete account success"
	return
}

// Get all token balance information for the specified CEX account
func (s *accountCexsrvc) GetAllTokenBalances(ctx context.Context, p *accountcex.GetAllTokenBalancesPayload) (res *accountcex.GetAllTokenBalancesResult, err error) {
	s.logger.Printf("accountCex.getAllTokenBalances - Account ID: %s", p.AccountID)

	// Create CEX account service instance
	cexAccountService := apps_cex_account.NewCexAccountService(s.logger)

	// Call service to get all token balances
	tokenBalances, err := cexAccountService.GetAllTokenBalances(ctx, p.AccountID)
	if err != nil {
		// Return error response
		return &accountcex.GetAllTokenBalancesResult{
			Code:    30001,
			Message: err.Error(),
			Result:  nil,
		}, nil
	}

	// Convert service returned token balance list to API response format
	apiTokenBalances := make([]*accountcex.CATokenBalance, len(tokenBalances))
	for i, tokenBalance := range tokenBalances {
		apiTokenBalances[i] = &accountcex.CATokenBalance{
			Asset:  tokenBalance.Asset,
			Free:   tokenBalance.Free,
			Locked: tokenBalance.Locked,
			Total:  tokenBalance.Total,
		}
	}

	// Return successful response
	return &accountcex.GetAllTokenBalancesResult{
		Code:    0,
		Message: "success",
		Result:  apiTokenBalances,
	}, nil
}
