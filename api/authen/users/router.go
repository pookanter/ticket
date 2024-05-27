package users

import (
	"ticket/pkg/apikit"
	"ticket/pkg/auth"

	"github.com/labstack/echo/v4"
)

func Router(api *apikit.API) {
	h := NewHandler(api)
	g := api.App.Group("/users")

	var guard echo.MiddlewareFunc
	guard = auth.New(api.Config).AuthMiddleware

	g.Use(guard)
	g.GET("/me", h.GetMe)
}
