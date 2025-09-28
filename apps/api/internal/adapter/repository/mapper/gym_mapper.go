package mapper

import (
	"gogym-api/internal/adapter/repository/record"
	"gogym-api/internal/domain/gym"
)

// ToGymEntity はGymRecordをGymドメインエンティティに変換する
func ToGymEntity(r *record.GymRecord) *gym.Gym {
	entity := &gym.Gym{
		BaseEntity: gym.BaseEntity{
			ID:        gym.ID(r.ID),
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		},
		Name:        r.Name,
		Description: r.Description,
		Location: gym.Location{
			Latitude:  r.Latitude,
			Longitude: r.Longitude,
		},
		Address:    r.Address,
		City:       r.City,
		Prefecture: r.Prefecture,
		PostalCode: r.PostalCode,
		IsActive:   r.IsActive,
	}

	// タグが存在する場合は変換する
	if len(r.Tags) > 0 {
		entity.Tags = make([]gym.Tag, len(r.Tags))
		for i, tagRecord := range r.Tags {
			entity.Tags[i] = *ToTagEntity(&tagRecord)
		}
	}

	return entity
}

// FromGymEntity はGymドメインエンティティをGymRecordに変換する
func FromGymEntity(g *gym.Gym) *record.GymRecord {
	gymRecord := &record.GymRecord{
		ID:          int64(g.ID),
		Name:        g.Name,
		Description: g.Description,
		Latitude:    g.Location.Latitude,
		Longitude:   g.Location.Longitude,
		Address:     g.Address,
		City:        g.City,
		Prefecture:  g.Prefecture,
		PostalCode:  g.PostalCode,
		IsActive:    g.IsActive,
		CreatedAt:   g.CreatedAt,
		UpdatedAt:   g.UpdatedAt,
	}

	// タグが存在する場合は変換する
	if len(g.Tags) > 0 {
		gymRecord.Tags = make([]record.TagRecord, len(g.Tags))
		for i, tagEntity := range g.Tags {
			gymRecord.Tags[i] = *FromTagEntity(&tagEntity)
		}
	}

	return gymRecord
}

// ToTagEntity はTagRecordをTagドメインエンティティに変換する
func ToTagEntity(r *record.TagRecord) *gym.Tag {
	return &gym.Tag{
		BaseEntity: gym.BaseEntity{
			ID:        gym.ID(r.ID),
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		},
		Name: r.Name,
	}
}

// FromTagEntity はTagドメインエンティティをTagRecordに変換する
func FromTagEntity(t *gym.Tag) *record.TagRecord {
	return &record.TagRecord{
		ID:        int64(t.ID),
		Name:      t.Name,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

// ToFavoriteEntity はFavoriteRecordをFavoriteドメインエンティティに変換する
func ToFavoriteEntity(r *record.FavoriteRecord) *gym.Favorite {
	return &gym.Favorite{
		BaseEntity: gym.BaseEntity{
			ID:        gym.ID(r.ID),
			CreatedAt: r.CreatedAt,
		},
		UserID: gym.ID(r.UserID),
		GymID:  gym.ID(r.GymID),
	}
}

// FromFavoriteEntity はFavoriteドメインエンティティをFavoriteRecordに変換する
func FromFavoriteEntity(f *gym.Favorite) *record.FavoriteRecord {
	return &record.FavoriteRecord{
		ID:        int64(f.ID),
		UserID:    int64(f.UserID),
		GymID:     int64(f.GymID),
		CreatedAt: f.CreatedAt,
	}
}
