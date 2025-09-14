package router

import (
	"gogym-api/internal/adapter/http/handler"

	"github.com/labstack/echo/v4"
)

func SetupReviewRoutes(e *echo.Echo, reviewHandler *handler.ReviewHandler) {
	gymGroup := e.Group("/api/v1/gyms")
	gymGroup.GET("/:gym_id/reviews", reviewHandler.GetReviews)
}
