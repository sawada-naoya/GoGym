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

// UseCase represents gym use cases
type UseCase struct {
	gymRepo gym.Repository
	tagRepo gym.TagRepository
	logger  *slog.Logger
}

// NewUseCase creates a new gym use case
func NewUseCase(gymRepo gym.Repository, tagRepo gym.TagRepository, logger *slog.Logger) *UseCase {
	return &UseCase{
		gymRepo: gymRepo,
		tagRepo: tagRepo,
		logger:  logger,
	}
}

// SearchGymRequest represents search gym input
type SearchGymRequest struct {
	Query      string
	Location   *common.Location
	RadiusM    *int
	Cursor     string
	Limit      int
}

// SearchGymsResponse represents search gym output
type SearchGymsResponse struct {
	Gyms       []gym.Gym
	NextCursor *string
	HasMore    bool
}

// SearchGyms searches gyms based on criteria
func (uc *UseCase) SearchGyms(ctx context.Context, req SearchGymRequest) (*SearchGymsResponse, error) {
	uc.logger.InfoContext(ctx, "searching gyms",
		"query", req.Query,
		"location", req.Location,
		"radius_m", req.RadiusM,
		"limit", req.Limit,
	)

	// Set default limit
	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 20
	}

	// Set default radius
	if req.RadiusM != nil && (*req.RadiusM < 100 || *req.RadiusM > 50000) {
		defaultRadius := 5000
		req.RadiusM = &defaultRadius
	}

	searchQuery := common.SearchQuery{
		Query:    req.Query,
		Location: req.Location,
		RadiusM:  req.RadiusM,
		Pagination: common.Pagination{
			Cursor: req.Cursor,
			Limit:  req.Limit,
		},
	}

	result, err := uc.gymRepo.Search(ctx, searchQuery)
	if err != nil {
		uc.logger.ErrorContext(ctx, "failed to search gyms", "error", err)
		return nil, common.NewDomainError(err, "search_failed", "failed to search gyms")
	}

	return &SearchGymsResponse{
		Gyms:       result.Items,
		NextCursor: result.NextCursor,
		HasMore:    result.HasMore,
	}, nil
}

// GetGym retrieves a gym by ID
func (uc *UseCase) GetGym(ctx context.Context, id common.ID) (*gym.Gym, error) {
	uc.logger.InfoContext(ctx, "getting gym", "gym_id", id)

	if id == 0 {
		return nil, common.NewDomainError(common.ErrInvalidInput, "invalid_gym_id", "gym ID is required")
	}

	foundGym, err := uc.gymRepo.FindByID(ctx, id)
	if err != nil {
		uc.logger.ErrorContext(ctx, "failed to get gym", "gym_id", id, "error", err)
		return nil, common.NewDomainError(err, "gym_not_found", "gym not found")
	}

	return foundGym, nil
}

// CreateGymRequest represents create gym input
type CreateGymRequest struct {
	Name        string
	Description *string
	Location    common.Location
	Address     string
	City        *string
	Prefecture  *string
	PostalCode  *string
	TagNames    []string
}

// CreateGym creates a new gym
func (uc *UseCase) CreateGym(ctx context.Context, req CreateGymRequest) (*gym.Gym, error) {
	uc.logger.InfoContext(ctx, "creating gym", "name", req.Name, "address", req.Address)

	// Create gym entity
	newGym, err := gym.NewGym(req.Name, req.Address, req.Location)
	if err != nil {
		uc.logger.ErrorContext(ctx, "invalid gym data", "error", err)
		return nil, err
	}

	// Set optional fields
	if req.Description != nil {
		newGym.SetDescription(*req.Description)
	}
	if req.City != nil {
		newGym.SetCity(*req.City)
	}
	if req.Prefecture != nil {
		newGym.SetPrefecture(*req.Prefecture)
	}
	if req.PostalCode != nil {
		newGym.SetPostalCode(*req.PostalCode)
	}

	// Handle tags
	if len(req.TagNames) > 0 {
		tags, err := uc.getOrCreateTags(ctx, req.TagNames)
		if err != nil {
			uc.logger.ErrorContext(ctx, "failed to handle tags", "error", err)
			return nil, err
		}
		newGym.Tags = tags
	}

	// Save gym
	if err := uc.gymRepo.Create(ctx, newGym); err != nil {
		uc.logger.ErrorContext(ctx, "failed to create gym", "error", err)
		return nil, common.NewDomainError(err, "create_failed", "failed to create gym")
	}

	uc.logger.InfoContext(ctx, "gym created successfully", "gym_id", newGym.ID)
	return newGym, nil
}

// getOrCreateTags retrieves existing tags or creates new ones
func (uc *UseCase) getOrCreateTags(ctx context.Context, tagNames []string) ([]gym.Tag, error) {
	// Find existing tags
	existing, err := uc.tagRepo.FindByNames(ctx, tagNames)
	if err != nil {
		return nil, common.NewDomainError(err, "tag_search_failed", "failed to search tags")
	}

	existingMap := make(map[string]gym.Tag)
	for _, tag := range existing {
		existingMap[tag.Name] = tag
	}

	// Create missing tags
	var newTags []gym.Tag
	var allTags []gym.Tag

	for _, name := range tagNames {
		if existingTag, exists := existingMap[name]; exists {
			allTags = append(allTags, existingTag)
		} else {
			newTag, err := gym.NewTag(name)
			if err != nil {
				return nil, err
			}
			newTags = append(newTags, *newTag)
			allTags = append(allTags, *newTag)
		}
	}

	// Save new tags if any
	if len(newTags) > 0 {
		if err := uc.tagRepo.CreateMany(ctx, newTags); err != nil {
			return nil, common.NewDomainError(err, "tag_create_failed", "failed to create tags")
		}
	}

	return allTags, nil
}

// GetAllTags retrieves all available tags
func (uc *UseCase) GetAllTags(ctx context.Context) ([]gym.Tag, error) {
	uc.logger.InfoContext(ctx, "getting all tags")

	tags, err := uc.tagRepo.FindAll(ctx)
	if err != nil {
		uc.logger.ErrorContext(ctx, "failed to get tags", "error", err)
		return nil, common.NewDomainError(err, "tags_fetch_failed", "failed to fetch tags")
	}

	return tags, nil
}