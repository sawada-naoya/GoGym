package handler

import (
	"gogym-api/internal/adapter/http/dto"
	su "gogym-api/internal/usecase/session"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SessionHandler struct {
	su su.UseCase
}

func NewSessionHandler(su su.UseCase) *SessionHandler {
	return &SessionHandler{su: su}
}

func (h *SessionHandler) Login(c echo.Context) error {
	ctx := c.Request().Context()
	slog.InfoContext(ctx, "Login Handler")

	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		slog.ErrorContext(ctx, "Failed to bind request", "error", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	// User認証
	err := h.su.Login(ctx, req)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to login user", "email", req.Email, "error", err)
		return c.JSON(http.StatusUnauthorized, err.Error())
	}

	// Session作成
	tokens, err := h.su.CreateSession(ctx, req.Email)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create session", "email", req.Email, "error", err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	slog.InfoContext(ctx, "User logged in successfully", "email", req.Email)
	return c.JSON(http.StatusOK, tokens)
}
