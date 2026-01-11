package util

import (
	"errors"
	"time"
)

const (
	JSTLocation = "Asia/Tokyo"
	DateLayout  = "2006-01-02"
)

var jstLoc *time.Location

func init() {
	var err error
	jstLoc, err = time.LoadLocation(JSTLocation)
	if err != nil {
		jstLoc = time.FixedZone("JST", 9*60*60)
	}
}

//
// ===== 基本変換 =====
//

// UTC → JST
func ToJST(t time.Time) time.Time {
	if t.IsZero() {
		return t
	}
	return t.In(jstLoc)
}

// JST/Local → UTC
func ToUTC(t time.Time) time.Time {
	if t.IsZero() {
		return t
	}
	return t.UTC()
}

//
// ===== フォーマット（表示専用） =====
//

// JSTの YYYY-MM-DD
func FormatJSTDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return ToJST(t).Format(DateLayout)
}

// JSTの YYYY-MM-DD HH:mm:ss
func FormatJSTDateTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return ToJST(t).Format("2006-01-02 15:04:05")
}

// JSTの HH:mm
func FormatJSTTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return ToJST(t).Format("15:04")
}

//
// ===== パース（入力専用） =====
//

// YYYY-MM-DD（JST）→ time.Time（JST 00:00）
func ParseJSTDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, errors.New("date is empty")
	}
	t, err := time.ParseInLocation(DateLayout, dateStr, jstLoc)
	if err != nil {
		return time.Time{}, err
	}
	return normalizeJSTDate(t), nil
}

// 空なら今日のJST日付を返す
func ParseJSTDateOrToday(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return TodayJST(), nil
	}
	return ParseJSTDate(dateStr)
}

//
// ===== 正規化（DB向け） =====
//

// JSTの「日付」を 00:00:00 に正規化
func normalizeJSTDate(t time.Time) time.Time {
	y, m, d := t.In(jstLoc).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, jstLoc)
}

// DB保存用：JST日付 → UTC
func NormalizeDateForDB(t time.Time) time.Time {
	if t.IsZero() {
		return t
	}
	return ToUTC(normalizeJSTDate(t))
}

//
// ===== 補助 =====
//

// 今日のJST日付（00:00）
func TodayJST() time.Time {
	now := time.Now().In(jstLoc)
	y, m, d := now.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, jstLoc)
}

// 現在時刻（JST）
func NowJST() time.Time {
	return time.Now().In(jstLoc)
}
