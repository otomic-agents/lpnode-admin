// Code generated by goa v3.11.0, DO NOT EDIT.
//
// chainConfig endpoints
//
// Command:
// $ goa gen admin-panel/design

package chainconfig

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Endpoints wraps the "chainConfig" service endpoints.
type Endpoints struct {
	SetChainList         goa.Endpoint
	DelChainList         goa.Endpoint
	ChainList            goa.Endpoint
	SetChainGasUsd       goa.Endpoint
	SetChainClientConfig goa.Endpoint
}

// NewEndpoints wraps the methods of the "chainConfig" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	return &Endpoints{
		SetChainList:         NewSetChainListEndpoint(s),
		DelChainList:         NewDelChainListEndpoint(s),
		ChainList:            NewChainListEndpoint(s),
		SetChainGasUsd:       NewSetChainGasUsdEndpoint(s),
		SetChainClientConfig: NewSetChainClientConfigEndpoint(s),
	}
}

// Use applies the given middleware to all the "chainConfig" service endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.SetChainList = m(e.SetChainList)
	e.DelChainList = m(e.DelChainList)
	e.ChainList = m(e.ChainList)
	e.SetChainGasUsd = m(e.SetChainGasUsd)
	e.SetChainClientConfig = m(e.SetChainClientConfig)
}

// NewSetChainListEndpoint returns an endpoint function that calls the method
// "setChainList" of service "chainConfig".
func NewSetChainListEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*SetChainListPayload)
		return s.SetChainList(ctx, p)
	}
}

// NewDelChainListEndpoint returns an endpoint function that calls the method
// "delChainList" of service "chainConfig".
func NewDelChainListEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*DelChainListPayload)
		return s.DelChainList(ctx, p)
	}
}

// NewChainListEndpoint returns an endpoint function that calls the method
// "chainList" of service "chainConfig".
func NewChainListEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return s.ChainList(ctx)
	}
}

// NewSetChainGasUsdEndpoint returns an endpoint function that calls the method
// "setChainGasUsd" of service "chainConfig".
func NewSetChainGasUsdEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*SetChainGasUsdPayload)
		return s.SetChainGasUsd(ctx, p)
	}
}

// NewSetChainClientConfigEndpoint returns an endpoint function that calls the
// method "setChainClientConfig" of service "chainConfig".
func NewSetChainClientConfigEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*SetChainClientConfigPayload)
		return s.SetChainClientConfig(ctx, p)
	}
}
