package model

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
