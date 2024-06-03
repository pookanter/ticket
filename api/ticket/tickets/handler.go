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

	status, err := qtx.GetStatus(ctx, db.GetStatusParams{
		ID:      uint32(statusID),
		BoardID: sql.NullInt32{Int32: int32(boardID), Valid: true},
		UserID:  sql.NullInt64{Int64: int64(claims.UserID), Valid: true},
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "status not found")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	count, err := qtx.CountTicketByStatusID(ctx, status.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = qtx.CreateTicket(ctx, db.CreateTicketParams{
		StatusID:    uint32(status.ID),
		Title:       null.NewString(body.Title, true),
		Description: null.NewString(body.Description, true),
		Contact:     null.NewString(body.Contact, true),
		SortOrder:   uint32(count),
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	ticket, err := qtx.GetLastInsertTicketByStatusID(ctx, uint32(status.ID))
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

	ticket, err := qtx.GetTicket(ctx, db.GetTicketParams{
		ID: ticketID,
		StatusID: sql.NullInt32{
			Int32: int32(statusID),
			Valid: true,
		},
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
			return echo.NewHTTPError(http.StatusNotFound, "ticket not found")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	isChanged := false
	ticketParam := db.UpdateTicketParams{
		StatusID:    ticket.StatusID,
		Title:       ticket.Title,
		Description: ticket.Description,
		Contact:     ticket.Contact,
		SortOrder:   ticket.SortOrder,
		ID:          ticket.ID,
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

	subctx, cancel := context.WithCancel(ctx)
	defer cancel()
	g, subctx := errgroup.WithContext(subctx)

	moveOut := body.StatusID != nil && ticket.StatusID != *body.StatusID

	if moveOut {
		if body.SortOrder == nil {
			return echo.NewHTTPError(http.StatusBadRequest, "sort order on is required on move to new status")
		}

		isChanged = true
		newStatus, err := qtx.GetStatus(ctx, db.GetStatusParams{
			ID: *body.StatusID,
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
				return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("moving failed: status %d not found", *body.StatusID))
			}

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		ticketParam.StatusID = newStatus.ID

		g.Go(func() error {
			oldSortOrder := ticket.SortOrder

			oldFriends, err := qtx.GetTicketsWithMinimumSortOrder(subctx, db.GetTicketsWithMinimumSortOrderParams{
				StatusID:           ticket.StatusID,
				SortOrder:          oldSortOrder,
				SortOrderDirection: "asc",
			})
			if err != nil {
				return err
			}

			for _, t := range oldFriends {
				if t.ID == ticket.ID {
					continue
				}

				err = qtx.UpdateTicketSortOrder(subctx, db.UpdateTicketSortOrderParams{
					SortOrder: oldSortOrder,
					ID:        t.ID,
				})
				if err != nil {
					return err
				}

				oldSortOrder--
			}

			return nil
		})
	}

	if body.SortOrder != nil {
		isChanged = true
		count, err := qtx.CountTicketByStatusID(ctx, ticket.StatusID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		var isOutOfRange bool

		if moveOut {
			isOutOfRange = *body.SortOrder > uint32(count+1)
		} else {
			isOutOfRange = *body.SortOrder > uint32(count)
		}

		if isOutOfRange {
			return echo.NewHTTPError(http.StatusBadRequest, "sort order out of range")
		}

		ticketParam.SortOrder = *body.SortOrder

		g.Go(func() error {
			friends, err := qtx.GetTicketsByStatusID(subctx, ticket.StatusID)
			if err != nil {
				return err
			}

			newSortOrder := 0
			for _, t := range friends {
				newSortOrder++
				if t.ID == ticket.ID {
					continue
				}

				err = qtx.UpdateTicketSortOrder(subctx, db.UpdateTicketSortOrderParams{
					SortOrder: uint32(newSortOrder),
					ID:        t.ID,
				})
				if err != nil {
					return err
				}
			}

			return nil
		})
	}

	if isChanged {
		err = qtx.UpdateTicket(ctx, ticketParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		err = g.Wait()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		err = tx.Commit()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	ticket, err = h.Queries.GetTicket(ctx, db.GetTicketParams{
		ID: ticketID,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, ticket)
}
