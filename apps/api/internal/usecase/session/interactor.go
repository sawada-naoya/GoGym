// 役割: セッションユースケースの実装（依存性注入とビジネスロジック）
// 受け取り: 外部サービスの実装（Repository, TokenService等）
// 処理: セッション管理のビジネスロジック実装
// 返却: UseCase インターフェースの実装
package session

import (
	"context"
	userUC "gogym-api/internal/usecase/user"

	"gogym-api/internal/adapter/dto"
	dom "gogym-api/internal/domain/user"
)

type interactor struct {
	// 外部依存関係
	jwt JWT
	ur  UserRepository
	idp IDProvider
	tp  TimeProvider
	ph  PasswordHasher
	// 他のユースケース
	uu userUC.UserUseCase
}

func NewInteractor(
	jwt JWT,
	ur UserRepository,
	idp IDProvider,
	tp TimeProvider,
	ph PasswordHasher,
	uu userUC.UserUseCase,
) SessionUseCase {
	return &interactor{
		jwt: jwt,
		ur:  ur,
		idp: idp,
		tp:  tp,
		ph:  ph,
		uu:  uu,
	}
}

func (i *interactor) Login(ctx context.Context, req dto.LoginRequest) error {

	// emailのバリデーション
	email, err := dom.NewEmail(req.Email)
	if err != nil {
		return err
	}

	// ユーザー検索
	user, err := i.ur.FindByEmail(ctx, email)
	if err != nil {
		return dom.NewDomainError("email_not_found")
	}
	if user == nil {
		return dom.NewDomainError("user_not_found")
	}

	// パスワード照合
	if err := i.ph.VerifyPassword(req.Password, user.PasswordHash); err != nil {
		return dom.NewDomainError("invalid_password")
	}
	return nil
}

func (i *interactor) CreateSession(ctx context.Context, email string) (dto.TokenPairResponse, error) {

	emailObj, err := dom.NewEmail(email)
	if err != nil {
		return dto.TokenPairResponse{}, err
	}

	user, err := i.ur.FindByEmail(ctx, emailObj)
	if err != nil {
		return dto.TokenPairResponse{}, dom.NewDomainError("email_not_found")
	}
	if user == nil {
		return dto.TokenPairResponse{}, dom.NewDomainError("user_not_found")
	}

	access, accessTTL, err := i.jwt.IssueAccess(user.ID)
	if err != nil {
		return dto.TokenPairResponse{}, err
	}

	refresh, _, _, exp, err := i.jwt.IssueRefresh(user.ID)
	if err != nil {
		return dto.TokenPairResponse{}, err
	}

	return dto.TokenPairResponse{
		AccessToken:  access,
		ExpiresIn:    int64(accessTTL.Seconds()),
		RefreshToken: refresh,
		RefreshExp:   exp.Unix(),
	}, nil
}
