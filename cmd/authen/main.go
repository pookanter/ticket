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

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Authen!")
	})

	fmt.Println("Starting Authen service..sssssss.")
	authen.UseRouter(e)

	e.Logger.Fatal(e.Start(os.Getenv("HOST") + ":" + os.Getenv("PORT")))
}
