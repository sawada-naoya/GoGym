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
		BaseEntity: user.BaseEntity{
			ID:        user.ID(r.ID),
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		},
		Email:        email,
		PasswordHash: r.PasswordHash,
		DisplayName:  r.DisplayName,
	}

	return userEntity, nil
}

// FromUserEntity はUserドメインエンティティをUserRecordに変換する
func FromUserEntity(u *user.User) *record.UserRecord {
	return &record.UserRecord{
		ID:           int64(u.ID),
		Email:        u.Email.String(),
		PasswordHash: u.PasswordHash,
		DisplayName:  u.DisplayName,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

// ToRefreshTokenEntity はRefreshTokenRecordをRefreshTokenドメインエンティティに変換する
func ToRefreshTokenEntity(r *record.RefreshTokenRecord) *user.RefreshToken {
	entity := &user.RefreshToken{
		ID:        user.ID(r.ID),
		UserID:    user.ID(r.UserID),
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
		ID:        int64(rt.ID),
		UserID:    int64(rt.UserID),
		TokenHash: rt.TokenHash,
		ExpiresAt: rt.ExpiresAt,
		CreatedAt: rt.CreatedAt,
	}

	if rt.User != nil {
		tokenRecord.User = FromUserEntity(rt.User)
	}

	return tokenRecord
}