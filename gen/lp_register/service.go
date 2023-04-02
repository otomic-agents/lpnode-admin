// Code generated by goa v3.11.0, DO NOT EDIT.
//
// lpRegister service
//
// Command:
// $ goa gen admin-panel/design

package lpregister

import (
	"context"
)

// 用于管理Lp到Client的注册
type Service interface {
	// RegisterAll implements registerAll.
	RegisterAll(context.Context) (res *RegisterAllResult, err error)
	// UnRegisterAll implements unRegisterAll.
	UnRegisterAll(context.Context) (res *UnRegisterAllResult, err error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "lpRegister"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [2]string{"registerAll", "unRegisterAll"}

// RegisterAllResult is the result type of the lpRegister service registerAll
// method.
type RegisterAllResult struct {
	Code    *int64
	Result  *string
	Message *string
}

// UnRegisterAllResult is the result type of the lpRegister service
// unRegisterAll method.
type UnRegisterAllResult struct {
	Code    *int64
	Result  *string
	Message *string
}
