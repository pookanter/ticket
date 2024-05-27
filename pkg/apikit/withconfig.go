package apikit

import "ticket/config"

type Config func(*API)

func WithAPI(c APIConfig) Config {
	return func(a *API) {
		a.cf.api = c
	}
}

func WithDB(c DBConfig) Config {
	return func(a *API) {
		a.cf.db = c
	}
}

func WithGlobal(c config.Config) Config {
	return func(a *API) {
		a.cf.global = c
	}
}

func WithCerts(c Certs) Config {
	return func(a *API) {
		a.cf.certs = c
	}
}
