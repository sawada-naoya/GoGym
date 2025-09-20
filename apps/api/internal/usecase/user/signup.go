package user

import (
	"context"
	dom "gogym-api/internal/domain/user"
)

// 具象は非公開にして、外には UseCase として返す
type interactor struct {
	repo   Repository
	hasher PasswordHasher
}

func NewInteractor(r Repository, h PasswordHasher) UseCase {
	return &interactor{repo: r, hasher: h}
}

func (i *interactor) SignUp(ctx context.Context, req SignUpRequest) (SignUpResult, error) {
	// 1. Email バリューオブジェクト作成（バリデーション含む）
	email, err := dom.NewEmail(req.Email)
	if err != nil {
		return SignUpResult{}, err
	}

	// 2. メールアドレスの重複チェック
	exists, err := i.repo.ExistsByEmail(ctx, email)
	if err != nil {
		return SignUpResult{}, err
	}
	if exists {
		return SignUpResult{}, dom.NewDomainError(dom.ErrAlreadyExists, "email_already_exists", "このメールアドレスは既に使用されています")
	}

	// 3. パスワードハッシュ化
	hashedPassword, err := i.hasher.HashPassword(req.Password)
	if err != nil {
		return SignUpResult{}, err
	}

	// 4. ユーザーエンティティ作成（IDは自動生成される）
	user, err := dom.NewUser(req.Name, email, hashedPassword)
	if err != nil {
		return SignUpResult{}, err
	}

	// 5. データベースに保存
	err = i.repo.Create(ctx, user)
	if err != nil {
		return SignUpResult{}, err
	}

	// 6. 結果を返す
	return SignUpResult{
		UserID: string(user.ID),
		Name:   user.Name,
		Email:  user.Email.String(),
	}, nil
}
