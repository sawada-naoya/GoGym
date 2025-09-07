// internal/domain/gym/model.go
// 役割: ジムドメインのエンティティモデル
// ジム集約ルート、タグエンティティ、および関連エンティティの定義
package gym

import (
	"gogym-api/internal/domain/common"
	"strings"
)

// Gym represents the gym aggregate root
type Gym struct {
	common.BaseEntity
	Name         string               `json:"name" gorm:"size:255;not null" validate:"required,max=255"`
	Description  *string              `json:"description" gorm:"type:text"`
	Location     common.Location      `json:"location" gorm:"embedded;embeddedPrefix:location_"`
	Address      string               `json:"address" gorm:"size:500;not null" validate:"required,max=500"`
	City         *string              `json:"city" gorm:"size:100"`
	Prefecture   *string              `json:"prefecture" gorm:"size:100"`
	PostalCode   *string              `json:"postal_code" gorm:"size:10"`
	Tags         []Tag                `json:"tags,omitempty" gorm:"many2many:gym_tags;"`
	AverageRating *float32            `json:"average_rating" gorm:"-"`
	ReviewCount   int                 `json:"review_count" gorm:"-"`
}

// TableName returns the table name for GORM
func (Gym) TableName() string {
	return "gyms"
}

// NewGym creates a new gym with validation
func NewGym(name, address string, location common.Location) (*Gym, error) {
	gym := &Gym{
		Name:     strings.TrimSpace(name),
		Address:  strings.TrimSpace(address),
		Location: location,
	}

	if err := gym.Validate(); err != nil {
		return nil, err
	}

	return gym, nil
}

// Validate validates gym data
func (g *Gym) Validate() error {
	if g.Name == "" {
		return common.NewDomainError(common.ErrInvalidInput, "invalid_name", "gym name is required")
	}

	if len(g.Name) > 255 {
		return common.NewDomainError(common.ErrInvalidInput, "invalid_name", "gym name too long")
	}

	if g.Address == "" {
		return common.NewDomainError(common.ErrInvalidInput, "invalid_address", "gym address is required")
	}

	if len(g.Address) > 500 {
		return common.NewDomainError(common.ErrInvalidInput, "invalid_address", "gym address too long")
	}

	if !g.Location.IsValid() {
		return common.NewDomainError(common.ErrInvalidLocation, "invalid_location", "invalid location coordinates")
	}

	return nil
}

// SetDescription sets the gym description
func (g *Gym) SetDescription(description string) {
	trimmed := strings.TrimSpace(description)
	if trimmed == "" {
		g.Description = nil
	} else {
		g.Description = &trimmed
	}
}

// SetCity sets the gym city
func (g *Gym) SetCity(city string) {
	trimmed := strings.TrimSpace(city)
	if trimmed == "" {
		g.City = nil
	} else {
		g.City = &trimmed
	}
}

// SetPrefecture sets the gym prefecture
func (g *Gym) SetPrefecture(prefecture string) {
	trimmed := strings.TrimSpace(prefecture)
	if trimmed == "" {
		g.Prefecture = nil
	} else {
		g.Prefecture = &trimmed
	}
}

// SetPostalCode sets the gym postal code
func (g *Gym) SetPostalCode(postalCode string) {
	trimmed := strings.TrimSpace(postalCode)
	if trimmed == "" {
		g.PostalCode = nil
	} else {
		g.PostalCode = &trimmed
	}
}

// Tag represents a gym tag entity
type Tag struct {
	common.BaseEntity
	Name string `json:"name" gorm:"uniqueIndex;size:50;not null" validate:"required,max=50"`
	Gyms []Gym  `json:"gyms,omitempty" gorm:"many2many:gym_tags;"`
}

// TableName returns the table name for GORM
func (Tag) TableName() string {
	return "tags"
}

// NewTag creates a new tag with validation
func NewTag(name string) (*Tag, error) {
	tag := &Tag{
		Name: strings.TrimSpace(name),
	}

	if err := tag.Validate(); err != nil {
		return nil, err
	}

	return tag, nil
}

// Validate validates tag data
func (t *Tag) Validate() error {
	if t.Name == "" {
		return common.NewDomainError(common.ErrInvalidInput, "invalid_name", "tag name is required")
	}

	if len(t.Name) > 50 {
		return common.NewDomainError(common.ErrInvalidInput, "invalid_name", "tag name too long")
	}

	return nil
}

// GymTag represents the many-to-many relationship between gyms and tags
type GymTag struct {
	GymID common.ID `json:"gym_id" gorm:"primaryKey"`
	TagID common.ID `json:"tag_id" gorm:"primaryKey"`
	Gym   *Gym      `json:"gym,omitempty" gorm:"foreignKey:GymID"`
	Tag   *Tag      `json:"tag,omitempty" gorm:"foreignKey:TagID"`
}

// TableName returns the table name for GORM
func (GymTag) TableName() string {
	return "gym_tags"
}

// Favorite represents a user's favorite gym
type Favorite struct {
	common.BaseEntity
	UserID common.ID `json:"user_id" gorm:"not null;index;uniqueIndex:unique_user_gym_favorite,priority:1"`
	GymID  common.ID `json:"gym_id" gorm:"not null;index;uniqueIndex:unique_user_gym_favorite,priority:2"`
}

// TableName returns the table name for GORM
func (Favorite) TableName() string {
	return "favorites"
}

// NewFavorite creates a new favorite with validation
func NewFavorite(userID, gymID common.ID) (*Favorite, error) {
	favorite := &Favorite{
		UserID: userID,
		GymID:  gymID,
	}

	if err := favorite.Validate(); err != nil {
		return nil, err
	}

	return favorite, nil
}

// Validate validates favorite data
func (f *Favorite) Validate() error {
	if f.UserID == 0 {
		return common.NewDomainError(common.ErrInvalidInput, "invalid_user_id", "user ID is required")
	}

	if f.GymID == 0 {
		return common.NewDomainError(common.ErrInvalidInput, "invalid_gym_id", "gym ID is required")
	}

	return nil
}