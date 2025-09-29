package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	ru "gogym-api/internal/usecase/review"

	"github.com/labstack/echo/v4"
)

type ReviewHandler struct {
	ru ru.ReviewUseCase
}

func NewReviewHandler(ru ru.ReviewUseCase) *ReviewHandler {
	return &ReviewHandler{
		ru: ru,
	}
}

// GET /api/v1/gyms/:gym_id/reviews
func (h *ReviewHandler) GetReviews(c echo.Context) error {
	ctx := c.Request().Context()
	slog.InfoContext(ctx, "GetReviews Handler")

	gymID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		slog.ErrorContext(ctx, "Invalid gym ID", "id", c.Param("id"), "error", err)
		return c.NoContent(http.StatusBadRequest)
	}

	cursor := c.QueryParam("cursor")
	limit := 10
	if v := c.QueryParam("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			limit = n
		}
	}

	response, err := h.ru.GetReviewsByGymID(ctx, gymID, cursor, limit)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get reviews", "gym_id", gymID, "error", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	slog.InfoContext(ctx, "Successfully retrieved reviews", "gym_id", gymID)
	return c.JSON(http.StatusOK, response)
}
