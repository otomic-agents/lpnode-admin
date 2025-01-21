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
	
	// MainLogic_RefreshLpWallet()
	return &mainLogicsrvc{logger}
}
func MainLogic_RefreshLpWallet() {
	index := 0
	go func() {
		if index == 0 {
			time.Sleep(time.Second * 60 * 1) //Refresh for the first time with an interval of 1 minute.
		} else {
			time.Sleep(time.Second * 60 * 20)
		}

		log.Println("auto refresh lp wallet")
		dwls := service.NewDexWalletLogicService()
		dwls.RefreshLpWallet()
	}()
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
