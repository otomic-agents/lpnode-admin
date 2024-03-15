package adminapiservice

import (
	lpregister "admin-panel/gen/lp_register"
	"admin-panel/service"
	"context"
	"fmt"
	"log"

	"github.com/aws/smithy-go/ptr"
)

// lpRegister service example implementation.
// The example methods log the requests and return zero values.
type lpRegistersrvc struct {
	logger *log.Logger
}

// NewLpRegister returns the lpRegister service implementation.
func NewLpRegister(logger *log.Logger) lpregister.Service {
	return &lpRegistersrvc{logger}
}

// RegisterAll implements registerAll.
func (s *lpRegistersrvc) RegisterAll(ctx context.Context) (res *lpregister.RegisterAllResult, err error) {
	res = &lpregister.RegisterAllResult{}
	cpls := service.NewCtrlPanelLogicService()
	lprls := service.NewLpRegisterLogicService()
	ret, err := cpls.GetInstallRowByInstallType("ammClient")
	if err != nil {
		return
	}
	log.Println("currently there are n clients that need to be registered", len(ret))
	go func() {
		for _, item := range ret {
			if item.RegisterClientStatus == 0 {
				register, regErr := lprls.RegisterItem(item.ID.Hex(), item.ServiceName, item.Name, item.ChainType, item.ChainId, item.Namespace)
				if regErr != nil {
					log.Println("register failed", regErr)
					// err = regErr
					continue
					// return
				}
				if !register {
					log.Println(fmt.Errorf("registration failed, id:%s", item.ID.Hex()))
					continue
					// return
				}
			}
		}
	}()

	res.Code = ptr.Int64(0)
	res.Message = ptr.String("")
	return
}

// UnRegisterAll implements unRegisterAll.
func (s *lpRegistersrvc) UnRegisterAll(ctx context.Context) (res *lpregister.UnRegisterAllResult, err error) {
	res = &lpregister.UnRegisterAllResult{}
	s.logger.Print("lpRegister.unRegisterAll")
	return
}
