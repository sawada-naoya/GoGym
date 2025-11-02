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

func (i *interactor) CreateWorkoutRecord(ctx context.Context, workout dom.WorkoutRecord) error {
	// Domain logic can be added here (validation, business rules, etc.)

	err := i.repo.Create(ctx, workout)
	if err != nil {
		return err
	}
	return nil
}
