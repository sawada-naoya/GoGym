package handler

import (
	gu "gogym-api/internal/usecase/gym"
)

type GymHandler struct {
	gu gu.GymUseCase
}

func NewGymHandler(gu gu.GymUseCase) *GymHandler {
	return &GymHandler{
		gu: gu,
	}
}
