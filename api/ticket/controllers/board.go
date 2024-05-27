package controllers

import (
	"ticket/pkg/apikit"
	"ticket/pkg/db"

	"github.com/labstack/echo/v4"
)

type BoardController struct {
	DB       *db.Queries
	DBConfig apikit.DBConfig
}

func NewBoardController(g *echo.Group, api *apikit.API) *BoardController {
	ctrl := &BoardController{
		DB:       api.DB,
		DBConfig: api.GetDBConfig(),
	}

	return ctrl
}

func (ctrl *BoardController) Get() echo.HandlerFunc {
	return func(c echo.Context) error {

		return c.JSON(200, "Get board")
	}
}
