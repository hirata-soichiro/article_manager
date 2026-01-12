package entity

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewArticle(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		url     string
		summary string
		tags    []string
		memo    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "正常系：全ての必須項目が有効",
			title:   "Go言語入門",
			url:     "https://example.com/go-introduction",
			summary: "Go言語の基本的な文法と特徴について解説します。",
			tags:    []string{"Go", "プログラミング"},
			memo:    "後で読む",
			wantErr: false,
		},
		{
			name:    "正常系：タグが空配列",
			title:   "Go言語入門",
			url:     "https://example.com/go-introduction",
			summary: "Go言語の基本的な文法と特徴について解説します。",
			tags:    []string{},
			memo:    "後で読む",
			wantErr: false,
		},
		{
			name:    "正常系：タグがnil",
			title:   "Go言語入門",
			url:     "https://example.com/go-introduction",
			summary: "Go言語の基本的な文法と特徴について解説します。",
			tags:    nil,
			memo:    "後で読む",
			wantErr: false,
		},
		{
			name:    "正常系：メモが空文字列",
			title:   "Go言語入門",
			url:     "https://example.com/go-introduction",
			summary: "Go言語の基本的な文法と特徴について解説します。",
			tags:    []string{"Go"},
			memo:    "",
			wantErr: false,
		},
		{
			name:    "正常系：タイトルが255文字ちょうど",
			title:   strings.Repeat("あ", 255),
			url:     "https://example.com/long-title",
			summary: "長いタイトルの記事です。",
			tags:    []string{"テスト"},
			memo:    "",
			wantErr: false,
		},
		{
			name:    "正常系：要約が1000文字ちょうど",
			title:   "Go言語入門",
			url:     "https://example.com/go-introduction",
			summary: strings.Repeat("あ", 1000),
			tags:    []string{"Go"},
			memo:    "",
			wantErr: false,
		},
		{
			name:    "正常系：タグが50文字ちょうど",
			title:   "Go言語入門",
			url:     "https://example.com/go-introduction",
			summary: "Go言語の基本的な文法と特徴について解説します。",
			tags:    []string{strings.Repeat("あ", 50)},
			memo:    "",
			wantErr: false,
		},
		{
			name:    "正常系：httpのURL",
			title:   "Go言語入門",
			url:     "http://example.com/go-introduction",
			summary: "Go言語の基本的な文法と特徴について解説します。",
			tags:    []string{"Go"},
			memo:    "",
			wantErr: false,
		},
		{
			name:    "正常系：複雑なURL",
			title:   "Go言語入門",
			url:     "https://example.com/path/to/article?id=123&lang=ja#section1",
			summary: "Go言語の基本的な文法と特徴について解説します。",
			tags:    []string{"Go"},
			memo:    "",
			wantErr: false,
		},
		{
			name:    "異常系：タイトルが空文字列",
			title:   "",
			url:     "https://example.com/go-introduction",
			summary: "Go言語の基本的な文法と特徴について解説します。",
			tags:    []string{"Go"},
			memo:    "",
			wantErr: true,
			errMsg:  "title is required",
		},
		{
			name:    "異常系：タイトルが255文字超過",
			title:   strings.Repeat("あ", 256),
			url:     "https://example.com/long-title",
			summary: "長いタイトルの記事です。",
			tags:    []string{"テスト"},
			memo:    "",
			wantErr: true,
			errMsg:  "title must be 255 characters or less",
		},
		{
			name:    "異常系：URLが空文字列",
			title:   "Go言語入門",
			url:     "",
			summary: "Go言語の基本的な文法と特徴について解説します。",
			tags:    []string{"Go"},
			memo:    "",
			wantErr: true,
			errMsg:  "url is required",
		},
		{
			name:    "異常系：URLが無効な形式（httpまたはhttpsで始まらない）",
			title:   "Go言語入門",
			url:     "ftp://example.com/go-introduction",
			summary: "Go言語の基本的な文法と特徴について解説します。",
			tags:    []string{"Go"},
			memo:    "",
			wantErr: true,
			errMsg:  "url must start with http:// or https://",
		},
		{
			name:    "異常系：URLが無効な形式（スキームなし）",
			title:   "Go言語入門",
			url:     "example.com/go-introduction",
			summary: "Go言語の基本的な文法と特徴について解説します。",
			tags:    []string{"Go"},
			memo:    "",
			wantErr: true,
			errMsg:  "url must start with http:// or https://",
		},
		{
			name:    "異常系：要約が空文字列",
			title:   "Go言語入門",
			url:     "https://example.com/go-introduction",
			summary: "",
			tags:    []string{"Go"},
			memo:    "",
			wantErr: true,
			errMsg:  "summary is required",
		},
		{
			name:    "異常系：要約が1000文字超過",
			title:   "Go言語入門",
			url:     "https://example.com/go-introduction",
			summary: strings.Repeat("あ", 1001),
			tags:    []string{"Go"},
			memo:    "",
			wantErr: true,
			errMsg:  "summary must be 1000 characters or less",
		},
		{
			name:    "異常系：タグの1つが50文字超過",
			title:   "Go言語入門",
			url:     "https://example.com/go-introduction",
			summary: "Go言語の基本的な文法と特徴について解説します。",
			tags:    []string{"Go", strings.Repeat("あ", 51)},
			memo:    "",
			wantErr: true,
			errMsg:  "each tag must be 50 characters or less",
		},
		{
			name:    "異常系：タグに空文字列が含まれる",
			title:   "Go言語入門",
			url:     "https://example.com/go-introduction",
			summary: "Go言語の基本的な文法と特徴について解説します。",
			tags:    []string{"Go", ""},
			memo:    "",
			wantErr: true,
			errMsg:  "tag cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beforeTime := time.Now()
			article, err := NewArticle(tt.title, tt.url, tt.summary, tt.tags, tt.memo)
			afterTime := time.Now()

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, article)
			} else {
				require.NoError(t, err)
				require.NotNil(t, article)
				assert.Equal(t, tt.title, article.Title)
				assert.Equal(t, tt.url, article.URL)
				assert.Equal(t, tt.summary, article.Summary)

				if tt.tags == nil {
					assert.Empty(t, article.Tags)
				} else {
					assert.Equal(t, tt.tags, article.Tags)
				}

				assert.Equal(t, tt.memo, article.Memo)

				assert.Equal(t, int64(0), article.ID)
				assert.True(t, article.CreatedAt.After(beforeTime) || article.CreatedAt.Equal(beforeTime))
				assert.True(t, article.CreatedAt.Before(afterTime) || article.CreatedAt.Equal(afterTime))
				assert.Equal(t, article.CreatedAt, article.UpdatedAt)
			}
		})
	}
}

