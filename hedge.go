package adminapiservice

import (
	hedge "admin-panel/gen/hedge"
	"context"
	"log"
)

// hedge service example implementation.
// The example methods log the requests and return zero values.
type hedgesrvc struct {
	logger *log.Logger
}

// NewHedge returns the hedge service implementation.
func NewHedge(logger *log.Logger) hedge.Service {
	return &hedgesrvc{logger}
}

// List implements list.
func (s *hedgesrvc) List(ctx context.Context) (res *hedge.ListResult, err error) {
	res = &hedge.ListResult{}
	s.logger.Print("hedge.list")
	return
}

// Edit implements edit.
func (s *hedgesrvc) Edit(ctx context.Context, p *hedge.EditPayload) (res *hedge.EditResult, err error) {
	res = &hedge.EditResult{}
	s.logger.Print("hedge.edit")
	return
}

// Del implements del.
func (s *hedgesrvc) Del(ctx context.Context, p *hedge.DelPayload) (res *hedge.DelResult, err error) {
	res = &hedge.DelResult{}
	s.logger.Print("hedge.del")
	return
}
