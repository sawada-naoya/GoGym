package router

import (
	"gogym-api/internal/adapter/handler"

	"github.com/labstack/echo/v4"
)

func SessionRoutes(e *echo.Group, sh *handler.SessionHandler) {
	e.POST("/sessions/login", sh.Login)
	e.POST("/sessions/refresh", sh.RefreshToken)
}
