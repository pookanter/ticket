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
		bw.Statuses = append(bw.Statuses, NewStatusWithRelated(status, t))
	}

	return bw
}

func NewStatusWithRelated(s db.Status, t []db.Ticket) StatusWithRelated {
	sw := StatusWithRelated{
		Status:  s,
		Tickets: []db.Ticket{},
	}

	for _, ticket := range t {
		if uint32(ticket.StatusID.Int32) == s.ID {
			sw.Tickets = append(sw.Tickets, ticket)
		}
	}

	return sw
}
