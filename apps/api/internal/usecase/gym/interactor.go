package gym

import (
	"log/slog"

	"gogym-api/internal/usecase/tag"
)

type interactor struct {
	repo    Repository
	tagRepo tag.Repository
	logger  *slog.Logger
}

func NewUseCase(repo Repository, tagRepo tag.Repository, logger *slog.Logger) UseCase {
	return &interactor{
		repo:    repo,
		tagRepo: tagRepo,
		logger:  logger,
	}
}