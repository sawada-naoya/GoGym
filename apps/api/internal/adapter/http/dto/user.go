package dto

// ユーザー認証関連DTO
type SignupRequest struct {
	Email       string `json:"email" validate:"required,email,max=255"`
	Password    string `json:"password" validate:"required,min=8,max=128"`
	DisplayName string `json:"display_name" validate:"required,max=100"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int          `json:"expires_in"`
	User         UserResponse `json:"user"`
}

type UserResponse struct {
	ID          int64  `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// S3プリサインURL関連DTO
type PresignRequest struct {
	Filename    string `json:"filename" validate:"required,max=255"`
	ContentType string `json:"content_type" validate:"required,oneof=image/jpeg image/png image/webp"`
}

type PresignResponse struct {
	PresignedURL string `json:"presigned_url"`
	PublicURL    string `json:"public_url"`
	ExpiresIn    int    `json:"expires_in"`
}

// エラーレスポンスDTO
type ErrorResponse struct {
	Error   string                 `json:"error"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}