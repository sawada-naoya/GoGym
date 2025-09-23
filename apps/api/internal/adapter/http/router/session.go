package router

import (
	"gogym-api/internal/adapter/http/handler"

	"github.com/labstack/echo/v4"
)

func SessionRoutes(e *echo.Echo, sh *handler.SessionHandler) {
	sessions := e.Group("/api/v1/sessions")

	sessions.POST("", sh.Login)
}
