package adminapiservice

import (
	authenticationlimiter "admin-panel/gen/authentication_limiter"
	"admin-panel/service"
	"context"
	"log"

	"github.com/aws/smithy-go/ptr"
)

// authenticationLimiter service example implementation.
// The example methods log the requests and return zero values.
type authenticationLimitersrvc struct {
	logger *log.Logger
}

// NewAuthenticationLimiter returns the authenticationLimiter service
// implementation.
func NewAuthenticationLimiter(logger *log.Logger) authenticationlimiter.Service {
	return &authenticationLimitersrvc{logger}
}

func (s *authenticationLimitersrvc) GetAuthenticationLimiter(ctx context.Context) (res *authenticationlimiter.GetAuthenticationLimiterResult, err error) {
	res = &authenticationlimiter.GetAuthenticationLimiterResult{}
	als := service.NewAuthenticationLimiterService()
	ret, err := als.Get()
	if err != nil {
		return
	}
	res.Code = ptr.Int64(0)
	res.Result = &ret.Data
	res.Message = ptr.String("")
	s.logger.Print("authenticationLimiter.getAuthenticationLimiter")
	return
}

func (s *authenticationLimitersrvc) SetAuthenticationLimiter(ctx context.Context, p *authenticationlimiter.SetAuthenticationLimiterPayload) (res *authenticationlimiter.SetAuthenticationLimiterResult, err error) {
	res = &authenticationlimiter.SetAuthenticationLimiterResult{}
	als := service.NewAuthenticationLimiterService()
	als.Set(ptr.ToString(p.AuthenticationLimiter))
	res.Code = ptr.Int64(0)
	res.Result = ptr.String("")
	res.Message = ptr.String("")

	s.logger.Print("authenticationLimiter.setAuthenticationLimiter")
	return
}

func (s *authenticationLimitersrvc) DelAuthenticationLimiter(ctx context.Context) (res *authenticationlimiter.DelAuthenticationLimiterResult, err error) {
	res = &authenticationlimiter.DelAuthenticationLimiterResult{}
	als := service.NewAuthenticationLimiterService()
	delCount, err := als.Del()
	if err != nil {
		return
	}
	res.Code = ptr.Int64(0)
	res.Result = ptr.Int64(delCount)
	res.Message = ptr.String("")
	s.logger.Print("authenticationLimiter.delAuthenticationLimiter")
	return
}
