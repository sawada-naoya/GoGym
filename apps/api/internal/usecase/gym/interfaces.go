// internal/usecase/gym/interfaces.go
// 役割: ジムユースケースとRepository interface（Application Layer）
// ビジネスロジックを実装し、ドメインエンティティとインターフェースにのみ依存する
package gym

import (
	"context"
	"gogym-api/internal/domain/gym"
	"log/slog"
)

type UseCase struct {
	gymRepo      Repository
	tagRepo      TagRepository
	favoriteRepo FavoriteRepository
	logger       *slog.Logger
}


func NewUseCase(gymRepo Repository, tagRepo TagRepository, logger *slog.Logger) *UseCase {
	return &UseCase{
		gymRepo:      gymRepo,
		tagRepo:      tagRepo,
		favoriteRepo: nil, // TODO: Add favorite repository when implemented
		logger:       logger,
	}
}

type Repository interface {
	FindByID(ctx context.Context, id gym.ID) (*gym.Gym, error)
	FindDetailByID(ctx context.Context, id gym.ID) (*gym.Gym, error) // With relations (reviews, amenities, etc.)
	Search(ctx context.Context, query gym.SearchQuery) (*gym.PaginatedResult[gym.Gym], error)
	Create(ctx context.Context, gym *gym.Gym) error
	Update(ctx context.Context, gym *gym.Gym) error
	Delete(ctx context.Context, id gym.ID) error
}

type TagRepository interface {
	FindAll(ctx context.Context) ([]gym.Tag, error)
	FindByIDs(ctx context.Context, ids []gym.ID) ([]gym.Tag, error)
	FindByNames(ctx context.Context, names []string) ([]gym.Tag, error)
	Create(ctx context.Context, tag *gym.Tag) error
	CreateMany(ctx context.Context, tags []gym.Tag) error
}

type FavoriteRepository interface {
	AddFavorite(ctx context.Context, userID gym.ID, gymID gym.ID) error
	RemoveFavorite(ctx context.Context, userID gym.ID, gymID gym.ID) error
	GetFavoriteGyms(ctx context.Context, userID gym.ID, pagination gym.Pagination) (*gym.PaginatedResult[gym.Gym], error)
	IsFavorite(ctx context.Context, userID gym.ID, gymID gym.ID) (bool, error)
}
