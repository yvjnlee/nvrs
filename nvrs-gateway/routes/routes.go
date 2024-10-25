package routes

import (
	"nvrs-gateway/handlers"
	"nvrs-gateway/middleware"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func InitRoutes(e *echo.Echo) {
	// Global rate limiting middleware
	e.Use(echoMiddleware.RateLimiterWithConfig(echoMiddleware.RateLimiterConfig{
		Store: echoMiddleware.NewRateLimiterMemoryStore(10),
	}))

	// Public routes
	e.GET("/health", handlers.HealthCheck)
	e.POST("/agents/register", handlers.RegisterAgent) // Public registration route

	// Protected routes (JWT secured)
	api := e.Group("/api", middleware.JWTMiddleware()) // Apply JWTMiddleware for all /api routes
	api.POST("/tasks", handlers.SubmitTask)
	api.POST("/agents/status", handlers.UpdateStatus)

	// New endpoints for querying tasks and retrieving status
	api.GET("/tasks", handlers.QueryTasks)        // Agent queries their tasks
	api.GET("/agents/status", handlers.GetStatus) // Agent retrieves their status
}
