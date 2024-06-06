package statuses

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"ticket/pkg/apikit"
	"ticket/pkg/auth"
	"ticket/pkg/db"

	"github.com/guregu/null/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/sync/errgroup"
)

type Handler struct {
	DB      *sql.DB
	Queries *db.Queries
	Auth    *auth.Auth
}

func New(api *apikit.API) *Handler {
	return &Handler{
		DB:      api.DB,
		Queries: db.New(api.DB),
		Auth:    auth.New(api.Config),
	}
}

func (h *Handler) CreateStatus(c echo.Context) error {
	claims := c.Get("claims").(*auth.Claims)

	boardID, err := strconv.ParseUint(c.Param("board_id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var body struct {
		Title string `json:"title" validate:"required,min=3,max=50"`
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

	board, err := h.Queries.GetBoard(ctx, db.GetBoardParams{
		ID:     uint32(boardID),
		UserID: claims.UserID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "board not found")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	count, err := h.Queries.CountStatusByBoardID(ctx, board.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	tx, err := h.DB.Begin()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer tx.Rollback()
	qtx := h.Queries.WithTx(tx)

	err = qtx.CreateStatus(ctx, db.CreateStatusParams{
		BoardID:   board.ID,
		Title:     null.NewString(body.Title, true),
		SortOrder: uint32(count + 1),
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	status, err := qtx.GetLastInsertStatus(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, db.NewStatusWithRelated(status, nil))
}

func (h *Handler) UpdateStatusPartial(c echo.Context) error {
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
		Title *string `json:"title" validate:"omitempty,min=3,max=50"`
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

	statusWithBoard, err := qtx.GetStatusWithBoard(ctx, db.GetStatusWithBoardParams{
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

	isChanged := false
	statusParams := db.UpdateStatusParams{
		ID:        statusWithBoard.Status.ID,
		Title:     statusWithBoard.Status.Title,
		SortOrder: statusWithBoard.Status.SortOrder,
	}

	if body.Title != nil {
		isChanged = true
		statusParams.Title = null.NewString(*body.Title, true)
	}

	if isChanged {
		err = qtx.UpdateStatus(ctx, statusParams)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		err = tx.Commit()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	status, err := h.Queries.GetStatus(ctx, db.GetStatusParams{
		ID: sql.NullInt32{Int32: int32(statusID), Valid: true},
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	tickets, err := h.Queries.GetTickets(ctx, db.GetTicketsParams{
		StatusIds:          []uint32{status.ID},
		SortOrderDirection: null.StringFrom("asc"),
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, db.NewStatusWithRelated(status, tickets))
}

func (h *Handler) SortStatusesOrder(c echo.Context) error {
	claims := c.Get("claims").(*auth.Claims)

	boardID, err := strconv.ParseUint(c.Param("board_id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var body struct {
		StatuseIDs []uint64 `json:"status_ids" validate:"required,dive"`
	}

	err = c.Bind(&body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = c.Validate(&body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var statusIDs []uint32
	statusIDMap := make(map[uint64]bool)
	for _, statusID := range body.StatuseIDs {
		if _, exists := statusIDMap[statusID]; exists {
			return echo.NewHTTPError(http.StatusBadRequest, "status id must be unique")
		}

		statusIDMap[statusID] = true
		statusIDs = append(statusIDs, uint32(statusID))
	}

	ctx := c.Request().Context()

	count, err := h.Queries.CountStatusWithBoard(ctx, db.CountStatusWithBoardParams{
		Ids:     statusIDs,
		BoardID: uint32(boardID),
		UserID:  claims.UserID,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if count != int64(len(statusIDs)) {
		return echo.NewHTTPError(http.StatusNotFound, "some status id not exist")
	}

	count, err = h.Queries.CountStatusWithBoardExclude(ctx, db.CountStatusWithBoardExcludeParams{
		Ids:     statusIDs,
		BoardID: uint32(boardID),
		UserID:  claims.UserID,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if count > 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "some status id is missing")
	}

	tx, err := h.DB.Begin()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer tx.Rollback()
	qtx := h.Queries.WithTx(tx)

	subctx1, cancel := context.WithCancel(ctx)
	g, subctx1 := errgroup.WithContext(subctx1)
	defer cancel()

	for i, statusID := range body.StatuseIDs {
		i := i
		s := statusID
		g.Go(func() error {
			err = qtx.UpdateStatusSortOrder(subctx1, db.UpdateStatusSortOrderParams{
				SortOrder: uint32(i + 1),
				ID:        uint32(s),
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

	subctx2, cancel := context.WithCancel(ctx)
	g, subctx2 = errgroup.WithContext(subctx2)
	defer cancel()
	chtickets := make(chan []db.Ticket)

	g.Go(func() error {
		tickets, err := h.Queries.GetTickets(subctx2, db.GetTicketsParams{
			StatusIds:          statusIDs,
			SortOrderDirection: null.StringFrom("asc"),
		})
		if err != nil {
			cancel()

			return err
		}

		chtickets <- tickets

		return nil
	})

	statuses, err := h.Queries.GetStatuses(ctx, db.GetStatusesParams{
		BoardID:            sql.NullInt32{Int32: int32(boardID), Valid: true},
		SortOrderDirection: null.StringFrom("asc"),
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var tickets []db.Ticket
	select {
	case <-subctx2.Done():
		return echo.NewHTTPError(http.StatusInternalServerError, subctx2.Err().Error())
	case tickets = <-chtickets:
	}

	err = g.Wait()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var statusesWithRelated []db.StatusWithRelated
	for _, status := range statuses {
		statusesWithRelated = append(statusesWithRelated, db.NewStatusWithRelated(status, tickets))
	}

	return c.JSON(http.StatusOK, statusesWithRelated)
}

func (h *Handler) BulkUpdateTicketOrderInStatuses(c echo.Context) error {
	claims := c.Get("claims").(*auth.Claims)

	boardID, err := strconv.ParseUint(c.Param("board_id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var body struct {
		Statuses []struct {
			ID        uint64   `json:"id" validate:"required"`
			TicketIDs []uint64 `json:"ticket_ids" validate:"required,dive"`
		} `json:"statuses" validate:"required,dive"`
	}

	err = c.Bind(&body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = c.Validate(&body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var statusIDs []uint32
	var ticketIDs []uint64
	statusIDMap := make(map[uint64]bool)
	ticketIDMap := make(map[uint64]bool)
	for _, status := range body.Statuses {
		for _, ID := range ticketIDs {
			if _, exists := ticketIDMap[ID]; exists {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("ticket id %d must be unique", ID))
			}

			ticketIDMap[ID] = true
		}

		if _, exists := statusIDMap[status.ID]; exists {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("status id %d must be unique", status.ID))
		}

		statusIDs = append(statusIDs, uint32(status.ID))
		ticketIDs = append(ticketIDs, status.TicketIDs...)
	}

	ctx := c.Request().Context()
	prectx, cancel := context.WithCancel(ctx)
	g, prectx := errgroup.WithContext(prectx)
	defer cancel()

	chtotal := make(chan int64)

	g.Go(func() error {
		total, err := h.Queries.CountTicketWithBoard(ctx, db.CountTicketWithBoardParams{
			Ids:     ticketIDs,
			BoardID: uint32(boardID),
			UserID:  claims.UserID,
		})
		if err != nil {
			cancel()

			return err
		}

		chtotal <- total

		return nil
	})

	total, err := h.Queries.CountStatusWithBoard(ctx, db.CountStatusWithBoardParams{
		Ids:     statusIDs,
		BoardID: uint32(boardID),
		UserID:  claims.UserID,
	})
	if err != nil {
		return err
	}

	if total != int64(len(statusIDs)) {
		return fmt.Errorf(fmt.Sprintf("ticket id %d not found", ticketIDs))
	}

	select {
	case <-prectx.Done():
		return echo.NewHTTPError(http.StatusInternalServerError, prectx.Err().Error())
	case total := <-chtotal:
		if total != int64(len(ticketIDs)) {
			return echo.NewHTTPError(http.StatusBadRequest, "some ticket's id does exist")
		}
	}

	err = g.Wait()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tx, err := h.DB.Begin()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer tx.Rollback()
	qtx := h.Queries.WithTx(tx)

	subctx, cancel := context.WithCancel(ctx)
	g, subctx = errgroup.WithContext(subctx)
	defer cancel()

	for _, s := range body.Statuses {
		statusID := s.ID

		for i, ID := range s.TicketIDs {
			i := i
			ID := ID
			g.Go(func() error {
				err = qtx.UpdateTicketSortOrderAndStatusID(subctx, db.UpdateTicketSortOrderAndStatusIDParams{
					StatusID:  uint32(statusID),
					SortOrder: uint32(i + 1),
					ID:        ID,
				})
				if err != nil {
					return err
				}

				return nil
			})
		}
	}

	err = g.Wait()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	posctx, cancel := context.WithCancel(ctx)
	g, posctx = errgroup.WithContext(posctx)
	defer cancel()

	chtickets := make(chan []db.Ticket)

	g.Go(func() error {
		tickets, err := h.Queries.GetTickets(posctx, db.GetTicketsParams{
			StatusIds:          statusIDs,
			SortOrderDirection: null.StringFrom("asc"),
		})
		if err != nil {
			cancel()

			return err
		}

		chtickets <- tickets

		return nil
	})

	statuses, err := h.Queries.GetStatuses(ctx, db.GetStatusesParams{
		Ids:                statusIDs,
		SortOrderDirection: null.StringFrom("asc"),
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var tickets []db.Ticket
	select {
	case <-posctx.Done():
		return echo.NewHTTPError(http.StatusInternalServerError, posctx.Err().Error())
	case tickets = <-chtickets:
	}

	err = g.Wait()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, db.NewStatusesWithRelated(statuses, tickets))
}
