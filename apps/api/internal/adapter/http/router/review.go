package router

import (
	"gogym-api/internal/adapter/http/handler"

	"github.com/labstack/echo/v4"
)

func SetupReviewRoutes(e *echo.Echo, reviewHandler *handler.ReviewHandler) {
	review := e.Group("/api/v1/gyms")
	review.GET("/:id/reviews", reviewHandler.GetReviews)
}
