package dbutil

import "ticket/pkg/db"

type StatusWithRelated struct {
	db.Status
	Tickets []db.Ticket `json:"tickets"`
}

type BoardWithRelated struct {
	db.Board
	Statuses []StatusWithRelated `json:"statuses"`
}

func NewBoardWithRelated(b db.Board, s []db.Status, t []db.Ticket) BoardWithRelated {
	bw := BoardWithRelated{
		Board:    b,
		Statuses: []StatusWithRelated{},
	}

	for _, status := range s {
		sw := StatusWithRelated{
			Status:  status,
			Tickets: []db.Ticket{},
		}

		for _, ticket := range t {
			if ticket.StatusID == status.ID {
				sw.Tickets = append(sw.Tickets, ticket)
			}
		}

		bw.Statuses = append(bw.Statuses, sw)
	}

	return bw
}
