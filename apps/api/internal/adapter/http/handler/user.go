package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	// TODO: Add user usecase
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// POST /api/v1/users/register
func (h *UserHandler) Register(c echo.Context) error {
	// TODO: User registration
	return c.JSON(http.StatusCreated, map[string]string{
		"message": "User registered",
	})
}

// POST /api/v1/users/login
func (h *UserHandler) Login(c echo.Context) error {
	// TODO: User login
	return c.JSON(http.StatusOK, map[string]string{
		"token": "jwt_token_here",
	})
}

// GET /api/v1/users/:id
func (h *UserHandler) GetUser(c echo.Context) error {
	// TODO: Get user profile
	return c.JSON(http.StatusOK, map[string]interface{}{})
}

// PUT /api/v1/users/:id
func (h *UserHandler) UpdateUser(c echo.Context) error {
	// TODO: Update user profile
	return c.JSON(http.StatusOK, map[string]string{
		"message": "User updated",
	})
}

// GET /api/v1/users/:id/favorites
func (h *UserHandler) GetUserFavorites(c echo.Context) error {
	// TODO: Get user's favorite gyms
	return c.JSON(http.StatusOK, map[string]interface{}{
		"favorites": []interface{}{},
	})
}

// GET /api/v1/users/:id/reviews
func (h *UserHandler) GetUserReviews(c echo.Context) error {
	// TODO: Get user's reviews
	return c.JSON(http.StatusOK, map[string]interface{}{
		"reviews": []interface{}{},
	})
}