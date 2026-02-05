package timeutil

import (
	"testing"
	"time"
)

func TestGetJST(t *testing.T) {
	t.Run("JSTタイムゾーンが正常にロードできる", func(t *testing.T) {
		location, err := GetJST()
		if err != nil {
			t.Fatalf("GetJST() returned error: %v", err)
		}
		if location == nil {
			t.Fatal("GetJST() returned nil location")
		}
		if location.String() != TimeZone {
			t.Errorf("Expected timezone %s, got %s", TimeZone, location.String())
		}
	})

	t.Run("複数回呼び出しても同じインスタンスを返す", func(t *testing.T) {
		location1, err1 := GetJST()
		location2, err2 := GetJST()

		if err1 != nil || err2 != nil {
			t.Fatalf("GetJST() returned error: %v, %v", err1, err2)
		}

		if location1 != location2 {
			t.Error("GetJST() returned different instances")
		}
	})
}

func TestFormatInJST(t *testing.T) {
	t.Run("UTC時刻をJSTに正しく変換する", func(t *testing.T) {
		utcTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		formatted, err := FormatInJST(utcTime)
		if err != nil {
			t.Fatalf("FormatInJST() returned error: %v", err)
		}

		expected := "2024-01-01 09:00:00"
		if formatted != expected {
			t.Errorf("Expected %s, got %s", expected, formatted)
		}
	})

	t.Run("異なる時刻でも正しく変換する", func(t *testing.T) {
		utcTime := time.Date(2024, 6, 15, 14, 30, 45, 0, time.UTC)
		formatted, err := FormatInJST(utcTime)
		if err != nil {
			t.Fatalf("FormatInJST() returned error: %v", err)
		}

		expected := "2024-06-15 23:30:45"
		if formatted != expected {
			t.Errorf("Expected %s, got %s", expected, formatted)
		}
	})
}

func TestMustFormatInJST(t *testing.T) {
	t.Run("正常にフォーマットできる", func(t *testing.T) {
		utcTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		formatted := MustFormatInJST(utcTime)

		expected := "2024-01-01 09:00:00"
		if formatted != expected {
			t.Errorf("Expected %s, got %s", expected, formatted)
		}
	})

	t.Run("有効なタイムゾーンでパニックしない", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("MustFormatInJST() panicked: %v", r)
			}
		}()

		utcTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		_ = MustFormatInJST(utcTime)
	})
}

func TestDateTimeFormat(t *testing.T) {
	t.Run("正しいフォーマットを使用する", func(t *testing.T) {
		if DateTimeFormat != "2006-01-02 15:04:05" {
			t.Errorf("DateTimeFormat is incorrect: %s", DateTimeFormat)
		}
	})
}

func TestTimeZone(t *testing.T) {
	t.Run("正しいタイムゾーンを使用する", func(t *testing.T) {
		if TimeZone != "Asia/Tokyo" {
			t.Errorf("TimeZone is incorrect: %s", TimeZone)
		}
	})
}
