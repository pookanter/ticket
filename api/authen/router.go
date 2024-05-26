package authen

import (
	"ticket/api/authen/internal/handler"

	"github.com/labstack/echo/v4"
)

type GenericResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func UseRouter(e *echo.Echo) {
	g := e.Group("/authen")
	handler.Index(g)
}
