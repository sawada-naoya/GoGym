package router

import (
	"gogym-api/internal/adapter/http/handler"

	"github.com/labstack/echo/v4"
)

func NewRouter(e *echo.Echo, gymHandler *handler.GymHandler, userHandler *handler.UserHandler, reviewHandler *handler.ReviewHandler, favoriteHandler *handler.FavoriteHandler) *echo.Echo {
	SetupHealthRoutes(e)
	SetupGymRoutes(e, gymHandler)
	SetupUserRoutes(e, userHandler)
	SetupReviewRoutes(e, reviewHandler)
	SetupFavoriteRoutes(e, favoriteHandler)
	return e
}
