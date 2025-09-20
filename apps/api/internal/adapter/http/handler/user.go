package handler

import (
	"net/http"

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
	var body struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	req := uc.SignUpRequest{
		Name: body.Name, Email: body.Email, Password: body.Password,
	}

	res, err := h.uc.SignUp(ctx, req)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"userId": res.UserID,
	})
}
