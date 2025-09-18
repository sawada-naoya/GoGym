package handler

import (
	"net/http"

	"gogym-api/internal/adapter/http/dto"
	"gogym-api/internal/usecase/user"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	uu user.UserUseCase
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// POST /api/v1/users/register
func (h *UserHandler) Register(c echo.Context) error {
	var req dto.RegisterUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	

	return c.JSON(http.StatusCreated, err.Erros())
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
