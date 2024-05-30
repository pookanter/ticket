package model

type Status struct {
	ID        uint32   `json:"id"`
	BoardID   uint32   `json:"board_id"`
	Title     string   `json:"title"`
	SortOrder uint32   `json:"sort_order"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	Tickets   []Ticket `json:"tickets,omitempty"`
}
