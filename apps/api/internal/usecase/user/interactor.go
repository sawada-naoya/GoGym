package user

import (
	"context"
	"crypto/rand"
	"time"

	"gogym-api/internal/adapter/dto"
	dom "gogym-api/internal/domain/user"

	"github.com/oklog/ulid/v2"
)

type interactor struct {
	repo   Repository
	hasher PasswordHasher
}

func NewInteractor(repo Repository, hasher PasswordHasher) UserUseCase {
	return &interactor{
		repo:   repo,
		hasher: hasher,
	}
}

// SignUp handles user registration
func (i *interactor) SignUp(ctx context.Context, req dto.SignUpRequest) error {
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
		return dom.NewDomainError("email_already_exists")
	}

	// パスワードハッシュ化
	hashedPassword, err := i.hasher.HashPassword(req.Password)
	if err != nil {
		return err
	}

	// ユーザーID(ulid)の生成
	t := time.Now()
	entropy := ulid.Monotonic(rand.Reader, 0)
	id := ulid.MustNew(ulid.Timestamp(t), entropy)

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
