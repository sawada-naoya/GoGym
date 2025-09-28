package router

import (
	"gogym-api/internal/adapter/handler"

	"github.com/labstack/echo/v4"
)

func ReviewRoutes(e *echo.Group, reviewHandler *handler.ReviewHandler) {
	review := e.Group("gyms")
	review.GET("/:id/reviews", reviewHandler.GetReviews)
}
