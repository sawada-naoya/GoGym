package router

import (
	"github.com/labstack/echo/v4"
	"gogym-api/internal/adapter/http/handler"
)

func SetupGymRoutes(e *echo.Echo, gymHandler *handler.GymHandler) {
	// パブリックなジム情報API
	api := e.Group("/api/v1")
	
	// おすすめジム取得
	api.GET("/gyms/recommended", gymHandler.GetRecommendedGyms)
	
	// ジム検索
	api.GET("/gyms", gymHandler.SearchGyms)
	
	// 特定ジム取得
	api.GET("/gyms/:id", gymHandler.GetGym)
}