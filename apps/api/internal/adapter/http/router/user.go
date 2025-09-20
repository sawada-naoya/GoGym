package router

import (
	"gogym-api/internal/adapter/http/handler"

	"github.com/labstack/echo/v4"
)

func SetupUserRoutes(e *echo.Echo, uh *handler.UserHandler) {
	userGroup := e.Group("/api/v1/users")

	userGroup.POST("/user", uh.SignUp)
}
