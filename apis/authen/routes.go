package authen

import (
	"ticket/apis"
	"ticket/apis/authen/controllers"
)

func Router(api *apis.API) {
	api.App.Group("/authen")
	controllers.Index(api)
}
