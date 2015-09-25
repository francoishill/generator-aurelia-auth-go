package DefaultAuthenticationService

import (
	. "<%= OWN_GO_IMPORT_PATH %>/Interface/Authentication"
	. "github.com/francoishill/golang-common-ddd/Interface/Security/Authentication"
	"net/http"

	. "<%= OWN_GO_IMPORT_PATH %>/Entities/User"
)

type service struct {
	BaseAuthenticationService
	UserRepository
}

func (s *service) GetUserFromRequest(r *http.Request) User {
	baseUser := s.BaseAuthenticationService.BaseGetUserFromRequest(r)
	return s.UserRepository.GetByUUID(baseUser.UUID())
}

func New(authenticationService BaseAuthenticationService, userRepository UserRepository) AuthenticationService {
	return &service{
		authenticationService,
		userRepository,
	}
}
