package models

import (
	"time"
)

// CurrentTime は現在時刻を返します
func CurrentTime() time.Time {
	return time.Now()
}

// FormatTime は時間を指定されたフォーマットで文字列に変換します
func FormatTime(t time.Time, format string) string {
	if format == "" {
		format = time.RFC3339
	}
	return t.Format(format)
}

// ParseTime は文字列を時間に変換します
func ParseTime(timeStr, format string) (time.Time, error) {
	if format == "" {
		format = time.RFC3339
	}
	return time.Parse(format, timeStr)
}

// ParseDate は日付文字列を時間に変換します
func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

// FormatDate は時間を日付文字列に変換します
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// IsZeroTime は時間がゼロ値かどうかを判定します
func IsZeroTime(t time.Time) bool {
	return t.IsZero()
}
