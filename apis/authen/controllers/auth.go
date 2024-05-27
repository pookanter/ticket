package controllers

import (
	"context"
	"database/sql"
	"net/http"
	"ticket/apis"
	"ticket/pkg/auth"
	"ticket/pkg/db"

	"github.com/labstack/echo/v4"
)

type AuthController struct {
	App      *echo.Echo
	Auth     *auth.Auth
	DB       *db.Queries
	DBConfig apis.DBConfig
}

func NewAuthController(g *echo.Group, api *apis.API) *AuthController {
	gcf := api.GetGlobalConfig()
	ctrl := &AuthController{
		App: api.App,
		Auth: auth.New(auth.AuthConfig{
			RSAKey:             api.GetPrivateKey(),
			AccessTokenExpire:  gcf.AccessTokenExpire,
			RefreshTokenExpire: gcf.RefreshTokenExpire,
		}),
		DB:       api.Db,
		DBConfig: api.GetDBConfig(),
	}

	g.POST("/sign-in", ctrl.SignIn())

	g.POST("/sign-up", ctrl.SignUp())

	g.POST("/refresh-token", ctrl.RefreshToken())

	return ctrl
}

func (ctrl *AuthController) SignIn() echo.HandlerFunc {
	return func(c echo.Context) error {
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

		ctx, cancel := context.WithTimeout(c.Request().Context(), ctrl.DBConfig.TimeOut)
		defer cancel()

		user, err := ctrl.DB.FindUserByEmail(ctx, sql.NullString{String: body.Email, Valid: true})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if user.ID == 0 {
			return echo.NewHTTPError(http.StatusNotFound, "user not found")
		}

		err = ctrl.Auth.ComparePassword(user.Password.String, body.Password)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid password")
		}

		payload := auth.TokenPayload{
			UserID: user.ID,
		}

		tokens, err := ctrl.Auth.GenerateTokens(payload)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, tokens)
	}
}

func (ctrl *AuthController) SignUp() echo.HandlerFunc {
	return func(c echo.Context) error {
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

		ctx, cancel := context.WithTimeout(c.Request().Context(), ctrl.DBConfig.TimeOut)
		defer cancel()

		user, err := ctrl.DB.FindUserByEmail(ctx, sql.NullString{String: body.Email, Valid: true})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if user.ID != 0 {
			return echo.NewHTTPError(http.StatusConflict, "user already exists")
		}

		hash, err := ctrl.Auth.HashPassword(body.Password)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		err = ctrl.DB.CreateUser(ctx, db.CreateUserParams{
			Name:     sql.NullString{String: body.Name, Valid: true},
			Lastname: sql.NullString{String: body.Lastname, Valid: true},
			Email:    sql.NullString{String: body.Email, Valid: true},
			Password: sql.NullString{String: hash, Valid: true},
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusCreated, apis.GenericResponse{
			Error:   false,
			Message: "user created",
		})
	}
}

func (ctrl *AuthController) RefreshToken() echo.HandlerFunc {
	return func(c echo.Context) error {
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

		claims, err := ctrl.Auth.ParseToken(body.RefreshToken)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
		}

		payload := auth.TokenPayload{
			UserID: claims.UserID,
		}

		tokens, err := ctrl.Auth.GenerateTokens(payload)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, tokens)
	}
}
