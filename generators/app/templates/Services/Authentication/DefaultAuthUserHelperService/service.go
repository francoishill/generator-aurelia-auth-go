package DefaultAuthUserHelperService

import (
	. "github.com/francoishill/golang-common-ddd/Interface/Misc/Errors"
	. "github.com/francoishill/golang-common-ddd/Interface/Misc/Validation"
	. "github.com/francoishill/golang-common-ddd/Interface/Security/Authentication"

	. "<%= OWN_GO_IMPORT_PATH %>/Entities/User"
)

type service struct {
	errors     ErrorsService
	validation BaseValidationService
	userRepo   UserRepository
}

func (s *service) BaseVerifyAndGetUserFromCredentials(email, username, password string) BaseUser {
	return s.userRepo.VerifyAndGetUserFromCredentials(email, password)
}

func (s *service) BaseGetUserWithUUID(uid interface{}) BaseUser {
	return s.userRepo.GetByUUID(uid)
}

func (s *service) BaseRegisterUser(email, username, password string) BaseUser {
	if !s.validation.IsEmail(email) {
		panic(s.errors.CreateHttpStatusClientError_BadRequest("The email is not valid"))
	}
	if s.validation.IsEmptyStringOrWhitespace(password) {
		panic(s.errors.CreateHttpStatusClientError_BadRequest("Password cannot be empty"))
	}
	return s.userRepo.Insert("", email, password)
}

func New(errorsService ErrorsService, validationService BaseValidationService, userRepository UserRepository) BaseAuthUserHelperService {
	return &service{
		errorsService,
		validationService,
		userRepository,
	}
}
