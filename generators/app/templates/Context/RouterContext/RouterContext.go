package RouterContext

import (
	. "github.com/francoishill/golang-common-ddd/Interface/Logger"
	. "github.com/francoishill/golang-common-ddd/Interface/Misc/Encryption"
	. "github.com/francoishill/golang-common-ddd/Interface/Misc/Errors"
	. "github.com/francoishill/golang-common-ddd/Interface/Misc/HttpRenderHelper"
	. "github.com/francoishill/golang-common-ddd/Interface/Misc/HttpRequestHelper"
	. "github.com/francoishill/golang-common-ddd/Interface/Misc/Validation"
	. "github.com/francoishill/golang-common-ddd/Interface/Security/Authentication"
	. "github.com/francoishill/golang-common-ddd/Interface/Security/Authentication/Jwt"
	. "github.com/francoishill/golang-common-ddd/Interface/Storage/DbStorage"
	. "github.com/francoishill/golang-common-ddd/Interface/Storage/KeyValueStorage"
	"gopkg.in/redis.v3"

	"github.com/francoishill/golang-common-ddd/Implementations/Logger/DefaultLogger"
	"github.com/francoishill/golang-common-ddd/Implementations/Misc/Encryption/DefaultEncryptionService"
	"github.com/francoishill/golang-common-ddd/Implementations/Misc/Errors/DefaultErrorsService"
	"github.com/francoishill/golang-common-ddd/Implementations/Misc/HttpRenderHelper/DefaultHttpRenderHelperService"
	"github.com/francoishill/golang-common-ddd/Implementations/Misc/HttpRequestHelper/DefaultHttpRequestHelperService"
	"github.com/francoishill/golang-common-ddd/Implementations/Misc/Validation/DefaultBaseValidationService"
	"github.com/francoishill/golang-common-ddd/Implementations/Security/Authentication/DefaultAuthenticationMiddleware"
	"github.com/francoishill/golang-common-ddd/Implementations/Security/Authentication/Jwt/DefaultJwtHelperService"
	"github.com/francoishill/golang-common-ddd/Implementations/Security/Authentication/Jwt/JWTAuthenticationService"
	"github.com/francoishill/golang-common-ddd/Implementations/Storage/DbStorage/DatabaseSqlxDbStorage"
	"github.com/francoishill/golang-common-ddd/Implementations/Storage/KeyValueStorage/RedisKeyValueStorage"

	"<%= OWN_GO_IMPORT_PATH %>/Authentication/DefaultAuthUserHelperService"
	"<%= OWN_GO_IMPORT_PATH %>/Authentication/DefaultAuthenticationService"
	. "<%= OWN_GO_IMPORT_PATH %>/Entities/User"
	. "<%= OWN_GO_IMPORT_PATH %>/Interface/Authentication"
	. "<%= OWN_GO_IMPORT_PATH %>/Repositories/User"
	. "<%= OWN_GO_IMPORT_PATH %>/Settings"
)

type MiscServices struct {
	HttpRenderHelperService
	HttpRequestHelperService
	ErrorsService
	EncryptionService
	BaseValidationService
}

type Middlewares struct {
	Auth AuthenticationMiddleware
}

type RouterContext struct {
	Logger

	dbStorage       DbStorage
	keyValueStorage KeyValueStorage

	UserRepository

	jwtHelperService      JwtHelperService
	authUserHelperService BaseAuthUserHelperService

	Settings

	AuthenticationService
	*Middlewares
	Misc *MiscServices
}

func (r *RouterContext) OnSettingsReloaded() {
	r.loadRepos()
	r.Logger.Debug("Router Context reloaded.")
}

