package statuses

import (
	"context"
	"database/sql"
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

	status, err := qtx.GetLastInsertStatusViewByBoardID(ctx, board.ID)
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
		Title     *string `json:"title" validate:"omitempty,min=3,max=50"`
		SortOrder *uint32 `json:"sort_order" validte:"omitempty,min=0"`
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

	status, err := qtx.GetStatus(ctx, db.GetStatusParams{
		ID: uint32(statusID),
		BoardID: sql.NullInt32{
			Int32: int32(boardID),
			Valid: true,
		},
		UserID: sql.NullInt64{
			Int64: int64(claims.UserID),
			Valid: true,
		},
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "status not found")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	isChanged := false
	statusParams := db.UpdateStatusParams{
		ID:        status.ID,
		Title:     status.Title,
		SortOrder: status.SortOrder,
	}

	if body.Title != nil {
		isChanged = true
		statusParams.Title = null.NewString(*body.Title, true)
	}

	subctx, cancel := context.WithCancel(ctx)
	defer cancel()
	g, subctx := errgroup.WithContext(subctx)

	if body.SortOrder != nil {
		isChanged = true
		count, err := qtx.CountStatusByBoardID(ctx, status.ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		var isOutOfRange = *body.SortOrder > uint32(count-1)

		if isOutOfRange {
			return echo.NewHTTPError(http.StatusBadRequest, "sort order out of range")
		}

		statusParams.SortOrder = *body.SortOrder

		g.Go(func() error {
			var newSortOrder uint32

			if status.SortOrder > statusParams.SortOrder {
				newSortOrder = statusParams.SortOrder
			} else if status.SortOrder < statusParams.SortOrder {
				newSortOrder = status.SortOrder - 1
			}

			friends, err := qtx.GetStatusesWithMinimumSortOrder(subctx, db.GetStatusesWithMinimumSortOrderParams{
				BoardID:            status.BoardID,
				SortOrder:          newSortOrder,
				SortOrderDirection: "asc",
			})
			if err != nil {
				return err
			}

			for _, s := range friends {
				if s.ID == status.ID {
					continue
				}

				newSortOrder++

				err = qtx.UpdateStatusSortOrder(subctx, db.UpdateStatusSortOrderParams{
					SortOrder: newSortOrder,
					ID:        s.ID,
				})
				if err != nil {
					return err
				}

			}

			return nil
		})
	}

	if isChanged {
		err = qtx.UpdateStatus(ctx, statusParams)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		err = g.Wait()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		err = tx.Commit()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	statusView, err := qtx.GetStatusView(ctx, uint32(statusID))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, statusView)
}
