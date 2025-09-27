package user

import (
	"context"
	"gogym-api/internal/adapter/http/dto"
)

type UseCase interface {
	SignUp(ctx context.Context, req dto.SignUpRequest) error
}
