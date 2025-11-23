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

	// 6部位のシードデータを作成
	defaultParts := []dom.WorkoutPart{
		{Name: "胸", Owner: &ownerULID},
		{Name: "肩", Owner: &ownerULID},
		{Name: "背中", Owner: &ownerULID},
		{Name: "腕", Owner: &ownerULID},
		{Name: "脚", Owner: &ownerULID},
		{Name: "その他", Owner: &ownerULID},
	}

	return i.repo.CreateWorkoutParts(ctx, userID, defaultParts)
}

func (i *interactor) CreateWorkoutExercise(ctx context.Context, userID string, exercises []dto.CreateWorkoutExerciseItem) error {
	// ユーザーIDをULIDに変換
	ownerULID := dom.ULID(userID)

	// DTOをドメインモデルに変換
	domainExercises := make([]dom.WorkoutExerciseRef, 0, len(exercises))
	for _, ex := range exercises {
		partID := dom.ID(ex.WorkoutPartID)
		exerciseRef := dom.WorkoutExerciseRef{
			Name:   ex.Name,
			PartID: &partID,
			Owner:  &ownerULID, // ユーザー作成の種目なのでOwnerを設定
		}

		// IDがある場合は設定（update対象）
		if ex.ID != nil {
			exerciseRef.ID = dom.ID(*ex.ID)
		}

		domainExercises = append(domainExercises, exerciseRef)
	}

	// リポジトリを通じてデータベースに保存（upsert）
	return i.repo.UpsertWorkoutExercises(ctx, userID, domainExercises)
}

func (i *interactor) DeleteWorkoutExercise(ctx context.Context, userID string, exerciseID int64) error {
	return i.repo.DeleteWorkoutExercise(ctx, userID, exerciseID)
}
