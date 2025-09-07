package router

import (
	"gogym-api/internal/adapter/http/handler"

	"github.com/labstack/echo/v4"
)

func SetupUserRoutes(e *echo.Echo, userHandler *handler.UserHandler) {
	userGroup := e.Group("/users")

	userGroup.POST("/signup", userHandler.Signup)
}
