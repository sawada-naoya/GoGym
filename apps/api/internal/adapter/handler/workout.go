package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"gogym-api/internal/adapter/dto"
	dom "gogym-api/internal/domain/workout"
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

	// dateクエリパラメータを取得（空文字列の場合はUseCaseで今日のJST日付を使用）
	date := c.QueryParam("date")

	response, err := h.wu.GetWorkoutRecords(ctx, userID, date)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get workout records", "userID", userID, "date", date, "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	slog.InfoContext(ctx, "Successfully retrieved workout records", "userID", userID, "date", date)
	return c.JSON(http.StatusOK, response)
}

func (h *WorkoutHandler) CreateWorkoutRecord(c echo.Context) error {
	ctx := c.Request().Context()
	slog.InfoContext(ctx, "CreateWorkoutRecord Handler")

	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		slog.ErrorContext(ctx, "User ID not found in context")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	var req dto.WorkoutRecordDTO
	err := c.Bind(&req)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to bind request", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	// dto → domain変換
	domainRecord, err := dto.WorkoutRecordDTOToDomain(&req)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to convert DTO to domain model", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Invalid request data: %v", err)})
	}

	// Set userID
	domainRecord.UserID = dom.ULID(userID)

	err = h.wu.CreateWorkoutRecord(ctx, *domainRecord)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create workout record", "userID", userID, "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	slog.InfoContext(ctx, "Successfully created workout record", "userID", userID)
	return c.JSON(http.StatusCreated, map[string]string{"message": "Workout record created successfully"})
}

func (h *WorkoutHandler) GetWorkoutParts(c echo.Context) error {
	ctx := c.Request().Context()
	slog.InfoContext(ctx, "GetWorkoutParts Handler")

	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		slog.ErrorContext(ctx, "User ID not found in context")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	parts, err := h.wu.GetWorkoutParts(ctx, userID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get workout parts", "userID", userID, "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	slog.InfoContext(ctx, "Successfully retrieved workout parts", "userID", userID, "count", len(parts))
	return c.JSON(http.StatusOK, parts)
}
