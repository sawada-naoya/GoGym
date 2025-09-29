// internal/domain/user/types.go
// 役割: ユーザードメインの基本型（Domain Layer）
// ユーザードメイン固有の基本型とエラー定義。GORM/JSONタグは一切なし
package user

import (
	"crypto/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

// ID はULIDベースのエンティティ識別子
type ID string

// GenerateID は新しいULIDを生成する
func GenerateID() ID {
	return ID(ulid.MustNew(ulid.Timestamp(time.Now()), rand.Reader).String())
}

// ParseID はULID文字列からIDを作成する（バリデーション付き）
func ParseID(s string) (ID, error) {
	if _, err := ulid.Parse(s); err != nil {
		return "", NewDomainError("invalid_id")
	}
	return ID(s), nil
}
