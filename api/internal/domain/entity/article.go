package entity

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"
)

// 記事エンティティ
type Article struct {
	ID        int64
	Title     string
	URL       string
	Summary   string
	Tags      []string
	Memo      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// 新しい記事の作成
func NewArticle(title, url, summary string, tags []string, memo string) (*Article, error) {
	if err := validateTitle(title); err != nil {
		return nil, err
	}
	if err := validateURL(url); err != nil {
		return nil, err
	}
	if err := validateSummary(summary); err != nil {
		return nil, err
	}
	if err := validateTags(tags); err != nil {
		return nil, err
	}

	if tags == nil {
		tags = []string{}
	}

	now := time.Now()
	article := &Article{
		ID:        0,
		Title:     title,
		URL:       url,
		Summary:   summary,
		Tags:      tags,
		Memo:      memo,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return article, nil
}

// 記事の内容を更新
func (a *Article) Update(title, url, summary string, tags []string, memo string) error {
	if err := validateTitle(title); err != nil {
		return err
	}
	if err := validateURL(url); err != nil {
		return err
	}
	if err := validateSummary(summary); err != nil {
		return err
	}
	if err := validateTags(tags); err != nil {
		return err
	}

	if tags == nil {
		tags = []string{}
	}

	a.Title = title
	a.URL = url
	a.Summary = summary
	a.Tags = tags
	a.Memo = memo
	a.UpdatedAt = time.Now()

	return nil
}

func validateTitle(title string) error {
	if title == "" {
		return errors.New("title is required")
	}
	if utf8.RuneCountInString(title) > 255 {
		return errors.New("title must be 255 characters or less")
	}
	return nil
}

func validateURL(url string) error {
	if url == "" {
		return errors.New("url is required")
	}
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return errors.New("url must start with http:// or https://")
	}
	return nil
}

func validateSummary(summary string) error {
	if summary == "" {
		return errors.New("summary is required")
	}
	if utf8.RuneCountInString(summary) > 1000 {
		return errors.New("summary must be 1000 characters or less")
	}
	return nil
}

func validateTags(tags []string) error {
	for _, tag := range tags {
		if tag == "" {
			return errors.New("tag cannot be empty")
		}
		if utf8.RuneCountInString(tag) > 50 {
			return errors.New("each tag must be 50 characters or less")
		}
	}
	return nil
}
