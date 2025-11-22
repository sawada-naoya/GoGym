package workout

import (
	"context"
	dom "gogym-api/internal/domain/workout"
)

type Repository interface {
	GetRecordsByDate(ctx context.Context, userID string, date string) (dom.WorkoutRecord, error)
	Create(ctx context.Context, workout dom.WorkoutRecord) error
	GetWorkoutParts(ctx context.Context, userID string) ([]dom.WorkoutPart, error)
	CreateWorkoutParts(ctx context.Context, userID string, parts []dom.WorkoutPart) error
	CountUserWorkoutParts(ctx context.Context, userID string) (int64, error)
	CreateWorkoutExercises(ctx context.Context, userID string, exercises []dom.WorkoutExerciseRef) error
}
