package usecases

import "ticket/pkg/db"

type Ticket struct {
	ID          uint64 `json:"id"`
	StatusID    uint32 `json:"status_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Contact     string `json:"contact"`
	SortOrder   uint32 `json:"sort_order"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type StatusFullDetail struct {
	ID        uint32   `json:"id"`
	BoardID   uint32   `json:"board_id"`
	Title     string   `json:"title"`
	SortOrder uint32   `json:"sort_order"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	Tickets   []Ticket `json:"tickets"`
}

type BoardFullDetail struct {
	ID        uint32             `json:"id"`
	UserID    uint64             `json:"user_id"`
	Title     string             `json:"title"`
	SortOrder uint32             `json:"sort_order"`
	CreatedAt string             `json:"created_at"`
	UpdatedAt string             `json:"updated_at"`
	Statuses  []StatusFullDetail `json:"statuses"`
}

func NewBoardFullDetail(b db.Board) BoardFullDetail {
	return BoardFullDetail{
		ID:        b.ID,
		UserID:    b.UserID,
		Title:     b.Title.String,
		SortOrder: uint32(b.SortOrder),
		CreatedAt: b.CreatedAt.Time.String(),
		UpdatedAt: b.UpdatedAt.Time.String(),
		Statuses:  []StatusFullDetail{},
	}
}

func NewStatusFullDetail(s db.Status) StatusFullDetail {
	return StatusFullDetail{
		ID:        s.ID,
		BoardID:   s.BoardID,
		Title:     s.Title.String,
		SortOrder: s.SortOrder,
		CreatedAt: s.CreatedAt.Time.String(),
		UpdatedAt: s.UpdatedAt.Time.String(),
		Tickets:   []Ticket{},
	}
}

func NewTicket(t db.Ticket) Ticket {
	return Ticket{
		ID:          t.ID,
		StatusID:    t.StatusID,
		Title:       t.Title.String,
		Description: t.Description.String,
		Contact:     t.Contact.String,
		SortOrder:   t.SortOrder,
		CreatedAt:   t.CreatedAt.Time.String(),
		UpdatedAt:   t.UpdatedAt.Time.String(),
	}
}
