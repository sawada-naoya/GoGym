package workout

import (
	"context"
	dom "gogym-api/internal/domain/workout"
)

type Repository interface {
	GetRecordsByDate(ctx context.Context, userID string, date string) (dom.WorkoutRecord, error)
	Create(ctx context.Context, workout dom.WorkoutRecord) error
}