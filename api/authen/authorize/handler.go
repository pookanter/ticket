package authorize

import (
	"context"
	"database/sql"
	"net/http"
	"ticket/pkg/apikit"
	"ticket/pkg/auth"
	"ticket/pkg/db"
	"time"

	"github.com/guregu/null/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/sync/errgroup"
)

type Handler struct {
	DB        *sql.DB
	Queries   *db.Queries
	DBTimeOut time.Duration
	Auth      *auth.Auth
}

func New(api *apikit.API) *Handler {
	return &Handler{
		DB:        api.DB,
		Queries:   db.New(api.DB),
		DBTimeOut: api.Config.DB().TimeOut,
		Auth:      auth.New(api.Config),
	}
}

func (h *Handler) SignIn(c echo.Context) error {
	var body struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	err := c.Bind(&body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = c.Validate(&body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), h.DBTimeOut)
	defer cancel()

	user, err := h.Queries.FindUserByEmail(ctx, null.NewString(body.Email, true))
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "user not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if user.ID == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	err = h.Auth.ComparePassword(user.Password.String, body.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid password")
	}

	payload := auth.TokenPayload{
		UserID: user.ID,
	}

	tokens, err := h.Auth.GenerateTokens(payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, tokens)

}

func (h *Handler) SignUp(c echo.Context) error {
	var body struct {
		Name     string `json:"name" validate:"required,min=3,max=100"`
		Lastname string `json:"lastname" validate:"required,min=3,max=100"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8,max=32"`
	}

	err := c.Bind(&body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = c.Validate(&body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), h.DBTimeOut)
	defer cancel()

	user, err := h.Queries.FindUserByEmail(ctx, null.NewString(body.Email, true))
	if err != nil && err != sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if user.ID != 0 {
		return echo.NewHTTPError(http.StatusConflict, "email already exists")
	}

	hash, err := h.Auth.HashPassword(body.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	tx, err := h.DB.Begin()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer tx.Rollback()
	qtx := h.Queries.WithTx(tx)

	err = qtx.CreateUser(ctx, db.CreateUserParams{
		Name:     null.NewString(body.Name, true),
		Lastname: null.NewString(body.Lastname, true),
		Email:    null.NewString(body.Email, true),
		Password: null.NewString(hash, true),
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	userID, err := qtx.GetLastInsertUserID(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = qtx.CreateBoard(ctx, db.CreateBoardParams{
		UserID: uint64(userID),
		Title:  null.NewString("My first board", true),
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	boardID, err := qtx.GetLastInsertBoardID(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	subctx, cancel := context.WithCancel(ctx)
	e, subctx := errgroup.WithContext(subctx)
	defer cancel()

	statusTitles := []string{"pending", "accepted", "resolved", "rejected"}
	for i, title := range statusTitles {
		i, title := i, title
		e.Go(func() error {
			err = qtx.CreateStatus(subctx, db.CreateStatusParams{
				BoardID:   uint32(boardID),
				Title:     null.NewString(title, true),
				SortOrder: uint32(i + 1),
			})
			if err != nil {
				cancel()

				return err
			}

			return nil
		})
	}

	err = e.Wait()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, apikit.GenericResponse[any]{
		Error:   false,
		Message: "user created",
	})
}

func (h *Handler) RefreshToken(c echo.Context) error {
	var body struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	err := c.Bind(&body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = c.Validate(&body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims, err := h.Auth.ParseToken(body.RefreshToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
	}

	payload := auth.TokenPayload{
		UserID: claims.UserID,
	}

	tokens, err := h.Auth.GenerateTokens(payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, tokens)
}
