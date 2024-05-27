package apikit

import (
	"ticket/config"
	"time"
)

type Option func(*API)

type APIConfig struct {
	Label string
	Host  string
	Port  int
}

func WithAPI(c APIConfig) Option {
	return func(a *API) {
		a.Config.api = c
	}
}

type DBConfig struct {
	Host     string
	Name     string
	User     string
	Password string
	TimeOut  time.Duration
}

func WithDB(c DBConfig) Option {
	return func(a *API) {
		a.Config.db = c
	}
}

func WithGlobal(c config.Config) Option {
	return func(a *API) {
		a.Config.global = c
	}
}

type Certs struct {
	PrivateKey []byte
	PublicKey  []byte
}

func WithCerts(c Certs) Option {
	return func(a *API) {
		a.Config.certs = c
	}
}
