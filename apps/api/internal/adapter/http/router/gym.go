package router

import (
	"gogym-api/internal/adapter/http/handler"

	"github.com/labstack/echo/v4"
)

func SetupGymRoutes(e *echo.Echo, g *handler.GymHandler) {
	// パブリックなジム情報API
	api := e.Group("/api/v1")

	// おすすめジム取得
	api.GET("/gyms/recommended", g.GetRecommendedGyms)

	// ジム検索
	api.GET("/gyms", g.SearchGyms)

	// 特定ジム取得
	api.GET("/gyms/:id", g.GetGym)
}
