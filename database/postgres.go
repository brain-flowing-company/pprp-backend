package database

import (
	"github.com/brain-flowing-company/pprp-backend/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(cfg *config.Config) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(cfg.DBUrl), &gorm.Config{
		TranslateError: true,
		Logger:         logger.Default.LogMode(logger.Silent),
	})
}

func CreateTestDatabase(cfg *config.Config) error {
	db, err := gorm.Open(postgres.Open(cfg.DBUrl), &gorm.Config{
		TranslateError: true,
		Logger:         logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	defer db.Exec("DROP DATABASE test_database")

	// Create the test database
	err = db.Exec("CREATE DATABASE test_database").Error
	if err != nil {
		return err
	}

	return nil
}