func TestArticle_Update(t *testing.T) {
	tests := []struct {
		name       string
		initial    *Article
		newTitle   string
		newURL     string
		newSummary string
		newTags    []string
		newMemo    string
		wantErr    bool
		errMsg     string
	}{
		{
			name: "正常系：全てのフィールドを更新",
			initial: &Article{
				ID:        1,
				Title:     "Go言語入門",
				URL:       "https://example.com/go-introduction",
				Summary:   "Go言語の基本的な文法と特徴について解説します。",
				Tags:      []string{"Go"},
				Memo:      "後で読む",
				CreatedAt: time.Now().Add(-24 * time.Hour),
				UpdatedAt: time.Now().Add(-24 * time.Hour),
			},
			newTitle:   "Go言語完全ガイド",
			newURL:     "https://example.com/go-complete-guide",
			newSummary: "Go言語の全てを網羅的に解説する完全ガイドです。",
			newTags:    []string{"Go", "完全ガイド"},
			newMemo:    "重要な記事",
			wantErr:    false,
		},
		{
			name: "正常系：メモを空文字列に更新",
			initial: &Article{
				ID:        1,
				Title:     "Go言語入門",
				URL:       "https://example.com/go-introduction",
				Summary:   "Go言語の基本的な文法と特徴について解説します。",
				Tags:      []string{"Go"},
				Memo:      "後で読む",
				CreatedAt: time.Now().Add(-24 * time.Hour),
				UpdatedAt: time.Now().Add(-24 * time.Hour),
			},
			newTitle:   "Go言語入門",
			newURL:     "https://example.com/go-introduction",
			newSummary: "Go言語の基本的な文法と特徴について解説します。",
			newTags:    []string{"Go"},
			newMemo:    "",
			wantErr:    false,
		},
		{
			name: "正常系：タグを空配列に更新",
			initial: &Article{
				ID:        1,
				Title:     "Go言語入門",
				URL:       "https://example.com/go-introduction",
				Summary:   "Go言語の基本的な文法と特徴について解説します。",
				Tags:      []string{"Go", "プログラミング"},
				Memo:      "",
				CreatedAt: time.Now().Add(-24 * time.Hour),
				UpdatedAt: time.Now().Add(-24 * time.Hour),
			},
			newTitle:   "Go言語入門",
			newURL:     "https://example.com/go-introduction",
			newSummary: "Go言語の基本的な文法と特徴について解説します。",
			newTags:    []string{},
			newMemo:    "",
			wantErr:    false,
		},
		{
			name: "正常系：タグをnilに更新",
			initial: &Article{
				ID:        1,
				Title:     "Go言語入門",
				URL:       "https://example.com/go-introduction",
				Summary:   "Go言語の基本的な文法と特徴について解説します。",
				Tags:      []string{"Go", "プログラミング"},
				Memo:      "",
				CreatedAt: time.Now().Add(-24 * time.Hour),
				UpdatedAt: time.Now().Add(-24 * time.Hour),
			},
			newTitle:   "Go言語入門",
			newURL:     "https://example.com/go-introduction",
			newSummary: "Go言語の基本的な文法と特徴について解説します。",
			newTags:    nil,
			newMemo:    "",
			wantErr:    false,
		},
		{
			name: "正常系：タイトルが255文字ちょうど",
			initial: &Article{
				ID:        1,
				Title:     "Go言語入門",
				URL:       "https://example.com/go-introduction",
				Summary:   "Go言語の基本的な文法と特徴について解説します。",
				Tags:      []string{"Go"},
				Memo:      "",
				CreatedAt: time.Now().Add(-24 * time.Hour),
				UpdatedAt: time.Now().Add(-24 * time.Hour),
			},
			newTitle:   strings.Repeat("あ", 255),
			newURL:     "https://example.com/long-title",
			newSummary: "長いタイトルの記事です。",
			newTags:    []string{"テスト"},
			newMemo:    "",
			wantErr:    false,
		},
		{
			name: "正常系：要約が1000文字ちょうど",
			initial: &Article{
				ID:        1,
				Title:     "Go言語入門",
				URL:       "https://example.com/go-introduction",
				Summary:   "Go言語の基本的な文法と特徴について解説します。",
				Tags:      []string{"Go"},
				Memo:      "",
				CreatedAt: time.Now().Add(-24 * time.Hour),
				UpdatedAt: time.Now().Add(-24 * time.Hour),
			},
			newTitle:   "Go言語入門",
			newURL:     "https://example.com/go-introduction",
			newSummary: strings.Repeat("あ", 1000),
			newTags:    []string{"Go"},
			newMemo:    "",
			wantErr:    false,
		},
		{
			name: "正常系：タグが50文字ちょうど",
			initial: &Article{
				ID:        1,
				Title:     "Go言語入門",
				URL:       "https://example.com/go-introduction",
				Summary:   "Go言語の基本的な文法と特徴について解説します。",
				Tags:      []string{"Go"},
				Memo:      "",
				CreatedAt: time.Now().Add(-24 * time.Hour),
				UpdatedAt: time.Now().Add(-24 * time.Hour),
			},
			newTitle:   "Go言語入門",
			newURL:     "https://example.com/go-introduction",
			newSummary: "Go言語の基本的な文法と特徴について解説します。",
			newTags:    []string{strings.Repeat("あ", 50)},
			newMemo:    "",
			wantErr:    false,
		},
		{
			name: "正常系：httpのURL",
			initial: &Article{
				ID:        1,
				Title:     "Go言語入門",
				URL:       "https://example.com/go-introduction",
				Summary:   "Go言語の基本的な文法と特徴について解説します。",
				Tags:      []string{"Go"},
				Memo:      "",
				CreatedAt: time.Now().Add(-24 * time.Hour),
				UpdatedAt: time.Now().Add(-24 * time.Hour),
			},
			newTitle:   "Go言語入門",
			newURL:     "http://example.com/go-introduction",
			newSummary: "Go言語の基本的な文法と特徴について解説します。",
			newTags:    []string{"Go"},
			newMemo:    "",
			wantErr:    false,
		},
		{
			name: "正常系：複雑なURL",
			initial: &Article{
				ID:        1,
				Title:     "Go言語入門",
				URL:       "https://example.com/go-introduction",
				Summary:   "Go言語の基本的な文法と特徴について解説します。",
				Tags:      []string{"Go"},
				Memo:      "",
				CreatedAt: time.Now().Add(-24 * time.Hour),
				UpdatedAt: time.Now().Add(-24 * time.Hour),
			},
			newTitle:   "Go言語入門",
			newURL:     "https://example.com/path/to/article?id=123&lang=ja#section1",
			newSummary: "Go言語の基本的な文法と特徴について解説します。",
			newTags:    []string{"Go"},
			newMemo:    "",
			wantErr:    false,
		},
		{
			name: "異常系：タイトルが空文字列",
			initial: &Article{
				ID:        1,
				Title:     "Go言語入門",
				URL:       "https://example.com/go-introduction",
				Summary:   "Go言語の基本的な文法と特徴について解説します。",
				Tags:      []string{"Go"},
				Memo:      "",
				CreatedAt: time.Now().Add(-24 * time.Hour),
				UpdatedAt: time.Now().Add(-24 * time.Hour),
			},
			newTitle:   "",
			newURL:     "https://example.com/go-introduction",
			newSummary: "Go言語の基本的な文法と特徴について解説します。",
			newTags:    []string{"Go"},
			newMemo:    "",
			wantErr:    true,
			errMsg:     "title is required",
		},
		{
			name: "異常系：タイトルが256文字超過",
			initial: &Article{
				ID:        1,
				Title:     "Go言語入門",
				URL:       "https://example.com/go-introduction",
				Summary:   "Go言語の基本的な文法と特徴について解説します。",
				Tags:      []string{"Go"},
				Memo:      "",
				CreatedAt: time.Now().Add(-24 * time.Hour),
				UpdatedAt: time.Now().Add(-24 * time.Hour),
			},
			newTitle:   strings.Repeat("あ", 256),
			newURL:     "https://example.com/long-title",
			newSummary: "長いタイトルの記事です。",
			newTags:    []string{"テスト"},
			newMemo:    "",
			wantErr:    true,
			errMsg:     "title must be 255 characters or less",
		},
		{
			name: "異常系：URLが空文字列",
			initial: &Article{
				ID:        1,
				Title:     "Go言語入門",
				URL:       "https://example.com/go-introduction",
				Summary:   "Go言語の基本的な文法と特徴について解説します。",
				Tags:      []string{"Go"},
				Memo:      "",
				CreatedAt: time.Now().Add(-24 * time.Hour),
				UpdatedAt: time.Now().Add(-24 * time.Hour),
			},
			newTitle:   "Go言語入門",
			newURL:     "",
			newSummary: "Go言語の基本的な文法と特徴について解説します。",
			newTags:    []string{"Go"},
			newMemo:    "",
			wantErr:    true,
			errMsg:     "url is required",
		},
		{
			name: "異常系：URLが無効な形式（httpまたはhttpsで始まらない）",
			initial: &Article{
				ID:        1,
				Title:     "Go言語入門",
				URL:       "https://example.com/go-introduction",
				Summary:   "Go言語の基本的な文法と特徴について解説します。",
				Tags:      []string{"Go"},
				Memo:      "",
				CreatedAt: time.Now().Add(-24 * time.Hour),
				UpdatedAt: time.Now().Add(-24 * time.Hour),
			},
			newTitle:   "Go言語入門",
			newURL:     "ftp://example.com/go-introduction",
			newSummary: "Go言語の基本的な文法と特徴について解説します。",
			newTags:    []string{"Go"},
			newMemo:    "",
			wantErr:    true,
			errMsg:     "url must start with http:// or https://",
		},
		{
			name: "異常系：URLが無効な形式（スキームなし）",
			initial: &Article{
				ID:        1,
				Title:     "Go言語入門",
				URL:       "https://example.com/go-introduction",
				Summary:   "Go言語の基本的な文法と特徴について解説します。",
				Tags:      []string{"Go"},
				Memo:      "",
				CreatedAt: time.Now().Add(-24 * time.Hour),
				UpdatedAt: time.Now().Add(-24 * time.Hour),
			},
			newTitle:   "Go言語入門",
			newURL:     "example.com/go-introduction",
			newSummary: "Go言語の基本的な文法と特徴について解説します。",
			newTags:    []string{"Go"},
			newMemo:    "",
			wantErr:    true,
			errMsg:     "url must start with http:// or https://",
		},
		{
			name: "異常系：要約が空文字列",
			initial: &Article{
				ID:        1,
				Title:     "Go言語入門",
				URL:       "https://example.com/go-introduction",
				Summary:   "Go言語の基本的な文法と特徴について解説します。",
				Tags:      []string{"Go"},
				Memo:      "",
				CreatedAt: time.Now().Add(-24 * time.Hour),
				UpdatedAt: time.Now().Add(-24 * time.Hour),
			},
			newTitle:   "Go言語入門",
			newURL:     "https://example.com/go-introduction",
			newSummary: "",
			newTags:    []string{"Go"},
			newMemo:    "",
			wantErr:    true,
			errMsg:     "summary is required",
		},
		{
			name: "異常系：要約が1001文字超過",
			initial: &Article{
				ID:        1,
				Title:     "Go言語入門",
				URL:       "https://example.com/go-introduction",
				Summary:   "Go言語の基本的な文法と特徴について解説します。",
				Tags:      []string{"Go"},
				Memo:      "",
				CreatedAt: time.Now().Add(-24 * time.Hour),
				UpdatedAt: time.Now().Add(-24 * time.Hour),
			},
			newTitle:   "Go言語入門",
			newURL:     "https://example.com/go-introduction",
			newSummary: strings.Repeat("あ", 1001),
			newTags:    []string{"Go"},
			newMemo:    "",
			wantErr:    true,
			errMsg:     "summary must be 1000 characters or less",
		},
		{
			name: "異常系：タグの1つが51文字超過",
			initial: &Article{
				ID:        1,
				Title:     "Go言語入門",
				URL:       "https://example.com/go-introduction",
				Summary:   "Go言語の基本的な文法と特徴について解説します。",
				Tags:      []string{"Go"},
				Memo:      "",
				CreatedAt: time.Now().Add(-24 * time.Hour),
				UpdatedAt: time.Now().Add(-24 * time.Hour),
			},
			newTitle:   "Go言語入門",
			newURL:     "https://example.com/go-introduction",
			newSummary: "Go言語の基本的な文法と特徴について解説します。",
			newTags:    []string{"Go", strings.Repeat("あ", 51)},
			newMemo:    "",
			wantErr:    true,
			errMsg:     "each tag must be 50 characters or less",
		},
		{
			name: "異常系：タグに空文字列が含まれる",
			initial: &Article{
				ID:        1,
				Title:     "Go言語入門",
				URL:       "https://example.com/go-introduction",
				Summary:   "Go言語の基本的な文法と特徴について解説します。",
				Tags:      []string{"Go"},
				Memo:      "",
				CreatedAt: time.Now().Add(-24 * time.Hour),
				UpdatedAt: time.Now().Add(-24 * time.Hour),
			},
			newTitle:   "Go言語入門",
			newURL:     "https://example.com/go-introduction",
			newSummary: "Go言語の基本的な文法と特徴について解説します。",
			newTags:    []string{"Go", ""},
			newMemo:    "",
			wantErr:    true,
			errMsg:     "tag cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beforeUpdate := time.Now()
			originalCreatedAt := tt.initial.CreatedAt
			originalUpdatedAt := tt.initial.UpdatedAt

			err := tt.initial.Update(tt.newTitle, tt.newURL, tt.newSummary, tt.newTags, tt.newMemo)
			afterUpdate := time.Now()

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Equal(t, originalUpdatedAt, tt.initial.UpdatedAt)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.newTitle, tt.initial.Title)
				assert.Equal(t, tt.newURL, tt.initial.URL)
				assert.Equal(t, tt.newSummary, tt.initial.Summary)

				if tt.newTags == nil {
					assert.Empty(t, tt.initial.Tags)
				} else {
					assert.Equal(t, tt.newTags, tt.initial.Tags)
				}

				assert.Equal(t, tt.newMemo, tt.initial.Memo)

				assert.Equal(t, originalCreatedAt, tt.initial.CreatedAt)
				assert.True(t, tt.initial.UpdatedAt.After(beforeUpdate) || tt.initial.UpdatedAt.Equal(beforeUpdate))
				assert.True(t, tt.initial.UpdatedAt.Before(afterUpdate) || tt.initial.UpdatedAt.Equal(afterUpdate))
			}
		})
	}
}
