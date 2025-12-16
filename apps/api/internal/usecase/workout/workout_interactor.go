package workout

import (
	"context"
	dto "gogym-api/internal/adapter/dto"
	dom "gogym-api/internal/domain/entities"
	"time"
)

type workoutInteractor struct {
	repo Repository
}

func NewWorkoutInteractor(repo Repository) WorkoutUseCase {
	return &workoutInteractor{
		repo: repo,
	}
}

func (i *workoutInteractor) GetWorkoutRecords(ctx context.Context, userID string, dateParam string) (dto.WorkoutRecordDTO, error) {
	// dateが空文字列の場合は今日のJST日付を使用
	date := dateParam
	if date == "" {
		jst, _ := time.LoadLocation("Asia/Tokyo")
		date = time.Now().In(jst).Format("2006-01-02")
	}

	records, err := i.repo.GetRecordsByDateAndPart(ctx, userID, date, nil)
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

func (i *workoutInteractor) CreateWorkoutRecord(ctx context.Context, workout dom.WorkoutRecord) error {
	// Domain logic can be added here (validation, business rules, etc.)

	// 同日同部位ならupsert、それ以外は新規作成
	err := i.repo.UpsertWorkoutRecord(ctx, workout)
	if err != nil {
		return err
	}
	return nil
}

func (i *workoutInteractor) UpsertWorkoutRecord(ctx context.Context, workout dom.WorkoutRecord) error {
	// Domain logic can be added here (validation, business rules, etc.)

	// 同日同部位ならupsert、それ以外は新規作成
	err := i.repo.UpsertWorkoutRecord(ctx, workout)
	if err != nil {
		return err
	}
	return nil
}

func (i *workoutInteractor) GetWorkoutParts(ctx context.Context, userID string) ([]dto.WorkoutPartListItemDTO, error) {
	parts, err := i.repo.GetWorkoutParts(ctx, userID)
	if err != nil {
		return nil, err
	}
	return dto.WorkoutPartsToDTO(parts), nil
}

func (i *workoutInteractor) SeedWorkoutParts(ctx context.Context, userID string) error {
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

func (i *workoutInteractor) CreateWorkoutExercise(ctx context.Context, userID string, exercises []dto.CreateWorkoutExerciseItem) error {
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

func (i *workoutInteractor) DeleteWorkoutExercise(ctx context.Context, userID string, exerciseID int64) error {
	return i.repo.DeleteWorkoutExercise(ctx, userID, exerciseID)
}

func (i *workoutInteractor) GetLastWorkoutRecord(ctx context.Context, userID string, exerciseID int64) (*dto.ExerciseDTO, error) {
	// 最後のワークアウトレコードを取得
	record, err := i.repo.GetLastWorkoutRecord(ctx, userID, exerciseID)
	if err != nil {
		return nil, err
	}

	// レコードが空の場合はnilを返す
	if record.ID == nil || len(record.Sets) == 0 {
		return nil, nil
	}

	// 該当するエクササイズIDのセットだけをフィルタリング
	var exerciseSets []dom.WorkoutSet
	var exerciseName string
	var workoutPartID *int64

	for _, set := range record.Sets {
		if int64(set.Exercise.ID) == exerciseID {
			exerciseSets = append(exerciseSets, set)
			if exerciseName == "" {
				exerciseName = set.Exercise.Name
				if set.Exercise.PartID != nil {
					partID := int64(*set.Exercise.PartID)
					workoutPartID = &partID
				}
			}
		}
	}

	// 該当するセットが見つからない場合
	if len(exerciseSets) == 0 {
		return nil, nil
	}

	// ExerciseDTOに変換して返す
	exerciseDTO := dto.ExerciseDTO{
		ID:            &exerciseID,
		Name:          exerciseName,
		WorkoutPartID: workoutPartID,
		Sets:          []dto.SetDTO{},
	}

	// セット情報を追加
	for _, set := range exerciseSets {
		var setID *int64
		if set.ID != nil {
			id := int64(*set.ID)
			setID = &id
		}

		weight := float64(set.Weight)
		reps := int(set.Reps)

		exerciseDTO.Sets = append(exerciseDTO.Sets, dto.SetDTO{
			ID:        setID,
			SetNumber: set.SetNumber,
			WeightKg:  &weight,
			Reps:      &reps,
			Note:      set.Note,
		})
	}

	return &exerciseDTO, nil
}
