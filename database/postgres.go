package database

import (
	"github.com/brain-flowing-company/pprp-backend/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(cfg *config.Config) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(cfg.DBUrl))
}
