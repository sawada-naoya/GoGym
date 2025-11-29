package util

import "time"

const (
	// JSTLocation は日本標準時のタイムゾーン
	JSTLocation = "Asia/Tokyo"
)

var jstLoc *time.Location

func init() {
	var err error
	jstLoc, err = time.LoadLocation(JSTLocation)
	if err != nil {
		// フォールバック: FixedZone を使用
		jstLoc = time.FixedZone("JST", 9*60*60)
	}
}

// ToJST converts UTC time to JST
// DBから取得したUTC時刻をJSTに変換してフロントエンドに返す
func ToJST(t time.Time) time.Time {
	if t.IsZero() {
		return t
	}
	return t.In(jstLoc)
}

// ToUTC converts JST time to UTC
// フロントエンドから受け取ったJST時刻をUTCに変換してDBに保存する
func ToUTC(t time.Time) time.Time {
	if t.IsZero() {
		return t
	}
	return t.UTC()
}

// FormatJSTDate formats a time as YYYY-MM-DD in JST
func FormatJSTDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return ToJST(t).Format("2006-01-02")
}

// FormatJSTDateTime formats a time as YYYY-MM-DD HH:mm:ss in JST
func FormatJSTDateTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return ToJST(t).Format("2006-01-02 15:04:05")
}

// FormatJSTTime formats a time as HH:mm in JST
func FormatJSTTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return ToJST(t).Format("15:04")
}
