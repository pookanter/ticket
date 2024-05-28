package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"ticket/config"
	"ticket/pkg/apikit"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/valyala/fasthttp"
)

func main() {
	cf, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}

	apikit.NewAPI(apikit.WithAPI(apikit.APIConfig{
		Label: "Gateway",
		Host:  cf.Services.Gateway.Host,
		Port:  cf.Services.Gateway.Port,
	}), apikit.WithGlobal(cf)).Use(middleware.CORS()).UseRouter(func(api *apikit.API) {
		api.App.Any("/:service/:path", func(c echo.Context) error {
			cf := api.Config.GLobal()
			req := fasthttp.AcquireRequest()
			req.Header.SetContentType("application/json")
			res := fasthttp.AcquireResponse()
			defer fasthttp.ReleaseRequest(req)
			defer fasthttp.ReleaseResponse(res)

			service := c.Param("service")
			path := c.Param("path")

			url := ""
			if service == "authen-service" {
				url = cf.Services.Authen.URL
			} else if service == "ticket-service" {
				url = cf.Services.Ticket.URL
			}

			if url != "" {
				req.Header.SetMethod(c.Request().Method)
				for name, values := range c.Request().Header {
					for _, value := range values {
						req.Header.Add(name, value)
					}
				}
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

				fmt.Println("res.StatusCode():", res.StatusCode())

				return c.JSONBlob(res.StatusCode(), res.Body())
			}

			return c.JSON(http.StatusNotFound, "Not Found")
		})

	}).Start()
}
