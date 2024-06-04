package boards

import (
	"database/sql"
	"net/http"
	"strconv"
	"ticket/pkg/apikit"
	"ticket/pkg/auth"
	"ticket/pkg/db"
	"ticket/pkg/dbutil"

	"github.com/guregu/null/v5"
	"github.com/labstack/echo/v4"
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

func (h *Handler) GetBoards(c echo.Context) error {
	claims := c.Get("claims").(*auth.Claims)
	boards, err := h.Queries.GetBoardsByUserID(c.Request().Context(), claims.UserID)
	if err != nil && err != sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, boards)
}

func (h *Handler) GetBoardByID(c echo.Context) error {
	claims := c.Get("claims").(*auth.Claims)

	boardID, err := strconv.ParseUint(c.Param("board_id"), 10, 32)
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

	statuses, err := h.Queries.GetStatuses(ctx, db.GetStatusesParams{
		BoardID:            sql.NullInt32{Int32: int32(board.ID), Valid: true},
		SortOrderDirection: null.StringFrom("asc"),
	})
	if err != nil && err != sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var statusIDs []null.Int32
	for _, s := range statuses {
		statusIDs = append(statusIDs, null.NewInt32(int32(s.ID), true))
	}

	tickets, err := h.Queries.GetTickets(ctx, db.GetTicketsParams{
		StatusIds:          statusIDs,
		SortOrderDirection: null.StringFrom("asc"),
	})
	if err != nil && err != sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, dbutil.NewBoardWithRelated(board, statuses, tickets))
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

	ctx := c.Request().Context()

	tx, err := h.DB.Begin()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer tx.Rollback()
	qtx := h.Queries.WithTx(tx)

	user, err := qtx.FindUserByID(ctx, claims.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "user not found")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	count, err := qtx.CountBoardByUserID(ctx, user.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = qtx.CreateBoard(ctx, db.CreateBoardParams{
		UserID:    user.ID,
		Title:     null.NewString(body.Title, true),
		SortOrder: uint32(count),
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	boardID, err := qtx.GetLastInsertBoardID(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	board, err := qtx.GetBoard(ctx, db.GetBoardParams{
		ID: uint32(boardID),
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, board)
}
