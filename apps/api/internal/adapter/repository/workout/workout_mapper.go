package workout

import (
	dw "gogym-api/internal/domain/entities/workout"
	dom "gogym-api/internal/domain/entities"
)

func ToEntity(rec *WorkoutRecord) *dw.WorkoutRecord {
	if rec == nil {
		return nil
	}

	var gymName *string
	if rec.Gym != nil {
		gymName = &rec.Gym.Name
	}

	domainRecord := &dw.WorkoutRecord{
		ID:            ptrInt64ToDomainID(int64(rec.ID)),
		UserID:        dom.ULID(rec.UserID),
		GymID:         int64PtrToDomainIDPtr(rec.GymID),
		GymName:       gymName,
		PerformedDate: rec.PerformedDate,
		StartedAt:     rec.StartedAt,
		EndedAt:       rec.EndedAt,
		Note:          rec.Note,
		Condition:     intPtrToConditionLevel(rec.ConditionLevel),
		DurationMin:   rec.DurationMinutes,
		Sets:          make([]dw.WorkoutSet, 0, len(rec.Sets)),
		CreatedAt:     rec.CreatedAt,
		UpdatedAt:     rec.UpdatedAt,
	}

	for _, s := range rec.Sets {
		domainRecord.Sets = append(domainRecord.Sets, WorkoutSetToDomain(&s))
	}

	return domainRecord
}

func WorkoutSetToDomain(s *WorkoutSet) dw.WorkoutSet {
	exerciseRef := dw.WorkoutExerciseRef{
		ID:     dom.ID(s.WorkoutExerciseID),
		Name:   s.Exercise.Name,
		PartID: intPtrToDomainIDPtr(s.Exercise.WorkoutPartID),
		Owner:  stringPtrToULIDPtr(s.Exercise.UserID),
	}

	return dw.WorkoutSet{
		ID:           ptrInt64ToDomainID(int64(s.ID)),
		Exercise:     exerciseRef,
		SetNumber:    s.SetNumber,
		Weight:       dw.WeightKg(s.WeightKg),
		Reps:         dw.Reps(s.Reps),
		EstimatedMax: s.EstimatedMax,
		Note:         s.Note,
		CreatedAt:    s.CreatedAt,
		UpdatedAt:    s.UpdatedAt,
	}
}

// FromEntity converts dw.WorkoutRecord to WorkoutRecord
func FromEntity(domainRecord *dw.WorkoutRecord) *WorkoutRecord {
	if domainRecord == nil {
		return nil
	}

	rec := &WorkoutRecord{
		UserID:          string(domainRecord.UserID),
		GymID:           domainIDPtrToInt64Ptr(domainRecord.GymID),
		PerformedDate:   domainRecord.PerformedDate,
		StartedAt:       domainRecord.StartedAt,
		EndedAt:         domainRecord.EndedAt,
		Note:            domainRecord.Note,
		ConditionLevel:  conditionLevelToIntPtr(domainRecord.Condition),
		DurationMinutes: domainRecord.DurationMin,
		Sets:            make([]WorkoutSet, 0, len(domainRecord.Sets)),
	}

	if domainRecord.ID != nil {
		rec.ID = int(*domainRecord.ID)
	}

	for _, domainSet := range domainRecord.Sets {
		rec.Sets = append(rec.Sets, WorkoutSetToRecord(&domainSet, 0))
	}

	return rec
}

func WorkoutSetToRecord(domainSet *dw.WorkoutSet, workoutRecordID int) WorkoutSet {
	recSet := WorkoutSet{
		WorkoutRecordID:   workoutRecordID,
		WorkoutExerciseID: int(domainSet.Exercise.ID),
		SetNumber:         domainSet.SetNumber,
		WeightKg:          float64(domainSet.Weight),
		Reps:              int(domainSet.Reps),
		EstimatedMax:      domainSet.EstimatedMax,
		Note:              domainSet.Note,
	}

	if domainSet.ID != nil {
		recSet.ID = int(*domainSet.ID)
	}

	return recSet
}

// WorkoutPartToDomain converts WorkoutPart to dw.WorkoutPart
func WorkoutPartToDomain(rec *WorkoutPart) *dw.WorkoutPart {
	if rec == nil {
		return nil
	}

	// Translationsを変換
	translations := make([]dw.WorkoutPartTranslation, 0, len(rec.Translations))
	for _, trans := range rec.Translations {
		translations = append(translations, dw.WorkoutPartTranslation{
			ID:            dom.ID(trans.ID),
			WorkoutPartID: dom.ID(trans.WorkoutPartID),
			Locale:        trans.Locale,
			Name:          trans.Name,
		})
	}

	// Exercisesを変換
	exercises := make([]dw.WorkoutExerciseRef, 0, len(rec.Exercises))
	for _, ex := range rec.Exercises {
		var partIDPtr *dom.ID
		if ex.WorkoutPartID != nil {
			pid := dom.ID(*ex.WorkoutPartID)
			partIDPtr = &pid
		}

		var ownerPtr *dom.ULID
		if ex.UserID != nil {
			owner := dom.ULID(*ex.UserID)
			ownerPtr = &owner
		}

		exercises = append(exercises, dw.WorkoutExerciseRef{
			ID:     dom.ID(ex.ID),
			Name:   ex.Name,
			PartID: partIDPtr,
			Owner:  ownerPtr,
		})
	}

	return &dw.WorkoutPart{
		ID:           dom.ID(rec.ID),
		Key:          rec.Key,
		Owner:        stringPtrToULIDPtr(rec.UserID),
		Translations: translations,
		Exercises:    exercises,
	}
}

// WorkoutPartsToDomain converts slice of WorkoutPart to slice of dw.WorkoutPart
func WorkoutPartsToDomain(recs []WorkoutPart) []dw.WorkoutPart {
	result := make([]dw.WorkoutPart, len(recs))
	for i, rec := range recs {
		result[i] = *WorkoutPartToDomain(&rec)
	}
	return result
}

// Helper functions

func ptrInt64ToDomainID(i int64) *dom.ID {
	if i == 0 {
		return nil
	}
	id := dom.ID(i)
	return &id
}

func intPtrToDomainIDPtr(i *int) *dom.ID {
	if i == nil {
		return nil
	}
	id := dom.ID(*i)
	return &id
}

func int64PtrToDomainIDPtr(i *int64) *dom.ID {
	if i == nil {
		return nil
	}
	id := dom.ID(*i)
	return &id
}

func domainIDPtrToInt64Ptr(id *dom.ID) *int64 {
	if id == nil {
		return nil
	}
	i := int64(*id)
	return &i
}

func stringPtrToULIDPtr(s *string) *dom.ULID {
	if s == nil {
		return nil
	}
	ulid := dom.ULID(*s)
	return &ulid
}

func intPtrToConditionLevel(i *int) dw.ConditionLevel {
	if i == nil {
		return dw.CondUnknown
	}
	switch *i {
	case 1:
		return dw.Cond1
	case 2:
		return dw.Cond2
	case 3:
		return dw.Cond3
	case 4:
		return dw.Cond4
	case 5:
		return dw.Cond5
	default:
		return dw.CondUnknown
	}
}

func conditionLevelToIntPtr(c dw.ConditionLevel) *int {
	if c == dw.CondUnknown {
		return nil
	}
	i := int(c)
	return &i
}
