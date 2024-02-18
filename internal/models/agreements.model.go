package models

import (
	"time"

	"github.com/google/uuid"
)

type Agreement struct {
	AgreementID   uuid.UUID  `json:"agreement_id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	PropertyID    uuid.UUID  `json:"property_id"`
	OwnerUserID   uuid.UUID  `json:"owner_user_id"`
	DwellerUserID uuid.UUID  `json:"dweller_user_id"`
	AgreementDate time.Time  `json:"agreement_date"`
	CreatedAt     *time.Time `gorm:"autoCreateTime"`
	UpdatedAt     *time.Time `gorm:"autoUpdateTime"`
	DeletedAt     *time.Time `gorm:"default:null"`
}

type CreatingAgreement struct {
	PropertyID    uuid.UUID `json:"property_id"`
	OwnerUserID   uuid.UUID `json:"owner_user_id"`
	DwellerUserID uuid.UUID `json:"dweller_user_id"`
	AgreementDate time.Time `json:"agreement_date"`
}

func (a Agreement) TableName() string {
	return "agreements"
}
