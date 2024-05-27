package apikit

import (
	"ticket/config"
)

type Configuration struct {
	api    APIConfig
	db     DBConfig
	global config.Config
	certs  Certs
}

func (cf *Configuration) API() APIConfig {
	return cf.api
}

func (cf *Configuration) DB() DBConfig {
	return cf.db
}

func (cf *Configuration) GLobal() config.Config {
	return cf.global
}

func (cf *Configuration) Certs() Certs {
	return cf.certs
}

func (cf *Configuration) PrivateKey() string {
	return cf.certs.PrivateKey
}

func (cf *Configuration) PublicKey() string {
	return cf.certs.PublicKey
}

func (cf *Configuration) AccessTokenExpire() int {
	return cf.global.AccessTokenExpire
}

func (cf *Configuration) RefreshTokenExpire() int {
	return cf.global.RefreshTokenExpire
}
