package routes

import (
	"user-management-api/internal/handler"
	"user-management-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App, userHandler *handler.UserHandler) {
	// Apply global middleware
	app.Use(middleware.RequestID())
	app.Use(middleware.Logger())
	
	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
		})
	})
	
	// API routes
	api := app.Group("/users")
	
	api.Post("/", userHandler.CreateUser)
	api.Get("/", userHandler.ListUsers)
	api.Get("/:id", userHandler.GetUserByID)
	api.Put("/:id", userHandler.UpdateUser)
	api.Delete("/:id", userHandler.DeleteUser)
}
