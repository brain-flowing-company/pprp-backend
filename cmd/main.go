package main

import (
	"fmt"
	"net/http"

	"github.com/brain-flowing-company/pprp-backend/config"
	"github.com/brain-flowing-company/pprp-backend/database"
	_ "github.com/brain-flowing-company/pprp-backend/docs"
	"github.com/brain-flowing-company/pprp-backend/internal/agreements"
	"github.com/brain-flowing-company/pprp-backend/internal/appointments"
	"github.com/brain-flowing-company/pprp-backend/internal/auth"
	"github.com/brain-flowing-company/pprp-backend/internal/google"
	"github.com/brain-flowing-company/pprp-backend/internal/greeting"
	"github.com/brain-flowing-company/pprp-backend/internal/property"
	"github.com/brain-flowing-company/pprp-backend/internal/users"
	"github.com/brain-flowing-company/pprp-backend/middleware"
	"github.com/brain-flowing-company/pprp-backend/storage"
	"github.com/gofiber/contrib/fiberzap"
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

	fmt.Println(cfg)

	db, err := database.New(cfg)
	if err != nil {
		panic(fmt.Sprintf("Could not establish connection with database with err: %v", err.Error()))
	}

	storage, err := storage.New(cfg)
	if err != nil {
		panic(fmt.Sprintf("Could not establish connection with AWS S3 with err: %v", err.Error()))
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
	propertyService := property.NewService(logger, propertyRepo)
	propertyHandler := property.NewHandler(propertyService)

	agreementsRepo := agreements.NewRepository(db)
	agreementsService := agreements.NewService(logger, agreementsRepo)
	agreementsHandler := agreements.NewHandler(agreementsService)

	usersRepo := users.NewRepository(db)
	usersService := users.NewService(logger, cfg, usersRepo, storage)
	usersHandler := users.NewHandler(usersService)

	googleRepo := google.NewRepository(db)
	googleService := google.NewService(logger, cfg, googleRepo)
	googleHandler := google.NewHandler(logger, cfg, googleService)

	// Initialize the repository, service, and handler
	authRepository := auth.NewRepository(db)
	authService := auth.NewService(logger, cfg, authRepository)
	authHandler := auth.NewHandler(cfg, authService)

	appointmentRepository := appointments.NewRepository(db)
	appointmentService := appointments.NewService(logger, appointmentRepository)
	appointmentHandler := appointments.NewHandler(appointmentService)

	mw := middleware.NewMiddleware(cfg)

	apiv1 := app.Group("/api/v1", mw.SessionMiddleware)

	apiv1.Post("/upload", func(c *fiber.Ctx) error {
		file, err := c.FormFile("profile")
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}

		fileReader, err := file.Open()
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}

		url, err := storage.Upload(fmt.Sprintf("profiles/%v", file.Filename), fileReader)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}

		return c.SendString(url)
	})

	apiv1.Get("/greeting", hwHandler.Greeting)
	apiv1.Get("/user/greeting", mw.AuthMiddlewareWrapper(hwHandler.UserGreeting))

	apiv1.Get("/property/:propertyId", propertyHandler.GetPropertyById)
	apiv1.Get("/properties", propertyHandler.GetAllProperties)
	apiv1.Get("/properties/search", propertyHandler.SeachProperties)

	apiv1.Get("/appointments/:appointmentId", appointmentHandler.GetAppointmentById)
	apiv1.Get("/appointments", appointmentHandler.GetAllAppointments)
	apiv1.Post("/appointments", appointmentHandler.CreateAppointments)
	apiv1.Delete("/appointments", appointmentHandler.DeleteAppointments)
	apiv1.Patch("/appointments/:appointmentId", appointmentHandler.UpdateAppointmentStatus)

	apiv1.Get("/users", usersHandler.GetAllUsers)
	apiv1.Get("/user/me", mw.AuthMiddlewareWrapper(usersHandler.GetCurrentUser))
	apiv1.Get("/user/me/registered", usersHandler.GetRegisteredType)
	apiv1.Get("/user/:userId", usersHandler.GetUserById)
	apiv1.Put("/user/:userId", usersHandler.UpdateUser)
	apiv1.Delete("/user/:userId", usersHandler.DeleteUser)

	apiv1.Post("/register", usersHandler.Register)
	apiv1.Post("/login", authHandler.Login)
	apiv1.Post("/logout", authHandler.Logout)

	apiv1.Get("/agreements", agreementsHandler.GetAllAgreements)
	apiv1.Get("/agreement/:agreementId", agreementsHandler.GetAgreementById)
	apiv1.Get("/user/:userId/agreements", agreementsHandler.GetAgreementsByOwnerId)
	apiv1.Get("/user/:userId/dwelling-agreements", agreementsHandler.GetAgreementsByDwellerId)
	apiv1.Post("/agreement", agreementsHandler.CreateAgreement)
	apiv1.Delete("/agreement/:agreementId", agreementsHandler.DeleteAgreement)

	apiv1.Get("/oauth/google", googleHandler.GoogleLogin)
	apiv1.Get("/oauth/callback", googleHandler.ExchangeToken)

	err = app.Listen(fmt.Sprintf(":%v", cfg.AppPort))
	if err != nil {
		panic(fmt.Sprintf("Server could not start with error: %v", err.Error()))
	}
}
