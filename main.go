package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/igris-hash/go-event-app/internal/handler"
	"github.com/igris-hash/go-event-app/internal/middleware"
	"github.com/igris-hash/go-event-app/internal/repository"
	"github.com/igris-hash/go-event-app/internal/service"
	"github.com/igris-hash/go-event-app/internal/utils"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "modernc.org/sqlite"

	_ "github.com/igris-hash/go-event-app/docs" // swagger docs
)

// @title Event Management API
// @version 1.0
// @description This is a simple event management service API
// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Initialize SQLite database
	db, err := utils.InitializeDB("events.db")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	eventRepo := repository.NewEventRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo)
	eventService := service.NewEventService(eventRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	eventHandler := handler.NewEventHandler(eventService)

	// Initialize router
	router := gin.Default()

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public routes
	router.POST("/users/register", userHandler.Register)
	router.POST("/users/login", userHandler.Login)

	router.GET("/events", eventHandler.ListEvents)
	router.GET("/events/:id", eventHandler.GetEvent)

	// Protected routes
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// User routes
		protected.GET("/users/me", userHandler.GetProfile)
		protected.PUT("/users/me", userHandler.UpdateProfile)

		// Event routes
		protected.POST("/events", eventHandler.CreateEvent)
		protected.PUT("/events/:id", eventHandler.UpdateEvent)
		protected.DELETE("/events/:id", eventHandler.DeleteEvent)
		protected.POST("/events/:id/register", eventHandler.RegisterForEvent)
		protected.POST("/events/:id/unregister", eventHandler.UnregisterFromEvent)
		protected.GET("/events/:id/registrations", eventHandler.GetRegisteredUsers)
	}

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
