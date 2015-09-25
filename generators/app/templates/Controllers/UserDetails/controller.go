package UserDetails

import (
	"net/http"

	. "<%= OWN_GO_IMPORT_PATH %>/Context/RouterContext"
)

type controller struct{}

func (c *controller) Get(w http.ResponseWriter, r *http.Request, ctx *RouterContext) {
	usr := ctx.AuthenticationService.GetUserFromRequest(r)

	ctx.Misc.HttpRenderHelperService.RenderJson(w, &struct {
		Id       int64
		FullName string
		Email    string
	}{
		usr.Id(),
		usr.FullName(),
		usr.Email(),
	})
}

func (c *controller) Post(w http.ResponseWriter, r *http.Request, ctx *RouterContext) {
	usr := ctx.AuthenticationService.GetUserFromRequest(r)

	tmpRequestObj := &struct {
		FullName string
	}{}
	ctx.Misc.HttpRequestHelperService.DecodeJsonRequest(r, tmpRequestObj)

	clonedUser := usr.Clone()
	clonedUser.SetFullName(tmpRequestObj.FullName)
	ctx.UserRepository.Update(clonedUser)
}

func NewUserDetailsController() *controller {
	return &controller{}
}
