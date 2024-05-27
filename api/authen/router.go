package authen

import (
	"ticket/api/authen/controllers"
	"ticket/pkg/apikit"
)

func Router(api *apikit.API) {
	controllers.NewAuthController(api.App.Group("/auth"), api)
}
