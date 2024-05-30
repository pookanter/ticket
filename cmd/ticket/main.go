package main

import (
	"ticket/api/ticket"
	"ticket/config"
	"ticket/pkg/apikit"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cf, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}

	pri, err := config.ReadPrivateKey(cf)
	if err != nil {
		panic(err)
	}

	pub, err := config.ReadPublicKey(cf)
	if err != nil {
		panic(err)
	}

	apikit.NewAPI(apikit.WithAPI(apikit.APIConfig{
		Label: "Ticket",
		Host:  cf.Services.Ticket.Host,
		Port:  cf.Services.Ticket.Port,
	}), apikit.WithDB(apikit.DBConfig{
		Host:     cf.Services.Database.Host,
		Name:     cf.Services.Database.Dbname,
		User:     cf.Services.Database.User,
		Password: cf.Services.Database.Password,
		TimeOut:  5 * time.Second,
	}), apikit.WithGlobal(cf), apikit.WithCerts(apikit.Certs{
		PrivateKey: pri,
		PublicKey:  pub,
	})).UseRouter(ticket.Router).Start()
}
