package handler

import (
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
