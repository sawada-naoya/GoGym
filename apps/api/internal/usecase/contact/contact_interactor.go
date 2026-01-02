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

func (i *contactInteractor) SendContact(ctx context.Context, email, message string, userID *string, ip, ua string) error {
	// ビジネスロジック: バリデーション
	email = strings.TrimSpace(email)
	message = strings.TrimSpace(message)

	if email == "" || message == "" {
		return errors.New("email and message are required")
	}
	if len(email) > 255 || len(message) > 2000 {
		return errors.New("email or message is too long")
	}

	// メッセージの整形（ビジネスロジック）
	userIDStr := "anonymous"
	if userID != nil && *userID != "" {
		userIDStr = *userID
	}

	text := "📩 Contact Form Submission\n\n" +
		"From: " + email + "\n" +
		"User ID: " + userIDStr + "\n" +
		"IP: " + ip + "\n" +
		"User Agent: " + ua + "\n\n" +
		"Message:\n" + message

	return i.sg.NotifyContact(ctx, text)
}
