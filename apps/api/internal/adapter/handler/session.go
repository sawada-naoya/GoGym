package handler

import (
	"gogym-api/internal/adapter/dto"
	su "gogym-api/internal/application/session"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SessionHandler struct {
	su su.SessionUseCase
}

func NewSessionHandler(su su.SessionUseCase) *SessionHandler {
	return &SessionHandler{su: su}
}

func (h *SessionHandler) Login(c echo.Context) error {
	ctx := c.Request().Context()

	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	// User認証
	err := h.su.Login(ctx, req)
	if err != nil {
		slog.Error("Login failed", "email", req.Email, "error", err.Error())
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
	}

	// Session作成
	tokens, err := h.su.CreateSession(ctx, req.Email)
	if err != nil {
		slog.Error("Failed to create session", "email", req.Email, "error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create session"})
	}

	return c.JSON(http.StatusOK, tokens)
}

func (h *SessionHandler) RefreshToken(c echo.Context) error {
	ctx := c.Request().Context()

	var req dto.RefreshRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	// トークンリフレッシュ
	tokens, err := h.su.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		slog.Error("Token refresh failed", "error", err.Error())
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid or expired refresh token"})
	}

	return c.JSON(http.StatusOK, tokens)
}
