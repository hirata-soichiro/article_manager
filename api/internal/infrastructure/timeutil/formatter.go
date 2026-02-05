package timeutil

import (
	"fmt"
	"sync"
	"time"
)

const (
	// 日時表示フォーマット
	DateTimeFormat = "2006-01-02 15:04:05"
	// アプリケーションのタイムゾーン
	TimeZone = "Asia/Tokyo"
)

var (
	jst  *time.Location
	once sync.Once
	err  error
)

// JSTタイムゾーンを取得（初回のみロード）
func GetJST() (*time.Location, error) {
	once.Do(func() {
		jst, err = time.LoadLocation(TimeZone)
	})
	return jst, err
}

// 時刻をJST形式でフォーマット
func FormatInJST(t time.Time) (string, error) {
	location, err := GetJST()
	if err != nil {
		return "", fmt.Errorf("failed to load JST timezone: %w", err)
	}
	return t.In(location).Format(DateTimeFormat), nil
}

// 時刻をJST形式でフォーマット（エラー時はパニック）
func MustFormatInJST(t time.Time) string {
	formatted, err := FormatInJST(t)
	if err != nil {
		panic(err)
	}
	return formatted
}
