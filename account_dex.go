package adminapiservice

import (
	accountdex "admin-panel/gen/account_dex"
	"context"
	"log"
)

// accountDex service example implementation.
// The example methods log the requests and return zero values.
type accountDexsrvc struct {
	logger *log.Logger
}

// NewAccountDex returns the accountDex service implementation.
func NewAccountDex(logger *log.Logger) accountdex.Service {
	return &accountDexsrvc{logger}
}

// WalletInfo implements walletInfo.
func (s *accountDexsrvc) WalletInfo(ctx context.Context, p *accountdex.WalletInfoPayload) (res *accountdex.WalletInfoResult, err error) {
	res = &accountdex.WalletInfoResult{}
	s.logger.Print("accountDex.walletInfo")
	return
}
