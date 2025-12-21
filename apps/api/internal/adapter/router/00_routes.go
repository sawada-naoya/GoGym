package router

import (
	"gogym-api/internal/adapter/handler"
	"gogym-api/internal/middleware"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(
	e *echo.Echo,
	gymHandler *handler.GymHandler,
	userHandler *handler.UserHandler,
	sessionHandler *handler.SessionHandler,
	workoutHandler *handler.WorkoutHandler,
	jwtSecret string,
) {
	v1 := e.Group("/api/v1")

	// 認証不要なルート
	UserRoutes(v1, userHandler)
	SessionRoutes(v1, sessionHandler)

	// 認証が必要なルート
	authMiddleware := middleware.AuthMiddleware(jwtSecret)
	authGroup := v1.Group("", authMiddleware)
	GymRoutes(authGroup, gymHandler)
	WorkoutRoutes(authGroup, workoutHandler)
}
