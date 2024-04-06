package models

import (
	"github.com/google/uuid"
)

type Reviews struct {
	ReviewId      uuid.UUID `json:"review_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	PropertyId    uuid.UUID `json:"property_id" gorm:"type:uuid;not null"`
	DwellerUserId uuid.UUID `json:"dweller_user_id" gorm:"type:uuid;not null"`
	Rating        int8      `json:"rating" gorm:"type:rating;not null"`
	Review        string    `json:"review" gorm:"type:text;default:null"`
	CommonModels
}

type RatingResponse struct {
	ReviewId      uuid.UUID `json:"review_id"`
	PropertyId    uuid.UUID `json:"property_id"`
	DwellerUserId uuid.UUID `json:"dweller_user_id"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Rating        int8      `json:"rating"`
	Review        string    `json:"review"`
	CommonModels
}

type UpdateRatingStatus struct {
	Rating int8   `json:"rating"`
	Review string `json:"review"`
}
