package DefaultSettings

import (
	"github.com/francoishill/golang-common-ddd/Implementations/Settings/AutoReloadingSettings"
	. "github.com/francoishill/golang-common-ddd/Interface/Settings/ReloadableSettings"
	"github.com/ian-kent/go-log/log"

	. "<%= OWN_GO_IMPORT_PATH %>/Settings"
)

type settings struct {
	c *config
	ReloadableSettings
}

func (s *settings) ValidateAndUseFile(filePath string) {
	cfg := loadConfigFile(filePath)
	cfg.Validate()
	s.c = cfg
	log.Info("Config loaded " + filePath)
}

func (s *settings) OnWatchReloadError(err error) {
	log.Error("Error: %+v", err)
}

func New(configPath string) Settings {
	s := &settings{}
	s.ReloadableSettings = AutoReloadingSettings.New(configPath, s)
	return s
}
