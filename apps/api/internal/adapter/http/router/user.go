package router

import (
	"gogym-api/internal/adapter/http/handler"
	"log/slog"

	"github.com/labstack/echo/v4"
)

func SetupUserRoutes(e *echo.Echo, uh *handler.UserHandler) {
	slog.Info("Setting up user routes")
	userGroup := e.Group("/api/v1/users")

	userGroup.POST("", uh.SignUp)
	slog.Info("User routes registered", "endpoint", "POST /api/v1/users")
}
