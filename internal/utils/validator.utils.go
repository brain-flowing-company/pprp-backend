package utils

import (
	"net/mail"

	"github.com/brain-flowing-company/pprp-backend/internal/enums"
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

func IsValidRating(rating string) bool {
	switch rating {
	case string(enums.Rating0), string(enums.Rating0_5), string(enums.Rating1), string(enums.Rating1_5), string(enums.Rating2), string(enums.Rating2_5), string(enums.Rating3), string(enums.Rating3_5), string(enums.Rating4), string(enums.Rating4_5), string(enums.Rating5):
		return true
	default:
		return false
	}
}
