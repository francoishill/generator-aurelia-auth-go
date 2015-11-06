package DefaultAuthorizationMiddleware

import (
	. "<%= OWN_GO_IMPORT_PATH %>/Interface/Authentication"
	. "<%= OWN_GO_IMPORT_PATH %>/Interface/Authorization"
)

type middleware struct {
	authentication AuthenticationService
	authorization  AuthorizationService
}

func New(authenticationService AuthenticationService, authorizationService AuthorizationService) AuthorizationMiddleware {
	return &middleware{
		authenticationService,
		authorizationService,
	}
}
