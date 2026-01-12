package workout

import (
	"context"
	dto "gogym-api/internal/adapter/dto"
	dom "gogym-api/internal/domain/entities"
	"time"
)

type WorkoutUseCase interface {
	GetWorkoutRecords(ctx context.Context, userID string, date time.Time) (dto.WorkoutRecordDTO, error)
	CreateWorkoutRecord(ctx context.Context, workout dom.WorkoutRecord) error
	UpsertWorkoutRecord(ctx context.Context, workout dom.WorkoutRecord) error
	GetWorkoutParts(ctx context.Context, userID string) ([]dto.WorkoutPartListItemDTO, error)
	SeedWorkoutParts(ctx context.Context, userID string) error
	CreateWorkoutExercise(ctx context.Context, userID string, exercises []dto.CreateWorkoutExerciseItem) error
	DeleteWorkoutExercise(ctx context.Context, userID string, exerciseID int64) error
	GetLastWorkoutRecord(ctx context.Context, userID string, exerciseID int64) (*dto.ExerciseDTO, error)

	ResolveGymIDFromName(ctx context.Context, userID string, gymName string) (dom.ID, error)
}
