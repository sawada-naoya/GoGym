package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ReviewHandler struct {
	// TODO: Add review usecase
}

func NewReviewHandler() *ReviewHandler {
	return &ReviewHandler{}
}

// GET /api/v1/gyms/:gym_id/reviews
func (h *ReviewHandler) GetReviews(c echo.Context) error {
	// TODO: Get reviews for gym
	return c.JSON(http.StatusOK, map[string]interface{}{
		"reviews": []interface{}{},
	})
}

// POST /api/v1/gyms/:gym_id/reviews
func (h *ReviewHandler) CreateReview(c echo.Context) error {
	// TODO: Create new review
	return c.JSON(http.StatusCreated, map[string]string{
		"message": "Review created",
	})
}

// GET /api/v1/reviews/:id
func (h *ReviewHandler) GetReview(c echo.Context) error {
	// TODO: Get review by ID
	return c.JSON(http.StatusOK, map[string]interface{}{})
}

// PUT /api/v1/reviews/:id
func (h *ReviewHandler) UpdateReview(c echo.Context) error {
	// TODO: Update review
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Review updated",
	})
}

// DELETE /api/v1/reviews/:id
func (h *ReviewHandler) DeleteReview(c echo.Context) error {
	// TODO: Delete review
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Review deleted",
	})
}