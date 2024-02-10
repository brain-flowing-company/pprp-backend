package main

import (
	"fmt"

	"github.com/brain-flowing-company/pprp-backend/config"
	"github.com/brain-flowing-company/pprp-backend/database"
	_ "github.com/brain-flowing-company/pprp-backend/docs"
	"github.com/brain-flowing-company/pprp-backend/internal/google"
	"github.com/brain-flowing-company/pprp-backend/internal/greeting"
	"github.com/brain-flowing-company/pprp-backend/internal/login"
	"github.com/brain-flowing-company/pprp-backend/internal/property"
	"github.com/brain-flowing-company/pprp-backend/internal/users"
	"github.com/brain-flowing-company/pprp-backend/middleware"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// @title        Bangkok Property Matchmaking Platform
// @description  Bangkok Property Matchmaking Platform API docs
// @version      1.0
// @host         localhost:8000
// @BasePath     /
func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Could not load environment variables from .env file: %v\n", err.Error())
	}

	cfg := &config.Config{}
	err = config.Load(cfg)
	if err != nil {
		panic(fmt.Sprintf("Could not load config with error: %v", err.Error()))
	}

	db, err := database.New(cfg)
	if err != nil {
		panic(fmt.Sprintf("Could not establish connection with database with err: %v", err.Error()))
	}

	app := fiber.New()

	var logger *zap.Logger
	if cfg.IsDevelopment() {
		logger = zap.Must(zap.NewDevelopment())
	} else {
		logger = zap.Must(zap.NewProduction())
	}

	app.Use(fiberzap.New(fiberzap.Config{
		Logger: logger,
	}))

	fmt.Println(cfg.AllowOrigin)

	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowOrigin,
		AllowCredentials: true,
	}))

	if cfg.IsDevelopment() {
		app.Get("/docs/*", swagger.HandlerDefault)
	}

	hwService := greeting.NewService()
	hwHandler := greeting.NewHandler(hwService)

	propertyRepo := property.NewRepository(db)
	propertyService := property.NewService(propertyRepo, logger)
	propertyHandler := property.NewHandler(propertyService)

	usersRepo := users.NewRepository(db)
	usersService := users.NewService(usersRepo, logger)
	usersHandler := users.NewHandler(usersService)

	googleRepo := google.NewRepository(db)
	googleService := google.NewService(googleRepo, cfg, logger)
	googleHandler := google.NewHandler(googleService, logger, cfg)

	// Initialize the repository, service, and handler
	loginRepository := login.NewRepository(db)
	loginService := login.NewService(loginRepository, cfg, logger)
	loginHandler := login.NewHandler(loginService, cfg, logger)

	mw := middleware.NewMiddleware(cfg)

	apiv1 := app.Group("/api/v1", mw.SessionMiddleware)

	apiv1.Get("/greeting", hwHandler.Greeting)
	apiv1.Get("/user/greeting", mw.AuthMiddlewareWrapper(hwHandler.UserGreeting))

	apiv1.Get("/property/:propertyId", propertyHandler.GetPropertyById)
	apiv1.Get("/properties", propertyHandler.GetAllProperties)

	apiv1.Get("/users", usersHandler.GetAllUsers)
	apiv1.Get("/user/me", mw.AuthMiddlewareWrapper(usersHandler.GetCurrentUser))
	apiv1.Get("/user/:userId", usersHandler.GetUserById)
	apiv1.Put("/user/:userId", usersHandler.UpdateUser)
	apiv1.Delete("/user/:userId", usersHandler.DeleteUser)

	apiv1.Post("/register", usersHandler.Register)
	apiv1.Post("/login", loginHandler.Login)

	apiv1.Get("/oauth/google", googleHandler.GoogleLogin)
	apiv1.Get("/oauth/callback", googleHandler.ExchangeToken)

	err = app.Listen(fmt.Sprintf(":%v", cfg.AppPort))
	if err != nil {
		panic(fmt.Sprintf("Server could not start with error: %v", err.Error()))
	}
}
