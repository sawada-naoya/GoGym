package router

import (
	"gogym-api/internal/adapter/handler"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(
	e *echo.Echo,
	gymHandler *handler.GymHandler,
	userHandler *handler.UserHandler,
	sessionHandler *handler.SessionHandler,
	workoutHandler *handler.WorkoutHandler,
) {
	v1 := e.Group("/api/v1")
	GymRoutes(v1, gymHandler)
	UserRoutes(v1, userHandler)
	SessionRoutes(v1, sessionHandler)
	WorkoutRoutes(v1, workoutHandler)
}
