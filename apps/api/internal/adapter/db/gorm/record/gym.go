// internal/adapter/db/gorm/record/gym.go
// 役割: ジム関連のGORMレコード構造体（Infrastructure Layer）
// DB行の形でGORMタグ付きstruct。ドメインエンティティとの変換はconverterで実行
package record

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// GymRecord はジムエンティティ用のGORMレコードを表す
type GymRecord struct {
	ID                int64       `gorm:"primaryKey;autoIncrement"`
	Name              string      `gorm:"size:255;not null"`
	Description       *string     `gorm:"type:text"`
	Location          string      `gorm:"column:location;type:point"`
	LocationLatitude  float64     `gorm:"-"` // 計算フィールド
	LocationLongitude float64     `gorm:"-"` // 計算フィールド
	Address           string      `gorm:"size:500;not null"`
	City              *string     `gorm:"size:100"`
	Prefecture        *string     `gorm:"size:100"`
	PostalCode        *string     `gorm:"size:10"`
	AverageRating     *float32    `gorm:"column:average_rating;type:decimal(3,2)"`
	ReviewCount       int         `gorm:"column:review_count;default:0"`
	CreatedAt         time.Time   `gorm:"autoCreateTime"`
	UpdatedAt         time.Time   `gorm:"autoUpdateTime"`
	Tags              []TagRecord `gorm:"many2many:gym_tags;foreignKey:ID;joinForeignKey:gym_id;References:ID;joinReferences:tag_id;"`
}

// TableName はGORM用のテーブル名を返す
func (GymRecord) TableName() string {
	return "gyms"
}

// AfterFind はGORM後処理フック - POINT型から座標を抽出
func (g *GymRecord) AfterFind(tx *gorm.DB) error {
	// POINT型文字列から座標を抽出 (例: "POINT(35.6598 139.7016)" - latitude longitude順)
	if g.Location != "" {
		// POINT(lat lon) 形式をパース
		locationStr := strings.TrimSpace(g.Location)
		if strings.HasPrefix(locationStr, "POINT(") && strings.HasSuffix(locationStr, ")") {
			coords := strings.TrimSuffix(strings.TrimPrefix(locationStr, "POINT("), ")")
			parts := strings.Split(coords, " ")
			if len(parts) == 2 {
				var latitude, longitude float64
				if _, err := fmt.Sscanf(parts[0], "%f", &latitude); err == nil {
					g.LocationLatitude = latitude
				}
				if _, err := fmt.Sscanf(parts[1], "%f", &longitude); err == nil {
					g.LocationLongitude = longitude
				}
			}
		}
	}
	return nil
}

// TagRecord はタグエンティティ用のGORMレコードを表す
type TagRecord struct {
	ID        int64       `gorm:"primaryKey;autoIncrement"`
	Name      string      `gorm:"uniqueIndex;size:50;not null"`
	CreatedAt time.Time   `gorm:"autoCreateTime"`
	UpdatedAt time.Time   `gorm:"autoUpdateTime"`
	Gyms      []GymRecord `gorm:"many2many:gym_tags;foreignKey:ID;joinForeignKey:tag_id;References:ID;joinReferences:gym_id;"`
}

// TableName はGORM用のテーブル名を返す
func (TagRecord) TableName() string {
	return "tags"
}

// GymTagRecord は多対多関係用のGORMレコードを表す
type GymTagRecord struct {
	GymID int64      `gorm:"primaryKey"`
	TagID int64      `gorm:"primaryKey"`
	Gym   *GymRecord `gorm:"foreignKey:GymID"`
	Tag   *TagRecord `gorm:"foreignKey:TagID"`
}

// TableName はGORM用のテーブル名を返す
func (GymTagRecord) TableName() string {
	return "gym_tags"
}

// FavoriteRecord はお気に入りエンティティ用のGORMレコードを表す
type FavoriteRecord struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	UserID    int64     `gorm:"not null;index;uniqueIndex:unique_user_gym_favorite,priority:1"`
	GymID     int64     `gorm:"not null;index;uniqueIndex:unique_user_gym_favorite,priority:2"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TableName はGORM用のテーブル名を返す
func (FavoriteRecord) TableName() string {
	return "favorites"
}
