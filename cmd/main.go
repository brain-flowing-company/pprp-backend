package main

import (
	"fmt"

	"github.com/brain-flowing-company/pprp-backend/config"
	"github.com/brain-flowing-company/pprp-backend/database"
	_ "github.com/brain-flowing-company/pprp-backend/docs"
	"github.com/brain-flowing-company/pprp-backend/internal/greeting"
	"github.com/brain-flowing-company/pprp-backend/internal/property"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
)

// @title        Bangkok Property Matchmaking Platform
// @description  Bangkok Property Matchmaking Platform API docs
// @version      1.0
// @host         localhost:3000
// @BasePath     /
func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Failed to load environment variables from .env file: %v\n", err.Error())
	}

	cfg := config.Config{}
	err = config.Load(&cfg)
	if err != nil {
		panic(fmt.Sprintf("Could not load config with error: %v", err.Error()))
	}

	db, err := database.New(&cfg)
	if err != nil {
		panic(fmt.Sprintf("Could not establish connection with database with err: %v", err.Error()))
	}

	app := fiber.New()

	app.Use(logger.New(logger.Config{
		TimeFormat: "02-01-2006 15:04:05",
		TimeZone:   "Asia/Bangkok",
	}))

	if cfg.IsDevelopment() {
		app.Get("/docs/*", swagger.HandlerDefault)
	}

	hwService := greeting.NewService()
	hwHandler := greeting.NewHandler(hwService)

	propertyRepo := property.NewRepository(db)
	propertyService := property.NewService(propertyRepo)
	propertyHandler := property.NewHandler(propertyService)

	app.Get("/greeting", hwHandler.Greeting)
	app.Get("/api/v1/property/:propertyId", propertyHandler.GetPropertyById)

	err = app.Listen(fmt.Sprintf(":%v", cfg.AppPort))
	if err != nil {
		panic(fmt.Sprintf("Server cannot start with error: %v", err.Error()))
	}
}
