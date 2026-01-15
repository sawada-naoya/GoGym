package workout

import (
	"context"
	"errors"
	"gogym-api/internal/util"
	"time"

	dto "gogym-api/internal/adapter/dto"
	dom "gogym-api/internal/domain/entities"
	gymUsecase "gogym-api/internal/application/gym"
)

type workoutInteractor struct {
	repo    Repository
	gymRepo gymUsecase.Repository
}

func NewWorkoutInteractor(repo Repository, gymRepo gymUsecase.Repository) WorkoutUseCase {
	return &workoutInteractor{
		repo:    repo,
		gymRepo: gymRepo,
	}
}

func (i *workoutInteractor) GetWorkoutRecords(ctx context.Context, userID string, date time.Time) (dto.WorkoutRecordDTO, error) {
	// repositoryからドメインエンティティを取得
	domainRecord, err := i.repo.GetRecordsByDate(ctx, userID, date)
	if err != nil {
		return dto.WorkoutRecordDTO{}, err
	}

	// レコードが空（IDがnil）の場合は、日付だけ設定したDTOを返す
	if domainRecord.ID == nil {
		return dto.WorkoutRecordDTO{
			PerformedDate: util.FormatJSTDate(date),
			Parts:         []dto.WorkoutPartGroupDTO{},
		}, nil
	}

	// ドメインエンティティをDTOに変換
	response := dto.WorkoutDomainToDTO(&domainRecord)
	if response == nil {
		return dto.WorkoutRecordDTO{}, errors.New("failed to convert domain record to DTO")
	}

	return *response, nil
}

func (i *workoutInteractor) CreateWorkoutRecord(ctx context.Context, workout dom.WorkoutRecord) error {
	// 同日同部位ならupsert、それ以外は新規作成
	err := i.repo.UpsertWorkoutRecord(ctx, workout)
	if err != nil {
		return err
	}
	return nil
}

func (i *workoutInteractor) UpsertWorkoutRecord(ctx context.Context, workout dom.WorkoutRecord) error {
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

	// 6部位のシードデータを作成（多言語対応）
	defaultParts := []dom.WorkoutPart{
		{
			Key:   "chest",
			Owner: &ownerULID,
			Translations: []dom.WorkoutPartTranslation{
				{Locale: "ja", Name: "胸"},
				{Locale: "en", Name: "Chest"},
			},
		},
		{
			Key:   "shoulders",
			Owner: &ownerULID,
			Translations: []dom.WorkoutPartTranslation{
				{Locale: "ja", Name: "肩"},
				{Locale: "en", Name: "Shoulders"},
			},
		},
		{
			Key:   "back",
			Owner: &ownerULID,
			Translations: []dom.WorkoutPartTranslation{
				{Locale: "ja", Name: "背中"},
				{Locale: "en", Name: "Back"},
			},
		},
		{
			Key:   "arms",
			Owner: &ownerULID,
			Translations: []dom.WorkoutPartTranslation{
				{Locale: "ja", Name: "腕"},
				{Locale: "en", Name: "Arms"},
			},
		},
		{
			Key:   "legs",
			Owner: &ownerULID,
			Translations: []dom.WorkoutPartTranslation{
				{Locale: "ja", Name: "脚"},
				{Locale: "en", Name: "Legs"},
			},
		},
		{
			Key:   "others",
			Owner: &ownerULID,
			Translations: []dom.WorkoutPartTranslation{
				{Locale: "ja", Name: "その他"},
				{Locale: "en", Name: "Others"},
			},
		},
	}

	return i.repo.CreateWorkoutParts(ctx, userID, defaultParts)
}

func (i *workoutInteractor) CreateWorkoutExercise(ctx context.Context, userID string, exercises []dto.CreateWorkoutExerciseItem) error {
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

// ResolveGymIDFromName resolves gym_name to gym_id (finds or creates)
func (i *workoutInteractor) ResolveGymIDFromName(ctx context.Context, userID string, gymName string) (dom.ID, error) {
	// Normalize gym name
	normalizedName := dom.NormalizeName(gymName)
	if normalizedName == "" {
		return 0, errors.New("gym name cannot be empty")
	}

	// Try to find existing gym
	gym, err := i.gymRepo.FindByNormalizedName(ctx, userID, normalizedName)
	if err == nil && gym != nil {
		return dom.ID(gym.ID), nil
	}

	// If not found, create new gym
	if errors.Is(err, gymUsecase.ErrNotFound) {
		gym, err = i.gymRepo.CreateGym(ctx, userID, gymName, normalizedName)
		if err != nil {
			return 0, err
		}
		return dom.ID(gym.ID), nil
	}

	return 0, err
}
