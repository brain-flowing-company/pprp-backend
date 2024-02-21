package models

import "time"

type EmailVerificationData struct {
	Email     string    `json:"email"`
	Code      string    `json:"code"`
	ExpiredAt time.Time `json:"expire_at"`
}

type EmailVerificationRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func (e EmailVerificationData) TableName() string {
	return "email_verification_data"
}
