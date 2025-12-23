package domain

import (
	"errors"
	"time"
)

type ULID string // users.id は ULID想定
type ID int64

type ConditionLevel uint8 // 1..5 を表す
const (
	CondUnknown ConditionLevel = 0
	Cond1       ConditionLevel = 1
	Cond2       ConditionLevel = 2
	Cond3       ConditionLevel = 3
	Cond4       ConditionLevel = 4
	Cond5       ConditionLevel = 5
)

type WeightKg float64
type Reps int

func (w WeightKg) Valid() bool { return w >= 0 }
func (r Reps) Valid() bool     { return r >= 0 }

type WorkoutPart struct {
	ID        ID
	Name      string
	Owner     *ULID                // nil ならプリセット、値があればユーザー作成
	Exercises []WorkoutExerciseRef // この部位に紐づく種目
}

type WorkoutExerciseRef struct {
	ID     ID
	Name   string
	PartID *ID
	Owner  *ULID // nil ならプリセット、値があればユーザー作成
}

type WorkoutSet struct {
	ID           *ID
	Exercise     WorkoutExerciseRef
	SetNumber    int
	Weight       WeightKg
	Reps         Reps
	EstimatedMax *float64
	Note         *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type WorkoutRecord struct {
	ID            *ID
	UserID        ULID
	GymID         *ID
	GymName       *string    // ジム名（表示用）
	PerformedDate time.Time  // DATE を day-start に固定、時刻は別扱い
	StartedAt     *time.Time // 実日時（PerformedDateに紐づけて作る）
	EndedAt       *time.Time
	Note          *string
	Condition     ConditionLevel
	DurationMin   *int // 派生値：Started/Endedから再計算
	Sets          []WorkoutSet
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewWorkoutRecord(user ULID, performedDate time.Time) (*WorkoutRecord, error) {
	if user == "" {
		return nil, errors.New("user required")
	}
	// performedDate は時刻00:00に正規化しておく設計が楽
	pd := time.Date(performedDate.Year(), performedDate.Month(), performedDate.Day(), 0, 0, 0, 0, performedDate.Location())
	return &WorkoutRecord{
		UserID:        user,
		PerformedDate: pd,
		Condition:     CondUnknown,
		Sets:          []WorkoutSet{},
	}, nil
}

func (r *WorkoutRecord) SetTimes(start, end *time.Time) error {
	if start != nil && end != nil && start.After(*end) {
		return errors.New("start must be <= end")
	}
	r.StartedAt, r.EndedAt = start, end
	r.recalcDuration()
	return nil
}

func (r *WorkoutRecord) recalcDuration() {
	if r.StartedAt != nil && r.EndedAt != nil {
		min := int(r.EndedAt.Sub(*r.StartedAt).Minutes())
		if min < 0 {
			min = 0
		}
		r.DurationMin = &min
	} else {
		r.DurationMin = nil
	}
}

func (r *WorkoutRecord) AddSet(s WorkoutSet) error {
	if s.SetNumber <= 0 {
		return errors.New("setNumber must be >= 1")
	}
	if !s.Weight.Valid() || !s.Reps.Valid() {
		return errors.New("invalid weight or reps")
	}
	// (exerciseID, setNumber) の一意性
	for _, cur := range r.Sets {
		if cur.Exercise.ID == s.Exercise.ID && cur.SetNumber == s.SetNumber {
			return errors.New("duplicate setNumber for the exercise")
		}
	}
	r.Sets = append(r.Sets, s)
	return nil
}

func (r *WorkoutRecord) ReorderSets(exerciseID ID) {
	// 同一 exercise の setNumber を 1..N に詰め直すユーティリティ（必要なら）
}
