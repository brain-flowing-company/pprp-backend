package models

import "time"

type Property struct {
	PropertyId            string
	Description           string
	ResidentialType       string
	ProjectName           string
	Address               string
	Alley                 string
	Street                string
	SubDistrict           string
	District              string
	Province              string
	Country               string
	PostalCode            string
	PropertyListTimestamp time.Time
	PropertyImages        []PropertyImage `gorm:"references:PropertyId"`
	SellingProperty       SellingProperty `gorm:"references:PropertyId"`
	RentingProperty       RentingProperty `gorm:"references:PropertyId"`
}

type PropertyImage struct {
	PropertyId string
	ImageUrl   string
}

type SellingProperty struct {
	PropertyId string
	Price      float64
	IsSold     bool
}

type RentingProperty struct {
	PropertyId    string
	PricePerMonth float64
	IsOccupied    bool
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
