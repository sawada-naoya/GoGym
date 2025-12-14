package gym

import (
	domain "gogym-api/internal/domain/entities"
)

// ToEntity converts GymRecord to domain entity
func ToEntity(r *GymRecord) *domain.Gym {
	if r == nil {
		return nil
	}

	gym := &domain.Gym{
		ID:        int(r.ID),
		Name:      r.Name,
		Latitude:  r.Latitude,
		Longitude: r.Longitude,
		SourceURL: r.SourceURL,
	}

	if r.PrimaryPhotoURL != nil {
		gym.PrimaryPhotoURL = *r.PrimaryPhotoURL
	}

	if r.PlaceID != nil {
		// PlaceIDは文字列だがdomainではint型の場合、適切な変換が必要
		// 現在のスキーマではVARCHAR(128)なので、そのまま保存する場合はdomainも変更が必要
		// ここでは一旦0として扱う（要確認）
		gym.PlaceID = 0
	}

	return gym
}

// FromEntity converts domain entity to GymRecord
func FromEntity(g *domain.Gym, createdBy string) *GymRecord {
	if g == nil {
		return nil
	}

	record := &GymRecord{
		ID:        int64(g.ID),
		Name:      g.Name,
		Latitude:  g.Latitude,
		Longitude: g.Longitude,
		SourceURL: g.SourceURL,
		CreatedBy: createdBy,
	}

	if g.PrimaryPhotoURL != "" {
		record.PrimaryPhotoURL = &g.PrimaryPhotoURL
	}

	// PlaceIDの変換（domain側がintの場合の仮実装）
	// 実際はdomain側をstringにするか、この変換ロジックを修正する必要がある
	if g.PlaceID != 0 {
		// placeIDStr := strconv.Itoa(g.PlaceID)
		// record.PlaceID = &placeIDStr
	}

	return record
}

// ToEntities converts slice of GymRecord to slice of domain entities
func ToEntities(records []*GymRecord) []*domain.Gym {
	entities := make([]*domain.Gym, 0, len(records))
	for _, r := range records {
		entities = append(entities, ToEntity(r))
	}
	return entities
}
