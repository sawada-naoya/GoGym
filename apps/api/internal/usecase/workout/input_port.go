package workout

import (
	"context"
	dom "gogym-api/internal/domain/workout"
)

type WorkoutUseCase interface {
	GetWorkoutRecords(ctx context.Context, userID string, date string) (dom.WorkoutRecord, error)
	CreateWorkoutRecord(ctx context.Context, workout dom.WorkoutRecord) error
}
