package entity

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"
)

// タグエンティティ
type Tag struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// 新しいタグの作成
func NewTag(name string) (*Tag, error) {
	if err := validateTagName(name); err != nil {
		return nil, err
	}

	now := time.Now()
	tag := &Tag{
		ID:        0,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return tag, nil
}

// タグ名のバリデーション
func validateTagName(name string) error {
	// 空文字チェック
	if name == "" {
		return errors.New("name is required")
	}

	// 空白のみチェック
	if strings.TrimSpace(name) == "" {
		return errors.New("name cannot be only whitespace")
	}

	// 文字数チェック
	if utf8.RuneCountInString(name) > 50 {
		return errors.New("name must be 50 characters or less")
	}

	return nil
}
