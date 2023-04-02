// Code generated by goa v3.11.0, DO NOT EDIT.
//
// orderCenter service
//
// Command:
// $ goa gen admin-panel/design

package ordercenter

import (
	"context"
)

// 用于管理orderCenter
type Service interface {
	// List implements list.
	List(context.Context, *ListPayload) (res *ListResult, err error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "orderCenter"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [1]string{"list"}

// ListPayload is the payload type of the orderCenter service list method.
type ListPayload struct {
	Status   *int64
	Page     *int64
	PageSize *int64
}

// ListResult is the result type of the orderCenter service list method.
type ListResult struct {
	Code    *int64
	Result  *OrderCenterRetResult
	Message *string
}

type OrderCenterOrderItem interface{}

type OrderCenterRetResult struct {
	List      []OrderCenterOrderItem
	PageCount *int64
}
