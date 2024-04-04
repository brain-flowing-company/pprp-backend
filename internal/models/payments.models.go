package models

import (
	"time"

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

type HistoryResponse struct {
	PaymentID        uuid.UUID             `json:"payment_id"`
	UserID           uuid.UUID             `json:"user_id"`
	Price            float64               `json:"price"`
	IsSuccess        bool                  `json:"is_success"`
	Name             string                `json:"name"`
	AgreementID      uuid.UUID             `json:"agreement_id"`
	PaymentMethod    enums.PaymentMethods  `json:"payment_method"`
	AgreementType    enums.AgreementTypes  `json:"agreement_type"`
	PropertyID       uuid.UUID             `json:"property_id"`
	OwnerUserID      uuid.UUID             `json:"owner_user_id"`
	DwellerUserID    uuid.UUID             `json:"dweller_user_id"`
	AgreementDate    time.Time             `json:"agreement_date"`
	Status           enums.AgreementStatus `json:"status"`
	DepositAmount    float64               `json:"deposit_amount"`
	PaymentPerMonth  float64               `json:"payment_per_month"`
	PaymentDuration  int                   `json:"payment_duration"`
	TotalPayment     float64               `json:"total_payment"`
	CancelledMessage string                `json:"cancelled_message"`
	CommonModels
}

// type Agreements struct {
// 	AgreementId   uuid.UUID `json:"agreement_id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
// 	AgreementType enums.AgreementTypes `json:"agreement_type" example:"SELLING"`
// 	PropertyId    uuid.UUID `json:"property_id" example:"00000000-0000-0000-0000-000000000000"`
// 	OwnerUserId   uuid.UUID `json:"owner_user_id" example:"00000000-0000-0000-0000-000000000000"`
// 	DwellerUserId uuid.UUID `json:"dweller_user_id" example:"00000000-0000-0000-0000-000000000000"`
// 	AgreementDate time.Time `json:"agreement_date" example:"2021-01-01T00:00:00Z"`
// 	Status enums.AgreementStatus `json:"status" example:"AWAITING_DEPOSIT"`
// 	DepositAmount float64 `json:"deposit_amount" example:"1000000"`
// 	PaymentPerMonth float64 `json:"payment_per_month" example:"1000000"`
// 	PaymentDuration int `json:"payment_duration" example:"12"`
// 	TotalPayment float64 `json:"total_payment" example:"12000000"`
// 	CancelledMessage string `json:"cancelled_message" example:"This is cancelled message."`
// 	CommonModels
// }
