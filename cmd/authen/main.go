package main

import (
	"fmt"
	"net/http"
	"ticket/api/authen"
	"ticket/config"

	"github.com/labstack/echo/v4"
)

func main() {
	config.Initialize()

	cf := config.GetConfig()
	e := echo.New()

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Authen, OK!")
	})

	fmt.Println("Starting Authen service...")
	authen.UseRouter(e)

	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%d", cf.Services.Authen.Host, cf.Services.Authen.Port)))
}
