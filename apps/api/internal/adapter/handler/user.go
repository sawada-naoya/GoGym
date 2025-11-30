package handler

import (
	"log/slog"
	"net/http"

	"gogym-api/internal/adapter/dto"
	uu "gogym-api/internal/usecase/user"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	uu       uu.UserUseCase
	validate *validator.Validate
}

func NewUserHandler(uu uu.UserUseCase) *UserHandler {
	return &UserHandler{
		uu:       uu,
		validate: validator.New(),
	}
}

// POST /api/v1/user
func (h *UserHandler) SignUp(c echo.Context) error {
	ctx := c.Request().Context()
	slog.InfoContext(ctx, "SignUp Handler")
	var req dto.SignUpRequest
	if err := c.Bind(&req); err != nil {
		slog.ErrorContext(ctx, "Failed to bind request", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	// バリデーション実行
	if err := h.validate.Struct(req); err != nil {
		slog.ErrorContext(ctx, "Validation failed", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "パスワードは8文字以上で、大文字・小文字・数字を含める必要があります"})
	}

	err := h.uu.SignUp(ctx, req)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to sign up user", "error", err)
		return c.JSON(http.StatusConflict, err.Error())
	}

	slog.InfoContext(ctx, "User signed up successfully")
	return c.NoContent(http.StatusCreated)
}
