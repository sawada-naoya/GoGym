package contact

import "context"

type ContactUseCase interface {
	SendContact(ctx context.Context, email, message string, userID *string, ip, ua string) error
}
