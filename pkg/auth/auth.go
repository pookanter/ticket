package auth

type AuthConfig struct {
	RSAKey             string
	AccessTokenExpire  int
	RefreshTokenExpire int
}

type Auth struct {
	config AuthConfig
}

func New(config AuthConfig) *Auth {
	return &Auth{config}
}
