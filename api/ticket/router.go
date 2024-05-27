package ticket

import (
	"ticket/pkg/apikit"
)

func Router(api *apikit.API) apikit.Router {
	return func(api *apikit.API) {
		api.App.Group("/boards")
	}
}
