package handler

import (
	"gogym-api/internal/usecase/contact"
	"log/slog"

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
	slog.InfoContext(ctx, "PostContact Handler")

	ip := c.RealIP()
	ua := c.Request().UserAgent()

	var req ContactRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(400, "invalid request body")
	}

	// ユーザーIDはnil許容なので、存在しない場合もある
	userID, _ := c.Get("user_id").(string)
	var userIDPtr *string
	if userID != "" {
		userIDPtr = &userID
	}

	err := h.cu.SendContact(ctx, req.Email, req.Message, userIDPtr, ip, ua)
	if err != nil {
		// エラーログは出力するが、通知失敗でもユーザーには成功を返す
		// （問い合わせ自体は受け付けたとみなす）
		c.Logger().Error("failed to send contact notification:", err)
	}

	return c.NoContent(204)
}
