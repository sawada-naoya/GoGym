// internal/adapter/db/gorm/gym_repo.go
// 役割: ジムドメイン用リポジトリ実装（Infrastructure Layer）
// record↔entity変換を行いDB操作する実装。ドメインエンティティとGORMレコード間の変換を担当
package gorm

import (
	"context"
	"fmt"
	"gogym-api/internal/adapter/db/gorm/record"
	"gogym-api/internal/domain/gym"
	gymUsecase "gogym-api/internal/usecase/gym"
	"gorm.io/gorm"
	"strings"
)

// gymRepository はgym.Repositoryインターフェースを実装する
type gymRepository struct {
	db *gorm.DB
}

// NewGymRepository は新しいジムリポジトリを作成する
func NewGymRepository(db *gorm.DB) gymUsecase.Repository {
	return &gymRepository{db: db}
}

// FindByID はIDでジムを検索する
func (r *gymRepository) FindByID(ctx context.Context, id gym.ID) (*gym.Gym, error) {
	var gymRecord record.GymRecord
	if err := r.db.WithContext(ctx).Preload("Tags").First(&gymRecord, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gym.NewDomainError(gym.ErrNotFound, "gym_not_found", "gym not found")
		}
		return nil, err
	}

	return ToGymEntity(&gymRecord), nil
}

// Search はクエリとページングでジムを検索する
func (r *gymRepository) Search(ctx context.Context, query gym.SearchQuery) (*gym.PaginatedResult[gym.Gym], error) {
	var gymRecords []record.GymRecord
	
	db := r.db.WithContext(ctx).Preload("Tags")

	// 検索フィルタを適用
	if query.Query != "" {
		searchTerm := "%" + strings.ToLower(query.Query) + "%"
		db = db.Where("LOWER(name) LIKE ? OR LOWER(address) LIKE ?", searchTerm, searchTerm)
	}

	// 位置フィルタが提供されている場合は適用
	if query.Location != nil && query.RadiusM != nil {
		// 簡単な距離計算（本番環境ではPostGISの使用を検討）
		lat := query.Location.Latitude
		lng := query.Location.Longitude
		radius := float64(*query.RadiusM) / 111000.0 // メートルを度に変換（近似値）
		
		db = db.Where("location_latitude BETWEEN ? AND ? AND location_longitude BETWEEN ? AND ?",
			lat-radius, lat+radius, lng-radius, lng+radius)
	}

	// カーソルベースページングを適用
	if query.Pagination.Cursor != "" {
		db = db.Where("id > ?", query.Pagination.Cursor)
	}

	// レコードがさらに存在するかチェックするためlimit + 1を設定
	limit := query.Pagination.Limit
	if limit <= 0 {
		limit = 20 // デフォルト制限
	}
	
	if err := db.Limit(limit + 1).Find(&gymRecords).Error; err != nil {
		return nil, err
	}

	// レコードをエンティティに変換
	entities := make([]gym.Gym, 0, len(gymRecords))
	hasMore := len(gymRecords) > limit
	
	recordsToProcess := gymRecords
	if hasMore {
		recordsToProcess = gymRecords[:limit]
	}

	for _, gymRecord := range recordsToProcess {
		entities = append(entities, *ToGymEntity(&gymRecord))
	}

	// 次のカーソルを決定
	var nextCursor string
	if hasMore && len(entities) > 0 {
		nextCursor = fmt.Sprintf("%d", entities[len(entities)-1].ID)
	}

	return &gym.PaginatedResult[gym.Gym]{
		Items:      entities,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

// Create は新しいジムを作成する
func (r *gymRepository) Create(ctx context.Context, gymEntity *gym.Gym) error {
	gymRecord := FromGymEntity(gymEntity)
	
	if err := r.db.WithContext(ctx).Create(gymRecord).Error; err != nil {
		return err
	}

	// 生成されたIDでエンティティを更新
	gymEntity.ID = gym.ID(gymRecord.ID)
	gymEntity.CreatedAt = gymRecord.CreatedAt
	gymEntity.UpdatedAt = gymRecord.UpdatedAt

	return nil
}

// Update は既存のジムを更新する
func (r *gymRepository) Update(ctx context.Context, gymEntity *gym.Gym) error {
	gymRecord := FromGymEntity(gymEntity)
	
	if err := r.db.WithContext(ctx).Save(gymRecord).Error; err != nil {
		return err
	}

	// 更新タイムスタンプでエンティティを更新
	gymEntity.UpdatedAt = gymRecord.UpdatedAt

	return nil
}

// Delete はIDでジムを削除する
func (r *gymRepository) Delete(ctx context.Context, id gym.ID) error {
	result := r.db.WithContext(ctx).Delete(&record.GymRecord{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gym.NewDomainError(gym.ErrNotFound, "gym_not_found", "gym not found")
	}

	return nil
}