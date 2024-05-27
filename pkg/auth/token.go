package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenPayload struct {
	UserID uint64 `json:"user_id"`
}

type Claims struct {
	UserID uint64 `json:"user_id"`
	jwt.RegisteredClaims
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type GenerateTokensConfig struct {
	AccessTokenExpire  int
	RefreshTokenExpire int
}

func (a *Auth) GenerateTokens(tokenPayload TokenPayload, cf GenerateTokensConfig) (Tokens, error) {
	accessToken, err := a.GenerateTokenString(tokenPayload, cf.AccessTokenExpire)
	if err != nil {
		return Tokens{}, err
	}

	refreshToken, err := a.GenerateTokenString(tokenPayload, cf.RefreshTokenExpire)
	if err != nil {
		return Tokens{}, err
	}

	return Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a *Auth) GenerateTokenString(tokenPayload TokenPayload, unixDuration int) (string, error) {
	claims := Claims{
		tokenPayload.UserID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(unixDuration))),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(a.config.RSAKey)
}

func (a *Auth) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwt.ParseECPublicKeyFromPEM([]byte(a.config.RSAKey))
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid claims")
}
