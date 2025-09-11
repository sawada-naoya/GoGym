package dto

import "gogym-api/internal/domain/gym"

// リクエストDTO - HTTPリクエストから受け取るデータ
type CreateGymRequest struct {
	Name        string       `json:"name" validate:"required,max=255"`
	Description *string      `json:"description,omitempty"`
	Location    gym.Location `json:"location" validate:"required"`
	Address     string       `json:"address" validate:"required,max=500"`
	City        *string      `json:"city,omitempty" validate:"omitempty,max=100"`
	Prefecture  *string      `json:"prefecture,omitempty" validate:"omitempty,max=100"`
	PostalCode  *string      `json:"postal_code,omitempty" validate:"omitempty,pattern=^[0-9]{3}-[0-9]{4}$"`
	TagNames    []string     `json:"tag_names,omitempty"`
}

type SearchGymRequest struct {
	Query   string   `json:"q,omitempty"`
	Lat     *float64 `json:"lat,omitempty" validate:"omitempty,min=-90,max=90"`
	Lon     *float64 `json:"lon,omitempty" validate:"omitempty,min=-180,max=180"`
	RadiusM *int     `json:"radius_m,omitempty" validate:"omitempty,min=100,max=50000"`
	Cursor  string   `json:"cursor,omitempty"`
	Limit   int      `json:"limit,omitempty" validate:"omitempty,min=1,max=100"`
}

// レスポンスDTO - HTTPレスポンスで返すデータ
type GymResponse struct {
	ID            int64         `json:"id"`
	Name          string        `json:"name"`
	Description   *string       `json:"description"`
	Location      gym.Location  `json:"location"`
	Address       string        `json:"address"`
	City          *string       `json:"city"`
	Prefecture    *string       `json:"prefecture"`
	PostalCode    *string       `json:"postal_code"`
	Tags          []TagResponse `json:"tags"`
	AverageRating *float32      `json:"average_rating"`
	ReviewCount   int           `json:"review_count"`
	CreatedAt     string        `json:"created_at"`
	UpdatedAt     string        `json:"updated_at"`
}

type TagResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type SearchGymResponse struct {
	Gyms       []GymResponse `json:"gyms"`
	NextCursor *string       `json:"next_cursor"`
	HasMore    bool          `json:"has_more"`
}
