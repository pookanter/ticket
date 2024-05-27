package users

import "ticket/pkg/apikit"

func Router(api *apikit.API) {
	h := NewHandler(api)

	api.App.Group("/users").GET("/me", h.GetMe, h.Auth.AuthMiddleware)
}
