package tag

import (
	"context"
	"gogym-api/internal/domain/gym"
)

// handler → usecase (if needed in the future)
type UseCase interface {
	GetTags(ctx context.Context) ([]gym.Tag, error)
}