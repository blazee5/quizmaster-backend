package auth

import (
	"crypto/sha256"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

const (
	signingKey = "7!qK&5pTg#r*Fz$@9W"
	tokenTTL   = time.Hour * 72
	salt       = "Xy@6#L9*Z!q2r$Pc"
)

type TokenClaims struct {
	jwt.RegisteredClaims
	UserID int `json:"user_id"`
	RoleID int `json:"role_id"`
}

func GenerateToken(userID, roleID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		userID,
		roleID,
	})

	return token.SignedString([]byte(signingKey))
}

func ParseToken(token string) (int, int, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	})

	if err != nil {
		return 0, 0, err
	}

	claims, ok := parsedToken.Claims.(*TokenClaims)
	if !ok {
		return 0, 0, err
	}

	return claims.UserID, claims.RoleID, nil
}

func GenerateHashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func GenerateNewTokenCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(72 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	}
}
