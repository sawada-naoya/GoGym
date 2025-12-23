package domain

import (
	"errors"
	"regexp"
	"strings"
)

var multiSpaceRegex = regexp.MustCompile(`\s+`)

type Gym struct {
	ID              int
	Name            string `validate:"required,max=255"`
	NormalizedName  string `validate:"required,max=255"`
	Latitude        float64
	Longitude       float64
	SourceURL       string
	PrimaryPhotoURL string
	PlaceID         int
}

// NormalizeName normalizes gym name for deduplication
// - Convert to lowercase
// - Trim whitespace
// - Replace consecutive whitespace with single space
func NormalizeName(name string) string {
	normalized := strings.TrimSpace(name)
	normalized = multiSpaceRegex.ReplaceAllString(normalized, " ")
	return strings.ToLower(normalized)
}

func NewGym(name, address string, latitude, longitude float64) (*Gym, error) {
	trimmedName := strings.TrimSpace(name)
	gym := &Gym{
		Name:           trimmedName,
		NormalizedName: NormalizeName(trimmedName),
		Latitude:       latitude,
		Longitude:      longitude,
	}

	if err := gym.Validate(); err != nil {
		return nil, err
	}

	return gym, nil
}

func (g *Gym) Validate() error {
	if g.Name == "" {
		return errors.New("invalid name")
	}

	if len(g.Name) > 255 {
		return errors.New("invalid name")
	}

	return nil
}
