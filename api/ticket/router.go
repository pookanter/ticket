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
	bg.PUT("/:board_id", b.UpdateBoardByID)

	s := statuses.New(api)

	sg := bg.Group("/:board_id/statuses")
	sg.POST("", s.CreateStatus)
	sg.PUT("/sort-orders", s.SortStatusesOrder)
	sg.PATCH("/:status_id", s.UpdateStatusPartial)
	sg.PUT("/tickets/bulk-reorder", s.BulkUpdateTicketOrderInStatuses)

	t := tickets.New(api)
	tg := sg.Group("/:status_id/tickets")
	tg.POST("", t.CreateTicket)
	tg.PUT("/sort-orders", t.SortTicketsOrder)
	tg.PATCH("/:ticket_id", t.UpdateTicketPartial)
}
