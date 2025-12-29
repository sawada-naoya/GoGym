package contact

import "context"

type SlackGateway interface {
	NotifyContact(ctx context.Context, text string) error
	NotifyError(ctx context.Context, text string) error
}
