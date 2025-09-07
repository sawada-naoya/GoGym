// internal/usecase/gym/interfaces.go
// 役割: ジムユースケースとRepository interface（Application Layer）
// ビジネスロジックを実装し、ドメインエンティティとインターフェースにのみ依存する
package gym

import (
	"context"
	"gogym-api/internal/domain/common"
	"gogym-api/internal/domain/gym"
	"log/slog"
)

// Repository interface for gym data access
type Repository interface {
	FindByID(ctx context.Context, id common.ID) (*gym.Gym, error)
	Search(ctx context.Context, query common.SearchQuery) (*common.PaginatedResult[gym.Gym], error)
	Create(ctx context.Context, gym *gym.Gym) error
	Update(ctx context.Context, gym *gym.Gym) error
	Delete(ctx context.Context, id common.ID) error
}

// TagRepository interface for tag data access
type TagRepository interface {
	FindAll(ctx context.Context) ([]gym.Tag, error)
	FindByIDs(ctx context.Context, ids []common.ID) ([]gym.Tag, error)
	FindByNames(ctx context.Context, names []string) ([]gym.Tag, error)
	Create(ctx context.Context, tag *gym.Tag) error
	CreateMany(ctx context.Context, tags []gym.Tag) error
}

// FavoriteRepository interface for favorite data access
type FavoriteRepository interface {
	AddFavorite(ctx context.Context, userID common.ID, gymID common.ID) error
	RemoveFavorite(ctx context.Context, userID common.ID, gymID common.ID) error
	GetFavoriteGyms(ctx context.Context, userID common.ID, pagination common.Pagination) (*common.PaginatedResult[gym.Gym], error)
	IsFavorite(ctx context.Context, userID common.ID, gymID common.ID) (bool, error)
}

// UseCase represents gym use cases
type UseCase struct {
	gymRepo      Repository
	tagRepo      TagRepository
	favoriteRepo FavoriteRepository
	logger       *slog.Logger
}

// NewUseCase creates a new gym use case
func NewUseCase(gymRepo Repository, tagRepo TagRepository, logger *slog.Logger) *UseCase {
	return &UseCase{
		gymRepo:      gymRepo,
		tagRepo:      tagRepo,
		favoriteRepo: nil, // TODO: Add favorite repository when implemented
		logger:       logger,
	}
}