package DefaultAuthorizationService

import (
	. "<%= OWN_GO_IMPORT_PATH %>/Interface/Authorization"
)

type service struct {
}

func New() AuthorizationService {
	return &service{}
}
