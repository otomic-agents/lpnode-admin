package adminapiservice

import (
	mainlogic "admin-panel/gen/main_logic"
	"admin-panel/service"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/smithy-go/ptr"
)

// mainLogic service example implementation.
// The example methods log the requests and return zero values.
type mainLogicsrvc struct {
	logger *log.Logger
}

// NewMainLogic returns the mainLogic service implementation.
func NewMainLogic(logger *log.Logger) mainlogic.Service {
	logger.Println("new Main.....")
	go func() {
		for {
			time.Sleep(time.Second * 60 * 5) // automatically register lpnode every five minutes
			cpls := service.NewCtrlPanelLogicService()
			ret, err := cpls.GetInstallRowByInstallType("ammClient")
			if err != nil {
				return
			}
			lprls := service.NewLpRegisterLogicService()
			log.Println("currently there are n clients that need to be registered", len(ret))
			for _, item := range ret {
				if item.RegisterClientStatus == 0 {
					register, regErr := lprls.RegisterItemWithoutRestart(item.ID.Hex(), item.ServiceName, item.Name, item.ChainType, item.ChainId, item.Namespace)
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
		}

	}()
	return &mainLogicsrvc{logger}
}

// MainLogic implements mainLogic.
func (s *mainLogicsrvc) MainLogic(ctx context.Context) (err error) {
	s.logger.Print("mainLogic.mainLogic")
	return
}
func (s *mainLogicsrvc) MainLogicLink(ctx context.Context) (res *mainlogic.MainLogicLinkResult, err error) {
	s.logger.Print("mainLogic.MainLogicLink")
	res = &mainlogic.MainLogicLinkResult{}
	res.Code = ptr.Int64(0)
	res.Data = ptr.String(fmt.Sprintf("%d", time.Now().Unix()))

	return
}
