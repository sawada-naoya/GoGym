package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	// TODO: ユースケースを追加
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// Signup handles user signup requests
func (h *UserHandler) Signup(c echo.Context) error {
	// TODO: サインアップ実装
	return c.JSON(http.StatusOK, map[string]string{
		"message": "User signup endpoint - coming soon",
	})
}