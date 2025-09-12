package entity

// 必要な場合のみimportを追加

// Gym はジムドメインエンティティ
type Gym struct {
	BaseEntity
	Name            string
	Description     string
	Address         string
	PhoneNumber     string
	Email           string
	Website         string
	Latitude        float64
	Longitude       float64
	Images          []string
	AverageRating   float32
	ReviewCount     int
	ViewCount       int
	MonthlyCost     *int
	InitialCost     *int
	DailyRate       *int
	HourlyRate      *int
	UserID          ID
	
	// Relations
	Reviews   []Review
	Favorites []Favorite
	Tags      []Tag
	Locations []Location
	User      *User
}

// Review はレビュードメインエンティティ
type Review struct {
	BaseEntity
	Title     string
	Content   string
	Rating    int       // 1-5 rating
	ImageURL  *string   // Optional image
	GymID     ID
	UserID    ID
	
	// Relations
	User *User
	Gym  *Gym
}

// Favorite はお気に入りドメインエンティティ
type Favorite struct {
	BaseEntity
	UserID ID
	GymID  ID
	
	// Relations
	User *User
	Gym  *Gym
}

// Tag はタグドメインエンティティ
type Tag struct {
	BaseEntity
	Name string
	
	// Relations
	Gyms []Gym
}

// Location はロケーションドメインエンティティ
type Location struct {
	BaseEntity
	Name      string
	Latitude  float64
	Longitude float64
	
	// Relations
	Gyms []Gym
}