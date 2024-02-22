package utils

import "math/rand"

func RandomInt(max int) int {
	return rand.Intn(max)
}

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[RandomInt(len(charset))]
	}
	return string(b)
}
