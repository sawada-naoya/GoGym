package gym

import (
	"context"
	"errors"

	dom "gogym-api/internal/domain/entities"
)

// ErrNotFound is returned when a gym is not found
var ErrNotFound = errors.New("gym not found")

type Repository interface {
	// FindByNormalizedName finds a gym by normalized name and creator
	FindByNormalizedName(ctx context.Context, createdBy string, normalizedName string) (*dom.Gym, error)
	// CreateGym creates a new gym
	CreateGym(ctx context.Context, createdBy string, name string, normalizedName string) (*dom.Gym, error)
}
