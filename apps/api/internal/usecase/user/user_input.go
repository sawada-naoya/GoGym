package user

import (
	"context"
	"gogym-api/internal/adapter/dto"
)

type UserUseCase interface {
	SignUp(ctx context.Context, req dto.SignUpRequest) error
}
