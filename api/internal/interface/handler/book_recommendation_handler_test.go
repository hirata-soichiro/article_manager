package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"article-manager/internal/domain/entity"
	"article-manager/internal/domain/repository"
	"article-manager/internal/domain/service"
	"article-manager/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// モック ArticleRepository
type mockArticleRepositoryForHandler struct {
	findAllFunc func(ctx context.Context) ([]*entity.Article, error)
}

func (m *mockArticleRepositoryForHandler) Create(ctx context.Context, article *entity.Article) (*entity.Article, error) {
	return nil, nil
}

func (m *mockArticleRepositoryForHandler) FindByID(ctx context.Context, id int64) (*entity.Article, error) {
	return nil, nil
}

func (m *mockArticleRepositoryForHandler) FindAll(ctx context.Context) ([]*entity.Article, error) {
	if m.findAllFunc != nil {
		return m.findAllFunc(ctx)
	}
	return []*entity.Article{}, nil
}

func (m *mockArticleRepositoryForHandler) Update(ctx context.Context, article *entity.Article) (*entity.Article, error) {
	return nil, nil
}

func (m *mockArticleRepositoryForHandler) Delete(ctx context.Context, id int64) error {
	return nil
}

func (m *mockArticleRepositoryForHandler) Search(ctx context.Context, keyword string) ([]*entity.Article, error) {
	return nil, nil
}

// モック BookRecommendationRepository
type mockBookRecommendationRepositoryForHandler struct {
	findLatestValidFunc func(ctx context.Context) (*entity.BookRecommendationCache, error)
	saveFunc            func(ctx context.Context, cache *entity.BookRecommendationCache) (*entity.BookRecommendationCache, error)
}

func (m *mockBookRecommendationRepositoryForHandler) FindLatestValid(ctx context.Context) (*entity.BookRecommendationCache, error) {
	if m.findLatestValidFunc != nil {
		return m.findLatestValidFunc(ctx)
	}
	return nil, errors.New("not found")
}

func (m *mockBookRecommendationRepositoryForHandler) Save(ctx context.Context, cache *entity.BookRecommendationCache) (*entity.BookRecommendationCache, error) {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, cache)
	}
	cache.ID = 1
	return cache, nil
}

// モック BookRecommendationService
type mockBookRecommendationServiceForHandler struct {
	recommendBooksFunc func(ctx context.Context, articles []*entity.Article) ([]entity.Book, error)
}

func (m *mockBookRecommendationServiceForHandler) RecommendBooks(ctx context.Context, articles []*entity.Article) ([]entity.Book, error) {
	if m.recommendBooksFunc != nil {
		return m.recommendBooksFunc(ctx, articles)
	}
	return []entity.Book{}, nil
}

// テスト用のハンドラをセットアップ
func setupBookRecommendationHandler(
	articleRepo repository.ArticleRepository,
	bookRecommendationRepo repository.BookRecommendationRepository,
	bookRecommendationService service.BookRecommendationService,
) *BookRecommendationHandler {
	uc := usecase.NewBookRecommendationUsecase(articleRepo, bookRecommendationRepo, bookRecommendationService)
	return NewBookRecommendationHandler(uc)
}

