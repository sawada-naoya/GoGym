package domain

import (
	"time"
)

type WeightKg float64

type Reps int

func (w WeightKg) Valid() bool { return w >= 0 }

func (r Reps) Valid() bool { return r >= 0 }

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
