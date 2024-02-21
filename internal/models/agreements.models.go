package models

import (
	"time"

	"github.com/google/uuid"
)

type Agreements struct {
	AgreementID   uuid.UUID `json:"agreement_id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	PropertyID    uuid.UUID `json:"property_id"`
	OwnerUserID   uuid.UUID `json:"owner_user_id"`
	DwellerUserID uuid.UUID `json:"dweller_user_id"`
	CommonModels
}

type CreatingAgreements struct {
	PropertyID    uuid.UUID `json:"property_id"`
	OwnerUserID   uuid.UUID `json:"owner_user_id"`
	DwellerUserID uuid.UUID `json:"dweller_user_id"`
	AgreementDate time.Time `json:"agreement_date"`
}

func (a Agreements) TableName() string {
	return "agreements"
}
