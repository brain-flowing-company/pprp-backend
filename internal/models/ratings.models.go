package models

import (
	"github.com/brain-flowing-company/pprp-backend/internal/enums"
	"github.com/google/uuid"
)

// CREATE TABLE review(
//     review_id UUID PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
//     property_id UUID REFERENCES properties(property_id) NOT NULL,
//     dweller_user_id UUID REFERENCES users(user_id) NOT NULL,
//     rating Rating NOT NULL,
//     review TEXT DEFAULT NULL,
//     created_at TIMESTAMP(0) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
//     updated_at TIMESTAMP(0) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
//     deleted_at TIMESTAMP(0) WITH TIME ZONE DEFAULT NULL
// );

type Review struct {
	ReviewId      uuid.UUID    `json:"review_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	PropertyId    uuid.UUID    `json:"property_id" gorm:"type:uuid;not null"`
	DwellerUserId uuid.UUID    `json:"dweller_user_id" gorm:"type:uuid;not null"`
	Rating        enums.Rating `json:"rating" gorm:"type:rating;not null"`
	Review        string       `json:"review" gorm:"type:text;default:null"`
	CommonModels
}
