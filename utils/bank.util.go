package utils

import "github.com/brain-flowing-company/pprp-backend/internal/models"

func ParseBankName(s string) models.BankName {
	return map[string]models.BankName{
		"KBANK": models.KBANK,
		"BBL":   models.BBL,
		"KTB":   models.KTB,
		"BAY":   models.BAY,
		"CIMB":  models.CIMB,
		"TTB":   models.TTB,
		"SCB":   models.SCB,
		"GSB":   models.GSB,
		"NULL":  models.NULL,
	}[s]
}
