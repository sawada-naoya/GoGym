package router

import (
	"github.com/labstack/echo/v4"
	"gogym-api/internal/adapter/http/handler"
)

func SetupUserRoutes(e *echo.Echo, userHandler *handler.UserHandler) {
	userGroup := e.Group("/users")
	
	userGroup.POST("/signup", userHandler.Signup)
}