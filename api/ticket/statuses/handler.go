package statuses

import (
	"database/sql"
	"net/http"
	"strconv"
	"ticket/pkg/apikit"
	"ticket/pkg/auth"
	"ticket/pkg/db"

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
			return echo.NewHTTPError(http.StatusBadRequest, "board not found")
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
