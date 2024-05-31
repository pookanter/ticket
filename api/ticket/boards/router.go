package boards

import (
	"ticket/pkg/apikit"
	"ticket/pkg/auth"
)

func Router(api *apikit.API) {
	h := NewHandler(api)
	g := api.App.Group("/boards")

	guard := auth.New(api.Config).AuthMiddleware

	g.Use(guard)
	g.GET("", h.GetBoards)
	g.POST("", h.CreateBoard)
}
