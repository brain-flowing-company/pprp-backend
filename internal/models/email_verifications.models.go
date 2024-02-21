package models

import "time"

type EmailVerificationCodes struct {
	Email     string    `json:"email"`
	Code      string    `json:"code"`
	ExpiredAt time.Time `json:"expire_at"`
}

type EmailVerificationRequests struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func (e EmailVerificationCodes) TableName() string {
	return "email_verification_codes"
}
