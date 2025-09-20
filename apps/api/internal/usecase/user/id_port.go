package user

import dom "gogym-api/internal/domain/user"

// IDProvider はIDを生成するインターフェース（実際にはドメインに移譲）
type IDProvider interface {
	Generate() dom.ID
}