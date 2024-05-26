package controller

import (
	"context"
	"database/sql"
	"net/http"
	"ticket/api"
	"ticket/pkg/auth"
	"ticket/pkg/db"

	"github.com/labstack/echo/v4"
)

func Index(a *api.API) *echo.Echo {
	g := a.App
	authen := auth.New(auth.AuthConfig{
		PrivateKey: a.GetGlobalConfig().PrivateKey,
		PublicKey:  a.GetGlobalConfig().PublicKey,
	})
	cf := api.GetAPI().GetGlobalConfig()

	g.POST("/sign-in", func(c echo.Context) error {
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

		ctx, cancel := context.WithTimeout(c.Request().Context(), a.GetDBConfig().TimeOut)
		defer cancel()

		user, err := a.Db.FindUserByEmail(ctx, sql.NullString{String: body.Email, Valid: true})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if user.ID == 0 {
			return echo.NewHTTPError(http.StatusNotFound, "user not found")
		}

		err = authen.ComparePassword(user.Password.String, body.Password)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid password")
		}

		payload := auth.TokenPayload{
			UserID: user.ID,
		}

		tokens, err := authen.GenerateTokens(payload, auth.GenerateTokensConfig{
			AccessTokenExpire:  cf.AccessTokenExpire,
			RefreshTokenExpire: cf.RefreshTokenExpire,
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, tokens)
	})

	g.POST("/sign-up", func(c echo.Context) error {
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

		ctx, cancel := context.WithTimeout(c.Request().Context(), a.GetDBConfig().TimeOut)
		defer cancel()

		user, err := a.Db.FindUserByEmail(ctx, sql.NullString{String: body.Email, Valid: true})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if user.ID != 0 {
			return echo.NewHTTPError(http.StatusConflict, "user already exists")
		}

		hash, err := authen.HashPassword(body.Password)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		err = a.Db.CreateUser(ctx, db.CreateUserParams{
			Name:     sql.NullString{String: body.Name, Valid: true},
			Lastname: sql.NullString{String: body.Lastname, Valid: true},
			Email:    sql.NullString{String: body.Email, Valid: true},
			Password: sql.NullString{String: hash, Valid: true},
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusCreated, api.GenericResponse{
			Error:   false,
			Message: "user created",
		})
	})

	g.POST("/refresh-token", func(c echo.Context) error {
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

		claims, err := authen.ParseToken(body.RefreshToken)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
		}

		payload := auth.TokenPayload{
			UserID: claims.UserID,
		}

		tokens, err := authen.GenerateTokens(payload, auth.GenerateTokensConfig{
			AccessTokenExpire:  cf.AccessTokenExpire,
			RefreshTokenExpire: cf.RefreshTokenExpire,
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, tokens)
	})

	return g
}
