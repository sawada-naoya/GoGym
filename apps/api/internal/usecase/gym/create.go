package gym

import (
	"context"
	"gogym-api/internal/domain/gym"
)

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
		return nil, gym.NewDomainErrorWithCause(err, "create_failed", "failed to create gym")
	}

	uc.logger.InfoContext(ctx, "gym created successfully", "gym_id", newGym.ID)
	return newGym, nil
}

// getOrCreateTags retrieves existing tags or creates new ones
func (uc *UseCase) getOrCreateTags(ctx context.Context, tagNames []string) ([]gym.Tag, error) {
	// Find existing tags
	existing, err := uc.tagRepo.FindByNames(ctx, tagNames)
	if err != nil {
		return nil, gym.NewDomainErrorWithCause(err, "tag_search_failed", "failed to search tags")
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
			return nil, gym.NewDomainErrorWithCause(err, "tag_create_failed", "failed to create tags")
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
		return nil, gym.NewDomainErrorWithCause(err, "tags_fetch_failed", "failed to fetch tags")
	}

	return tags, nil
}