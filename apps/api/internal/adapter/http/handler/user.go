package handler

import (
	"log/slog"
	"net/http"

	httpError "gogym-api/internal/adapter/http/error"
	uc "gogym-api/internal/usecase/user"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	uc uc.UseCase
}

func NewUserHandler(usecase uc.UseCase) *UserHandler {
	return &UserHandler{uc: usecase}
}

// POST /api/v1/user
func (h *UserHandler) SignUp(c echo.Context) error {
	ctx := c.Request().Context()
	slog.InfoContext(ctx, "SignUp Handler")
	var req uc.SignUpRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err := h.uc.SignUp(ctx, req)
	if err != nil {
		return c.JSON(http.StatusConflict, httpError.ErrorResponse{
			Code:    "email_already_exists",
			Message: err.Error(),
		})
	}

	return c.NoContent(http.StatusCreated)
}
