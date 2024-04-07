package utils

import (
	"github.com/go-playground/validator/v10"
)

func PropertyTypeValidator(propertyType validator.FieldLevel) bool {
	switch propertyType.Field().String() {
	case "CONDOMINIUM", "APARTMENT", "SEMI-DETACHED_HOUSE", "HOUSE", "SERVICED_APARTMENT", "TOWNHOUSE":
		return true
	}
	return false
}

func PostalCodeValidator(postalCode validator.FieldLevel) bool {
	return len(postalCode.Field().String()) == 5
}

func FurnishingValidator(furnishing validator.FieldLevel) bool {
	switch furnishing.Field().String() {
	case "UNFURNISHED", "PARTIALLY_FURNISHED", "FULLY_FURNISHED":
		return true
	}
	return false
}

func FloorSizeUnitValidator(floorSizeUnit validator.FieldLevel) bool {
	switch floorSizeUnit.Field().String() {
	case "SQM", "SQFT":
		return true
	}
	return false
}

func NewPropertyValidator() (*validator.Validate, error) {
	v := validator.New()
	if err := v.RegisterValidation("property_type", PropertyTypeValidator); err != nil {
		return nil, err
	}

	if err := v.RegisterValidation("postal_code", PostalCodeValidator); err != nil {
		return nil, err
	}

	if err := v.RegisterValidation("furnishing", FurnishingValidator); err != nil {
		return nil, err
	}

	if err := v.RegisterValidation("fs_unit", FloorSizeUnitValidator); err != nil {
		return nil, err
	}

	return v, nil
}
