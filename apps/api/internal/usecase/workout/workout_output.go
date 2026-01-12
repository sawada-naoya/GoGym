package workout

import (
	"context"
	dom "gogym-api/internal/domain/entities"
	"time"
)

type Repository interface {
	GetRecordsByDate(ctx context.Context, userID string, date time.Time) (dom.WorkoutRecord, error)
	CreateWorkoutRecord(ctx context.Context, workout dom.WorkoutRecord) error
	UpsertWorkoutRecord(ctx context.Context, workout dom.WorkoutRecord) error
	GetWorkoutParts(ctx context.Context, userID string) ([]dom.WorkoutPart, error)
	CreateWorkoutParts(ctx context.Context, userID string, parts []dom.WorkoutPart) error
	CountUserWorkoutParts(ctx context.Context, userID string) (int64, error)
	UpsertWorkoutExercises(ctx context.Context, userID string, exercises []dom.WorkoutExerciseRef) error
	DeleteWorkoutExercise(ctx context.Context, userID string, exerciseID int64) error
	GetLastWorkoutRecord(ctx context.Context, userID string, exerciseID int64) (dom.WorkoutRecord, error)
}
