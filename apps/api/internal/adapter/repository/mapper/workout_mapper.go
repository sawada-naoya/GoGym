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
		ID:     dom.ID(s.WorkoutExerciseID),
		Name:   s.Exercise.Name,
		PartID: intPtrToDomainIDPtr(s.Exercise.WorkoutPartID),
		Owner:  stringPtrToULIDPtr(s.Exercise.UserID),
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

// WorkoutRecordToRecord converts domain.WorkoutRecord to record.WorkoutRecord
func WorkoutRecordToRecord(domainRecord *dom.WorkoutRecord) *record.WorkoutRecord {
	if domainRecord == nil {
		return nil
	}

	rec := &record.WorkoutRecord{
		UserID:          string(domainRecord.UserID),
		PerformedDate:   domainRecord.PerformedDate,
		StartedAt:       domainRecord.StartedAt,
		EndedAt:         domainRecord.EndedAt,
		Place:           domainRecord.Place,
		Note:            domainRecord.Note,
		ConditionLevel:  conditionLevelToIntPtr(domainRecord.Condition),
		DurationMinutes: domainRecord.DurationMin,
		Sets:            make([]record.WorkoutSet, 0, len(domainRecord.Sets)),
	}

	if domainRecord.ID != nil {
		rec.ID = int(*domainRecord.ID)
	}

	for _, domainSet := range domainRecord.Sets {
		rec.Sets = append(rec.Sets, WorkoutSetToRecord(&domainSet, 0))
	}

	return rec
}

// WorkoutSetToRecord converts domain.WorkoutSet to record.WorkoutSet
func WorkoutSetToRecord(domainSet *dom.WorkoutSet, workoutRecordID int) record.WorkoutSet {
	recSet := record.WorkoutSet{
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

func conditionLevelToIntPtr(c dom.ConditionLevel) *int {
	if c == dom.CondUnknown {
		return nil
	}
	i := int(c)
	return &i
}

// WorkoutPartToDomain converts record.WorkoutPart to domain.WorkoutPart
func WorkoutPartToDomain(rec *record.WorkoutPart) *dom.WorkoutPart {
	if rec == nil {
		return nil
	}
	return &dom.WorkoutPart{
		ID:    dom.ID(rec.ID),
		Name:  rec.Name,
		Owner: stringPtrToULIDPtr(rec.UserID),
	}
}

// WorkoutPartsToDomain converts slice of record.WorkoutPart to slice of domain.WorkoutPart
func WorkoutPartsToDomain(recs []record.WorkoutPart) []dom.WorkoutPart {
	result := make([]dom.WorkoutPart, len(recs))
	for i, rec := range recs {
		result[i] = *WorkoutPartToDomain(&rec)
	}
	return result
}
