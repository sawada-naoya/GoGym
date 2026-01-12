package dto

import (
	"fmt"
	"gogym-api/internal/util"
	"sort"
	"time"

	dom "gogym-api/internal/domain/entities"
)

type WorkoutPartListItemDTO struct {
	ID           int64                        `json:"id"`
	Key          string                       `json:"key"`
	Translations []WorkoutPartTranslationDTO  `json:"translations"`
	Exercises    []WorkoutExerciseListItemDTO `json:"exercises"`
}

type WorkoutPartTranslationDTO struct {
	Locale string `json:"locale"`
	Name   string `json:"name"`
}

type WorkoutRecordDTO struct {
	ID             *int64  `json:"id,omitempty"`
	PerformedDate  string  `json:"performed_date"`
	StartedAt      *string `json:"started_at,omitempty"`
	EndedAt        *string `json:"ended_at,omitempty"`
	GymID          *int64  `json:"gym_id,omitempty"`
	GymName        *string `json:"gym_name,omitempty"`
	Note           *string `json:"note,omitempty"`
	ConditionLevel *int    `json:"condition_level,omitempty"`

	Parts []WorkoutPartGroupDTO `json:"parts"`
}

type WorkoutPartGroupDTO struct {
	ID           int64                       `json:"id"`
	Key          string                      `json:"key"`
	Translations []WorkoutPartTranslationDTO `json:"translations"`
	Exercises    []ExerciseDTO               `json:"exercises"`
}

type WorkoutPartDTO struct {
	ID     *int64  `json:"id,omitempty"`
	Name   *string `json:"name,omitempty"`
	Source *string `json:"source,omitempty"` // "preset" | "custom" | null
}

type ExerciseDTO struct {
	ID            *int64   `json:"id,omitempty"`
	Name          string   `json:"name"`
	WorkoutPartID *int64   `json:"workout_part_id,omitempty"`
	Sets          []SetDTO `json:"sets"`
}

type SetDTO struct {
	ID        *int64   `json:"id,omitempty"`
	SetNumber int      `json:"set_number"`
	WeightKg  *float64 `json:"weight_kg,omitempty"` // 空文字→null→nil→層内で検証
	Reps      *int     `json:"reps,omitempty"`
	Note      *string  `json:"note,omitempty"`
}

type CreateWorkoutExerciseRequest struct {
	Exercises []CreateWorkoutExerciseItem `json:"exercises"`
}

type CreateWorkoutExerciseItem struct {
	ID            *int64 `json:"id,omitempty"` // nil = insert, value = update
	Name          string `json:"name"`
	WorkoutPartID int64  `json:"workout_part_id"`
}

type WorkoutExerciseListItemDTO struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	WorkoutPartID *int64 `json:"workout_part_id,omitempty"`
}

