package router

import (
	"gogym-api/internal/adapter/http/handler"

	"github.com/labstack/echo/v4"
)

func SetupFavoriteRoutes(e *echo.Echo, favoriteHandler *handler.FavoriteHandler) {
	// Nested under gyms
	gym := e.Group("/api/v1/gyms")
	gym.POST("/:gym_id/favorite", favoriteHandler.CreateFavorite)
	gym.DELETE("/:gym_id/favorite", favoriteHandler.DeleteFavorite)

	// Direct favorites access
	favorite := e.Group("/api/v1/favorites")
	favorite.GET("", favoriteHandler.GetUserFavorites)
}
