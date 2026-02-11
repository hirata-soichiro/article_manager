package external

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"article-manager/internal/domain/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SearchBookのテスト
func TestSearchBook(t *testing.T) {
	t.Run("正常系：書籍を検索できる（ISBN-10あり）", func(t *testing.T) {
		// モックサーバーのセットアップ
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, "/volumes")
			assert.Contains(t, r.URL.Query().Get("q"), "リーダブルコード")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"items": [
					{
						"volumeInfo": {
							"title": "リーダブルコード",
							"authors": ["Dustin Boswell", "Trevor Foucher"],
							"industryIdentifiers": [
								{
									"type": "ISBN_10",
									"identifier": "4873115655"
								},
								{
									"type": "ISBN_13",
									"identifier": "9784873115658"
								}
							]
						}
					}
				]
			}`))
		}))
		defer mockServer.Close()

		config := &GoogleBooksConfig{
			APIKey:        "test-api-key",
			BaseURL:       mockServer.URL,
			Timeout:       10 * time.Second,
			MaxRetries:    3,
			RetryWaitTime: 1 * time.Second,
		}

		client := NewGoogleBooksClient(config)
		book, err := client.SearchBook(context.Background(), "リーダブルコード", "Dustin Boswell")

		require.NoError(t, err)
		assert.NotNil(t, book)
		assert.Equal(t, "リーダブルコード", book.Title)
		assert.Equal(t, "Dustin Boswell, Trevor Foucher", book.Author)
		assert.Equal(t, "4873115655", book.ISBN)
		assert.Equal(t, "9784873115658", book.ISBN13)
		assert.Equal(t, "https://www.amazon.co.jp/dp/4873115655", book.PurchaseLinks.Amazon)
		assert.Contains(t, book.PurchaseLinks.Rakuten, "9784873115658")
	})

	t.Run("正常系：書籍を検索できる（ISBN-10なし、ISBN-13のみ）", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"items": [
					{
						"volumeInfo": {
							"title": "Clean Architecture",
							"authors": ["Robert C. Martin"],
							"industryIdentifiers": [
								{
									"type": "ISBN_13",
									"identifier": "9780134494166"
								}
							]
						}
					}
				]
			}`))
		}))
		defer mockServer.Close()

		config := DefaultGoogleBooksConfig("")
		config.BaseURL = mockServer.URL
		client := NewGoogleBooksClient(config)

		book, err := client.SearchBook(context.Background(), "Clean Architecture", "Robert C. Martin")

		require.NoError(t, err)
		assert.NotNil(t, book)
		assert.Equal(t, "Clean Architecture", book.Title)
		assert.Equal(t, "Robert C. Martin", book.Author)
		assert.Equal(t, "9780134494166", book.ISBN)
		assert.Equal(t, "9780134494166", book.ISBN13)
		assert.Equal(t, "https://www.amazon.co.jp/dp/9780134494166", book.PurchaseLinks.Amazon)
		assert.Contains(t, book.PurchaseLinks.Rakuten, "9780134494166")
	})

	t.Run("正常系：検索結果が0件の場合、基本情報のみ返す", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"items": []
			}`))
		}))
		defer mockServer.Close()

		config := DefaultGoogleBooksConfig("")
		config.BaseURL = mockServer.URL
		client := NewGoogleBooksClient(config)

		book, err := client.SearchBook(context.Background(), "存在しない書籍", "Unknown Author")

		require.NoError(t, err)
		assert.NotNil(t, book)
		assert.Equal(t, "存在しない書籍", book.Title)
		assert.Equal(t, "Unknown Author", book.Author)
		assert.Equal(t, "", book.ISBN)
		assert.Equal(t, "", book.ISBN13)
		assert.Equal(t, "", book.PurchaseLinks.Amazon)
		assert.Equal(t, "", book.PurchaseLinks.Rakuten)
	})

	t.Run("正常系：著者名が複数の場合、カンマ区切りで結合される", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"items": [
					{
						"volumeInfo": {
							"title": "プログラミング言語Go",
							"authors": ["Alan A. A. Donovan", "Brian W. Kernighan"],
							"industryIdentifiers": [
								{
									"type": "ISBN_13",
									"identifier": "9784621300251"
								}
							]
						}
					}
				]
			}`))
		}))
		defer mockServer.Close()

		config := DefaultGoogleBooksConfig("")
		config.BaseURL = mockServer.URL
		client := NewGoogleBooksClient(config)

		book, err := client.SearchBook(context.Background(), "プログラミング言語Go", "")

		require.NoError(t, err)
		assert.Equal(t, "Alan A. A. Donovan, Brian W. Kernighan", book.Author)
	})

	t.Run("正常系：著者名を指定しない場合も検索できる", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			query := r.URL.Query().Get("q")
			// URLデコード済みのクエリが渡される
			assert.Contains(t, query, "Go言語入門")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"items": [
					{
						"volumeInfo": {
							"title": "Go言語入門",
							"authors": ["Test Author"],
							"industryIdentifiers": [
								{
									"type": "ISBN_10",
									"identifier": "1234567890"
								}
							]
						}
					}
				]
			}`))
		}))
		defer mockServer.Close()

		config := DefaultGoogleBooksConfig("")
		config.BaseURL = mockServer.URL
		client := NewGoogleBooksClient(config)

		book, err := client.SearchBook(context.Background(), "Go言語入門", "")

		require.NoError(t, err)
		assert.NotNil(t, book)
		assert.Equal(t, "Go言語入門", book.Title)
	})

	t.Run("異常系：タイトルが空の場合エラー", func(t *testing.T) {
		config := DefaultGoogleBooksConfig("")
		client := NewGoogleBooksClient(config)

		book, err := client.SearchBook(context.Background(), "", "Author")

		require.Error(t, err)
		assert.Nil(t, book)

		bookErr, ok := err.(*service.BookRecommendationError)
		require.True(t, ok)
		assert.Equal(t, service.ErrCodeBooksAPIError, bookErr.Code)
		assert.Contains(t, bookErr.Message, "title is required")
	})

	t.Run("異常系：400 Bad Request", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{
				"error": {
					"code": 400,
					"message": "Invalid query"
				}
			}`))
		}))
		defer mockServer.Close()

		config := DefaultGoogleBooksConfig("")
		config.BaseURL = mockServer.URL
		client := NewGoogleBooksClient(config)

		book, err := client.SearchBook(context.Background(), "Test Book", "")

		require.Error(t, err)
		assert.Nil(t, book)

		bookErr, ok := err.(*service.BookRecommendationError)
		require.True(t, ok)
		assert.Equal(t, service.ErrCodeBooksAPIError, bookErr.Code)
		assert.Contains(t, bookErr.Message, "bad request")
	})

	t.Run("異常系：401 Unauthorized（無効なAPIキー）", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{
				"error": {
					"code": 401,
					"message": "Invalid API key"
				}
			}`))
		}))
		defer mockServer.Close()

		config := DefaultGoogleBooksConfig("invalid-key")
		config.BaseURL = mockServer.URL
		client := NewGoogleBooksClient(config)

		book, err := client.SearchBook(context.Background(), "Test Book", "")

		require.Error(t, err)
		assert.Nil(t, book)

		bookErr, ok := err.(*service.BookRecommendationError)
		require.True(t, ok)
		assert.Equal(t, service.ErrCodeBooksAPIError, bookErr.Code)
		assert.Contains(t, bookErr.Message, "invalid API key")
	})

	t.Run("異常系：429 Too Many Requests（レート制限）", func(t *testing.T) {
		attemptCount := 0
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attemptCount++
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{
				"error": {
					"code": 429,
					"message": "Rate limit exceeded"
				}
			}`))
		}))
		defer mockServer.Close()

		config := DefaultGoogleBooksConfig("")
		config.BaseURL = mockServer.URL
		config.MaxRetries = 1 // 1回リトライ
		config.RetryWaitTime = 10 * time.Millisecond
		client := NewGoogleBooksClient(config)

		book, err := client.SearchBook(context.Background(), "Test Book", "")

		require.Error(t, err)
		assert.Nil(t, book)

		bookErr, ok := err.(*service.BookRecommendationError)
		require.True(t, ok)
		assert.Equal(t, service.ErrCodeBooksAPIError, bookErr.Code)
		// リトライ後のエラーメッセージ
		assert.Contains(t, bookErr.Message, "failed after")
		assert.Equal(t, 2, attemptCount) // 最初の試行 + 1回リトライ
	})

	t.Run("異常系：500 Internal Server Error", func(t *testing.T) {
		attemptCount := 0
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attemptCount++
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{
				"error": {
					"code": 500,
					"message": "Internal server error"
				}
			}`))
		}))
		defer mockServer.Close()

		config := DefaultGoogleBooksConfig("")
		config.BaseURL = mockServer.URL
		config.MaxRetries = 1
		config.RetryWaitTime = 10 * time.Millisecond
		client := NewGoogleBooksClient(config)

		book, err := client.SearchBook(context.Background(), "Test Book", "")

		require.Error(t, err)
		assert.Nil(t, book)

		bookErr, ok := err.(*service.BookRecommendationError)
		require.True(t, ok)
		assert.Equal(t, service.ErrCodeBooksAPIError, bookErr.Code)
		// リトライ後のエラーメッセージ
		assert.Contains(t, bookErr.Message, "failed after")
		assert.Equal(t, 2, attemptCount) // 最初の試行 + 1回リトライ
	})

	t.Run("異常系：不正なJSONレスポンス", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{invalid json`))
		}))
		defer mockServer.Close()

		config := DefaultGoogleBooksConfig("")
		config.BaseURL = mockServer.URL
		client := NewGoogleBooksClient(config)

		book, err := client.SearchBook(context.Background(), "Test Book", "")

		require.Error(t, err)
		assert.Nil(t, book)

		bookErr, ok := err.(*service.BookRecommendationError)
		require.True(t, ok)
		assert.Equal(t, service.ErrCodeBooksAPIError, bookErr.Code)
		assert.Contains(t, bookErr.Message, "failed to parse response")
	})

	t.Run("異常系：コンテキストキャンセル", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(2 * time.Second)
			w.WriteHeader(http.StatusOK)
		}))
		defer mockServer.Close()

		config := DefaultGoogleBooksConfig("")
		config.BaseURL = mockServer.URL
		config.Timeout = 100 * time.Millisecond
		client := NewGoogleBooksClient(config)

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		book, err := client.SearchBook(ctx, "Test Book", "")

		require.Error(t, err)
		assert.Nil(t, book)
	})

	t.Run("正常系：リトライ機能が動作する（500エラー後に成功）", func(t *testing.T) {
		attemptCount := 0
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attemptCount++
			if attemptCount == 1 {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"items": [
					{
						"volumeInfo": {
							"title": "Test Book",
							"authors": ["Test Author"],
							"industryIdentifiers": []
						}
					}
				]
			}`))
		}))
		defer mockServer.Close()

		config := DefaultGoogleBooksConfig("")
		config.BaseURL = mockServer.URL
		config.MaxRetries = 2
		config.RetryWaitTime = 10 * time.Millisecond
		client := NewGoogleBooksClient(config)

		book, err := client.SearchBook(context.Background(), "Test Book", "")

		require.NoError(t, err)
		assert.NotNil(t, book)
		assert.Equal(t, "Test Book", book.Title)
		assert.Equal(t, 2, attemptCount) // 1回目失敗、2回目成功
	})
}
