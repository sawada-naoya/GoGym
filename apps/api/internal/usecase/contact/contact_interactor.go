package contact

import (
	"context"
	"errors"
	"strings"
)

type contactInteractor struct {
	sg SlackGateway
}

func NewContactInteractor(sg SlackGateway) ContactUseCase {
	return &contactInteractor{
		sg: sg,
	}
}

type SendContactInput struct {
	Email   string
	Message string
	UserID  *string
	IP      string
	UA      string
}

func (i *contactInteractor) SendContact(ctx context.Context, in SendContactInput) error {
	email := strings.TrimSpace(in.Email)
	msg := strings.TrimSpace(in.Message)

	if email == "" || msg == "" {
		return errors.New("email and message are required")
	}
	if len(email) > 255 || len(msg) > 2000 {
		return errors.New("email or message is too long")
	}

	// 「何を送るか」はusecaseで決める
	text := "📩 Contact\nEmail: " + email + "\n\n" + msg
	return i.sg.NotifyContact(ctx, text)
}

func (i *contactInteractor) SendError(ctx context.Context, text string) error {
	t := strings.TrimSpace(text)
	if t == "" {
		return errors.New("text is required")
	}
	if len(t) > 4000 {
		return errors.New("text is too long")
	}
	return i.sg.NotifyError(ctx, t)
}
