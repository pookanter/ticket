package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"ticket/config"

	"github.com/labstack/echo/v4"
	"github.com/valyala/fasthttp"
)

func main() {
	config.Initialize()

	cf := config.GetConfig()
	e := echo.New()

	e.Any("/:service/:path", handler)

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Gateway, OK!")
	})

	fmt.Println("Starting Gateway service...")

	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%d", cf.Services.Gateway.Host, cf.Services.Gateway.Port)))
}

func handler(c echo.Context) error {
	cf := config.GetConfig()
	req := fasthttp.AcquireRequest()
	req.Header.SetContentType("application/json")
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	service := c.Param("service")
	path := c.Param("path")

	fmt.Println("service:", service == "authen")

	url := ""
	if service == "authen" {
		url = cf.Services.Authen.URL
	}

	if url != "" {
		req.Header.SetContentType("application/json; charset=UTF-8")
		req.Header.SetMethod(c.Request().Method)
		bodyBytes, err := io.ReadAll(c.Request().Body)
		defer c.Request().Body.Close()
		if err != nil {
			log.Fatal(err)
		}

		req.SetBody(bodyBytes)

		req.SetRequestURI(fmt.Sprintf("%s/%s", url, path))

		err = fasthttp.Do(req, res)
		if err != nil {
			log.Println("fashttp failed to do request")
			log.Print(err)

			return c.String(http.StatusInternalServerError, "Internal Server Error")
		}

		return c.String(res.StatusCode(), string(res.Body()))
	}

	return c.String(http.StatusNotFound, "Not Found")
}
