package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/negroni"
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	. "github.com/francoishill/golang-web-dry/middleware/accessloggingmiddleware"
	. "github.com/francoishill/golang-web-dry/middleware/recoverymiddleware"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/thoas/stats"
	"gopkg.in/tylerb/graceful.v1"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	. "github.com/francoishill/golang-common-ddd/Interface/Misc/Errors/ClientError"

	. "<%= OWN_GO_IMPORT_PATH %>/Context/RouterContext"
	"<%= OWN_GO_IMPORT_PATH %>/Routers"
	. "<%= OWN_GO_IMPORT_PATH %>/Routers/Setup"
	"<%= OWN_GO_IMPORT_PATH %>/Settings/DefaultSettings"
)

func getBaseUrlWithoutSlash(url string) string {
	returnUrl := url
	for len(returnUrl) > 0 && returnUrl[len(returnUrl)-1] == '/' {
		returnUrl = returnUrl[0 : len(returnUrl)-1]
	}
	return returnUrl
}

func getNegroniHandlers(ctx *RouterContext, router *mux.Router) []negroni.Handler {
	tmpArray := []negroni.Handler{}

	if frontendUrl := ctx.Settings.ServerFrontendUrl(); strings.TrimSpace(frontendUrl) != "" {
		tmpArray = append(tmpArray, cors.New(cors.Options{
			AllowedOrigins: []string{getBaseUrlWithoutSlash(frontendUrl)},
			AllowedHeaders: []string{"*"},
		}))
	}

	if ctx.Settings.IsDevMode() {
		tmpArray = append(tmpArray, NewAccessLoggingMiddleware(NewSimpleAccessInfoHandler(
			func(info *StartAccessInfo) {
				ctx.Logger.Debug(
					"START:     %s %s, RemAddr: %s, RemIP: %s, Proxies: %s",
					info.HttpMethod, info.RequestURI, info.RemoteAddr, info.RemoteIP, info.Proxies)
			},
			func(info *EndAccessInfo) {
				if !info.GotPanic {
					ctx.Logger.Debug(
						"END [OK]:  %s %s, Status: %d %s, Duration: %s",
						info.HttpMethod, info.RequestURI, info.Status, info.StatusText, info.Duration)
				} else {
					ctx.Logger.Debug(
						"END [ERR]: %s %s, Status: %d %s, Duration: %s",
						info.HttpMethod, info.RequestURI, info.Status, info.StatusText, info.Duration)
				}
			},
		)))

		middleware := stats.New()
		router.HandleFunc("/stats.json", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			b, _ := json.Marshal(middleware.Data())
			w.Write(b)
		})
		tmpArray = append(tmpArray, middleware)
	}

	tmpArray = append(tmpArray, NewRecovery(func(errDetails *RecoveredErrorDetails) *RecoveryResponse {
		switch errObj := errDetails.OriginalError.(type) {
		case *ClientError:
			//No logging for this error, this is client side only
			return &RecoveryResponse{
				errObj.StatusCode,
				errObj.StatusText,
			}
		default:
			ctx.Logger.Error("ERROR recovered: %s\nStack:\n%s", errDetails.Error, errDetails.StackTrace)
			return nil
		}
	}))

	return tmpArray
}

func setupApiV1Routes(ctx *RouterContext, router *mux.Router) {
	baseRouterMiddleware := []http.HandlerFunc{}

	apiV1Router := router.PathPrefix("/api/v1").Subrouter()

	routers := Routers.GetRouters(ctx)
	RegisterRouters(ctx, apiV1Router, baseRouterMiddleware, routers)
}

type notFoundHandler struct {
	ctx *RouterContext
}

func (n *notFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	panic(n.ctx.Misc.ErrorsService.CreateHttpStatusClientError_NotFound("Page not found"))
}

func createAndRegisterRoutersHandler(ctx *RouterContext) http.Handler {
	mainRouter := mux.NewRouter().StrictSlash(true)
	setupApiV1Routes(ctx, mainRouter)

	mainRouter.NotFoundHandler = &notFoundHandler{ctx}

	n := negroni.New(getNegroniHandlers(ctx, mainRouter)...)
	n.UseHandler(context.ClearHandler(mainRouter))
	return n
}

type registerRoutesSettingsObserver struct {
	ctx    *RouterContext
	server *graceful.Server
}

func (r *registerRoutesSettingsObserver) OnSettingsReloaded() {
	r.ctx.Logger.Debug("Settings reloaded, now re-registering routers")
	r.server.Handler = nil
	r.server.Handler = createAndRegisterRoutersHandler(r.ctx)
}

func main() {
	fmt.Println("")
	fmt.Println("")
	fmt.Println("-------------------------STARTUP OF API V1 SERVER-------------------------")

	if len(os.Args) < 2 {
		panic("Cannot startup, the first argument must be the config path")
	}

	settings := DefaultSettings.New(os.Args[1])
	ctx := NewRouterContext(settings)
	settings.SubscribeReloadObserver(ctx)

	backendUrl, err := url.Parse(settings.ServerBackendUrl())
	CheckError(err)

	ctx.Logger.Info("Now serving on %s", settings.ServerBackendUrl())

	var gracefulTimeout time.Duration = 0
	if !settings.IsDevMode() {
		gracefulTimeout = 5 * time.Second
	}

	srv := &graceful.Server{
		Timeout: gracefulTimeout,
		Server:  &http.Server{Addr: backendUrl.Host, Handler: createAndRegisterRoutersHandler(ctx)},
	}

	r := &registerRoutesSettingsObserver{ctx, srv}
	settings.SubscribeReloadObserver(r)

	err = srv.ListenAndServe()
	CheckError(err)
}
