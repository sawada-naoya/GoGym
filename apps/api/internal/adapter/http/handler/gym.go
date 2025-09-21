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
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	foundGym, err := h.gu.GetGym(ctx, gym.ID(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	response := gu.ToGymResponse(*foundGym)
	return c.JSON(http.StatusOK, response)
}


func (h *GymHandler) GetRecommendedGyms(c echo.Context) error {
	ctx := c.Request().Context()
	slog.InfoContext(ctx, "GetRecommendedGyms Handler")

	gyms, err := h.gu.RecommendGyms(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := make([]gu.GymResponse, len(gyms))
	for i, gym := range gyms {
		response[i] = gu.ToGymResponse(gym)
	}

	return c.JSON(http.StatusOK, response)
}
