package boards

import (
	"database/sql"
	"net/http"
	"strconv"
	"ticket/pkg/apikit"
	"ticket/pkg/auth"
	"ticket/pkg/db"
	"ticket/pkg/usecases"

	"github.com/guregu/null"
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

	dbboard, err := h.Queries.GetBoard(ctx, db.GetBoardParams{
		ID:     uint32(boardID),
		UserID: claims.UserID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "board not found")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	board := usecases.NewBoardFullDetail(dbboard)

	statuses, err := h.Queries.GetStatusesByBoardID(ctx, dbboard.ID)
	if err != nil && err != sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	for i, s := range statuses {
		status := usecases.NewStatusFullDetail(s)

		tickets, err := h.Queries.GetTicketsByStatusID(ctx, s.ID)
		if err != nil && err != sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		for _, t := range tickets {
			t.SortOrder = uint32(i)
			status.Tickets = append(status.Tickets, usecases.NewTicket(t))
		}

		board.Statuses = append(board.Statuses, status)
	}

	return c.JSON(http.StatusOK, board)
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

	board, err := qtx.GetLastInsertBoardViewByUserID(ctx, user.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, board)
}
