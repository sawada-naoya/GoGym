// internal/adapter/db/gorm/record/converter.go
// 役割: Entity ↔ Record 変換ユーティリティ（Infrastructure Layer）
// ドメインエンティティとGORMレコード間の双方向変換を担当。クリーンアーキテクチャの境界を維持
package record

import (
	"gogym-api/internal/domain/common"
	"gogym-api/internal/domain/gym"
	"gogym-api/internal/domain/user"
)

// ユーザードメイン変換

// ToUserEntity はUserRecordをUserドメインエンティティに変換する
func (r *UserRecord) ToEntity() (*user.User, error) {
	email, err := user.NewEmail(r.Email)
	if err != nil {
		return nil, err
	}

	userEntity := &user.User{
		BaseEntity: common.BaseEntity{
			ID:        common.ID(r.ID),
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
func FromUserEntity(u *user.User) *UserRecord {
	return &UserRecord{
		ID:           int64(u.ID),
		Email:        u.Email.String(),
		PasswordHash: u.PasswordHash,
		DisplayName:  u.DisplayName,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

// ToRefreshTokenEntity はRefreshTokenRecordをRefreshTokenドメインエンティティに変換する
func (r *RefreshTokenRecord) ToEntity() *user.RefreshToken {
	entity := &user.RefreshToken{
		ID:        common.ID(r.ID),
		UserID:    common.ID(r.UserID),
		TokenHash: r.TokenHash,
		ExpiresAt: r.ExpiresAt,
		CreatedAt: r.CreatedAt,
	}

	if r.User != nil {
		if userEntity, err := r.User.ToEntity(); err == nil {
			entity.User = userEntity
		}
	}

	return entity
}

// FromRefreshTokenEntity はRefreshTokenドメインエンティティをRefreshTokenRecordに変換する
func FromRefreshTokenEntity(rt *user.RefreshToken) *RefreshTokenRecord {
	record := &RefreshTokenRecord{
		ID:        int64(rt.ID),
		UserID:    int64(rt.UserID),
		TokenHash: rt.TokenHash,
		ExpiresAt: rt.ExpiresAt,
		CreatedAt: rt.CreatedAt,
	}

	if rt.User != nil {
		record.User = FromUserEntity(rt.User)
	}

	return record
}

// ジムドメイン変換

// ToGymEntity はGymRecordをGymドメインエンティティに変換する
func (r *GymRecord) ToEntity() *gym.Gym {
	entity := &gym.Gym{
		BaseEntity: common.BaseEntity{
			ID:        common.ID(r.ID),
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		},
		Name:        r.Name,
		Description: r.Description,
		Location: common.Location{
			Latitude:  r.LocationLatitude,
			Longitude: r.LocationLongitude,
		},
		Address:    r.Address,
		City:       r.City,
		Prefecture: r.Prefecture,
		PostalCode: r.PostalCode,
	}

	// タグが存在する場合は変換する
	if len(r.Tags) > 0 {
		entity.Tags = make([]gym.Tag, len(r.Tags))
		for i, tagRecord := range r.Tags {
			entity.Tags[i] = *tagRecord.ToTagEntity()
		}
	}

	return entity
}

// FromGymEntity はGymドメインエンティティをGymRecordに変換する
func FromGymEntity(g *gym.Gym) *GymRecord {
	record := &GymRecord{
		ID:                int64(g.ID),
		Name:              g.Name,
		Description:       g.Description,
		LocationLatitude:  g.Location.Latitude,
		LocationLongitude: g.Location.Longitude,
		Address:           g.Address,
		City:              g.City,
		Prefecture:        g.Prefecture,
		PostalCode:        g.PostalCode,
		CreatedAt:         g.CreatedAt,
		UpdatedAt:         g.UpdatedAt,
	}

	// タグが存在する場合は変換する
	if len(g.Tags) > 0 {
		record.Tags = make([]TagRecord, len(g.Tags))
		for i, tagEntity := range g.Tags {
			record.Tags[i] = *FromTagEntity(&tagEntity)
		}
	}

	return record
}

// ToTagEntity はTagRecordをTagドメインエンティティに変換する
func (r *TagRecord) ToTagEntity() *gym.Tag {
	return &gym.Tag{
		BaseEntity: common.BaseEntity{
			ID:        common.ID(r.ID),
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		},
		Name: r.Name,
	}
}

// FromTagEntity はTagドメインエンティティをTagRecordに変換する
func FromTagEntity(t *gym.Tag) *TagRecord {
	return &TagRecord{
		ID:        int64(t.ID),
		Name:      t.Name,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

// ToFavoriteEntity はFavoriteRecordをFavoriteドメインエンティティに変換する
func (r *FavoriteRecord) ToEntity() *gym.Favorite {
	return &gym.Favorite{
		BaseEntity: common.BaseEntity{
			ID:        common.ID(r.ID),
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		},
		UserID: common.ID(r.UserID),
		GymID:  common.ID(r.GymID),
	}
}

// FromFavoriteEntity はFavoriteドメインエンティティをFavoriteRecordに変換する
func FromFavoriteEntity(f *gym.Favorite) *FavoriteRecord {
	return &FavoriteRecord{
		ID:        int64(f.ID),
		UserID:    int64(f.UserID),
		GymID:     int64(f.GymID),
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}