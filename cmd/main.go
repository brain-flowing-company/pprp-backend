package main

import (
	"fmt"

	"github.com/brain-flowing-company/pprp-backend/config"
	"github.com/brain-flowing-company/pprp-backend/database"
	_ "github.com/brain-flowing-company/pprp-backend/docs"
	"github.com/brain-flowing-company/pprp-backend/internal/core/agreements"
	"github.com/brain-flowing-company/pprp-backend/internal/core/appointments"
	"github.com/brain-flowing-company/pprp-backend/internal/core/auth"
	"github.com/brain-flowing-company/pprp-backend/internal/core/chats"
	"github.com/brain-flowing-company/pprp-backend/internal/core/emails"
	"github.com/brain-flowing-company/pprp-backend/internal/core/google"
	"github.com/brain-flowing-company/pprp-backend/internal/core/greetings"
	"github.com/brain-flowing-company/pprp-backend/internal/core/payments"
	"github.com/brain-flowing-company/pprp-backend/internal/core/properties"
	"github.com/brain-flowing-company/pprp-backend/internal/core/ratings"
	"github.com/brain-flowing-company/pprp-backend/internal/core/users"
	"github.com/brain-flowing-company/pprp-backend/internal/middleware"
	"github.com/brain-flowing-company/pprp-backend/storage"
	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/contrib/websocket"
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

		// test websocket
		app.Static("/", "./internal/core/chats/public/")
	}

	hwService := greetings.NewService()
	hwHandler := greetings.NewHandler(hwService)

	propertyRepo := properties.NewRepository(db)
	propertyService := properties.NewService(logger, propertyRepo, storage)
	propertyHandler := properties.NewHandler(propertyService)

	usersRepo := users.NewRepository(db)
	usersService := users.NewService(logger, cfg, usersRepo, storage)
	usersHandler := users.NewHandler(logger, cfg, usersService)

	googleRepo := google.NewRepository(db)
	googleService := google.NewService(logger, cfg, googleRepo)
	googleHandler := google.NewHandler(logger, cfg, googleService)

	emailRepository := emails.NewRepository(db)
	emailService := emails.NewService(logger, cfg, emailRepository)
	emailHandler := emails.NewHandler(logger, cfg, emailService)

	// Initialize the repository, service, and handler
	authRepository := auth.NewRepository(db)
	authService := auth.NewService(logger, cfg, authRepository, googleService, emailService)
	authHandler := auth.NewHandler(cfg, authService)

	chatRepository := chats.NewRepository(db)
	chatService := chats.NewService(logger, chatRepository)
	hub := chats.NewHub(chatService)
	chatHandler := chats.NewHandler(logger, cfg, hub, chatService)

	appointmentRepository := appointments.NewRepository(db)
	appointmentService := appointments.NewService(logger, appointmentRepository)
	appointmentHandler := appointments.NewHandler(hub, appointmentService)

	agreementsRepo := agreements.NewRepository(db)
	agreementsService := agreements.NewService(logger, agreementsRepo)
	agreementsHandler := agreements.NewHandler(hub, agreementsService)

	paymentsRepository := payments.NewRepository(db)
	paymentsService := payments.NewService(logger, paymentsRepository, cfg)
	paymentsHandler := payments.NewHandler(cfg, paymentsService)

	ratingsRepository := ratings.NewRepository(db)
	ratingsService := ratings.NewService(ratingsRepository, logger, cfg)
	ratingsHandler := ratings.NewHandler(ratingsService)

	mw := middleware.NewMiddleware(cfg)

	apiv1 := app.Group("/api/v1", mw.SessionMiddleware)

	apiv1.Get("/checkout", mw.WithAuthentication(paymentsHandler.CreatePayment))
	apiv1.Get("/payments", mw.WithAuthentication(paymentsHandler.GetPaymentByUserId))
	apiv1.Get("/payments/history", mw.WithAuthentication(paymentsHandler.GetHistoryPaymentByUserId))

	apiv1.Get("/greeting", hwHandler.Greeting)
	apiv1.Get("/user/greeting", mw.WithAuthentication(hwHandler.UserGreeting))

	apiv1.Get("/properties/:propertyId", propertyHandler.GetPropertyById)
	apiv1.Get("/properties", propertyHandler.GetAllProperties)
	apiv1.Get("/user/me/properties", mw.WithAuthentication(propertyHandler.GetMyProperties))
	apiv1.Post("/properties", mw.WithOwnerAccess(propertyHandler.CreateProperty))
	apiv1.Patch("/properties/:propertyId", mw.WithOwnerAccess(propertyHandler.UpdatePropertyById))
	apiv1.Delete("/properties/:propertyId", mw.WithOwnerAccess(propertyHandler.DeletePropertyById))
	apiv1.Post("/properties/favorites/:propertyId", mw.WithAuthentication(propertyHandler.AddFavoriteProperty))
	apiv1.Delete("/properties/favorites/:propertyId", mw.WithAuthentication(propertyHandler.RemoveFavoriteProperty))
	apiv1.Get("/user/me/favorites", mw.WithAuthentication(propertyHandler.GetMyFavoriteProperties))
	apiv1.Get("/top10properties", propertyHandler.GetTop10Properties)

	apiv1.Get("/appointments", mw.WithAuthentication(appointmentHandler.GetAllAppointments))
	apiv1.Get("/appointments/:appointmentId", mw.WithAuthentication(appointmentHandler.GetAppointmentById))
	apiv1.Get("/user/me/appointments", mw.WithAuthentication(appointmentHandler.GetMyAppointments))
	apiv1.Post("/appointments", mw.WithAuthentication(appointmentHandler.CreateAppointment))
	apiv1.Delete("/appointments", mw.WithAuthentication(appointmentHandler.DeleteAppointment))
	apiv1.Patch("/appointments/:appointmentId", mw.WithAuthentication(appointmentHandler.UpdateAppointmentStatus))

	apiv1.Get("/users", usersHandler.GetAllUsers)
	apiv1.Get("/user/me/personal-information", mw.WithAuthentication(usersHandler.GetCurrentUser))
	apiv1.Get("/user/me/financial-information", mw.WithAuthentication(usersHandler.GetUserFinancialInformation))
	apiv1.Get("/user/me/registered", usersHandler.GetRegisteredType)
	apiv1.Get("/user/:userId", usersHandler.GetUserById)
	apiv1.Put("/user/me/personal-information", mw.WithAuthentication(usersHandler.UpdateUser))
	apiv1.Put("/user/me/financial-information", mw.WithAuthentication(usersHandler.UpdateUserFinancialInformation))
	apiv1.Post("/user/me/verify", mw.WithAuthentication(usersHandler.VerifyCitizenId))
	apiv1.Delete("/user/:userId", mw.WithAuthentication(usersHandler.DeleteUser))

	apiv1.Post("/register", usersHandler.Register)
	apiv1.Post("/login", authHandler.Login)
	apiv1.Post("/logout", authHandler.Logout)

	apiv1.Get("/agreements", mw.WithAuthentication(agreementsHandler.GetAllAgreements))
	apiv1.Get("/agreements/:agreementId", mw.WithAuthentication(agreementsHandler.GetAgreementById))
	apiv1.Get("/user/me/agreements", mw.WithAuthentication(agreementsHandler.GetMyAgreements))
	apiv1.Post("/agreements", mw.WithOwnerAccess(agreementsHandler.CreateAgreement))
	apiv1.Delete("/agreements/:agreementId", mw.WithOwnerAccess(agreementsHandler.DeleteAgreement))
	apiv1.Patch("/agreements/:agreementId", mw.WithOwnerAccess(agreementsHandler.UpdateAgreementStatus))

	apiv1.Get("/oauth/google", googleHandler.GoogleLogin)
	apiv1.Post("/email", emailHandler.SendVerificationEmail)
	apiv1.Get("/auth/callback", authHandler.Callback)

	apiv1.Get("/chats", mw.WithAuthentication(chatHandler.GetAllChats))
	apiv1.Get("/chats/:recvUserId", mw.WithAuthentication(chatHandler.GetMessagesInChat))

	apiv1.Post("/ratings", mw.WithAuthentication(ratingsHandler.CreateRating))
	apiv1.Get("/ratings/:propertyId", mw.WithAuthentication(ratingsHandler.GetRatingByPropertyId))
	apiv1.Get("/ratings", mw.WithAuthentication(ratingsHandler.GetAllRatings))
	apiv1.Get("/ratings/sorted/:propertyId", mw.WithAuthentication(ratingsHandler.GetRatingByPropertyIdSortedByRating))
	apiv1.Get("/ratings/newest/:propertyId", mw.WithAuthentication(ratingsHandler.GetRatingByPropertyIdSortedByNewest))
	apiv1.Patch("/ratings/:ratingId", mw.WithAuthentication(ratingsHandler.UpdateRatingStatus))
	apiv1.Delete("/ratings/:ratingId", mw.WithAuthentication(ratingsHandler.DeleteRating))

	ws := app.Group("/ws")
	ws.Get("/chats", websocket.New(chatHandler.OpenConnection))

	err = app.Listen(fmt.Sprintf(":%v", cfg.AppPort))
	if err != nil {
		panic(fmt.Sprintf("Server could not start with error: %v", err.Error()))
	}
}
