package workout

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	dom "gogym-api/internal/domain/entities"
	dw "gogym-api/internal/domain/entities/workout"
)

func TestWorkoutInteractor_SeedWorkoutParts(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := NewMockRepository(ctrl)

	uc := NewWorkoutInteractor(repo, nil)

	ctx := context.Background()
	userID := "01FGZ9K6TV3J5ZZZQX6Z9X6K7W" // ULID

	t.Run("正常系: 部位が未登録の場合、シードデータを登録する", func(t *testing.T) {
		t.Parallel()

		repo.EXPECT().
			CountUserWorkoutParts(gomock.Any(), userID).
			Return(int64(0), nil)

		repo.EXPECT().
			CreateWorkoutParts(gomock.Any(), userID, gomock.Any()).
			DoAndReturn(func(ctx context.Context, _ string, parts []dw.WorkoutPart) error {
				require.Len(t, parts, 6)

				keys := []string{parts[0].Key, parts[1].Key, parts[2].Key, parts[3].Key, parts[4].Key, parts[5].Key}
				require.Equal(t, []string{"chest", "shoulders", "back", "arms", "legs", "others"}, keys)

				for _, p := range parts {
					require.NotNil(t, p.Owner)
					require.Len(t, p.Translations, 2)
				}
				return nil
			})

		err := uc.SeedWorkoutParts(ctx, userID)
		require.NoError(t, err)
	})

	t.Run("正常系: 部位が既に登録されている場合、何もしない", func(t *testing.T) {
		t.Parallel()

		repo.EXPECT().
			CountUserWorkoutParts(gomock.Any(), userID).
			Return(int64(6), nil)

		// CreateWorkoutPartsは呼ばれないはず（モックなし）

		err := uc.SeedWorkoutParts(ctx, userID)
		require.NoError(t, err)
	})
}

func TestWorkoutInteractor_GetLastWorkoutRecord(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := NewMockRepository(ctrl)

	uc := NewWorkoutInteractor(repo, nil)

	ctx := context.Background()
	userID := "01FGZ9K6TV3J5ZZZQX6Z9X6K7W" // ULID
	exerciseID := int64(10)

	t.Run("異常系: レコードが存在しない場合、nilを返す", func(t *testing.T) {
		t.Parallel()

		repo.EXPECT().
			GetLastWorkoutRecord(gomock.Any(), userID, exerciseID).
			Return(dw.WorkoutRecord{}, nil)

		record, err := uc.GetLastWorkoutRecord(ctx, userID, exerciseID)
		require.NoError(t, err)
		require.Nil(t, record)
	})

	t.Run("異常系: セットが存在しない場合、nilを返す", func(t *testing.T) {
		t.Parallel()

		id := dom.ID(1)
		repo.EXPECT().
			GetLastWorkoutRecord(gomock.Any(), userID, exerciseID).
			Return(dw.WorkoutRecord{
				ID:   &id,
				Sets: nil,
			}, nil)

		record, err := uc.GetLastWorkoutRecord(ctx, userID, exerciseID)
		require.NoError(t, err)
		require.Nil(t, record)
	})

	t.Run("正常系: レコードが存在する場合、DTOを返す", func(t *testing.T) {
		t.Parallel()

		recordID := dom.ID(1)
		partID := dom.ID(3)
		setID100 := dom.ID(100)
		setID101 := dom.ID(101)

		domainRecord := dw.WorkoutRecord{
			ID: &recordID,
			Sets: []dw.WorkoutSet{
				{
					ID:        &setID100,
					SetNumber: 1,
					Weight:    80,
					Reps:      10,
					Exercise: dw.WorkoutExerciseRef{
						ID:     dom.ID(exerciseID),
						Name:   "Bench Press",
						PartID: &partID,
					},
				},
				{
					ID:        &setID101,
					SetNumber: 2,
					Weight:    100,
					Reps:      8,
					Exercise: dw.WorkoutExerciseRef{
						ID:   dom.ID(999),
						Name: "Squat",
					},
				},
			},
		}

		repo.EXPECT().
			GetLastWorkoutRecord(gomock.Any(), userID, exerciseID).
			Return(domainRecord, nil)

		result, err := uc.GetLastWorkoutRecord(ctx, userID, exerciseID)
		require.NoError(t, err)
		require.NotNil(t, result)

		require.NotNil(t, result.ID)
		require.Equal(t, exerciseID, *result.ID)
		require.Equal(t, "Bench Press", result.Name)
		require.NotNil(t, result.WorkoutPartID)
		require.Equal(t, int64(partID), *result.WorkoutPartID)
		require.Len(t, result.Sets, 1)
		require.Equal(t, 1, result.Sets[0].SetNumber)
		require.NotNil(t, result.Sets[0].WeightKg)
		require.Equal(t, 80.0, *result.Sets[0].WeightKg)
		require.NotNil(t, result.Sets[0].Reps)
		require.Equal(t, 10, *result.Sets[0].Reps)
	})
}

func ptrID(v int64) *dom.ID {
	id := dom.ID(v)
	return &id
}
