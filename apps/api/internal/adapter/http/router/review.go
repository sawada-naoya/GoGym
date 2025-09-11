package router

import (
	"gogym-api/internal/adapter/http/handler"

	"github.com/labstack/echo/v4"
)

func SetupReviewRoutes(e *echo.Echo, reviewHandler *handler.ReviewHandler) {
	// Nested under gyms
	gymGroup := e.Group("/api/v1/gyms")
	gymGroup.GET("/:gym_id/reviews", reviewHandler.GetReviews)
	gymGroup.POST("/:gym_id/reviews", reviewHandler.CreateReview)
	
	// Direct review access
	reviewGroup := e.Group("/api/v1/reviews")
	reviewGroup.GET("/:id", reviewHandler.GetReview)
	reviewGroup.PUT("/:id", reviewHandler.UpdateReview)
	reviewGroup.DELETE("/:id", reviewHandler.DeleteReview)
}