// GetBookRecommendationsのテスト
func TestGetBookRecommendationsHandler(t *testing.T) {
	t.Run("正常系：書籍推薦を取得できる（キャッシュあり）", func(t *testing.T) {
		// モックのセットアップ
		now := time.Now()
		expectedCache := &entity.BookRecommendationCache{
			ID: 1,
			Books: []entity.Book{
				{
					Title: "リーダブルコード",
					PurchaseLinks: entity.PurchaseLinks{
						Amazon:  "https://www.amazon.co.jp/dp/4873115655",
						Rakuten: "https://books.rakuten.co.jp/rb/11753651/",
					},
				},
				{
					Title: "Clean Code",
					PurchaseLinks: entity.PurchaseLinks{
						Amazon: "https://www.amazon.co.jp/dp/0132350884",
					},
				},
			},
			GeneratedAt: now.Add(-1 * time.Hour),
			ExpiresAt:   now.Add(23 * time.Hour),
		}

		mockBookRecommendationRepo := &mockBookRecommendationRepositoryForHandler{
			findLatestValidFunc: func(ctx context.Context) (*entity.BookRecommendationCache, error) {
				return expectedCache, nil
			},
		}

		handler := setupBookRecommendationHandler(
			&mockArticleRepositoryForHandler{},
			mockBookRecommendationRepo,
			&mockBookRecommendationServiceForHandler{},
		)

		// リクエスト作成
		req := httptest.NewRequest(http.MethodGet, "/api/book-recommendations", nil)
		rec := httptest.NewRecorder()

		// ハンドラ実行
		handler.GetBookRecommendations(rec, req)

		// レスポンス検証
		require.Equal(t, http.StatusOK, rec.Code)

		var response BookRecommendationResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		// 書籍数の確認
		assert.Equal(t, 2, len(response.Books))

		// 1冊目の確認
		assert.Equal(t, "リーダブルコード", response.Books[0].Title)
		assert.Equal(t, "https://www.amazon.co.jp/dp/4873115655", response.Books[0].PurchaseLinks.Amazon)
		assert.Equal(t, "https://books.rakuten.co.jp/rb/11753651/", response.Books[0].PurchaseLinks.Rakuten)

		// 2冊目の確認
		assert.Equal(t, "Clean Code", response.Books[1].Title)
		assert.Equal(t, "https://www.amazon.co.jp/dp/0132350884", response.Books[1].PurchaseLinks.Amazon)

		// キャッシュフラグとタイムスタンプの確認
		assert.True(t, response.Cached)
		assert.NotNil(t, response.GeneratedAt)
		assert.NotNil(t, response.ExpiresAt)
	})

	t.Run("正常系：記事が0件の場合、空の配列を返す", func(t *testing.T) {
		mockArticleRepo := &mockArticleRepositoryForHandler{
			findAllFunc: func(ctx context.Context) ([]*entity.Article, error) {
				return []*entity.Article{}, nil
			},
		}

		mockBookRecommendationRepo := &mockBookRecommendationRepositoryForHandler{
			findLatestValidFunc: func(ctx context.Context) (*entity.BookRecommendationCache, error) {
				return nil, errors.New("not found")
			},
		}

		handler := setupBookRecommendationHandler(
			mockArticleRepo,
			mockBookRecommendationRepo,
			&mockBookRecommendationServiceForHandler{},
		)

		req := httptest.NewRequest(http.MethodGet, "/api/book-recommendations", nil)
		rec := httptest.NewRecorder()

		handler.GetBookRecommendations(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		var response BookRecommendationResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, 0, len(response.Books))
		assert.False(t, response.Cached)
		assert.Nil(t, response.GeneratedAt)
		assert.Nil(t, response.ExpiresAt)
	})

	t.Run("正常系：新規生成された推薦を返す（キャッシュなし）", func(t *testing.T) {
		mockArticleRepo := &mockArticleRepositoryForHandler{
			findAllFunc: func(ctx context.Context) ([]*entity.Article, error) {
				return []*entity.Article{
					{
						ID:      1,
						Title:   "Go言語入門",
						URL:     "https://example.com/go",
						Summary: "Go言語の基本",
						Tags:    []string{"Go"},
					},
				}, nil
			},
		}

		mockBookRecommendationRepo := &mockBookRecommendationRepositoryForHandler{
			findLatestValidFunc: func(ctx context.Context) (*entity.BookRecommendationCache, error) {
				return nil, errors.New("not found")
			},
		}

		mockService := &mockBookRecommendationServiceForHandler{
			recommendBooksFunc: func(ctx context.Context, articles []*entity.Article) ([]entity.Book, error) {
				return []entity.Book{
					{
						Title: "プログラミング言語Go",
						PurchaseLinks: entity.PurchaseLinks{
							Amazon: "https://www.amazon.co.jp/dp/4621300253",
						},
					},
				}, nil
			},
		}

		handler := setupBookRecommendationHandler(mockArticleRepo, mockBookRecommendationRepo, mockService)

		req := httptest.NewRequest(http.MethodGet, "/api/book-recommendations", nil)
		rec := httptest.NewRecorder()

		handler.GetBookRecommendations(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		var response BookRecommendationResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, 1, len(response.Books))
		assert.Equal(t, "プログラミング言語Go", response.Books[0].Title)
		assert.True(t, response.Cached)
		assert.NotNil(t, response.GeneratedAt)
		assert.NotNil(t, response.ExpiresAt)
	})

	t.Run("正常系：複数の書籍推薦を返す", func(t *testing.T) {
		now := time.Now()
		expectedCache := &entity.BookRecommendationCache{
			ID: 1,
			Books: []entity.Book{
				{
					Title: "書籍1",
					PurchaseLinks: entity.PurchaseLinks{
						Amazon: "https://www.amazon.co.jp/dp/1234567890",
					},
				},
				{
					Title: "書籍2",
					PurchaseLinks: entity.PurchaseLinks{
						Amazon:  "https://www.amazon.co.jp/dp/0987654321",
						Rakuten: "https://books.rakuten.co.jp/rb/test/",
					},
				},
				{
					Title: "書籍3",
					PurchaseLinks: entity.PurchaseLinks{
						Amazon: "https://www.amazon.co.jp/dp/1111111111",
					},
				},
			},
			GeneratedAt: now.Add(-1 * time.Hour),
			ExpiresAt:   now.Add(23 * time.Hour),
		}

		mockBookRecommendationRepo := &mockBookRecommendationRepositoryForHandler{
			findLatestValidFunc: func(ctx context.Context) (*entity.BookRecommendationCache, error) {
				return expectedCache, nil
			},
		}

		handler := setupBookRecommendationHandler(
			&mockArticleRepositoryForHandler{},
			mockBookRecommendationRepo,
			&mockBookRecommendationServiceForHandler{},
		)

		req := httptest.NewRequest(http.MethodGet, "/api/book-recommendations", nil)
		rec := httptest.NewRecorder()

		handler.GetBookRecommendations(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		var response BookRecommendationResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, 3, len(response.Books))
		assert.Equal(t, "書籍1", response.Books[0].Title)
		assert.Equal(t, "書籍2", response.Books[1].Title)
		assert.Equal(t, "書籍3", response.Books[2].Title)
	})

	t.Run("正常系：ISBNなしの書籍も正しく返す", func(t *testing.T) {
		now := time.Now()
		expectedCache := &entity.BookRecommendationCache{
			ID: 1,
			Books: []entity.Book{
				{
					Title: "ISBNなし書籍",
					PurchaseLinks: entity.PurchaseLinks{
						Amazon: "",
					},
				},
			},
			GeneratedAt: now,
			ExpiresAt:   now.Add(24 * time.Hour),
		}

		mockBookRecommendationRepo := &mockBookRecommendationRepositoryForHandler{
			findLatestValidFunc: func(ctx context.Context) (*entity.BookRecommendationCache, error) {
				return expectedCache, nil
			},
		}

		handler := setupBookRecommendationHandler(
			&mockArticleRepositoryForHandler{},
			mockBookRecommendationRepo,
			&mockBookRecommendationServiceForHandler{},
		)

		req := httptest.NewRequest(http.MethodGet, "/api/book-recommendations", nil)
		rec := httptest.NewRecorder()

		handler.GetBookRecommendations(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		var response BookRecommendationResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, 1, len(response.Books))
		assert.Equal(t, "ISBNなし書籍", response.Books[0].Title)
		assert.Equal(t, "", response.Books[0].PurchaseLinks.Amazon)
	})

	t.Run("異常系：ユースケースがエラーを返す（データベースエラー）", func(t *testing.T) {
		mockArticleRepo := &mockArticleRepositoryForHandler{
			findAllFunc: func(ctx context.Context) ([]*entity.Article, error) {
				return nil, errors.New("database connection error")
			},
		}

		mockBookRecommendationRepo := &mockBookRecommendationRepositoryForHandler{
			findLatestValidFunc: func(ctx context.Context) (*entity.BookRecommendationCache, error) {
				return nil, errors.New("not found")
			},
		}

		handler := setupBookRecommendationHandler(
			mockArticleRepo,
			mockBookRecommendationRepo,
			&mockBookRecommendationServiceForHandler{},
		)

		req := httptest.NewRequest(http.MethodGet, "/api/book-recommendations", nil)
		rec := httptest.NewRecorder()

		handler.GetBookRecommendations(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "error")
	})

	t.Run("異常系：ユースケースがAI APIエラーを返す", func(t *testing.T) {
		mockArticleRepo := &mockArticleRepositoryForHandler{
			findAllFunc: func(ctx context.Context) ([]*entity.Article, error) {
				return []*entity.Article{
					{
						ID:      1,
						Title:   "記事",
						URL:     "https://example.com",
						Summary: "要約",
						Tags:    []string{"Go"},
					},
				}, nil
			},
		}

		mockBookRecommendationRepo := &mockBookRecommendationRepositoryForHandler{
			findLatestValidFunc: func(ctx context.Context) (*entity.BookRecommendationCache, error) {
				return nil, errors.New("not found")
			},
		}

		mockService := &mockBookRecommendationServiceForHandler{
			recommendBooksFunc: func(ctx context.Context, articles []*entity.Article) ([]entity.Book, error) {
				return nil, &service.BookRecommendationError{
					Code:    service.ErrCodeAIError,
					Message: "AI API error",
					Err:     errors.New("rate limit exceeded"),
				}
			},
		}

		handler := setupBookRecommendationHandler(mockArticleRepo, mockBookRecommendationRepo, mockService)

		req := httptest.NewRequest(http.MethodGet, "/api/book-recommendations", nil)
		rec := httptest.NewRecorder()

		handler.GetBookRecommendations(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "error")
	})

	t.Run("異常系：キャッシュ保存時にエラーが発生", func(t *testing.T) {
		mockArticleRepo := &mockArticleRepositoryForHandler{
			findAllFunc: func(ctx context.Context) ([]*entity.Article, error) {
				return []*entity.Article{
					{
						ID:      1,
						Title:   "記事",
						URL:     "https://example.com",
						Summary: "要約",
						Tags:    []string{"Go"},
					},
				}, nil
			},
		}

		mockBookRecommendationRepo := &mockBookRecommendationRepositoryForHandler{
			findLatestValidFunc: func(ctx context.Context) (*entity.BookRecommendationCache, error) {
				return nil, errors.New("not found")
			},
			saveFunc: func(ctx context.Context, cache *entity.BookRecommendationCache) (*entity.BookRecommendationCache, error) {
				return nil, errors.New("database write error")
			},
		}

		mockService := &mockBookRecommendationServiceForHandler{
			recommendBooksFunc: func(ctx context.Context, articles []*entity.Article) ([]entity.Book, error) {
				return []entity.Book{
					{
						Title: "書籍タイトル",
						PurchaseLinks: entity.PurchaseLinks{
							Amazon: "https://www.amazon.co.jp/dp/1234567890",
						},
					},
				}, nil
			},
		}

		handler := setupBookRecommendationHandler(mockArticleRepo, mockBookRecommendationRepo, mockService)

		req := httptest.NewRequest(http.MethodGet, "/api/book-recommendations", nil)
		rec := httptest.NewRecorder()

		handler.GetBookRecommendations(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "error")
	})
}
