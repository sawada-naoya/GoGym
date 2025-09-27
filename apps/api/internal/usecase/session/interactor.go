// 役割: セッションユースケースの実装（依存性注入とビジネスロジック）
// 受け取り: 外部サービスの実装（Repository, TokenService等）
// 処理: セッション管理のビジネスロジック実装
// 返却: UseCase インターフェースの実装
package session

import (
	userUC "gogym-api/internal/usecase/user"
)

type interactor struct {
	// 外部依存関係
	jwt JWT
	ur  UserRepository
	idp IDProvider
	tp  TimeProvider
	ph  PasswordHasher
	// 他のユースケース
	uu userUC.UseCase
}

func NewInteractor(
	jwt JWT,
	ur UserRepository,
	idp IDProvider,
	tp TimeProvider,
	ph PasswordHasher,
	uu userUC.UseCase,
) UseCase {
	return &interactor{
		jwt: jwt,
		ur:  ur,
		idp: idp,
		tp:  tp,
		ph:  ph,
		uu:  uu,
	}
}
