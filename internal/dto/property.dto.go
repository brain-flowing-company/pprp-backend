package dto

import "github.com/brain-flowing-company/pprp-backend/internal/models"

type GetPropertyByIdResponse struct {
	Property models.Property `json:"property"`
}
