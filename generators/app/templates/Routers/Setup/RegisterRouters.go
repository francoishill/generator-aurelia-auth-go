package Setup

import (
	"github.com/gorilla/mux"
	"net/http"

	. "<%= OWN_GO_IMPORT_PATH %>/Context/RouterContext"
)

func setupRouters(ctx *RouterContext, router *mux.Router, parentMiddleWare []http.HandlerFunc, routers []*Router) {
	if len(routers) == 0 {
		return
	}

	for _, rd := range routers {
		combinedMiddleWareHandlers := []http.HandlerFunc{}
		combinedMiddleWareHandlers = append(combinedMiddleWareHandlers, parentMiddleWare...)
		combinedMiddleWareHandlers = append(combinedMiddleWareHandlers, rd.middlewares...)

		for method, h := range GetControllerMethods(rd.controller) {
			muxRoute := router.Handle(rd.urlPart, newNegroniMiddleware(ctx, combinedMiddleWareHandlers, h))
			muxRoute.Methods(method)
		}

		var subRouterToUse *mux.Router
		if rd.urlPart != "" {
			subRouterToUse = router.PathPrefix(rd.urlPart).Subrouter()
		} else {
			subRouterToUse = router
		}
		setupRouters(ctx, subRouterToUse, combinedMiddleWareHandlers, rd.subRouters)
	}
}

func RegisterRouters(ctx *RouterContext, router *mux.Router, baseMiddleWare []http.HandlerFunc, routers []*Router) {
	//Routers
	setupRouters(ctx, router, baseMiddleWare, routers)
}
