package utils

import "github.com/brain-flowing-company/pprp-backend/internal/models"

func ParseRegisteredType(s string) models.RegisteredType {
	val, _ := map[string]models.RegisteredType{
		"EMAIL":  models.EMAIL,
		"GOOGLE": models.GOOGLE,
	}[s]

	return val
}
