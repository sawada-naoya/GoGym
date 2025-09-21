package error

// ErrorResponse ログ用のエラーレスポンス
type ErrorResponse struct {
	Code    string `json:"code"`    // ログ用コード
	Message string `json:"message"` // ログ用メッセージ
}
