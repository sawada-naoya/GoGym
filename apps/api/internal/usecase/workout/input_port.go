package workout

import (
	"context"
	dto "gogym-api/internal/adapter/dto"
	dom "gogym-api/internal/domain/workout"
)

type WorkoutUseCase interface {
	GetWorkoutRecords(ctx context.Context, userID string, date string) (dto.WorkoutRecordDTO, error)
	CreateWorkoutRecord(ctx context.Context, workout dom.WorkoutRecord) error
	GetWorkoutParts(ctx context.Context, userID string) ([]dto.WorkoutPartListItemDTO, error)
	SeedWorkoutParts(ctx context.Context, userID string) error
	CreateWorkoutExercise(ctx context.Context, userID string, exercises []dto.CreateWorkoutExerciseItem) error
}
