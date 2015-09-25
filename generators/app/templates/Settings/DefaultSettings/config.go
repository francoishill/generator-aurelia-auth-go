package DefaultSettings

import (
	. "github.com/francoishill/golang-web-dry/errors/checkerror"

	"gopkg.in/sconf/ini.v0"
	"gopkg.in/sconf/sconf.v0"
)

type config struct {
	Common struct {
		DevMode bool
		UseMock bool
	}
	Server struct {
		FrontendUrl string
		BackendUrl  string
	}
	Database struct {
		MysqlDataSource    string
		MysqlMigrationsDir string
	}
	Redis struct {
		HostAndPort string
		Password    string
		DB          int64
	}
	Jwt struct {
		PrivateKeyFilePath   string
		PublicKeyFilePath    string
		ExpirationDeltaHours int
	}
}

func (c *config) Validate() {
	//Nothing for now
}

func loadConfigFile(configPath string) *config {
	cfg := &config{}
	err := sconf.Must(cfg).Read(ini.File(configPath))
	CheckError(err)
	return cfg
}
