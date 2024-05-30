package model

type Board struct {
	ID        uint32   `json:"id"`
	UserID    uint64   `json:"user_id"`
	Title     string   `json:"title"`
	SortOrder uint32   `json:"sort_order"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	Statuses  []Status `json:"statuses"`
}
