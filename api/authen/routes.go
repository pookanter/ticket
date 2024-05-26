package authen

import (
	"ticket/api"
	"ticket/api/authen/internal/controller"
)

type AuthenAPI struct {
	*api.API
}

func Router(a *api.API) {
	a.App.Group("/authen")
	controller.Index(a)
}
