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

func (a *Auth) GenerateTokens(tokenPayload TokenPayload) (Tokens, error) {
	accessToken, err := a.GenerateTokenString(tokenPayload, a.config.AccessTokenExpire)
	fmt.Println(err)
	if err != nil {
		return Tokens{}, err
	}

	refreshToken, err := a.GenerateTokenString(tokenPayload, a.config.RefreshTokenExpire)
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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(unixDuration))),
		},
	}

	rsaPrivatekey, err := jwt.ParseRSAPrivateKeyFromPEM(a.config.PrivateKey)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(rsaPrivatekey)
}

func (a *Auth) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return jwt.ParseRSAPublicKeyFromPEM(a.config.PublicKey)
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}
