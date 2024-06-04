package tickets

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"ticket/pkg/apikit"
	"ticket/pkg/auth"
	"ticket/pkg/db"

	"github.com/guregu/null"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	DB       *sql.DB
	DBConfig apikit.DBConfig
	Queries  *db.Queries
	Auth     *auth.Auth
}

func New(api *apikit.API) *Handler {
	return &Handler{
		DB:       api.DB,
		DBConfig: api.Config.DB(),
		Queries:  db.New(api.DB),
		Auth:     auth.New(api.Config),
	}
}

func (h *Handler) CreateTicket(c echo.Context) error {
	claims := c.Get("claims").(*auth.Claims)

	boardID, err := strconv.ParseUint(c.Param("board_id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	statusID, err := strconv.ParseUint(c.Param("status_id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var body struct {
		Title       string `json:"title" validate:"required,min=3,max=100"`
		Description string `json:"description" validate:"required,min=3,max=500"`
		Contact     string `json:"contact" validate:"required,min=3,max=100"`
	}

	err = c.Bind(&body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = c.Validate(&body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	tx, err := h.DB.Begin()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer tx.Rollback()
	qtx := h.Queries.WithTx(tx)

	status, err := qtx.GetStatusWithBoard(ctx, db.GetStatusWithBoardParams{
		ID:      uint32(statusID),
		BoardID: uint32(boardID),
		UserID:  claims.UserID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "status not found")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	count, err := qtx.CountTicketByStatusID(ctx, status.Status.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = qtx.CreateTicket(ctx, db.CreateTicketParams{
		StatusID:    uint32(status.Status.ID),
		Title:       null.NewString(body.Title, true),
		Description: null.NewString(body.Description, true),
		Contact:     null.NewString(body.Contact, true),
		SortOrder:   uint32(count),
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	ticket, err := qtx.GetLastInsertTicketByStatusID(ctx, uint32(status.Status.ID))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, ticket)
}

func (h *Handler) UpdateTicketPartial(c echo.Context) error {
	claims := c.Get("claims").(*auth.Claims)

	boardID, err := strconv.ParseUint(c.Param("board_id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	statusID, err := strconv.ParseUint(c.Param("status_id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ticketID, err := strconv.ParseUint(c.Param("ticket_id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var body struct {
		Title       *string `json:"title" validate:"omitempty,min=3,max=100"`
		Description *string `json:"description" validate:"omitempty,min=3,max=500"`
		Contact     *string `json:"contact" validate:"omitempty,min=3,max=100"`
		SortOrder   *uint32 `json:"sort_order" validte:"omitempty,min=0"`
		StatusID    *uint32 `json:"status_id" validate:"omitempty,min=0"`
	}

	err = c.Bind(&body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = c.Validate(&body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	tx, err := h.DB.Begin()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer tx.Rollback()
	qtx := h.Queries.WithTx(tx)

	ticket, err := qtx.GetTicketWithBoard(ctx, db.GetTicketWithBoardParams{
		ID:      ticketID,
		BoardID: uint32(boardID),
		UserID:  claims.UserID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "ticket not found")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if statusID != uint64(ticket.Ticket.StatusID) {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("status_id is not match, expected: %d", ticket.Ticket.StatusID))
	}

	isChanged := false
	ticketParam := db.UpdateTicketParams{
		StatusID:    ticket.Ticket.StatusID,
		Title:       ticket.Ticket.Title,
		Description: ticket.Ticket.Description,
		Contact:     ticket.Ticket.Contact,
		SortOrder:   ticket.Ticket.SortOrder,
		ID:          ticket.Ticket.ID,
	}

	if body.Title != nil {
		isChanged = true
		ticketParam.Title = null.NewString(*body.Title, true)
	}

	if body.Description != nil {
		isChanged = true
		ticketParam.Description = null.NewString(*body.Description, true)
	}

	if body.Contact != nil {
		isChanged = true
		ticketParam.Contact = null.NewString(*body.Contact, true)
	}

	if isChanged {
		err = qtx.UpdateTicket(ctx, ticketParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		err = tx.Commit()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	t, err := h.Queries.GetTicketByID(ctx, ticketID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, t)
}
