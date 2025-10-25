package workout

import (
	"context"
	dom "gogym-api/internal/domain/workout"
)

type interactor struct {
	repo Repository
}

func NewInteractor(repo Repository) WorkoutUseCase {
	return &interactor{
		repo: repo,
	}
}

func (i *interactor) GetWorkoutRecords(ctx context.Context, userID string, date string) (dom.WorkoutRecord, error) {
	records, err := i.repo.GetRecordsByDate(ctx, userID, date)
	if err != nil {
		return dom.WorkoutRecord{}, err
	}
	return records, nil
}