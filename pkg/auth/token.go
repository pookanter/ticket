package auth

import (
	"fmt"
	"ticket/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenPayload struct {
	UserID int `json:"user_id"`
}

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func GenerateTokens(tokenPayload TokenPayload) (Tokens, error) {
	cf := config.GetConfig()
	accessToken, err := GenerateTokenString(tokenPayload, cf.AccessTokenExpire)
	if err != nil {
		return Tokens{}, err
	}

	refreshToken, err := GenerateTokenString(tokenPayload, cf.RefreshTokenExpire)
	if err != nil {
		return Tokens{}, err
	}

	return Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func GenerateTokenString(tokenPayload TokenPayload, unixDuration int) (string, error) {
	cf := config.GetConfig()
	claims := Claims{
		tokenPayload.UserID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(unixDuration))),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(cf.PrivateKey)
}

func ParseToken(tokenString string) (*Claims, error) {
	cf := config.GetConfig()
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwt.ParseECPublicKeyFromPEM([]byte(cf.PublicKey))
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid claims")
}
