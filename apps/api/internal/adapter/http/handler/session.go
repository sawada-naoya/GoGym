package handler

import (
	"gogym-api/internal/adapter/http/dto"
	httpError "gogym-api/internal/adapter/http/error"
	su "gogym-api/internal/usecase/session"
	uu "gogym-api/internal/usecase/user"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SessionHandler struct {
	su su.UseCase
	uu uu.UseCase
}

func NewSessionHandler(su su.UseCase, uu uu.UseCase) *SessionHandler {
	return &SessionHandler{su: su, uu: uu}
}

func (h *SessionHandler) Login(c echo.Context) error {
	ctx := c.Request().Context()
	slog.InfoContext(ctx, "Login Handler")

	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err := h.uu.Login(ctx, req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, httpError.ErrorResponse{
			Code:    "invalid_credentials",
			Message: "Invalid email or password",
		})
	}

	return c.NoContent(http.StatusOK)
}
