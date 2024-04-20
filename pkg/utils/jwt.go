package utils

import (
	"errors"
	"time"

	"github.com/Hank-Kuo/chat-app/config"
	"github.com/Hank-Kuo/chat-app/pkg/customError"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Username string
	Email    string
	UserID   string
	jwt.RegisteredClaims
}

func GetJwt(cfg *config.Config, userId, username, email string) (string, error) {
	now := time.Now().In(cfg.Server.Location)
	t := time.Duration(cfg.Server.AccessJwtExpireTime) * time.Second

	claims := &Claims{
		Username: username,
		Email:    email,
		UserID:   userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(now.Add(t).Unix(), 0)),
			IssuedAt:  jwt.NewNumericDate(time.Unix(now.Unix(), 0)),
			Issuer:    "hank-kuo",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Server.AccessJwtSecret))
}

func ValidJwt(cfg *config.Config, tokenStr string) (*Claims, error) {

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.Server.AccessJwtSecret), nil
	})

	switch {
	case token.Valid:
		if claims, ok := token.Claims.(*Claims); ok {
			return claims, nil
		}
	case errors.Is(err, jwt.ErrTokenMalformed) || errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return nil, customError.ErrInvalidJWTClaims
	case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
		return nil, customError.ErrExpiredJWTError
	}

	return nil, customError.ErrInvalidJWTClaims
}
