package handler

import (
	gu "gogym-api/internal/application/gym"
)

type GymHandler struct {
	gu gu.GymUseCase
}

func NewGymHandler(gu gu.GymUseCase) *GymHandler {
	return &GymHandler{
		gu: gu,
	}
}
