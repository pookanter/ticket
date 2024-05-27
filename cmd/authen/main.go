package main

import (
	"ticket/apis"
	"ticket/apis/authen"
	"ticket/config"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func main() {
	config.Initialize()
	cf := config.GetConfig()

	authen.Router(apis.Start(echo.New(), apis.Config{
		APIConfig: apis.APIConfig{
			Label: "Authen",
			Host:  cf.Services.Authen.Host,
			Port:  cf.Services.Authen.Port,
		},
		DBConfig: &apis.DBConfig{
			Host:     cf.Services.Database.Host,
			Name:     cf.Services.Database.Dbname,
			User:     cf.Services.Database.User,
			Password: cf.Services.Database.Password,
			TimeOut:  5 * time.Second,
		},
	}))
}
