package handler

import (
	"fmt"
	"gogym-api/internal/adapter/dto"
	"gogym-api/internal/util"
	"log/slog"
	"net/http"

	dom "gogym-api/internal/domain/entities"
	wu "gogym-api/internal/application/workout"

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

	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		slog.ErrorContext(ctx, "User ID not found in context")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	dateStr := c.QueryParam("date")
	date, err := util.ParseJSTDateOrToday(dateStr)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to parse date", "dateStr", dateStr, "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid date format"})
	}

	response, err := h.wu.GetWorkoutRecords(ctx, userID, date)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get workout records", "userID", userID, "date", date, "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

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

	// gym_name から gym_id を解決（gym_name優先、なければgym_idをそのまま使用）
	if req.GymName != nil && *req.GymName != "" {
		gymID, err := h.wu.ResolveGymIDFromName(ctx, userID, *req.GymName)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to resolve gym_id from gym_name", "userID", userID, "gymName", *req.GymName, "error", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to resolve gym"})
		}
		domainRecord.GymID = &gymID
	}

	err = h.wu.CreateWorkoutRecord(ctx, *domainRecord)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create workout record", "userID", userID, "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Workout record created successfully"})
}

func (h *WorkoutHandler) UpdateWorkoutRecord(c echo.Context) error {
	ctx := c.Request().Context()
	slog.InfoContext(ctx, "UpdateWorkoutRecord Handler")

	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		slog.ErrorContext(ctx, "User ID not found in context")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	// パスパラメータから record ID を取得（使用しない可能性もあるが、一応取得）
	// recordID := c.Param("id")

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

	// gym_name から gym_id を解決（gym_name優先、なければgym_idをそのまま使用）
	if req.GymName != nil && *req.GymName != "" {
		gymID, err := h.wu.ResolveGymIDFromName(ctx, userID, *req.GymName)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to resolve gym_id from gym_name", "userID", userID, "gymName", *req.GymName, "error", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to resolve gym"})
		}
		domainRecord.GymID = &gymID
	}

	// Upsert を使用（同日同部位の場合は更新、異なる場合は新規作成）
	err = h.wu.UpsertWorkoutRecord(ctx, *domainRecord)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to update workout record", "userID", userID, "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Workout record updated successfully"})
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

func (h *WorkoutHandler) SeedWorkoutParts(c echo.Context) error {
	ctx := c.Request().Context()
	slog.InfoContext(ctx, "SeedWorkoutParts Handler")

	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		slog.ErrorContext(ctx, "User ID not found in context")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	err := h.wu.SeedWorkoutParts(ctx, userID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to seed workout parts", "userID", userID, "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Workout parts seeded successfully"})
}

func (h *WorkoutHandler) CreateWorkoutExercise(c echo.Context) error {
	ctx := c.Request().Context()
	slog.InfoContext(ctx, "CreateWorkoutExercise Handler")

	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		slog.ErrorContext(ctx, "User ID not found in context")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	var req dto.CreateWorkoutExerciseRequest
	err := c.Bind(&req)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to bind request", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	err = h.wu.CreateWorkoutExercise(ctx, userID, req.Exercises)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create workout exercises", "userID", userID, "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Workout exercises created successfully"})
}

func (h *WorkoutHandler) DeleteWorkoutExercise(c echo.Context) error {
	ctx := c.Request().Context()
	slog.InfoContext(ctx, "DeleteWorkoutExercise Handler")

	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		slog.ErrorContext(ctx, "User ID not found in context")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	exerciseIDStr := c.Param("id")
	var exerciseID int64
	if _, err := fmt.Sscanf(exerciseIDStr, "%d", &exerciseID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid exercise ID format"})
	}

	err := h.wu.DeleteWorkoutExercise(ctx, userID, exerciseID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to delete workout exercise", "userID", userID, "exerciseID", exerciseID, "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Workout exercise deleted successfully"})
}

func (h *WorkoutHandler) GetLastWorkoutRecord(c echo.Context) error {
	ctx := c.Request().Context()
	slog.InfoContext(ctx, "GetLastWorkoutRecord Handler")

	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		slog.ErrorContext(ctx, "User ID not found in context")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	exerciseIDStr := c.Param("id")
	var exerciseID int64
	if _, err := fmt.Sscanf(exerciseIDStr, "%d", &exerciseID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid exercise ID format"})
	}

	response, err := h.wu.GetLastWorkoutRecord(ctx, userID, exerciseID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get last workout record", "userID", userID, "exerciseID", exerciseID, "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, response)
}
