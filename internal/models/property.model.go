package models

import "time"

type Property struct {
	PropertyId            string          `json:"property_id"              example:"f38f80b3-f326-4825-9afc-ebc331626875"`
	Description           string          `json:"description"              example:"Et sequi dolor praes"`
	ResidentialType       string          `json:"residential_type"         example:"Condo"`
	ProjectName           string          `json:"project_name"             example:"Supalai"`
	Address               string          `json:"address"                  example:"123/4"`
	Alley                 string          `json:"alley"                    example:"Pattaya Nua 78"`
	Street                string          `json:"street"                   example:"Pattaya"`
	SubDistrict           string          `json:"sub_district"             example:"Bang Bon"`
	District              string          `json:"district"                 example:"Bang Phli"`
	Province              string          `json:"province"                 example:"Pattaya"`
	Country               string          `json:"country"                  example:"Thailand"`
	PostalCode            string          `json:"postal_code"              example:"69096"`
	PropertyListTimestamp time.Time       `json:"property_list_timestamp"  example:"2023-08-12T11:37:08Z"`
	PropertyImages        []PropertyImage `gorm:"references:PropertyId" json:"images"`
	SellingProperty       SellingProperty `gorm:"references:PropertyId" json:"selling"`
	RentingProperty       RentingProperty `gorm:"references:PropertyId" json:"renting"`
}

type PropertyImage struct {
	PropertyId string `json:"-"`
	ImageUrl   string `json:"url" example:"https://image_url.com/abcd"`
}

type SellingProperty struct {
	PropertyId string  `json:"-"`
	Price      float64 `json:"price"   example:"12345.67"`
	IsSold     bool    `json:"is_sold" example:"true"`
}

type RentingProperty struct {
	PropertyId    string  `json:"-"`
	PricePerMonth float64 `json:"price_per_month" example:"12345.67"`
	IsOccupied    bool    `json:"is_occupied"     example:"true"`
}

func (p Property) TableName() string {
	return "property"
}

func (p PropertyImage) TableName() string {
	return "property_image"
}

func (p SellingProperty) TableName() string {
	return "selling_property"
}

func (p RentingProperty) TableName() string {
	return "renting_property"
}