func (r *RouterContext) loadRepos() {
	r.Logger = r.getLogger()
	r.dbStorage = r.getStorage()
	r.keyValueStorage = r.getKeyValueStorage()

	r.Misc = &MiscServices{
		r.getHttpRenderHelperService(),
		r.getHttpRequestHelperService(),
		r.getErrorsService(),
		r.getEncryptionService(),
		r.getBaseValidationService(),
	}

	r.UserRepository = r.getUserRepository()

	r.jwtHelperService = r.getJwtHelperService()
	r.authUserHelperService = r.getAuthUserHelperService()

	r.AuthenticationService = r.getAuthenticationService()
	r.Middlewares = &Middlewares{
		r.getAuthenticationMiddleware(),
	}
}

func (r *RouterContext) getLogger() Logger {
	logFileName := "rolling-log.log"
	if r.Settings.UseMock() {
		return DefaultLogger.New(logFileName, "[MOCK]", r.Settings.IsDevMode())
	} else {
		return DefaultLogger.New(logFileName, "", r.Settings.IsDevMode())
	}
}

func (r *RouterContext) getStorage() DbStorage {
	d := DatabaseSqlxDbStorage.New("mysql", r.Settings.MysqlDataSource(), r.Settings.MysqlMigrationsDir())
	numMigrationsApplied := d.Migrate()
	r.Logger.Info("Applied %d migrations", numMigrationsApplied)
	return d
}

func (r *RouterContext) getKeyValueStorage() KeyValueStorage {
	redisOptions := &redis.Options{
		Addr:     r.Settings.RedisHostAndPort(),
		Password: r.Settings.RedisPassword(),
		DB:       r.Settings.RedisDB(),
	}
	return RedisKeyValueStorage.New(redisOptions)
}

func (r *RouterContext) getUserRepository() UserRepository {
	if r.Settings.UseMock() {
		return NewMockUserRepository()
	} else {
		return NewDbUserRepository(r.Logger, r.Misc.ErrorsService, r.Misc.EncryptionService, r.dbStorage)
	}
}

func (r *RouterContext) getJwtHelperService() JwtHelperService {
	return DefaultJwtHelperService.New(r.keyValueStorage)
}

func (r *RouterContext) getAuthUserHelperService() BaseAuthUserHelperService {
	return DefaultAuthUserHelperService.New(r.Misc.ErrorsService, r.Misc.BaseValidationService, r.UserRepository)
}

func (r *RouterContext) getAuthenticationService() AuthenticationService {
	baseAuthenticationService := JWTAuthenticationService.New(
		r.Logger,
		r.Settings.JwtPrivateKeyBytes(), r.Settings.JwtPublicKeyBytes(), r.Settings.JwtExpirationDeltaHours(),
		r.Misc.ErrorsService, r.Misc.HttpRequestHelperService,
		r.jwtHelperService, r.authUserHelperService)
	return DefaultAuthenticationService.New(baseAuthenticationService, r.UserRepository)
}

func (r *RouterContext) getAuthenticationMiddleware() AuthenticationMiddleware {
	return DefaultAuthenticationMiddleware.New(r.AuthenticationService)
}

func (r *RouterContext) getHttpRenderHelperService() HttpRenderHelperService {
	isDevelopment := r.Settings.IsDevMode()
	indentJSON := r.Settings.IsDevMode()
	templateExtensions := []string{".gohtml", ".gotmpl"}
	return DefaultHttpRenderHelperService.New(isDevelopment, indentJSON, templateExtensions)
}

func (r *RouterContext) getHttpRequestHelperService() HttpRequestHelperService {
	return DefaultHttpRequestHelperService.New()
}

func (r *RouterContext) getErrorsService() ErrorsService {
	return DefaultErrorsService.New()
}

func (r *RouterContext) getEncryptionService() EncryptionService {
	return DefaultEncryptionService.New()
}

func (r *RouterContext) getBaseValidationService() BaseValidationService {
	return DefaultBaseValidationService.New()
}

func NewRouterContext(settings Settings) *RouterContext {
	rc := &RouterContext{Settings: settings}
	rc.loadRepos()
	return rc
}
