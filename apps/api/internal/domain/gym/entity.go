// internal/domain/gym/entity.go
// 役割: ジムドメインのEntity/VO（Domain Layer）
// ビジネスルールと不変条件を持つ純粋なドメインオブジェクト。GORM/JSONタグは一切なし
package gym

import (
	"strings"
)

// Gym はジムの集約ルートを表す
type Gym struct {
	BaseEntity
	Name          string `validate:"required,max=255"`
	Description   *string
	Location      Location
	Address       string `validate:"required,max=500"`
	City          *string
	Prefecture    *string
	PostalCode    *string
	IsActive      bool
	Tags          []Tag
	AverageRating *float32
	ReviewCount   int
}

// NewGym は検証付きで新しいジムを作成する
func NewGym(name, address string, location Location) (*Gym, error) {
	gym := &Gym{
		Name:     strings.TrimSpace(name),
		Address:  strings.TrimSpace(address),
		Location: location,
		IsActive: true, // デフォルトで有効
	}

	if err := gym.Validate(); err != nil {
		return nil, err
	}

	return gym, nil
}

// Validate はジムデータを検証する
func (g *Gym) Validate() error {
	if g.Name == "" {
		return NewDomainError("invalid_name")
	}

	if len(g.Name) > 255 {
		return NewDomainError("invalid_name")
	}

	if g.Address == "" {
		return NewDomainError("invalid_address")
	}

	if len(g.Address) > 500 {
		return NewDomainError("invalid_address")
	}

	if !g.Location.IsValid() {
		return NewDomainError("invalid_location")
	}

	return nil
}

// SetDescription はジムの説明を設定する
func (g *Gym) SetDescription(description string) {
	trimmed := strings.TrimSpace(description)
	if trimmed == "" {
		g.Description = nil
	} else {
		g.Description = &trimmed
	}
}

// SetCity はジムの都市を設定する
func (g *Gym) SetCity(city string) {
	trimmed := strings.TrimSpace(city)
	if trimmed == "" {
		g.City = nil
	} else {
		g.City = &trimmed
	}
}

// SetPrefecture はジムの都道府県を設定する
func (g *Gym) SetPrefecture(prefecture string) {
	trimmed := strings.TrimSpace(prefecture)
	if trimmed == "" {
		g.Prefecture = nil
	} else {
		g.Prefecture = &trimmed
	}
}

// SetPostalCode はジムの郵便番号を設定する
func (g *Gym) SetPostalCode(postalCode string) {
	trimmed := strings.TrimSpace(postalCode)
	if trimmed == "" {
		g.PostalCode = nil
	} else {
		g.PostalCode = &trimmed
	}
}

// Tag はジムタグエンティティを表す
type Tag struct {
	BaseEntity
	Name string `validate:"required,max=50"`
	Gyms []Gym
}

// NewTag は検証付きで新しいタグを作成する
func NewTag(name string) (*Tag, error) {
	tag := &Tag{
		Name: strings.TrimSpace(name),
	}

	if err := tag.Validate(); err != nil {
		return nil, err
	}

	return tag, nil
}

// Validate はタグデータを検証する
func (t *Tag) Validate() error {
	if t.Name == "" {
		return NewDomainError("invalid_name")
	}

	if len(t.Name) > 50 {
		return NewDomainError("invalid_name")
	}

	return nil
}

// GymTag はジムとタグの多対多関係を表す
type GymTag struct {
	GymID ID
	TagID ID
	Gym   *Gym
	Tag   *Tag
}

// Favorite はユーザーのお気に入りジムを表す
type Favorite struct {
	BaseEntity
	UserID ID
	GymID  ID
}

// NewFavorite は検証付きで新しいお気に入りを作成する
func NewFavorite(userID, gymID ID) (*Favorite, error) {
	favorite := &Favorite{
		UserID: userID,
		GymID:  gymID,
	}

	if err := favorite.Validate(); err != nil {
		return nil, err
	}

	return favorite, nil
}

// Validate はお気に入りデータを検証する
func (f *Favorite) Validate() error {
	if f.UserID == 0 {
		return NewDomainError("invalid_user_id")
	}

	if f.GymID == 0 {
		return NewDomainError("invalid_gym_id")
	}

	return nil
}
