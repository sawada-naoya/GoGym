package handler

import (
	"gogym-api/internal/usecase/contact"

	"github.com/labstack/echo/v4"
)

type ContactHandler struct {
	cu contact.ContactUseCase
}

func NewContactHandler(cu contact.ContactUseCase) *ContactHandler {
	return &ContactHandler{
		cu: cu,
	}
}

type ContactRequest struct {
	Email   string `json:"email"`
	Message string `json:"message"`
}

func (h *ContactHandler) PostContact(c echo.Context) error {
	ctx := c.Request().Context()

	var req ContactRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(400, "invalid request body")
	}

	in := contact.SendContactInput{
		Email:   req.Email,
		Message: req.Message,
		UserID:  nil, // 未ログイン時の問い合わせ対応
		IP:      c.RealIP(),
		UA:      c.Request().UserAgent(),
	}

	err := h.cu.SendContact(ctx, in)
	if err != nil {
		// エラーログは出力するが、通知失敗でもユーザーには成功を返す
		// （問い合わせ自体は受け付けたとみなす）
		c.Logger().Error("failed to send contact notification:", err)
	}

	return c.NoContent(204)
}
