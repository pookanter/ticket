package auth

type AuthConfig struct {
	RSAKey string
}

type Auth struct {
	config AuthConfig
}

func New(config AuthConfig) *Auth {
	return &Auth{config}
}
