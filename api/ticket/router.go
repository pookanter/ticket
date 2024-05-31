package ticket

import (
	"ticket/api/ticket/boards"
	"ticket/pkg/apikit"
	"ticket/pkg/auth"
)

func Router(api *apikit.API) {
	b := boards.New(api)

	api.App.Use(auth.Middleware(api.Config))
	api.App.GET("/boards", b.GetBoards)
	api.App.POST("/boards", b.CreateBoard)
}
