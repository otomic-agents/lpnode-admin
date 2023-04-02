package adminapiservice

import (
	mainlogic "admin-panel/gen/main_logic"
	"context"
	"log"
)

// mainLogic service example implementation.
// The example methods log the requests and return zero values.
type mainLogicsrvc struct {
	logger *log.Logger
}

// NewMainLogic returns the mainLogic service implementation.
func NewMainLogic(logger *log.Logger) mainlogic.Service {
	logger.Println("new Main.....")
	return &mainLogicsrvc{logger}
}

// MainLogic implements mainLogic.
func (s *mainLogicsrvc) MainLogic(ctx context.Context) (err error) {
	s.logger.Print("mainLogic.mainLogic")
	return
}
