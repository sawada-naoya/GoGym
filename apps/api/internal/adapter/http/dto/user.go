package dto

import "time"

type RegisterUserRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
}

type RegisterUserResponse struct {
	ID          int64     `json:"id"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relations (TODO: implement)
	Reviews   []ReviewResponse   `json:"reviews,omitempty"`
	Favorites []FavoriteResponse `json:"favorites,omitempty"`
}
