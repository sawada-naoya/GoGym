package handler

import (
	"log/slog"
	"net/http"

	"gogym-api/internal/adapter/http/dto"
	httpError "gogym-api/internal/adapter/http/error"
	uu "gogym-api/internal/usecase/user"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	uu uu.UseCase
}

func NewUserHandler(uu uu.UseCase) *UserHandler {
	return &UserHandler{uu: uu}
}

// POST /api/v1/user
func (h *UserHandler) SignUp(c echo.Context) error {
	ctx := c.Request().Context()
	slog.InfoContext(ctx, "SignUp Handler")
	var req dto.SignUpRequest
	if err := c.Bind(&req); err != nil {
		slog.ErrorContext(ctx, "Failed to bind request", "error", err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err := h.uu.SignUp(ctx, req)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to sign up user", "error", err)
		return c.JSON(http.StatusConflict, httpError.ErrorResponse{
			Code:    "email_already_exists",
			Message: err.Error(),
		})
	}

	slog.InfoContext(ctx, "User signed up successfully")
	return c.NoContent(http.StatusCreated)
}
