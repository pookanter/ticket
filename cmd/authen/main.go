package main

import (
	"ticket/api"
	"ticket/api/authen"
	"ticket/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func main() {
	config.Initialize()
	cf := config.GetConfig()

	api.Start(echo.New(), api.Config{
		APIConfig: api.APIConfig{
			Label: "Authen",
			Host:  cf.Services.Authen.Host,
			Port:  cf.Services.Authen.Port,
		},
		DBConfig: &api.DBConfig{
			Host:     cf.Services.Database.Host,
			Name:     cf.Services.Database.Dbname,
			User:     cf.Services.Database.User,
			Password: cf.Services.Database.Password,
		},
		Ctx: authen.New(),
	})
}
