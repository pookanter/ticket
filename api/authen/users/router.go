package users

import (
	"ticket/pkg/apikit"
	"ticket/pkg/auth"
)

func Router(api *apikit.API) {
	h := NewHandler(api)
	g := api.App.Group("/users")

	guard := auth.New(api.Config).AuthMiddleware

	g.GET("/me", h.GetMe, guard)
}
