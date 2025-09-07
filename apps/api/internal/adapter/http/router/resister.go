package router

import (
	"gogym-api/internal/adapter/http/handler"

	"github.com/labstack/echo/v4"
)

func NewRouter(e *echo.Echo, gymHandler *handler.GymHandler, userHandler *handler.UserHandler) *echo.Echo {
	SetupHealthRoutes(e)
	SetupGymRoutes(e, gymHandler)
	SetupUserRoutes(e, userHandler)
	return e
}
