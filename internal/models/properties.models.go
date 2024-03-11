package models

import (
	"github.com/brain-flowing-company/pprp-backend/internal/enums"
	"github.com/google/uuid"
)

type Properties struct {
	PropertyId          uuid.UUID            `gorm:"type:uuid;unique;primaryKey;default:uuid_generate_v4()"`
	OwnerId             uuid.UUID            `json:"owner_id"                 example:"123e4567-e89b-12d3-a456-426614174000"`
	PropertyName        string               `json:"property_name"            example:"Supalai"`
	PropertyDescription string               `json:"property_description"              example:"Et sequi dolor praes"`
	PropertyType        enums.PropertyTypes  `json:"property_type"            example:"CONDOMINIUM"`
	Address             string               `json:"address"                  example:"123/4"`
	Alley               string               `json:"alley" gorm:"default:null" example:"Pattaya Nua 78"`
	Street              string               `json:"street"                   example:"Pattaya"`
	SubDistrict         string               `json:"sub_district"             example:"Bang Bon"`
	District            string               `json:"district"                 example:"Bang Phli"`
	Province            string               `json:"province"                 example:"Pattaya"`
	Country             string               `json:"country"                  example:"Thailand"`
	PostalCode          string               `json:"postal_code"              example:"69096"`
	Bedrooms            int64                `json:"bedrooms"                 example:"3"`
	Bathrooms           int64                `json:"bathrooms"                example:"2"`
	Furnishing          enums.Furnishing     `json:"furnishing"               example:"UNFURNISHED"`
	Floor               int64                `json:"floor"                    example:"5"`
	FloorSize           float64              `json:"floor_size"               example:"123.45"`
	FloorSizeUnit       enums.FloorSizeUnits `json:"floor_size_unit" gorm:"default:SQM" example:"SQM"`
	UnitNumber          int64                `json:"unit_number"              example:"123"`
	PropertyImages      []PropertyImages     `gorm:"foreignKey:PropertyId; references:PropertyId" json:"images"`
	SellingProperty     SellingProperties    `gorm:"foreignKey:PropertyId; references:PropertyId" json:"selling"`
	RentingProperty     RentingProperties    `gorm:"foreignKey:PropertyId; references:PropertyId" json:"renting"`
	CommonModels
}

type PropertyImages struct {
	PropertyId uuid.UUID `json:"-"`
	ImageUrl   string    `json:"url" example:"https://image_url.com/abcd"`
	CommonModels
}

type SellingProperties struct {
	PropertyId uuid.UUID `json:"-"`
	Price      float64   `json:"price"   example:"12345.67"`
	IsSold     bool      `json:"is_sold" example:"true"`
	CommonModels
}

type RentingProperties struct {
	PropertyId    uuid.UUID `json:"-"`
	PricePerMonth float64   `json:"price_per_month" example:"12345.67"`
	IsOccupied    bool      `json:"is_occupied"     example:"true"`
	CommonModels
}

type FavoriteProperties struct {
	PropertyId uuid.UUID `json:"-"`
	UserId     uuid.UUID `json:"-"`
}

func (p Properties) TableName() string {
	return "properties"
}

func (p PropertyImages) TableName() string {
	return "property_images"
}

func (p SellingProperties) TableName() string {
	return "selling_properties"
}

func (p RentingProperties) TableName() string {
	return "renting_properties"
}

func (p FavoriteProperties) TableName() string {
	return "favorite_properties"
}
