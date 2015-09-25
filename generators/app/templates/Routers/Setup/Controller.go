package Setup

import (
	"net/http"

	. "<%= OWN_GO_IMPORT_PATH %>/Context/RouterContext"
)

type Controller interface{}

type optionsHandler interface {
	Options(w http.ResponseWriter, r *http.Request, ctx *RouterContext)
}

type getHandler interface {
	Get(w http.ResponseWriter, r *http.Request, ctx *RouterContext)
}

type headHandler interface {
	Head(w http.ResponseWriter, r *http.Request, ctx *RouterContext)
}

type postHandler interface {
	Post(w http.ResponseWriter, r *http.Request, ctx *RouterContext)
}

type putHandler interface {
	Put(w http.ResponseWriter, r *http.Request, ctx *RouterContext)
}

type deleteHandler interface {
	Delete(w http.ResponseWriter, r *http.Request, ctx *RouterContext)
}

func GetControllerMethods(controller Controller) map[string]ControllerMethod {
	m := make(map[string]ControllerMethod)

	cnt := 0

	if o, ok := controller.(optionsHandler); ok {
		cnt++
		m["OPTIONS"] = o.Options
	}

	if g, ok := controller.(getHandler); ok {
		cnt++
		m["GET"] = g.Get
	}

	if h, ok := controller.(headHandler); ok {
		cnt++
		m["HEAD"] = h.Head
	}

	if h, ok := controller.(postHandler); ok {
		cnt++
		m["POST"] = h.Post
	}

	if h, ok := controller.(putHandler); ok {
		cnt++
		m["PUT"] = h.Put
	}

	if h, ok := controller.(deleteHandler); ok {
		cnt++
		m["DELETE"] = h.Delete
	}

	if cnt == 0 {
		panic("Controller must have at least one exposed method 'Get', 'Put', etc.")
	}

	return m
}
