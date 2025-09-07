package router

import (
	"github.com/labstack/echo/v4"
	"gogym-api/internal/adapter/http/handler"
)

func NewRouter(e *echo.Echo, gymHandler *handler.GymHandler, userHandler *handler.UserHandler) {
	SetupHealthRoutes(e)
	SetupGymRoutes(e, gymHandler)
	SetupUserRoutes(e, userHandler)
}
