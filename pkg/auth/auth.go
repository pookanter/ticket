package auth

type AuthConfig struct {
	PrivateKey         string
	AccessTokenExpire  int
	RefreshTokenExpire int
}

type Auth struct {
	config AuthConfig
}

type Configurer interface {
	PrivateKey() string
	AccessTokenExpire() int
	RefreshTokenExpire() int
}

func New(c Configurer) *Auth {
	return &Auth{
		config: AuthConfig{
			PrivateKey:         c.PrivateKey(),
			AccessTokenExpire:  c.AccessTokenExpire(),
			RefreshTokenExpire: c.RefreshTokenExpire(),
		},
	}
}
