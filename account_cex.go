package adminapiservice

import (
	accountcex "admin-panel/gen/account_cex"
	"context"
	"log"
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

// WalletInfo implements walletInfo.
func (s *accountCexsrvc) WalletInfo(ctx context.Context) (res *accountcex.WalletInfoResult, err error) {
	res = &accountcex.WalletInfoResult{}
	s.logger.Print("accountCex.walletInfo")
	return
}
