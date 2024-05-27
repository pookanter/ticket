package auth

type AuthConfig struct {
	PrivateKey         []byte
	PublicKey          []byte
	AccessTokenExpire  int
	RefreshTokenExpire int
}

type Auth struct {
	config AuthConfig
}

type Configurer interface {
	PrivateKey() []byte
	PublicKey() []byte
	AccessTokenExpire() int
	RefreshTokenExpire() int
}

func New(c Configurer) *Auth {
	return &Auth{
		config: AuthConfig{
			PrivateKey:         c.PrivateKey(),
			PublicKey:          c.PublicKey(),
			AccessTokenExpire:  c.AccessTokenExpire(),
			RefreshTokenExpire: c.RefreshTokenExpire(),
		},
	}
}
