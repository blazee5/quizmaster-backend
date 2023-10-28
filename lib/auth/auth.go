package auth

import (
	"crypto/sha256"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const (
	signingKey = "7!qK&5pTg#r*Fz$@9W"
	tokenTTL   = 72 * time.Hour
	salt       = "Xy@6#L9*Z!q2r$Pc"
)

type tokenClaims struct {
	jwt.RegisteredClaims
	UserId int `json:"user_id"`
}

func GenerateToken(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		userId,
	})

	return token.SignedString([]byte(signingKey))
}

func ParseToken(token string) (int, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := parsedToken.Claims.(*tokenClaims)
	if !ok {
		return 0, err
	}

	return claims.UserId, nil
}

func GenerateHashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
