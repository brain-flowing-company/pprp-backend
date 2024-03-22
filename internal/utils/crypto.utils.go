package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(plainPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func ComparePassword(hashedPassword string, password string) bool {
	// return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		fmt.Println("uilt error: ", err)
		return false
	}
	return true
}
