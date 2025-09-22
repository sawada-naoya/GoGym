package user

import (
	"context"
	"gogym-api/internal/adapter/http/dto"
)

// handler → usecase
type UseCase interface {
	SignUp(ctx context.Context, req dto.SignUpRequest) error
}
