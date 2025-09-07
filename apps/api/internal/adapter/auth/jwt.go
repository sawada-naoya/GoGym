package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims はJWTのクレーム（payload）を表す構造体
// ユーザーIDとスコープ情報を含み、JWT標準クレームも継承
type JWTClaims struct {
	UserID int64    `json:"user_id"` // ユーザーの一意識別子
	Scopes []string `json:"scopes"`  // ユーザーの権限スコープ
	jwt.RegisteredClaims             // JWT標準クレーム（iss, exp, iatなど）
}

// TokenVerifier はトークンの検証を行うインターフェース
// JWTトークンの署名検証と有効性チェックを担当
type TokenVerifier interface {
	Verify(token string) (*JWTClaims, error)
}

// TokenIssuer はトークンの発行を行うインターフェース
// アクセストークンとリフレッシュトークンの生成を担当
type TokenIssuer interface {
	IssueAccess(uid int64, scopes []string) (string, error)
	IssueRefresh(uid int64, jti string) (string, error)
}

// JWTService はJWTトークンの発行と検証を行うサービス
// HMAC-SHA256を使用してトークンの署名と検証を実行
type JWTService struct {
	secret     []byte        // トークン署名用の秘密鍵
	issuer     string        // トークン発行者識別子
	accessTTL  time.Duration // アクセストークンの有効期間
	refreshTTL time.Duration // リフレッシュトークンの有効期間
}

// NewJWT はJWTServiceの新しいインスタンスを作成
func NewJWT(secret, issuer string, accessTTL, refreshTTL time.Duration) *JWTService {
	return &JWTService{
		secret:     []byte(secret),
		issuer:     issuer,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

// IssueAccess はアクセストークンを発行
// uid: ユーザーID, scopes: 権限スコープ
func (j *JWTService) IssueAccess(uid int64, scopes []string) (string, error) {
	now := time.Now()
	claims := JWTClaims{
		UserID: uid,
		Scopes: scopes,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.accessTTL)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

// IssueRefresh はリフレッシュトークンを発行
// uid: ユーザーID, jti: JWT ID（一意識別子）
func (j *JWTService) IssueRefresh(uid int64, jti string) (string, error) {
	now := time.Now()
	claims := JWTClaims{
		UserID: uid,
		Scopes: []string{"refresh"}, // リフレッシュ専用スコープ
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti, // JWT ID
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.refreshTTL)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

// Verify はJWTトークンを検証してクレームを返す
// トークンの署名検証、有効期限チェック、形式検証を実行
func (j *JWTService) Verify(token string) (*JWTClaims, error) {
	// トークンをパースしてクレームを取得
	parsedToken, err := jwt.ParseWithClaims(
		token,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// 署名方式の確認
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return j.secret, nil
		},
	)
	
	if err != nil {
		return nil, err
	}
	
	// トークンの有効性確認
	if !parsedToken.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	
	// クレームの型アサーション
	claims, ok := parsedToken.Claims.(*JWTClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}
	
	return claims, nil
}
