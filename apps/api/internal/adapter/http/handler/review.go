package handler

import (
	"net/http"
	"strconv"

	"gogym-api/internal/adapter/http/dto"
	"gogym-api/internal/usecase/review"

	"github.com/labstack/echo/v4"
)

type ReviewHandler struct {
	ru *review.UseCase
}

func NewReviewHandler(ru *review.UseCase) *ReviewHandler {
	return &ReviewHandler{
		ru: ru,
	}
}

// GET /api/v1/gyms/:gym_id/reviews
func (h *ReviewHandler) GetReviews(c echo.Context) error {
	ctx := c.Request().Context()
	gymID, err := strconv.ParseInt(c.Param("gym_id"), 10, 64)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	cursor := c.QueryParam("cursor")
	limit := 10
	if v := c.QueryParam("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			limit = n
		}
	}

	reviews, next, err := h.ru.GetReviewsByGymID(ctx, gymID, cursor, limit)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	reviewDTOs := make([]dto.ReviewResponse, len(reviews))
	for i, review := range reviews {
		reviewDTOs[i] = dto.FromReviewEntity(&review)
	}

	var nextCursor *string
	if next != "" {
		nextCursor = &next
	}

	response := dto.ReviewListResponse{
		Reviews:    reviewDTOs,
		NextCursor: nextCursor,
	}

	return c.JSON(http.StatusOK, response)
}
