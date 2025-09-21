package user

import (
	"context"
	dom "gogym-api/internal/domain/user"
	"log/slog"
	"time"
)

// SignUp handles user registration
func (i *interactor) SignUp(ctx context.Context, req SignUpRequest) error {
	slog.InfoContext(ctx, "SignUp UseCase", "Name", req.Name, "Email", req.Email)
	// emailのバリデーション
	email, err := dom.NewEmail(req.Email)
	if err != nil {
		return err
	}

	// メールアドレスの重複チェック
	exists, err := i.repo.ExistsByEmail(ctx, email)
	if err != nil {
		return err
	}
	if exists {
		return dom.NewDomainError(dom.ErrAlreadyExists, "email_already_exists", "このメールアドレスは既に使用されています")
	}

	// パスワードハッシュ化
	hashedPassword, err := i.hasher.HashPassword(req.Password)
	if err != nil {
		return err
	}

	// ユーザーIDの生成
	id := i.idProvider.NewUserID()

	now := time.Now()

	// ユーザーエンティティの生成
	user, err := dom.NewUser(id, req.Name, email, hashedPassword, now)
	if err != nil {
		return err
	}

	// データベースに保存
	err = i.repo.Create(ctx, user)
	if err != nil {
		return err
	}

	return nil
}
