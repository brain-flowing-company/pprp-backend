package models

import (
	"github.com/brain-flowing-company/pprp-backend/internal/enums"
	"github.com/google/uuid"
)

type Payments struct {
	PaymentId     uuid.UUID            `json:"payment_id" `
	UserId        uuid.UUID            `json:"user_id" `
	Price         float64              `json:"price" `
	IsSuccess     bool                 `json:"is_success"`
	Name          string               `json:"name"`
	AgreementId   uuid.UUID            `json:"agreement_id" `
	PaymentMethod enums.PaymentMethods `json:"payment_method" `
	CommonModels
}

type MyPaymentsResponse struct {
	Payments []Payments `json:"payments"`
}
