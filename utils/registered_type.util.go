package utils

import "github.com/brain-flowing-company/pprp-backend/internal/models"

func ParseRegisteredType(s string) models.RegisteredType {
	return map[string]models.RegisteredType{
		"EMAIL":  models.EMAIL,
		"GOOGLE": models.GOOGLE,
	}[s]
}
