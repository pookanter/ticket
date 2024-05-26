package authen

import (
	"ticket/api"
	"ticket/api/authen/controller"
)

func Router(a *api.API) {
	a.App.Group("/authen")
	controller.Index(a)
}
