package models

import "time"

type Property struct {
	PropertyId            string          `json:"property_id"`
	Description           string          `json:"description"`
	ResidentialType       string          `json:"residential_type"`
	ProjectName           string          `json:"project_name"`
	Address               string          `json:"address"`
	Alley                 string          `json:"alley"`
	Street                string          `json:"street"`
	SubDistrict           string          `json:"sub_district"`
	District              string          `json:"district"`
	Province              string          `json:"province"`
	Country               string          `json:"country"`
	PostalCode            string          `json:"postal_code"`
	PropertyListTimestamp time.Time       `json:"property_list_time_stamp"`
	PropertyImages        []PropertyImage `gorm:"references:PropertyId" json:"images"`
	SellingProperty       SellingProperty `gorm:"references:PropertyId" json:"selling"`
	RentingProperty       RentingProperty `gorm:"references:PropertyId" json:"renting"`
}

type PropertyImage struct {
	PropertyId string `json:"-"`
	ImageUrl   string `json:"url"`
}

type SellingProperty struct {
	PropertyId string  `json:"-"`
	Price      float64 `json:"price"`
	IsSold     bool    `json:"is_sold"`
}

type RentingProperty struct {
	PropertyId    string  `json:"-"`
	PricePerMonth float64 `json:"price_per_month"`
	IsOccupied    bool    `json:"is_occupied"`
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
