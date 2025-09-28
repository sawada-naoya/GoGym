package router

import (
	"gogym-api/internal/adapter/handler"

	"github.com/labstack/echo/v4"
)

func SessionRoutes(e *echo.Group, sh *handler.SessionHandler) {
	sessions := e.Group("/sessions")
	sessions.POST("/login", sh.Login)
}
