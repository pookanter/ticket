package auth

import (
	"strings"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		bearer := c.Request().Header["Authorization"]

		s := strings.Split(bearer[0], " ")

		if len(s) != 2 {
			return echo.ErrUnauthorized
		}

		claims, err := ParseToken(s[1])
		if err != nil || claims == nil {
			return echo.ErrUnauthorized
		}

		c.Set("claims", claims)

		return next(c)

	}
}
