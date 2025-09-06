package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gogym-api/internal/adapter/http/dto"
	"gogym-api/internal/domain/common"
)

// GymHandler handles gym-related HTTP requests
type GymHandler struct {
	// Add any dependencies like services here
}

// NewGymHandler creates a new GymHandler instance
func NewGymHandler() *GymHandler {
	return &GymHandler{}
}

