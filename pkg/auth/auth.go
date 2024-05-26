package auth

type AuthConfig struct {
	PrivateKey string
	PublicKey  string
}

type Auth struct {
	config AuthConfig
}

func New(config AuthConfig) *Auth {
	return &Auth{config}
}
