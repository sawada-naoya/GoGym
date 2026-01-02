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
	contactHandler *handler.ContactHandler,
	jwtSecret string,
) {
	v1 := e.Group("/api/v1")

	// ヘルスチェック
	v1.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	// 認証不要なルート
	UserRoutes(v1, userHandler)
	SessionRoutes(v1, sessionHandler)
	ContactRoutes(v1, contactHandler)

	// 認証が必要なルート
	authMiddleware := middleware.AuthMiddleware(jwtSecret)
	authGroup := v1.Group("", authMiddleware)
	GymRoutes(authGroup, gymHandler)
	WorkoutRoutes(authGroup, workoutHandler)
}
