package apis

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"ticket/config"
	"ticket/pkg/db"
	"time"

	"github.com/labstack/echo/v4"
)

type GenericResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type APIConfig struct {
	Label string
	Host  string
	Port  int
}

type DBConfig struct {
	Host     string
	Name     string
	User     string
	Password string
	TimeOut  time.Duration
}

type Config struct {
	APIConfig APIConfig
	DBConfig  *DBConfig
}

type API struct {
	config       *Config
	globalConfig *config.Config
	Db           *db.Queries
	App          *echo.Echo
}

var a *API

func Start(e *echo.Echo, cf Config) *API {
	a = &API{
		config: &cf,
	}

	if cf.DBConfig != nil {
		dsname := fmt.Sprintf("%s:%s@%s?parseTime=true", cf.DBConfig.User, cf.DBConfig.Password, cf.DBConfig.Name)
		sqldb, err := sql.Open(cf.DBConfig.Host, dsname)
		if err != nil {
			log.Fatal(err)
		}

		a.Db = db.New(sqldb)
	}

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, fmt.Sprintf("%s, OK!", cf.APIConfig.Label))
	})

	fmt.Printf("Starting API %s...\n", cf.APIConfig.Label)

	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%d", cf.APIConfig.Host, cf.APIConfig.Port)))

	return a
}

func (a *API) GetAPIConfig() APIConfig {
	return a.config.APIConfig
}

func (a *API) GetDBConfig() DBConfig {
	return *a.config.DBConfig
}

func (a *API) GetGlobalConfig() config.Config {
	return *a.globalConfig
}

func GetAPI() *API {
	return a
}
