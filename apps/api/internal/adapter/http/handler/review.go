package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	ru "gogym-api/internal/usecase/review"

	"github.com/labstack/echo/v4"
)

type ReviewHandler struct {
	ru ru.UseCase
}

func NewReviewHandler(ru ru.UseCase) *ReviewHandler {
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
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, response)
}
