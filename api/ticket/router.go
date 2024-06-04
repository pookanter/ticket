package ticket

import (
	"ticket/api/ticket/boards"
	"ticket/api/ticket/statuses"
	"ticket/api/ticket/tickets"
	"ticket/pkg/apikit"
	"ticket/pkg/auth"
)

func Router(api *apikit.API) {
	b := boards.New(api)
	guard := auth.Middleware(api.Config)

	bg := api.App.Group("/boards", guard)
	bg.GET("", b.GetBoards)
	bg.GET("/:board_id", b.GetBoardByID)
	bg.POST("", b.CreateBoard)

	s := statuses.New(api)

	sg := bg.Group("/:board_id/statuses")
	sg.POST("", s.CreateStatus)
	sg.PATCH("/:status_id", s.UpdateStatusPartial)

	t := tickets.New(api)
	tg := sg.Group("/:status_id/tickets")
	tg.POST("", t.CreateTicket)
	tg.PUT("/sort_order", t.SortTicketsOrder)
	tg.PATCH("/:ticket_id", t.UpdateTicketPartial)
}
