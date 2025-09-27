package router

import (
	"gogym-api/internal/adapter/http/handler"
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(e *echo.Echo, gymHandler *handler.GymHandler, userHandler *handler.UserHandler, reviewHandler *handler.ReviewHandler, favoriteHandler *handler.FavoriteHandler, sessionHandler *handler.SessionHandler) *echo.Echo {
	// CORSミドルウェア
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", "http://localhost:3001", "http://localhost:3003"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"*"},
	}))

	// リクエストログを追加
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogMethod: true,
		LogLatency: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			slog.Info("HTTP Request",
				"method", values.Method,
				"uri", values.URI,
				"status", values.Status,
				"latency", values.Latency.String(),
			)
			return nil
		},
	}))

	SetupHealthRoutes(e)
	SetupGymRoutes(e, gymHandler)
	SetupUserRoutes(e, userHandler)
	SetupReviewRoutes(e, reviewHandler)
	SetupFavoriteRoutes(e, favoriteHandler)
	SessionRoutes(e, sessionHandler)
	return e
}
