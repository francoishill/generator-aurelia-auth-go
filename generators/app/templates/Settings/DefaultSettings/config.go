package DefaultSettings

import (
	. "github.com/francoishill/golang-web-dry/errors/checkerror"

	"gopkg.in/sconf/ini.v0"
	"gopkg.in/sconf/sconf.v0"

	"fmt"
	"github.com/francoishill/goangi2/utils/osUtils"
	"strings"
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

func (c *config) assertNotBlank(str, msgIfBlank string) {
	if strings.TrimSpace(str) == "" {
		panic(msgIfBlank)
	}
}

func (c *config) assertDirExists(dir, msgIfBlank string) {
	if !osUtils.DirectoryExists(dir) {
		panic(msgIfBlank)
	}
}

func (c *config) assertFileExists(file, msgIfBlank string) {
	if !osUtils.FileExists(file) {
		panic(msgIfBlank)
	}
}

func (c *config) Validate() {
	defer func() {
		if r := recover(); r != nil {
			panic("ERROR validating CONFIG: " + fmt.Sprintf("%+v", r))
		}
	}()

	c.assertNotBlank(c.Server.FrontendUrl, "Server FrontendUrl cannot be blank")
	c.assertNotBlank(c.Server.BackendUrl, "Server BackendUrl cannot be blank")
	c.assertNotBlank(c.Database.MysqlDataSource, "Database MysqlDataSource cannot be blank")
	c.assertDirExists(c.Database.MysqlMigrationsDir, fmt.Sprintf("Database MysqlMigrationsDir ('%s') must be an existing dir", c.Database.MysqlMigrationsDir))
	c.assertNotBlank(c.Redis.HostAndPort, "Redis HostAndPort cannot be blank")
	c.assertFileExists(c.Jwt.PrivateKeyFilePath, fmt.Sprintf("Jwt PrivateKeyFilePath ('%s') must be an existing file", c.Jwt.PrivateKeyFilePath))
	c.assertFileExists(c.Jwt.PublicKeyFilePath, fmt.Sprintf("Jwt PublicKeyFilePath ('%s') must be an existing file", c.Jwt.PublicKeyFilePath))
}

func loadConfigFile(configPath string) *config {
	cfg := &config{}
	err := sconf.Must(cfg).Read(ini.File(configPath))
	CheckError(err)
	return cfg
}
