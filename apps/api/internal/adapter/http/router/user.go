package router

import (
	"gogym-api/internal/adapter/http/handler"

	"github.com/labstack/echo/v4"
)

func SetupUserRoutes(e *echo.Echo, uh *handler.UserHandler) {
	user := e.Group("/api/v1/users")

	user.POST("", uh.SignUp)
}
