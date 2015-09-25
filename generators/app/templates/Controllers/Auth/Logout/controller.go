package Logout

import (
	"net/http"

	. "<%= OWN_GO_IMPORT_PATH %>/Context/RouterContext"
)

type controller struct{}

func (c *controller) Post(w http.ResponseWriter, r *http.Request, ctx *RouterContext) {
	ctx.AuthenticationService.BaseLogoutHandler(w, r)
}

func NewAuthLogoutController() *controller {
	return &controller{}
}
