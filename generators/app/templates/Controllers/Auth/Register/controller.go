package Register

import (
	"net/http"

	. "<%= OWN_GO_IMPORT_PATH %>/Context/RouterContext"
)

type controller struct{}

func (c *controller) Post(w http.ResponseWriter, r *http.Request, ctx *RouterContext) {
	ctx.AuthenticationService.BaseRegisterHandler(w, r)
}

func NewAuthRegisterController() *controller {
	return &controller{}
}
