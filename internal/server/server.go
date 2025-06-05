package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/igris-hash/go-event-app/config"
	"github.com/igris-hash/go-event-app/internal/handler"
	"github.com/igris-hash/go-event-app/internal/repository"
	"github.com/igris-hash/go-event-app/internal/service"
	"github.com/igris-hash/go-event-app/pkg/logger"
	"github.com/igris-hash/go-event-app/pkg/middleware"
)

// Server holds the HTTP server and its dependencies
type Server struct {
	config *config.Config
	logger *logger.Logger
	db     *sql.DB
	Router *gin.Engine
	server *http.Server
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config, logger *logger.Logger, db *sql.DB) *Server {
	server := &Server{
		config: cfg,
		logger: logger,
		db:     db,
		Router: gin.New(),
	}

	server.setupRouter()
	return server
}

// setupRouter configures the HTTP router
func (s *Server) setupRouter() {
	// Middleware
	s.Router.Use(gin.Recovery())
	s.Router.Use(middleware.Logger(s.logger))
	s.Router.Use(cors.New(cors.Config{
		AllowOrigins: s.config.CorsOrigins,
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	// Initialize repositories
	userRepo := repository.NewUserRepository(s.db)
	eventRepo := repository.NewEventRepository(s.db)

	// Initialize services
	userService := service.NewUserService(userRepo)
	eventService := service.NewEventService(eventRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	eventHandler := handler.NewEventHandler(eventService)

	// API v1 routes
	v1 := s.Router.Group("/api/v1")
	{
		// User routes
		users := v1.Group("/users")
		{
			users.POST("/register", userHandler.Register)
			users.POST("/login", userHandler.Login)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(s.config.JWTSecret))
		{
			// Event routes
			events := protected.Group("/events")
			{
				events.POST("", eventHandler.CreateEvent)
				events.GET("", eventHandler.GetEvents)
				events.GET("/:id", eventHandler.GetEvent)
				events.PUT("/:id", eventHandler.UpdateEvent)
				events.DELETE("/:id", eventHandler.DeleteEvent)
				events.POST("/:id/register", eventHandler.RegisterForEvent)
				events.DELETE("/:id/register", eventHandler.CancelRegistration)
			}

			// User profile
			users.GET("/me", userHandler.GetProfile)
			users.PUT("/me", userHandler.UpdateProfile)
		}
	}
}

// Run starts the HTTP server
func (s *Server) Run() error {
	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%s", s.config.ServerPort),
		Handler: s.Router,
	}

	s.logger.Info("starting server", "port", s.config.ServerPort)
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
