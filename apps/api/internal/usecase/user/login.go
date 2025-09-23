// 役割: ユーザー認証（email/password照合）
// 受け取り: LoginRequest（email, password）
// 処理: メールアドレス検索、パスワード照合
// 返却: 認証結果（成功時はnil）、エラー
package user

import (
	"context"
	"log/slog"

	"gogym-api/internal/adapter/http/dto"
	dom "gogym-api/internal/domain/user"
)

func (i *interactor) Login(ctx context.Context, req dto.LoginRequest) error {
	slog.InfoContext(ctx, "Login UseCase", "Email", req.Email)

	// emailのバリデーション
	email, err := dom.NewEmail(req.Email)
	if err != nil {
		return err
	}

	// ユーザー検索
	user, err := i.repo.FindByEmail(ctx, email)
	if err != nil {
		return dom.NewDomainError(dom.ErrNotFound, "user_not_found", "ユーザーが見つかりません")
	}
	if user == nil {
		return dom.NewDomainError(dom.ErrNotFound, "user_not_found", "ユーザーが見つかりません")
	}

	// パスワード照合
	if err := i.hasher.VerifyPassword(req.Password, user.PasswordHash); err != nil {
		return dom.NewDomainError(dom.ErrUnauthorized, "invalid_credentials", "メールアドレスまたはパスワードが間違っています")
	}
	return nil
}