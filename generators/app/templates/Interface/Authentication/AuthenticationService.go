package Authentication

import (
	. "github.com/francoishill/golang-common-ddd/Interface/Security/Authentication"
	"net/http"

	. "<%= OWN_GO_IMPORT_PATH %>/Entities/User"
)

type AuthenticationService interface {
	BaseAuthenticationService

	GetUserFromRequest(r *http.Request) User
}
