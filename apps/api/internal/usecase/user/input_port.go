package user

import (
	"context"
)

// handler → usecase
type UseCase interface {
	SignUp(ctx context.Context, req SignUpRequest) error
}
