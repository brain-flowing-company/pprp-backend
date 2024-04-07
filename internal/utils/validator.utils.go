package utils

import (
	"net/mail"

	"github.com/google/uuid"
)

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsValidPassword(password string) bool {
	hasNumber := false
	for _, char := range password {
		if char >= '0' && char <= '9' {
			hasNumber = true
			break
		}
	}
	return len(password) >= 8 && hasNumber
}

func IsValidEmailVerificationCode(code string) bool {
	return len(code) == 10 && code[:4] == "SCK-"
}

func IsValidRating(rating int8) bool {
	if rating < 0 || rating > 5 {
		return false
	}
	return true
}
