package router

import (
	"gogym-api/internal/adapter/handler"
	"gogym-api/internal/adapter/server"
	"gogym-api/internal/configs"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(
	httpCfg configs.HTTPConfig,
	gymHandler *handler.GymHandler,
	userHandler *handler.UserHandler,
	reviewHandler *handler.ReviewHandler,
	sessionHandler *handler.SessionHandler,
) *echo.Echo {
	e := server.NewEcho(httpCfg)
	v1 := e.Group("/api/v1")

	GymRoutes(v1, gymHandler)
	UserRoutes(v1, userHandler)
	ReviewRoutes(v1, reviewHandler)
	SessionRoutes(v1, sessionHandler)
	return e
}
