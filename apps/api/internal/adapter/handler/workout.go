package handler

import (
	"log/slog"
	"net/http"

	"gogym-api/internal/adapter/dto"
	wu "gogym-api/internal/usecase/workout"

	"github.com/labstack/echo/v4"
)

type WorkoutHandler struct {
	wu wu.WorkoutUseCase
}

func NewWorkoutHandler(wu wu.WorkoutUseCase) *WorkoutHandler {
	return &WorkoutHandler{
		wu: wu,
	}
}

func (h *WorkoutHandler) GetWorkoutRecords(c echo.Context) error {
	ctx := c.Request().Context()
	slog.InfoContext(ctx, "GetWorkoutRecords Handler")

	// ユーザーIDを取得（認証ミドルウェアでコンテキストにセットされている想定）
	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		slog.ErrorContext(ctx, "User ID not found in context")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	date := c.QueryParam("date")
	if date == "" {
		slog.ErrorContext(ctx, "Date query parameter is required")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Date query parameter is required"})
	}

	domainRecord, err := h.wu.GetWorkoutRecords(ctx, userID, date)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get workout records", "userID", userID, "date", date, "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Domain → DTO変換
	response := dto.WorkoutRecordToDTO(&domainRecord)

	slog.InfoContext(ctx, "Successfully retrieved workout records", "userID", userID, "date", date)
	return c.JSON(http.StatusOK, response)
}
