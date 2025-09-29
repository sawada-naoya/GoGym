package router

import (
	"gogym-api/internal/adapter/handler"

	"github.com/labstack/echo/v4"
)

func GymRoutes(e *echo.Group, g *handler.GymHandler) {
	gym := e.Group("/gyms")
	// おすすめジム取得
	gym.GET("/recommended", g.GetRecommendedGyms)

	// 特定ジム取得
	gym.GET("/:id", g.GetGym)

}