func WorkoutDomainToDTO(record *dom.WorkoutRecord) *WorkoutRecordDTO {
	if record == nil {
		return nil
	}

	var id *int64
	if record.ID != nil {
		i := int64(*record.ID)
		id = &i
	}

	var startedAt *string
	if record.StartedAt != nil && !record.StartedAt.IsZero() {
		s := util.FormatJSTTime(*record.StartedAt)
		startedAt = &s
	}

	var endedAt *string
	if record.EndedAt != nil && !record.EndedAt.IsZero() {
		s := util.FormatJSTTime(*record.EndedAt)
		endedAt = &s
	}

	var gymID *int64
	if record.GymID != nil {
		gid := int64(*record.GymID)
		gymID = &gid
	}

	var conditionLevel *int
	if record.Condition != dom.CondUnknown {
		cl := int(record.Condition)
		conditionLevel = &cl
	}

	out := &WorkoutRecordDTO{
		ID:             id,
		PerformedDate:  util.FormatJSTDate(record.PerformedDate),
		StartedAt:      startedAt,
		EndedAt:        endedAt,
		GymID:          gymID,
		GymName:        record.GymName,
		Note:           record.Note,
		ConditionLevel: conditionLevel,
		Parts:          []WorkoutPartGroupDTO{},
	}

	// partId -> partGroup (部位情報は種目から取得)
	partMap := map[int64]*WorkoutPartGroupDTO{}
	// partId -> exerciseId -> exerciseDTO
	exMap := map[int64]map[int64]*ExerciseDTO{}

	for _, set := range record.Sets {
		ex := set.Exercise

		// 部位IDがない場合はスキップ
		if ex.PartID == nil {
			continue
		}

		pid := int64(*ex.PartID)

		// 部位グループが存在しない場合は作成（簡易版、翻訳データなし）
		if _, ok := partMap[pid]; !ok {
			p := &WorkoutPartGroupDTO{
				ID:           pid,
				Key:          "", // 部位のKeyはセット情報からは取得できない
				Translations: []WorkoutPartTranslationDTO{},
				Exercises:    []ExerciseDTO{},
			}
			partMap[pid] = p
			exMap[pid] = map[int64]*ExerciseDTO{}
		}

		eid := int64(ex.ID)
		if _, ok := exMap[pid][eid]; !ok {
			exID := eid
			partID := pid
			exMap[pid][eid] = &ExerciseDTO{
				ID:            &exID,
				Name:          ex.Name,
				WorkoutPartID: &partID,
				Sets:          []SetDTO{},
			}
		}

		var setID *int64
		if set.ID != nil {
			sid := int64(*set.ID)
			setID = &sid
		}

		weight := float64(set.Weight)
		reps := int(set.Reps)

		exMap[pid][eid].Sets = append(exMap[pid][eid].Sets, SetDTO{
			ID:        setID,
			SetNumber: set.SetNumber,
			WeightKg:  &weight,
			Reps:      &reps,
			Note:      set.Note,
		})
	}

	partIDs := make([]int64, 0, len(partMap))
	for pid := range partMap {
		partIDs = append(partIDs, pid)
	}
	sort.Slice(partIDs, func(i, j int) bool { return partIDs[i] < partIDs[j] })

	for _, pid := range partIDs {
		p := partMap[pid]

		// exercises を slice 化
		exIDs := make([]int64, 0, len(exMap[pid]))
		for eid := range exMap[pid] {
			exIDs = append(exIDs, eid)
		}
		sort.Slice(exIDs, func(i, j int) bool { return exIDs[i] < exIDs[j] })

		for _, eid := range exIDs {
			e := exMap[pid][eid]
			// set_number で並べる
			sort.Slice(e.Sets, func(i, j int) bool { return e.Sets[i].SetNumber < e.Sets[j].SetNumber })
			p.Exercises = append(p.Exercises, *e)
		}

		out.Parts = append(out.Parts, *p)
	}

	return out
}

func domainIDToInt64Ptr(id *dom.ID) *int64 {
	if id == nil {
		return nil
	}
	i := int64(*id)
	return &i
}

func stringPtr(s string) *string {
	return &s
}

func stringPtrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func conditionLevelToIntPtr(c dom.ConditionLevel) *int {
	if c == dom.CondUnknown {
		return nil
	}
	i := int(c)
	return &i
}

