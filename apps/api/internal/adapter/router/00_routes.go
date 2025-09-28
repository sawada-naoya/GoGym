package router

import (
	"gogym-api/internal/adapter/handler"
	"gogym-api/internal/adapter/server"
	"gogym-api/internal/configs"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(
	e *echo.Echo,
	gymHandler *handler.GymHandler,
	userHandler *handler.UserHandler,
	reviewHandler *handler.ReviewHandler,
	sessionHandler *handler.SessionHandler,
) {
	v1 := e.Group("/api/v1")

	GymRoutes(v1, gymHandler)
	UserRoutes(v1, userHandler)
	ReviewRoutes(v1, reviewHandler)
	SessionRoutes(v1, sessionHandler)
}

func BuildEcho(
	httpCfg configs.HTTPConfig,
	gym *handler.GymHandler,
	user *handler.UserHandler,
	review *handler.ReviewHandler,
	session *handler.SessionHandler,
) *echo.Echo {
	e := server.NewEcho(httpCfg)
	RegisterRoutes(e, gym, user, review, session)
	return e
}
