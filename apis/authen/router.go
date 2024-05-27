package authen

import (
	"ticket/apis"
	"ticket/apis/authen/controllers"
)

func Router(api *apis.API) {
	controllers.NewAuthController(api.App.Group("/auth"), api)
}
