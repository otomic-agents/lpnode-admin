package adminapiservice

import (
	settings "admin-panel/gen/settings"
	"admin-panel/service"
	"admin-panel/types"
	"admin-panel/utils"
	"context"
	"log"

	"github.com/aws/smithy-go/ptr"
	"github.com/pkg/errors"
)

// settings service example implementation.
// The example methods log the requests and return zero values.
type settingssrvc struct {
	logger *log.Logger
}

// NewSettings returns the settings service implementation.
func NewSettings(logger *log.Logger) settings.Service {
	return &settingssrvc{logger}
}

// Settings implements settings.
func (s *settingssrvc) Settings(ctx context.Context, p *settings.SettingsPayload) (res *settings.SettingsResult, err error) {
	res = &settings.SettingsResult{}
	settingsConfig := &types.SettingsConfig{
		RelayUri: p.RelayURI,
	}
	log.Println(settingsConfig)
	saved, err := service.NewSettingsService().SaveSetting(settingsConfig)
	if !saved {
		err = errors.WithMessage(utils.GetNoEmptyError(err), "SaveSetting err:")
		log.Println("SaveSetting err")
		return
	}
	res.Code = ptr.Int64(0)
	res.Message = ptr.String("ok")
	s.logger.Print("settings.settings")
	return
}
