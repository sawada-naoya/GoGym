package user

// SignUpRequest はユーザー登録のリクエスト
type SignUpRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignUpResult はユーザー登録の結果
type SignUpResult struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}
