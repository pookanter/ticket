package tickets

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"ticket/pkg/apikit"
	"ticket/pkg/auth"
	"ticket/pkg/db"

	"github.com/guregu/null"
	"github.com/labstack/echo/v4"
	"golang.org/x/sync/errgroup"
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

func (h *Handler) SortTicketsOrder(c echo.Context) error {
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
		Tickets []struct {
			ID uint64 `json:"id" validate:"required"`
		} `json:"tickets" validate:"required,dive"`
	}

	err = c.Bind(&body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = c.Validate(&body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var ticketIDs []uint64
	ticketIDMap := make(map[uint64]bool)
	for _, ticket := range body.Tickets {
		if _, exists := ticketIDMap[ticket.ID]; exists {
			return echo.NewHTTPError(http.StatusBadRequest, "ticket id must be unique")
		}

		ticketIDMap[ticket.ID] = true
		ticketIDs = append(ticketIDs, ticket.ID)
	}

	ctx := c.Request().Context()
	tx, err := h.DB.Begin()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer tx.Rollback()
	qtx := h.Queries.WithTx(tx)

	subctx1, cancel := context.WithCancel(ctx)
	g, subctx1 := errgroup.WithContext(subctx1)
	defer cancel()

	chstatus := make(chan db.GetStatusWithBoardRow, len(body.Tickets))

	g.Go(func() error {
		status, err := qtx.GetStatusWithBoard(subctx1, db.GetStatusWithBoardParams{
			ID:      uint32(statusID),
			BoardID: uint32(boardID),
			UserID:  claims.UserID,
		})

		if err != nil {
			return err
		}

		chstatus <- status

		return nil
	})

	ticketWithBoards, err := qtx.GetTicketsWithBoard(subctx1, db.GetTicketsWithBoardParams{
		BoardID: uint32(boardID),
		UserID:  claims.UserID,
		Ids:     ticketIDs,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if len(ticketWithBoards) != len(ticketIDs) {
		for _, t := range ticketWithBoards {
			if _, exists := ticketIDMap[t.Ticket.ID]; !exists {
				return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("The ticket with ID %d was not found on the board with ID %d, or the board does not exist.", t.Ticket.ID, boardID))
			}
		}

		return echo.NewHTTPError(http.StatusNotFound, "ticket not found")
	}

	var status db.GetStatusWithBoardRow
	select {
	case <-subctx1.Done():
		return echo.NewHTTPError(http.StatusInternalServerError, subctx1.Err().Error())
	case status = <-chstatus:
	}

	err = g.Wait()
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("The status with ID %d was not found on the board with ID %d, or the board does not exist.", statusID, boardID))
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	subctx2, cancel := context.WithCancel(ctx)
	g, subctx2 = errgroup.WithContext(subctx2)
	defer cancel()

	for i, t := range body.Tickets {
		g.Go(func() error {
			err = qtx.UpdateTicketSortOrderAndStatusID(subctx2, db.UpdateTicketSortOrderAndStatusIDParams{
				StatusID:  status.Status.ID,
				SortOrder: uint32(i + 1),
				ID:        t.ID,
			})
			if err != nil {
				return err
			}

			return nil
		})
	}

	err = g.Wait()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	tickets, err := qtx.GetTickets(ctx, db.GetTicketsParams{
		StatusIds:          []uint32{uint32(statusID)},
		SortOrderDirection: null.StringFrom("ASC"),
	})

	return c.JSON(http.StatusOK, tickets)
}
