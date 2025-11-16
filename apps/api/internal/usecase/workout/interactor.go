package workout

import (
	"context"
	dto "gogym-api/internal/adapter/dto"
	dom "gogym-api/internal/domain/workout"
	"time"
)

type interactor struct {
	repo Repository
}

func NewInteractor(repo Repository) WorkoutUseCase {
	return &interactor{
		repo: repo,
	}
}

func (i *interactor) GetWorkoutRecords(ctx context.Context, userID string, dateParam string) (dto.WorkoutRecordDTO, error) {
	// dateが空文字列の場合は今日のJST日付を使用
	date := dateParam
	if date == "" {
		jst, _ := time.LoadLocation("Asia/Tokyo")
		date = time.Now().In(jst).Format("2006-01-02")
	}

	records, err := i.repo.GetRecordsByDate(ctx, userID, date)
	if err != nil {
		return dto.WorkoutRecordDTO{}, err
	}

	// レコードが空（IDがnil）の場合は、日付だけ設定したDTOを返す
	if records.ID == nil {
		return dto.WorkoutRecordDTO{
			PerformedDate: date,
			Place:         "",
			Exercises:     []dto.ExerciseDTO{},
		}, nil
	}

	response := dto.WorkoutRecordToDTO(&records)
	return *response, nil
}

func (i *interactor) CreateWorkoutRecord(ctx context.Context, workout dom.WorkoutRecord) error {
	// Domain logic can be added here (validation, business rules, etc.)

	err := i.repo.Create(ctx, workout)
	if err != nil {
		return err
	}
	return nil
}

func (i *interactor) GetWorkoutParts(ctx context.Context, userID string) ([]dto.WorkoutPartListItemDTO, error) {
	parts, err := i.repo.GetWorkoutParts(ctx, userID)
	if err != nil {
		return nil, err
	}
	return dto.WorkoutPartsToDTO(parts), nil
}

func (i *interactor) SeedWorkoutParts(ctx context.Context, userID string) error {
	// すでにユーザーの部位が存在するかチェック
	count, err := i.repo.CountUserWorkoutParts(ctx, userID)
	if err != nil {
		return err
	}

	// すでに存在する場合は何もしない（idempotent）
	if count > 0 {
		return nil
	}

	// ULID型に変換
	ownerULID := dom.ULID(userID)

	// 5部位のシードデータを作成
	defaultParts := []dom.WorkoutPart{
		{Name: "胸", IsDefault: true, Owner: &ownerULID},
		{Name: "肩", IsDefault: true, Owner: &ownerULID},
		{Name: "背中", IsDefault: true, Owner: &ownerULID},
		{Name: "腕", IsDefault: true, Owner: &ownerULID},
		{Name: "脚", IsDefault: true, Owner: &ownerULID},
	}

	return i.repo.CreateWorkoutParts(ctx, userID, defaultParts)
}
