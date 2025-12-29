package contact

import "context"

type ContactUseCase interface {
	SendContact(ctx context.Context, in SendContactInput) error
}
