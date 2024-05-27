package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"ticket/apis"
	"ticket/config"

	"github.com/labstack/echo/v4"
	"github.com/valyala/fasthttp"
)

func main() {
	config.Initialize()

	cf := config.GetConfig()
	e := echo.New()

	e.Any("/:service/:path", handler)

	apis.Start(e, apis.Config{
		APIConfig: apis.APIConfig{
			Label: "Gateway",
			Host:  cf.Services.Authen.Host,
			Port:  cf.Services.Authen.Port,
		},
	})

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
	if service == "authen-service" {
		url = cf.Services.Authen.URL
	} else if service == "ticket-service" {
		url = cf.Services.Ticket.URL
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

			return c.JSON(http.StatusInternalServerError, "Internal Server Error")
		}

		if res.StatusCode() == http.StatusInternalServerError {
			log.Printf("Internal Server Error: %s", string(res.Body()))

			return c.JSON(http.StatusInternalServerError, "Internal Server Error")
		}

		return c.JSON(res.StatusCode(), string(res.Body()))
	}

	return c.JSON(http.StatusNotFound, "Not Found")
}
