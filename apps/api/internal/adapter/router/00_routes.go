package router

import (
	"gogym-api/internal/adapter/handler"
	"gogym-api/internal/adapter/server"
	"gogym-api/internal/configs"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(
	httpCfg configs.HTTPConfig,
	authCfg configs.AuthConfig,
	gymHandler *handler.GymHandler,
	userHandler *handler.UserHandler,
	sessionHandler *handler.SessionHandler,
	workoutHandler *handler.WorkoutHandler,
) *echo.Echo {
	e := server.NewEcho(httpCfg)

	// Health check endpoint (for Docker, Render, Fly.io)
	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	v1 := e.Group("/api/v1")

	// MVPに必要な機能のみ有効化
	GymRoutes(v1, gymHandler)           // ジム登録・検索
	UserRoutes(v1, userHandler)         // ユーザー登録
	SessionRoutes(v1, sessionHandler)   // 認証
	WorkoutRoutes(v1, workoutHandler, authCfg) // トレーニング記録

	// 将来実装予定の機能（一旦無効化）
	// ReviewRoutes(v1, reviewHandler)  // レビュー機能

	return e
}
