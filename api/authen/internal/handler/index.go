package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Index(g *echo.Group) *echo.Group {
	g.POST("/sign-in", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"error":   false,
			"message": "sign-in",
		})
	})

	g.POST("/sign-up", func(c echo.Context) error {
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

		// database.Queries

		return c.JSON(200, map[string]interface{}{
			"error":   false,
			"message": "sign-up",
		})
	})

	g.POST("/refresh-token", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"error":   false,
			"message": "refresh-token",
		})
	})

	return g
}
