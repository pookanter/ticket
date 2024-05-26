package main

import (
	"fmt"
	"net/http"
	"os"
	"ticket/api/authen"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Authen, OK!")
	})

	fmt.Println("Starting Authen service...")
	authen.UseRouter(e)

	e.Logger.Fatal(e.Start(os.Getenv("HOST") + ":" + os.Getenv("PORT")))
}
