package http

import (
	"github.com/labstack/echo/v4"
)

func NewRouter(e *echo.Echo, gymHandler *handler.GymHandler,userHandler *handler.UserHandler)  {
		// Middleware
	e.Use(LoggerMiddleware)

	// Routes
	e.GET("/healthz", HealthCheck)

	e.GET("/gyms/search", gymHandler.Search)
	e.POST("/users/signup", userHandler.Signup)
}