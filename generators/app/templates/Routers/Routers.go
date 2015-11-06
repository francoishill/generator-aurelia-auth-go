package Routers

import (
	. "<%= OWN_GO_IMPORT_PATH %>/Context/RouterContext"
	. "<%= OWN_GO_IMPORT_PATH %>/Controllers/Auth/Login"
	. "<%= OWN_GO_IMPORT_PATH %>/Controllers/Auth/Logout"
	. "<%= OWN_GO_IMPORT_PATH %>/Controllers/Auth/Register"
	. "<%= OWN_GO_IMPORT_PATH %>/Controllers/UserDetails"
	. "<%= OWN_GO_IMPORT_PATH %>/Routers/Setup"
)

func GetRouters(ctx *RouterContext) []*Router {
	return []*Router{
		NewRouterBuilder("/register").
			SetController(NewAuthRegisterController()).
			Build(),
		NewRouterBuilder("/login").
			SetController(NewAuthLoginController()).
			Build(),
		NewRouterBuilder("/logout").
			SetController(NewAuthLogoutController()).
			Build(),

		NewRouterBuilder("/user-details").
			AddMiddlewares(ctx.Middlewares.Authentication.CheckAuthentication).
			SetController(NewUserDetailsController()).
			Build(),
	}
}
