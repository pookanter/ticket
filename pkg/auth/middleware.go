package auth

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
)

type AuthContext struct {
	Claims *Claims
	echo.Context
}

func Middleware(c Configurer) echo.MiddlewareFunc {
	a := New(c)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			bearer := c.Request().Header.Get("Authorization")

			s := strings.Split(bearer, " ")

			if len(s) != 2 {
				return echo.ErrUnauthorized
			}

			token := s[1]

			claims, err := a.ParseToken(token)
			if err != nil || claims == nil {
				fmt.Println("err:", err)
				return echo.ErrUnauthorized
			}

			c.Set("claims", claims)

			fmt.Println("Authrized")

			return next(c)
		}
	}
}
