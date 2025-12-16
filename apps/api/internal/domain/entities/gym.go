package domain

import (
	"errors"
	"strings"
)

type Gym struct {
	ID              int
	Name            string `validate:"required,max=255"`
	Latitude        float64
	Longitude       float64
	SourceURL       string
	PrimaryPhotoURL string
	PlaceID         int
}

func NewGym(name, address string, latitude, longitude float64) (*Gym, error) {
	gym := &Gym{
		Name:      strings.TrimSpace(name),
		Latitude:  latitude,
		Longitude: longitude,
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
