package boards

import (
	"database/sql"
	"net/http"
	"ticket/pkg/apikit"
	"ticket/pkg/auth"
	"ticket/pkg/db"
	"ticket/pkg/model"
	"ticket/pkg/util"
	"time"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	DB        *db.Queries
	DBTimeOut time.Duration
	Auth      *auth.Auth
}

func NewHandler(api *apikit.API) *Handler {
	return &Handler{
		DB:        db.New(api.DB),
		DBTimeOut: api.Config.DB().TimeOut,
		Auth:      auth.New(api.Config),
	}
}

func (h *Handler) GetBoards(c echo.Context) error {
	claims := c.Get("claims").(*auth.Claims)

	dbBoards, err := h.DB.GetBoardsByUserId(c.Request().Context(), claims.UserID)
	if err != nil && err != sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	boards := make([]model.Board, len(dbBoards))
	for _, b := range dbBoards {
		dbStatuses, err := h.DB.GetStatusesByBoardId(c.Request().Context(), b.ID)
		if err != nil && err != sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		board := model.Board{
			ID:        b.ID,
			UserID:    b.UserID,
			Title:     b.Title.String,
			SortOrder: b.SortOrder,
			CreatedAt: b.CreatedAt.Time.Format(util.TimeFormat),
			UpdatedAt: b.UpdatedAt.Time.Format(util.TimeFormat),
			Statuses:  make([]model.Status, len(dbStatuses)),
		}

		for _, s := range dbStatuses {
			dbTickets, err := h.DB.GetTicketsByStatusId(c.Request().Context(), s.ID)
			if err != nil && err != sql.ErrNoRows {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			status := model.Status{
				ID:        s.ID,
				BoardID:   s.BoardID,
				Title:     s.Title.String,
				SortOrder: s.SortOrder,
				CreatedAt: s.CreatedAt.Time.Format(util.TimeFormat),
				UpdatedAt: s.UpdatedAt.Time.Format(util.TimeFormat),
				Tickets:   make([]model.Ticket, len(dbTickets)),
			}

			for _, t := range dbTickets {
				ticket := model.Ticket{
					ID:          t.ID,
					StatusID:    t.StatusID,
					Title:       t.Title.String,
					Description: t.Description.String,
					Contact:     t.Contact.String,
					SortOrder:   t.SortOrder,
					CreatedAt:   t.CreatedAt.Time.Format(util.TimeFormat),
					UpdatedAt:   t.UpdatedAt.Time.Format(util.TimeFormat),
				}

				status.Tickets = append(status.Tickets, ticket)
			}

			board.Statuses = append(board.Statuses, status)
		}

		boards = append(boards, board)
	}

	return c.JSON(http.StatusOK, boards)
}
