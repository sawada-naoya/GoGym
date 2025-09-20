// internal/adapter/db/gorm/user_mapper.go
// 役割: User Entity ↔ Record 変換ユーティリティ（Infrastructure Layer）
// ユーザードメインエンティティとGORMレコード間の双方向変換を担当
package gorm

import (
	"gogym-api/internal/adapter/db/gorm/record"
	"gogym-api/internal/domain/user"
)

// ToUserEntity はUserRecordをUserドメインエンティティに変換する
func ToUserEntity(r *record.UserRecord) (*user.User, error) {
	email, err := user.NewEmail(r.Email)
	if err != nil {
		return nil, err
	}

	userEntity := &user.User{
		ID:           user.ID(r.ID), // ULID文字列をそのまま変換
		Name:         r.Name,
		Email:        email,
		PasswordHash: r.PasswordHash,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}

	return userEntity, nil
}

// FromUserEntity はUserドメインエンティティをUserRecordに変換する
func FromUserEntity(u *user.User) *record.UserRecord {
	return &record.UserRecord{
		ID:           string(u.ID), // ULIDをそのまま文字列として格納
		Email:        u.Email.String(),
		PasswordHash: u.PasswordHash,
		Name:         u.Name, // DisplayName→Nameに統一
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

// ToRefreshTokenEntity はRefreshTokenRecordをRefreshTokenドメインエンティティに変換する
func ToRefreshTokenEntity(r *record.RefreshTokenRecord) *user.RefreshToken {
	entity := &user.RefreshToken{
		ID:        user.ID(r.ID),        // ULID文字列をそのまま変換
		UserID:    user.ID(r.UserID),    // ULID文字列をそのまま変換
		TokenHash: r.TokenHash,
		ExpiresAt: r.ExpiresAt,
		CreatedAt: r.CreatedAt,
	}

	if r.User != nil {
		if userEntity, err := ToUserEntity(r.User); err == nil {
			entity.User = userEntity
		}
	}

	return entity
}

// FromRefreshTokenEntity はRefreshTokenドメインエンティティをRefreshTokenRecordに変換する
func FromRefreshTokenEntity(rt *user.RefreshToken) *record.RefreshTokenRecord {
	tokenRecord := &record.RefreshTokenRecord{
		ID:        string(rt.ID),     // ULIDをそのまま文字列として格納
		UserID:    string(rt.UserID), // ULIDをそのまま文字列として格納
		TokenHash: rt.TokenHash,
		ExpiresAt: rt.ExpiresAt,
		CreatedAt: rt.CreatedAt,
	}

	if rt.User != nil {
		tokenRecord.User = FromUserEntity(rt.User)
	}

	return tokenRecord
}