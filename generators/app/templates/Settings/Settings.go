package Settings

import (
	. "github.com/francoishill/golang-common-ddd/Interface/Settings/ReloadableSettings"
)

type storage interface {
	MysqlDataSource() string
	MysqlMigrationsDir() string

	RedisHostAndPort() string
	RedisPassword() string
	RedisDB() int64
}

type jwt interface {
	JwtPrivateKeyBytes() []byte
	JwtPublicKeyBytes() []byte
	JwtExpirationDeltaHours() int
}

type Settings interface {
	ReloadableSettings
	storage
	jwt

	IsDevMode() bool
	UseMock() bool

	ServerFrontendUrl() string
	ServerBackendUrl() string
}
