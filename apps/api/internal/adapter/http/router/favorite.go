package router

import (
	"gogym-api/internal/adapter/http/handler"

	"github.com/labstack/echo/v4"
)

func SetupFavoriteRoutes(e *echo.Echo, favoriteHandler *handler.FavoriteHandler) {
	// Nested under gyms
	gymGroup := e.Group("/api/v1/gyms")
	gymGroup.POST("/:gym_id/favorite", favoriteHandler.CreateFavorite)
	gymGroup.DELETE("/:gym_id/favorite", favoriteHandler.DeleteFavorite)
	
	// Direct favorites access
	favGroup := e.Group("/api/v1/favorites")
	favGroup.GET("", favoriteHandler.GetUserFavorites)
}