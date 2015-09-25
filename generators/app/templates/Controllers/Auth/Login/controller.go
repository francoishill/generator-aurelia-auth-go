package Login

import (
	"net/http"

	. "<%= OWN_GO_IMPORT_PATH %>/Context/RouterContext"
)

type controller struct{}

func (c *controller) Post(w http.ResponseWriter, r *http.Request, ctx *RouterContext) {
	ctx.AuthenticationService.BaseLoginHandler(w, r)
}

func NewAuthLoginController() *controller {
	return &controller{}
}