// WorkoutRecordDTOToDomain converts WorkoutRecordDTO to domain.WorkoutRecord
func WorkoutRecordDTOToDomain(dto *WorkoutRecordDTO) (*dom.WorkoutRecord, error) {
	if dto == nil {
		return nil, fmt.Errorf("dto is nil")
	}

	// Parse performedDate as UTC midnight to avoid timezone shift
	// "2025-11-25" should be stored as "2025-11-25 00:00:00 UTC" regardless of JST
	performedDate, err := time.Parse("2006-01-02", dto.PerformedDate)
	if err != nil {
		return nil, fmt.Errorf("invalid performedDate format: %w", err)
	}

	// Ensure it's UTC midnight
	performedDateUTC := time.Date(performedDate.Year(), performedDate.Month(), performedDate.Day(), 0, 0, 0, 0, time.UTC)

	// Create WorkoutRecord with placeholder userID (will be set by handler/usecase)
	record := &dom.WorkoutRecord{
		UserID:        dom.ULID(""), // will be set by handler
		PerformedDate: performedDateUTC,
		Condition:     dom.CondUnknown,
		Sets:          []dom.WorkoutSet{},
	}

	// Set ID if exists
	if dto.ID != nil {
		id := dom.ID(*dto.ID)
		record.ID = &id
	}

	// Parse and set times
	var startedAt, endedAt *time.Time
	if dto.StartedAt != nil {
		t, err := parseTimeWithDate(performedDate, *dto.StartedAt)
		if err != nil {
			return nil, fmt.Errorf("invalid startedAt format: %w", err)
		}
		startedAt = &t
	}
	if dto.EndedAt != nil {
		t, err := parseTimeWithDate(performedDate, *dto.EndedAt)
		if err != nil {
			return nil, fmt.Errorf("invalid endedAt format: %w", err)
		}
		endedAt = &t
	}
	if err := record.SetTimes(startedAt, endedAt); err != nil {
		return nil, fmt.Errorf("invalid times: %w", err)
	}

	record.Note = dto.Note
	if dto.ConditionLevel != nil {
		record.Condition = dom.ConditionLevel(*dto.ConditionLevel)
	}
	if dto.GymID != nil {
		gymID := dom.ID(*dto.GymID)
		record.GymID = &gymID
	}

	for _, exercise := range dto.Parts[0].Exercises {
		exerciseRef := dom.WorkoutExerciseRef{
			Name: exercise.Name,
		}
		if exercise.ID != nil {
			exerciseRef.ID = dom.ID(*exercise.ID)
		}
		if exercise.WorkoutPartID != nil {
			partID := dom.ID(*exercise.WorkoutPartID)
			exerciseRef.PartID = &partID
		}
		// Owner は exercise.ID の有無で判定される想定
		// 既存の exercise.ID がない場合は新規作成されるため、Owner は usecase 層で設定される

		// Convert sets
		for _, setDTO := range exercise.Sets {
			// Skip empty sets (both weight and reps are nil)
			if setDTO.WeightKg == nil || setDTO.Reps == nil {
				continue
			}

			workoutSet := dom.WorkoutSet{
				Exercise:  exerciseRef,
				SetNumber: setDTO.SetNumber,
				Weight:    dom.WeightKg(*setDTO.WeightKg),
				Reps:      dom.Reps(*setDTO.Reps),
				Note:      setDTO.Note,
			}

			if setDTO.ID != nil {
				id := dom.ID(*setDTO.ID)
				workoutSet.ID = &id
			}

			// Add set to record
			if err := record.AddSet(workoutSet); err != nil {
				return nil, fmt.Errorf("failed to add set: %w", err)
			}
		}
	}

	return record, nil
}

// parseTimeWithDate combines a date and HH:mm time string
// フロントから受け取った HH:mm をそのまま UTC として保存（タイムゾーン変換しない）
func parseTimeWithDate(date time.Time, hhmmStr string) (time.Time, error) {
	t, err := time.Parse("15:04", hhmmStr)
	if err != nil {
		return time.Time{}, err
	}

	// タイムゾーン変換せず、そのままUTCとして保存
	utcTime := time.Date(date.Year(), date.Month(), date.Day(), t.Hour(), t.Minute(), 0, 0, time.UTC)

	return utcTime, nil
}

// WorkoutPartToDTO converts domain.WorkoutPart to WorkoutPartListItemDTO
func WorkoutPartToDTO(part *dom.WorkoutPart) *WorkoutPartListItemDTO {
	if part == nil {
		return nil
	}

	// Translationsを変換
	translations := make([]WorkoutPartTranslationDTO, 0, len(part.Translations))
	for _, trans := range part.Translations {
		translations = append(translations, WorkoutPartTranslationDTO{
			Locale: trans.Locale,
			Name:   trans.Name,
		})
	}

	// Exercisesを変換
	exercises := make([]WorkoutExerciseListItemDTO, 0, len(part.Exercises))
	for _, ex := range part.Exercises {
		var partIDPtr *int64
		if ex.PartID != nil {
			pid := int64(*ex.PartID)
			partIDPtr = &pid
		}

		exercises = append(exercises, WorkoutExerciseListItemDTO{
			ID:            int64(ex.ID),
			Name:          ex.Name,
			WorkoutPartID: partIDPtr,
		})
	}

	return &WorkoutPartListItemDTO{
		ID:           int64(part.ID),
		Key:          part.Key,
		Translations: translations,
		Exercises:    exercises,
	}
}

// WorkoutPartsToDTO converts slice of domain.WorkoutPart to slice of WorkoutPartListItemDTO
func WorkoutPartsToDTO(parts []dom.WorkoutPart) []WorkoutPartListItemDTO {
	result := make([]WorkoutPartListItemDTO, len(parts))
	for i, part := range parts {
		result[i] = *WorkoutPartToDTO(&part)
	}
	return result
}
