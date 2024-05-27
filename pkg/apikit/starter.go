package apikit

import (
	"context"
	"fmt"
	"net/http"
	"ticket/config"
	"ticket/pkg/db"

	"github.com/labstack/echo/v4"
)

type GenericResponse[T any] struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

type APIConfig struct {
	Label string
	Host  string
	Port  int
}

type Certs struct {
	PrivateKey string
	PublicKey  string
}

type Configuration struct {
	api    APIConfig
	db     DBConfig
	global config.Config
	certs  Certs
}

type Router func(api *API)

type API struct {
	cf  Configuration
	rt  []Router
	DB  *db.Queries
	App *echo.Echo
}

func NewAPI(configs ...Config) *API {
	api := &API{}

	for _, c := range configs {
		c(api)
	}

	return api
}

func (api *API) Use(routers ...Router) *API {
	api.rt = append(api.rt, routers...)

	return api
}

func (api *API) Start(e *echo.Echo) {
	if isDBConfigValid(api.cf.db) {
		go func() {
			dbcf := api.cf.db
			ctx, cancel := context.WithTimeout(context.Background(), dbcf.TimeOut)
			defer cancel()

			db, err := ConnectDBContext(ctx, dbcf)
			if err != nil {
				fmt.Printf("Error connecting to database: %v\n", err.Error())

				return
			}

			fmt.Printf("Connected to database %s\n", dbcf.Name)

			api.DB = db
		}()
	}

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, fmt.Sprintf("%s, OK!", api.cf.api.Label))
	})

	for _, r := range api.rt {
		r(api)
	}

	fmt.Printf("Starting API %s...\n", api.cf.api.Label)

	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%d", api.cf.api.Host, api.cf.api.Port)))
}

func isDBConfigValid(dbcf DBConfig) bool {
	return dbcf.Host != "" && dbcf.Name != "" && dbcf.User != "" && dbcf.Password != ""
}

func (api *API) GetAPIConfig() APIConfig {
	return api.cf.api
}

func (api *API) GetDBConfig() DBConfig {
	return api.cf.db
}

func (api *API) GetGlobalConfig() config.Config {
	return api.cf.global
}

func (api *API) GetCerts() Certs {
	return api.cf.certs
}
