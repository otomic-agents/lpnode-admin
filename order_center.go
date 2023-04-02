package adminapiservice

import (
	ordercenter "admin-panel/gen/order_center"
	"admin-panel/service"
	"admin-panel/types"
	"context"
	"log"

	"github.com/aws/smithy-go/ptr"
	"go.mongodb.org/mongo-driver/bson"
)

// orderCenter service example implementation.
// The example methods log the requests and return zero values.
type orderCentersrvc struct {
	logger *log.Logger
}

// NewOrderCenter returns the orderCenter service implementation.
func NewOrderCenter(logger *log.Logger) ordercenter.Service {
	return &orderCentersrvc{logger}
}

// List implements list.
func (s *orderCentersrvc) List(ctx context.Context, p *ordercenter.ListPayload) (res *ordercenter.ListResult, err error) {
	res = &ordercenter.ListResult{Result: &ordercenter.OrderCenterRetResult{}}
	res.Result.List = make([]ordercenter.OrderCenterOrderItem, 0)

	ocls := &service.OrderCenterLogicService{}
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
	var orderList []types.CenterOrder
	if ptr.ToInt64(p.Status) == 1 {
		orderList, err = ocls.AllRedis()
	} else {
		orderList, pageCount, err = ocls.All(queryOption, bson.M{})
		if err != nil {
			return
		}
	}
	log.Println(pageCount, "🦕🦕🦕🦕🦕🦕")
	for _, v := range orderList {
		res.Result.List = append(res.Result.List, v)
	}
	res.Code = ptr.Int64(0)
	res.Message = ptr.String("")
	res.Result.PageCount = ptr.Int64(pageCount)
	s.logger.Print("orderCenter.list")
	return
}
