// 役割: セッション管理用のHTTP DTO（Data Transfer Object）
// 受け取り: HTTPリクエスト/レスポンスのJSON
// 処理: バリデーション、ドメインエンティティとの相互変換
// 返却: ドメイン層とHTTP層間のデータ変換結果
package dto

type TokenPairResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	RefreshExp   int64  `json:"refresh_exp,omitempty"`
}
