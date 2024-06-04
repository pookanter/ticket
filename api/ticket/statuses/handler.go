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

	"github.com/guregu/null"
	"github.com/labstack/echo/v4"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/status"
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
	tx, err := h.DB.Begin()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer tx.Rollback()
	qtx := h.Queries.WithTx(tx)

	board, err := qtx.GetBoard(ctx, db.GetBoardParams{
		ID:     uint32(boardID),
		UserID: claims.UserID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "board not found")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	count, err := qtx.CountStatusByBoardID(ctx, board.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = qtx.CreateStatus(ctx, db.CreateStatusParams{
		BoardID:   board.ID,
		Title:     null.NewString(body.Title, true),
		SortOrder: uint32(count),
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	statusID, err := qtx.GetLastInsertStatusID(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	status, err := qtx.GetStatus(ctx, db.GetStatusParams{
		ID: sql.NullInt32{Int32: int32(statusID), Valid: true},
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, status)
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

	status, err := qtx.GetStatus(ctx, db.GetStatusParams{
		ID: sql.NullInt32{Int32: int32(statusID), Valid: true},
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, status)
}

func (h *Handler) SortStatusesOrder(c echo.Context) error {
	claims := c.Get("claims").(*auth.Claims)

	boardID, err := strconv.ParseUint(c.Param("board_id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var body struct {
		Statuses []struct {
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

	var statusIDs []uint64
	statusIDMap := make(map[uint64]bool)
	for _, status := range body.Statuses {
		if _, exists := statusIDMap[status.ID]; exists {
			return echo.NewHTTPError(http.StatusBadRequest, "status id must be unique")
		}

		statusIDMap[status.ID] = true
		statusIDs = append(statusIDs, status.ID)
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

	stausesWithBoard, err := qtx.GetStatusesWithBoard(subctx1, db.GetStatusesWithBoardParams{
		BoardID: uint32(boardID),
		UserID:  claims.UserID,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if len(stausesWithBoard) != len(statusIDs) {
		for _, s := range stausesWithBoard {
			if _, exists := statusIDMap[uint64(s.Status.ID)]; !exists {
				return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("The status with ID %d was not found on the board with ID %d, or the board does not exist.", s.Status.ID, boardID))
			}
		}

		return echo.NewHTTPError(http.StatusNotFound, "ticket not found")
	}

	err = g.Wait()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	subctx2, cancel := context.WithCancel(ctx)
	g, subctx2 = errgroup.WithContext(subctx2)
	defer cancel()

	for i, t := range body.Statuses {
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
