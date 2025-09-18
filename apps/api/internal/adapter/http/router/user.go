package router

import (
	"gogym-api/internal/adapter/http/handler"

	"github.com/labstack/echo/v4"
)

func SetupUserRoutes(e *echo.Echo, userHandler *handler.UserHandler) {
	userGroup := e.Group("/api/v1/users")

	userGroup.POST("/register", userHandler.Register)
	userGroup.POST("/login", userHandler.Login)

	userGroup.GET("/:id", userHandler.GetUser)
	userGroup.PUT("/:id", userHandler.UpdateUser)

	userGroup.GET("/:id/favorites", userHandler.GetUserFavorites)
	userGroup.GET("/:id/reviews", userHandler.GetUserReviews)
}
