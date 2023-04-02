package adminapiservice

import (
	ammordercenter "admin-panel/gen/amm_order_center"
	"admin-panel/service"
	"context"
	"log"

	"github.com/aws/smithy-go/ptr"
	"go.mongodb.org/mongo-driver/bson"
)

// ammOrderCenter service example implementation.
// The example methods log the requests and return zero values.
type ammOrderCentersrvc struct {
	logger *log.Logger
}

// NewAmmOrderCenter returns the ammOrderCenter service implementation.
func NewAmmOrderCenter(logger *log.Logger) ammordercenter.Service {
	return &ammOrderCentersrvc{logger}
}

// List implements list.
func (s *ammOrderCentersrvc) List(ctx context.Context, p *ammordercenter.ListPayload) (res *ammordercenter.ListResult, err error) {
	res = &ammordercenter.ListResult{Result: &ammordercenter.AmmOrderCenterRetResult{}}
	res.Result.List = make([]ammordercenter.OrderCenterAmmOrderItem, 0)
	ammocls := service.AmmOrderCenterLogicService{}
	queryOption := struct {
		Page     int64
		PageSize int64
		Status   int64
	}{
		Page:     1,
		PageSize: 20,
		Status:   0,
	}
	if ptr.ToInt64(p.Page) > 0 {
		queryOption.Page = ptr.ToInt64(p.Page)
	}
	if ptr.ToInt64(p.PageSize) > 0 {
		queryOption.PageSize = ptr.ToInt64(p.PageSize)
	}
	var pageCount int64 = 1
	orderList, pageCount, err := ammocls.All(queryOption, p.AmmName, bson.M{})
	for _, v := range orderList {
		res.Result.List = append(res.Result.List, v)
	}
	res.Code = ptr.Int64(0)
	res.Message = ptr.String("")
	res.Result.PageCount = &pageCount
	s.logger.Print("ammOrderCenter.list")
	return
}
