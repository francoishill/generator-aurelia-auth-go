package DefaultSettings

import (
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"io/ioutil"
)

func (s *settings) mustReadFile(filePath string) []byte {
	bytes, err := ioutil.ReadFile(filePath)
	CheckError(err)
	return bytes
}

func (s *settings) IsDevMode() bool           { return s.c.Common.DevMode }
func (s *settings) ServerFrontendUrl() string { return s.c.Server.FrontendUrl }
func (s *settings) ServerBackendUrl() string  { return s.c.Server.BackendUrl }
func (s *settings) UseMock() bool             { return s.c.Common.UseMock }

func (s *settings) MysqlDataSource() string    { return s.c.Database.MysqlDataSource }
func (s *settings) MysqlMigrationsDir() string { return s.c.Database.MysqlMigrationsDir }

func (s *settings) RedisHostAndPort() string { return s.c.Redis.HostAndPort }
func (s *settings) RedisPassword() string    { return s.c.Redis.Password }
func (s *settings) RedisDB() int64           { return s.c.Redis.DB }

func (s *settings) JwtPrivateKeyBytes() []byte   { return s.mustReadFile(s.c.Jwt.PrivateKeyFilePath) }
func (s *settings) JwtPublicKeyBytes() []byte    { return s.mustReadFile(s.c.Jwt.PublicKeyFilePath) }
func (s *settings) JwtExpirationDeltaHours() int { return s.c.Jwt.ExpirationDeltaHours }
