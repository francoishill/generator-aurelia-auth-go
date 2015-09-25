package Setup

import (
	"github.com/codegangsta/negroni"
	"net/http"

	. "<%= OWN_GO_IMPORT_PATH %>/Context/RouterContext"
)

type NegroniHandlerFunc func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)

func newNegroniMiddleware(ctx *RouterContext, middleWareFuncs []http.HandlerFunc, controllerMethod ControllerMethod) *negroni.Negroni {
	negroniHandlers := []negroni.Handler{}
	for ind, _ := range middleWareFuncs {
		negroniHandlers = append(negroniHandlers, negroni.Wrap(middleWareFuncs[ind]))
	}
	negroniHandlers = append(negroniHandlers, negroni.Wrap(createHttpHandlerFunc(ctx, controllerMethod)))
	return negroni.New(negroniHandlers...)
}

func createHttpHandlerFunc(ctx *RouterContext, fn ControllerMethod) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, ctx)
	})
}
