package boards

import (
	"database/sql"
	"net/http"
	"ticket/pkg/apikit"
	"ticket/pkg/auth"
	"ticket/pkg/db"
	"ticket/pkg/model"
	"ticket/pkg/util"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	DB   *db.Queries
	Auth *auth.Auth
}

func New(api *apikit.API) *Handler {
	return &Handler{
		DB:   db.New(api.DB),
		Auth: auth.New(api.Config),
	}
}

func (h *Handler) GetBoards(c echo.Context) error {
	claims := c.Get("claims").(*auth.Claims)

	dbBoards, err := h.DB.GetBoardsByUserId(c.Request().Context(), claims.UserID)
	if err != nil && err != sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	boards := make([]model.Board, len(dbBoards))
	for i, b := range dbBoards {
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

		for j, s := range dbStatuses {
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

			for k, t := range dbTickets {
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

				status.Tickets[k] = ticket
			}

			board.Statuses[j] = status
		}

		boards[i] = board
	}

	return c.JSON(http.StatusOK, boards)
}

func (h *Handler) CreateBoard(c echo.Context) error {
	claims := c.Get("claims").(*auth.Claims)

	var body struct {
		Title string `json:"title" validate:"required,min=3,max=100"`
	}

	err := c.Bind(&body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = c.Validate(&body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	dbUser, err := h.DB.FindUserByID(c.Request().Context(), claims.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "user not found")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	sortOrder := uint32(0)
	lasbDbBoard, err := h.DB.GetLastBoardByUserId(c.Request().Context(), dbUser.ID)
	if err == nil {
		sortOrder = lasbDbBoard.SortOrder + 1
	} else if err != sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = h.DB.CreateBoard(c.Request().Context(), db.CreateBoardParams{
		UserID:    dbUser.ID,
		Title:     sql.NullString{String: body.Title, Valid: true},
		SortOrder: sortOrder,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	dbBoard, err := h.DB.GetLastCreatedBoardByUserId(c.Request().Context(), dbUser.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "board not found")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	board := model.Board{
		ID:        dbBoard.ID,
		UserID:    dbBoard.UserID,
		Title:     dbBoard.Title.String,
		SortOrder: dbBoard.SortOrder,
		CreatedAt: dbBoard.CreatedAt.Time.Format(util.TimeFormat),
		UpdatedAt: dbBoard.UpdatedAt.Time.Format(util.TimeFormat),
	}

	return c.JSON(http.StatusCreated, board)
}
