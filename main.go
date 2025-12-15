package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"user-management-api/config"
	"user-management-api/internal/handler"
	"user-management-api/internal/logger"
	"user-management-api/internal/repository"
	"user-management-api/internal/routes"
	"user-management-api/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg := config.Load()
	
	// Initialize logger
	if err := logger.Init(cfg.Server.Environment); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()
	
	logger.Info("Starting User Management API",
		zap.String("environment", cfg.Server.Environment),
		zap.String("port", cfg.Server.Port),
	)
	
	// Connect to database
	db, err := sql.Open("postgres", cfg.Database.ConnectionString())
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()
	
	// Test database connection
	if err := db.Ping(); err != nil {
		logger.Error("Failed to ping database", zap.Error(err))
		log.Fatalf("Database ping failed: %v", err)
	}
	
	logger.Info("Database connection established")
	
	// Initialize layers
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	
	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			
			logger.Error("Request error",
				zap.Error(err),
				zap.Int("status", code),
				zap.String("path", c.Path()),
			)
			
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})
	
	// Add panic recovery middleware
	app.Use(recover.New())
	
	// Setup routes
	routes.Setup(app, userHandler)
	
	// Start server in a goroutine
	go func() {
		addr := fmt.Sprintf(":%s", cfg.Server.Port)
		logger.Info("Server starting", zap.String("address", addr))
		if err := app.Listen(addr); err != nil {
			logger.Error("Server failed to start", zap.Error(err))
			log.Fatalf("Server error: %v", err)
		}
	}()
	
	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	
	logger.Info("Shutting down server...")
	
	if err := app.Shutdown(); err != nil {
		logger.Error("Server shutdown error", zap.Error(err))
	}
	
	logger.Info("Server stopped gracefully")
}
