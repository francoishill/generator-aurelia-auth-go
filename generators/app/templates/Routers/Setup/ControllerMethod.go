package Setup

import (
	"net/http"

	. "<%= OWN_GO_IMPORT_PATH %>/Context/RouterContext"
)

type ControllerMethod func(w http.ResponseWriter, r *http.Request, ctx *RouterContext)
