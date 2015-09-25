package Setup

import (
	"net/http"
)

type Router struct {
	urlPart     string
	middlewares []http.HandlerFunc
	controller  Controller
	subRouters  []*Router
}
