package authen

import (
	"context"
	"database/sql"
	"net/http"
	"ticket/pkg/apikit"
	"ticket/pkg/auth"
	"ticket/pkg/db"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	App      *echo.Echo
	Auth     *auth.Auth
	DB       *db.Queries
	DBConfig apikit.DBConfig
}

func NewHandler(api *apikit.API) *Handler {
	return &Handler{
		App: api.App,
		Auth: auth.New(auth.AuthConfig{
			RSAKey:             api.GetCerts().PrivateKey,
			AccessTokenExpire:  api.GetGlobalConfig().AccessTokenExpire,
			RefreshTokenExpire: api.GetGlobalConfig().RefreshTokenExpire,
		}),
		DB:       api.DB,
		DBConfig: api.GetDBConfig(),
	}
}

func (h *Handler) SignIn(c echo.Context) error {
	var body struct {
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

	ctx, cancel := context.WithTimeout(c.Request().Context(), h.DBConfig.TimeOut)
	defer cancel()

	user, err := h.DB.FindUserByEmail(ctx, sql.NullString{String: body.Email, Valid: true})
	if err != nil {
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

	ctx, cancel := context.WithTimeout(c.Request().Context(), h.DBConfig.TimeOut)
	defer cancel()

	user, err := h.DB.FindUserByEmail(ctx, sql.NullString{String: body.Email, Valid: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if user.ID != 0 {
		return echo.NewHTTPError(http.StatusConflict, "user already exists")
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
