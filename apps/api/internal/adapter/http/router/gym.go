package router

import (
	"gogym-api/internal/adapter/http/handler"

	"github.com/labstack/echo/v4"
)

func SetupGymRoutes(e *echo.Echo, g *handler.GymHandler) {

	gym := e.Group("/api/v1")

	// おすすめジム取得
	gym.GET("/gyms/recommended", g.GetRecommendedGyms)

	// 特定ジム取得
	gym.GET("/gym/:id", g.GetGym)

}
