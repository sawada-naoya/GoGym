package mapper

import (
	"gogym-api/internal/adapter/repository/record"
	dom "gogym-api/internal/domain/workout"
)

// RecordToDomain converts record.WorkoutRecord to domain.WorkoutRecord
func WorkoutRecordToDomain(rec *record.WorkoutRecord) *dom.WorkoutRecord {
	if rec == nil {
		return nil
	}

	domainRecord := &dom.WorkoutRecord{
		ID:            ptrInt64ToDomainID(int64(rec.ID)),
		UserID:        dom.ULID(rec.UserID),
		PerformedDate: rec.PerformedDate,
		StartedAt:     rec.StartedAt,
		EndedAt:       rec.EndedAt,
		Place:         rec.Place,
		Note:          rec.Note,
		Condition:     intPtrToConditionLevel(rec.ConditionLevel),
		DurationMin:   rec.DurationMinutes,
		Sets:          make([]dom.WorkoutSet, 0, len(rec.Sets)),
		CreatedAt:     rec.CreatedAt,
		UpdatedAt:     rec.UpdatedAt,
	}

	for _, s := range rec.Sets {
		domainRecord.Sets = append(domainRecord.Sets, WorkoutSetToDomain(&s))
	}

	return domainRecord
}

// WorkoutSetToDomain converts record.WorkoutSet to domain.WorkoutSet
func WorkoutSetToDomain(s *record.WorkoutSet) dom.WorkoutSet {
	exerciseRef := dom.WorkoutExerciseRef{
		ID:       dom.ID(s.WorkoutExerciseID),
		Name:     s.Exercise.Name,
		PartID:   intPtrToDomainIDPtr(s.Exercise.WorkoutPartID),
		IsPreset: s.Exercise.UserID == nil,
		Owner:    stringPtrToULIDPtr(s.Exercise.UserID),
	}

	return dom.WorkoutSet{
		ID:           ptrInt64ToDomainID(int64(s.ID)),
		Exercise:     exerciseRef,
		SetNumber:    s.SetNumber,
		Weight:       dom.WeightKg(s.WeightKg),
		Reps:         dom.Reps(s.Reps),
		EstimatedMax: s.EstimatedMax,
		Note:         s.Note,
		CreatedAt:    s.CreatedAt,
		UpdatedAt:    s.UpdatedAt,
	}
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

func stringPtrToULIDPtr(s *string) *dom.ULID {
	if s == nil {
		return nil
	}
	ulid := dom.ULID(*s)
	return &ulid
}

func intPtrToConditionLevel(i *int) dom.ConditionLevel {
	if i == nil {
		return dom.CondUnknown
	}
	switch *i {
	case 1:
		return dom.Cond1
	case 2:
		return dom.Cond2
	case 3:
		return dom.Cond3
	case 4:
		return dom.Cond4
	case 5:
		return dom.Cond5
	default:
		return dom.CondUnknown
	}
}
