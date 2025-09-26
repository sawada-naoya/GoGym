package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"gogym-api/internal/domain/gym"
	gu "gogym-api/internal/usecase/gym"

	"github.com/labstack/echo/v4"
)

type GymHandler struct {
	gu gu.UseCase
}

func NewGymHandler(gu gu.UseCase) *GymHandler {
	return &GymHandler{
		gu: gu,
	}
}

func (h *GymHandler) GetGym(c echo.Context) error {
	ctx := c.Request().Context()
	slog.InfoContext(ctx, "GetGymDetail Handler")

	param := c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		slog.ErrorContext(ctx, "Invalid gym ID", "id", param, "error", err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	response, err := h.gu.GetGym(ctx, gym.ID(id))
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get gym detail", "id", id, "error", err)
		return c.JSON(http.StatusNotFound, err.Error())
	}

	slog.InfoContext(ctx, "Successfully retrieved gym detail", "id", id)
	return c.JSON(http.StatusOK, response)
}

func (h *GymHandler) GetRecommendedGyms(c echo.Context) error {
	ctx := c.Request().Context()
	slog.InfoContext(ctx, "GetRecommendedGyms Handler")

	responses, err := h.gu.RecommendGyms(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get recommended gyms", "error", err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	slog.InfoContext(ctx, "Successfully retrieved recommended gyms")
	return c.JSON(http.StatusOK, responses)
}
