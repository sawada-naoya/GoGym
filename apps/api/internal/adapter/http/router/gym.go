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

	// ジム検索・オートコンプリート
	api.GET("/gyms", g.SearchGyms)
	api.GET("/gyms/autocomplete", g.AutocompleteGyms)

	// 特定ジム取得
	api.GET("/gym/:id", g.GetGym)
	api.GET("/gyms/:id/images", g.GetGymImages)

	// ジム管理（認証必要）
	api.POST("/gyms", g.CreateGym)
	api.PUT("/gyms/:id", g.UpdateGym)
	api.DELETE("/gyms/:id", g.DeleteGym)
}
