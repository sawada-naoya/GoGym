package user

import (
	domain "gogym-api/internal/domain/entities"

	"github.com/oklog/ulid/v2"
)

// ToEntity converts User record to domain entity
func ToEntity(r *User) (*domain.User, error) {
	if r == nil {
		return nil, nil
	}

	id, err := ulid.Parse(r.ID)
	if err != nil {
		return nil, err
	}

	return domain.NewUser(
		id,
		r.Name,
		r.Email,
		r.PasswordHash,
		r.CreatedAt,
	), nil
}

// FromEntity converts domain entity to User record
func FromEntity(u *domain.User) *User {
	if u == nil {
		return nil
	}

	return &User{
		ID:           u.ID.String(),
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Name:         u.Name,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

// ToEntities converts slice of User records to slice of domain entities
func ToEntities(records []*User) ([]*domain.User, error) {
	entities := make([]*domain.User, 0, len(records))
	for _, r := range records {
		entity, err := ToEntity(r)
		if err != nil {
			return nil, err
		}
		entities = append(entities, entity)
	}
	return entities, nil
}
