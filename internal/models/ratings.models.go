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

type Reviews struct {
	ReviewId      uuid.UUID    `json:"review_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	PropertyId    uuid.UUID    `json:"property_id" gorm:"type:uuid;not null"`
	DwellerUserId uuid.UUID    `json:"dweller_user_id" gorm:"type:uuid;not null"`
	Rating        enums.Rating `json:"rating" gorm:"type:rating;not null"`
	Review        string       `json:"review" gorm:"type:text;default:null"`
	CommonModels
}

type RatingResponse struct {
	ReviewId      uuid.UUID    `json:"review_id"`
	PropertyId    uuid.UUID    `json:"property_id"`
	DwellerUserId uuid.UUID    `json:"dweller_user_id"`
	FirstName     string       `json:"first_name"`
	LastName      string       `json:"last_name"`
	Rating        enums.Rating `json:"rating"`
	Review        string       `json:"review"`
	CommonModels
}

// type Users struct {
// 	UserId          uuid.UUID             `json:"user_id"                      gorm:"default:uuid_generate_v4()"`
// 	RegisteredType  enums.RegisteredTypes `json:"registered_type"              example:"EMAIL"`
// 	Email           string                `json:"email"                        form:"email"                        gorm:"unique" example:"email@email.com"`
// 	Password        string                `json:"password"                     form:"password"                     gorm:"default:null" example:"password1234"`
// 	FirstName       string                `json:"first_name"                   form:"first_name"                   example:"John"`
// 	LastName        string                `json:"last_name"                    form:"last_name"                    example:"Doe"`
// 	PhoneNumber     string                `json:"phone_number"                 form:"phone_number"                 gorm:"unique" example:"0812345678"`
// 	ProfileImageUrl string                `json:"profile_image_url"            form:"profile_image_url"            gorm:"default:null" example:"https://image_url.com/abcd"`
// 	IsVerified      bool                  `json:"is_verified"                  gorm:"default:null" example:"false"`
// 	CommonModels
// }
