package authorize

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"ticket/pkg/apikit"
	"ticket/pkg/auth"
	"ticket/pkg/db"
	"time"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	DB        *db.Queries
	DBTimeOut time.Duration
	Auth      *auth.Auth
}

func New(api *apikit.API) *Handler {
	return &Handler{
		DB:        db.New(api.DB),
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

	user, err := h.DB.FindUserByEmail(ctx, sql.NullString{String: body.Email, Valid: true})
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

	user, err := h.DB.FindUserByEmail(ctx, sql.NullString{String: body.Email, Valid: true})
	if err != nil && err != sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	fmt.Printf("user: %+v\n", user)

	if user.ID != 0 {
		return echo.NewHTTPError(http.StatusConflict, "email already exists")
	}

	hash, err := h.Auth.HashPassword(body.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = h.DB.CreateUser(ctx, db.CreateUserParams{
		Name:     sql.NullString{String: body.Name, Valid: true},
		Lastname: sql.NullString{String: body.Lastname, Valid: true},
		Email:    sql.NullString{String: body.Email, Valid: true},
		Password: sql.NullString{String: hash, Valid: true},
	})
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
