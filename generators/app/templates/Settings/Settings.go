package Settings

import (
	. "github.com/francoishill/golang-common-ddd/Interface/Settings/ReloadableSettings"
)

type Settings interface {
	ReloadableSettings

	IsDevMode() bool

	ServerFrontendUrl() string
	ServerBackendUrl() string

	UseMock() bool

	MysqlDataSource() string
	MysqlMigrationsDir() string

	RedisHostAndPort() string
	RedisPassword() string
	RedisDB() int64

	JwtPrivateKeyBytes() []byte
	JwtPublicKeyBytes() []byte
	JwtExpirationDeltaHours() int
}
