package repository

import (
	"context"
	"gogym-api/internal/adapter/repository/mapper"
	"gogym-api/internal/adapter/repository/record"
	"gogym-api/internal/domain/gym"
	tagUsecase "gogym-api/internal/usecase/tag"

	"gorm.io/gorm"
)

// tagRepository はtag.Repositoryインターフェースを実装する
type tagRepository struct {
	db *gorm.DB
}

// NewTagRepository は新しいタグリポジトリを作成する
func NewTagRepository(db *gorm.DB) tagUsecase.Repository {
	return &tagRepository{db: db}
}

// FindAll はすべてのタグを検索する
func (r *tagRepository) FindAll(ctx context.Context) ([]gym.Tag, error) {
	var tagRecords []record.TagRecord
	if err := r.db.WithContext(ctx).Find(&tagRecords).Error; err != nil {
		return nil, err
	}

	tags := make([]gym.Tag, len(tagRecords))
	for i, tagRecord := range tagRecords {
		tags[i] = *mapper.ToTagEntity(&tagRecord)
	}

	return tags, nil
}

// FindByIDs はIDリストでタグを検索する
func (r *tagRepository) FindByIDs(ctx context.Context, ids []gym.ID) ([]gym.Tag, error) {
	var tagRecords []record.TagRecord
	int64IDs := make([]int64, len(ids))
	for i, id := range ids {
		int64IDs[i] = int64(id)
	}

	if err := r.db.WithContext(ctx).Where("id IN ?", int64IDs).Find(&tagRecords).Error; err != nil {
		return nil, err
	}

	tags := make([]gym.Tag, len(tagRecords))
	for i, tagRecord := range tagRecords {
		tags[i] = *mapper.ToTagEntity(&tagRecord)
	}

	return tags, nil
}

// FindByNames は名前リストでタグを検索する
func (r *tagRepository) FindByNames(ctx context.Context, names []string) ([]gym.Tag, error) {
	var tagRecords []record.TagRecord
	if err := r.db.WithContext(ctx).Where("name IN ?", names).Find(&tagRecords).Error; err != nil {
		return nil, err
	}

	tags := make([]gym.Tag, len(tagRecords))
	for i, tagRecord := range tagRecords {
		tags[i] = *mapper.ToTagEntity(&tagRecord)
	}

	return tags, nil
}

// Create は新しいタグを作成する
func (r *tagRepository) Create(ctx context.Context, tagEntity *gym.Tag) error {
	tagRecord := mapper.FromTagEntity(tagEntity)

	if err := r.db.WithContext(ctx).Create(tagRecord).Error; err != nil {
		return err
	}

	// 生成されたIDでエンティティを更新
	tagEntity.ID = gym.ID(tagRecord.ID)
	tagEntity.CreatedAt = tagRecord.CreatedAt
	tagEntity.UpdatedAt = tagRecord.UpdatedAt

	return nil
}

// CreateMany は複数のタグを作成する
func (r *tagRepository) CreateMany(ctx context.Context, tagEntities []gym.Tag) error {
	if len(tagEntities) == 0 {
		return nil
	}

	tagRecords := make([]record.TagRecord, len(tagEntities))
	for i, tagEntity := range tagEntities {
		tagRecords[i] = *mapper.FromTagEntity(&tagEntity)
	}

	if err := r.db.WithContext(ctx).Create(&tagRecords).Error; err != nil {
		return err
	}

	// 生成されたIDとタイムスタンプでエンティティを更新
	for i, tagRecord := range tagRecords {
		tagEntities[i].ID = gym.ID(tagRecord.ID)
		tagEntities[i].CreatedAt = tagRecord.CreatedAt
		tagEntities[i].UpdatedAt = tagRecord.UpdatedAt
	}

	return nil
}
