package authen

import (
	"ticket/api/authen/internal/handler"

	"github.com/labstack/echo/v4"
)

type AuthenAPI struct {
}

func New() *AuthenAPI {
	return &AuthenAPI{}
}

func (r *AuthenAPI) Router(e *echo.Echo) {
	g := e.Group("/authen")
	handler.Index(g)
}
