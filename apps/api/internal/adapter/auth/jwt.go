// internal/auth/jwt.go
package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid/v2"
)

type Service struct {
	accessSecret  []byte
	refreshSecret []byte
	issuer        string
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

// New 生成
func New(accessSecret, refreshSecret []byte, issuer string, accessTTL, refreshTTL time.Duration) *Service {
	return &Service{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
		issuer:        issuer,
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
	}
}

// IssueAccess アクセストークン発行。sub=userID を入れる
func (s *Service) IssueAccess(userID string) (token string, ttl time.Duration, err error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub": userID,
		"iss": s.issuer,
		"iat": now.Unix(),
		"exp": now.Add(s.accessTTL).Unix(),
		"typ": "access",
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, err := t.SignedString(s.accessSecret)
	return str, s.accessTTL, err
}

// IssueRefresh リフレッシュトークン発行。jti=ULID を入れる
func (s *Service) IssueRefresh(userID string) (token string, ttl time.Duration, jti string, exp time.Time, err error) {
	now := time.Now()
	jti = ulid.Make().String()
	exp = now.Add(s.refreshTTL)
	claims := jwt.MapClaims{
		"sub": userID,
		"jti": jti,
		"iss": s.issuer,
		"iat": now.Unix(),
		"exp": exp.Unix(),
		"typ": "refresh",
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, err := t.SignedString(s.refreshSecret)
	return str, s.refreshTTL, jti, exp, err
}

// RefreshClaims 検証後に必要な最小情報だけ返す
type RefreshClaims struct {
	JTI    string
	UserID string
	Exp    time.Time
}

// ParseRefresh リフレッシュトークン検証。署名/期限/typ を確認して JTI/USER を返す
func (s *Service) ParseRefresh(tokenStr string) (RefreshClaims, error) {
	tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected alg")
		}
		return s.refreshSecret, nil
	})
	if err != nil || !tok.Valid {
		return RefreshClaims{}, errors.New("invalid refresh token")
	}

	mc, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return RefreshClaims{}, errors.New("invalid claims")
	}
	if mc["typ"] != "refresh" {
		return RefreshClaims{}, errors.New("wrong token type")
	}

	// 期限チェック（Parseがやるが、追加で厳格化）
	expUnix, ok := mc["exp"].(float64)
	if !ok {
		return RefreshClaims{}, errors.New("no exp")
	}
	exp := time.Unix(int64(expUnix), 0)
	if time.Now().After(exp) {
		return RefreshClaims{}, errors.New("expired")
	}

	userID, _ := mc["sub"].(string)
	jti, _ := mc["jti"].(string)
	if userID == "" || jti == "" {
		return RefreshClaims{}, errors.New("missing sub or jti")
	}

	return RefreshClaims{JTI: jti, UserID: userID, Exp: exp}, nil
}
