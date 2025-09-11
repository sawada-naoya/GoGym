package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type FavoriteHandler struct {
	// TODO: Add favorite usecase
}

func NewFavoriteHandler() *FavoriteHandler {
	return &FavoriteHandler{}
}

// POST /api/v1/gyms/:gym_id/favorite
func (h *FavoriteHandler) CreateFavorite(c echo.Context) error {
	// TODO: Add gym to favorites
	return c.JSON(http.StatusCreated, map[string]string{
		"message": "Added to favorites",
	})
}

// DELETE /api/v1/gyms/:gym_id/favorite
func (h *FavoriteHandler) DeleteFavorite(c echo.Context) error {
	// TODO: Remove gym from favorites
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Removed from favorites",
	})
}

// GET /api/v1/favorites
func (h *FavoriteHandler) GetUserFavorites(c echo.Context) error {
	// TODO: Get current user's favorites
	return c.JSON(http.StatusOK, map[string]interface{}{
		"favorites": []interface{}{},
	})
}