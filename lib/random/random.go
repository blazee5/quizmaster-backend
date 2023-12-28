package random

import (
	"math/rand"
	"time"
)

func GenerateVerificationCode(length int) string {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, length)
	for i := range code {
		code[i] = characters[random.Intn(len(characters))]
	}
	return string(code)
}
