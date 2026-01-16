package workout

import (
	"context"
	"time"

	dw "gogym-api/internal/domain/entities/workout"
)

type Repository interface {
	GetRecordsByDate(ctx context.Context, userID string, date time.Time) (dw.WorkoutRecord, error)
	CreateWorkoutRecord(ctx context.Context, workout dw.WorkoutRecord) error
	UpsertWorkoutRecord(ctx context.Context, workout dw.WorkoutRecord) error
	GetWorkoutParts(ctx context.Context, userID string) ([]dw.WorkoutPart, error)
	CreateWorkoutParts(ctx context.Context, userID string, parts []dw.WorkoutPart) error
	CountUserWorkoutParts(ctx context.Context, userID string) (int64, error)
	UpsertWorkoutExercises(ctx context.Context, userID string, exercises []dw.WorkoutExerciseRef) error
	DeleteWorkoutExercise(ctx context.Context, userID string, exerciseID int64) error
	GetLastWorkoutRecord(ctx context.Context, userID string, exerciseID int64) (dw.WorkoutRecord, error)
}
