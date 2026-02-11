package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"article-manager/internal/domain/entity"
	"article-manager/internal/domain/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// モック BookRecommendationRepository
type mockBookRecommendationRepository struct {
	findLatestValidFunc func(ctx context.Context) (*entity.BookRecommendationCache, error)
	saveFunc            func(ctx context.Context, cache *entity.BookRecommendationCache) (*entity.BookRecommendationCache, error)
}

func (m *mockBookRecommendationRepository) FindLatestValid(ctx context.Context) (*entity.BookRecommendationCache, error) {
	return m.findLatestValidFunc(ctx)
}

func (m *mockBookRecommendationRepository) Save(ctx context.Context, cache *entity.BookRecommendationCache) (*entity.BookRecommendationCache, error) {
	return m.saveFunc(ctx, cache)
}

// モック BookRecommendationService
type mockBookRecommendationService struct {
	recommendBooksFunc func(ctx context.Context, articles []*entity.Article) ([]entity.Book, error)
}

func (m *mockBookRecommendationService) RecommendBooks(ctx context.Context, articles []*entity.Article) ([]entity.Book, error) {
	return m.recommendBooksFunc(ctx, articles)
}

// GetBookRecommendationsのテスト
func TestGetBookRecommendations(t *testing.T) {
	t.Run("正常系：キャッシュが有効な場合、キャッシュを返す", func(t *testing.T) {
		// 有効なキャッシュを準備
		now := time.Now()
		expectedCache := &entity.BookRecommendationCache{
			ID: 1,
			Books: []entity.Book{
				{
					Title:  "リーダブルコード",
					Author: "Dustin Boswell",
					ISBN:   "9784873115658",
					PurchaseLinks: entity.PurchaseLinks{
						Amazon:  "https://www.amazon.co.jp/dp/4873115655",
						Rakuten: "https://books.rakuten.co.jp/rb/11753651/",
					},
				},
			},
			GeneratedAt: now.Add(-1 * time.Hour),
			ExpiresAt:   now.Add(23 * time.Hour),
		}

		mockBookRecommendationRepo := &mockBookRecommendationRepository{
			findLatestValidFunc: func(ctx context.Context) (*entity.BookRecommendationCache, error) {
				return expectedCache, nil
			},
		}

		mockArticleRepo := &mockArticleRepository{}
		mockService := &mockBookRecommendationService{}

		usecase := NewBookRecommendationUsecase(mockArticleRepo, mockBookRecommendationRepo, mockService)

		result, err := usecase.GetBookRecommendations(context.Background())

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(1), result.ID)
		assert.Equal(t, 1, len(result.Books))
		assert.Equal(t, "リーダブルコード", result.Books[0].Title)
		assert.True(t, result.IsValid())
	})

	t.Run("正常系：キャッシュが存在しない場合、新規に生成する", func(t *testing.T) {
		// キャッシュなし
		mockBookRecommendationRepo := &mockBookRecommendationRepository{
			findLatestValidFunc: func(ctx context.Context) (*entity.BookRecommendationCache, error) {
				return nil, errors.New("not found")
			},
			saveFunc: func(ctx context.Context, cache *entity.BookRecommendationCache) (*entity.BookRecommendationCache, error) {
				cache.ID = 1
				return cache, nil
			},
		}

		// 記事リストを返す
		mockArticleRepo := &mockArticleRepository{
			findAllFunc: func(ctx context.Context) ([]*entity.Article, error) {
				return []*entity.Article{
					{
						ID:      1,
						Title:   "Go言語入門",
						URL:     "https://example.com/go",
						Summary: "Go言語の基本を学ぶ",
						Tags:    []string{"Go"},
					},
				}, nil
			},
		}

		// AIが書籍を推薦
		mockService := &mockBookRecommendationService{
			recommendBooksFunc: func(ctx context.Context, articles []*entity.Article) ([]entity.Book, error) {
				return []entity.Book{
					{
						Title:  "プログラミング言語Go",
						Author: "Alan A. A. Donovan",
						ISBN:   "9784621300251",
						PurchaseLinks: entity.PurchaseLinks{
							Amazon: "https://www.amazon.co.jp/dp/4621300253",
						},
					},
				}, nil
			},
		}

		usecase := NewBookRecommendationUsecase(mockArticleRepo, mockBookRecommendationRepo, mockService)

		result, err := usecase.GetBookRecommendations(context.Background())

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(1), result.ID)
		assert.Equal(t, 1, len(result.Books))
		assert.Equal(t, "プログラミング言語Go", result.Books[0].Title)
		assert.True(t, result.IsValid())
	})

	t.Run("正常系：キャッシュの有効期限が切れている場合、新規に生成する", func(t *testing.T) {
		// 有効期限切れのキャッシュ
		expiredCache := &entity.BookRecommendationCache{
			ID: 1,
			Books: []entity.Book{
				{
					Title:  "古い書籍",
					Author: "Author Name",
					ISBN:   "1234567890",
					PurchaseLinks: entity.PurchaseLinks{
						Amazon: "https://www.amazon.co.jp/dp/1234567890",
					},
				},
			},
			GeneratedAt: time.Now().Add(-25 * time.Hour),
			ExpiresAt:   time.Now().Add(-1 * time.Hour),
		}

		mockBookRecommendationRepo := &mockBookRecommendationRepository{
			findLatestValidFunc: func(ctx context.Context) (*entity.BookRecommendationCache, error) {
				return expiredCache, nil
			},
			saveFunc: func(ctx context.Context, cache *entity.BookRecommendationCache) (*entity.BookRecommendationCache, error) {
				cache.ID = 2
				return cache, nil
			},
		}

		mockArticleRepo := &mockArticleRepository{
			findAllFunc: func(ctx context.Context) ([]*entity.Article, error) {
				return []*entity.Article{
					{
						ID:      1,
						Title:   "新しい記事",
						URL:     "https://example.com/new",
						Summary: "新しい記事の内容",
						Tags:    []string{"Tech"},
					},
				}, nil
			},
		}

		mockService := &mockBookRecommendationService{
			recommendBooksFunc: func(ctx context.Context, articles []*entity.Article) ([]entity.Book, error) {
				return []entity.Book{
					{
						Title:  "新しい書籍",
						Author: "New Author",
						ISBN:   "9876543210",
						PurchaseLinks: entity.PurchaseLinks{
							Amazon: "https://www.amazon.co.jp/dp/9876543210",
						},
					},
				}, nil
			},
		}

		usecase := NewBookRecommendationUsecase(mockArticleRepo, mockBookRecommendationRepo, mockService)

		result, err := usecase.GetBookRecommendations(context.Background())

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(2), result.ID)
		assert.Equal(t, 1, len(result.Books))
		assert.Equal(t, "新しい書籍", result.Books[0].Title)
		assert.True(t, result.IsValid())
	})

	t.Run("正常系：記事が0件の場合、空の推薦を返す", func(t *testing.T) {
		mockBookRecommendationRepo := &mockBookRecommendationRepository{
			findLatestValidFunc: func(ctx context.Context) (*entity.BookRecommendationCache, error) {
				return nil, errors.New("not found")
			},
		}

		mockArticleRepo := &mockArticleRepository{
			findAllFunc: func(ctx context.Context) ([]*entity.Article, error) {
				return []*entity.Article{}, nil
			},
		}

		mockService := &mockBookRecommendationService{}

		usecase := NewBookRecommendationUsecase(mockArticleRepo, mockBookRecommendationRepo, mockService)

		result, err := usecase.GetBookRecommendations(context.Background())

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(0), result.ID)
		assert.Equal(t, 0, len(result.Books))
	})

	t.Run("異常系：記事取得時にエラーが発生", func(t *testing.T) {
		mockBookRecommendationRepo := &mockBookRecommendationRepository{
			findLatestValidFunc: func(ctx context.Context) (*entity.BookRecommendationCache, error) {
				return nil, errors.New("cache not found")
			},
		}

		mockArticleRepo := &mockArticleRepository{
			findAllFunc: func(ctx context.Context) ([]*entity.Article, error) {
				return nil, errors.New("database connection error")
			},
		}

		mockService := &mockBookRecommendationService{}

		usecase := NewBookRecommendationUsecase(mockArticleRepo, mockBookRecommendationRepo, mockService)

		result, err := usecase.GetBookRecommendations(context.Background())

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database connection error")
	})

	t.Run("異常系：AI書籍推薦時にエラーが発生", func(t *testing.T) {
		mockBookRecommendationRepo := &mockBookRecommendationRepository{
			findLatestValidFunc: func(ctx context.Context) (*entity.BookRecommendationCache, error) {
				return nil, errors.New("cache not found")
			},
		}

		mockArticleRepo := &mockArticleRepository{
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

		mockService := &mockBookRecommendationService{
			recommendBooksFunc: func(ctx context.Context, articles []*entity.Article) ([]entity.Book, error) {
				return nil, &service.BookRecommendationError{
					Code:    service.ErrCodeAIError,
					Message: "AI API error",
					Err:     errors.New("rate limit exceeded"),
				}
			},
		}

		usecase := NewBookRecommendationUsecase(mockArticleRepo, mockBookRecommendationRepo, mockService)

		result, err := usecase.GetBookRecommendations(context.Background())

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "AI API error")
	})

	t.Run("異常系：キャッシュ保存時にエラーが発生", func(t *testing.T) {
		mockBookRecommendationRepo := &mockBookRecommendationRepository{
			findLatestValidFunc: func(ctx context.Context) (*entity.BookRecommendationCache, error) {
				return nil, errors.New("cache not found")
			},
			saveFunc: func(ctx context.Context, cache *entity.BookRecommendationCache) (*entity.BookRecommendationCache, error) {
				return nil, errors.New("database write error")
			},
		}

		mockArticleRepo := &mockArticleRepository{
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

		mockService := &mockBookRecommendationService{
			recommendBooksFunc: func(ctx context.Context, articles []*entity.Article) ([]entity.Book, error) {
				return []entity.Book{
					{
						Title:  "書籍タイトル",
						Author: "著者名",
						ISBN:   "1234567890",
						PurchaseLinks: entity.PurchaseLinks{
							Amazon: "https://www.amazon.co.jp/dp/1234567890",
						},
					},
				}, nil
			},
		}

		usecase := NewBookRecommendationUsecase(mockArticleRepo, mockBookRecommendationRepo, mockService)

		result, err := usecase.GetBookRecommendations(context.Background())

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database write error")
	})

	t.Run("正常系：複数の書籍を推薦する", func(t *testing.T) {
		mockBookRecommendationRepo := &mockBookRecommendationRepository{
			findLatestValidFunc: func(ctx context.Context) (*entity.BookRecommendationCache, error) {
				return nil, errors.New("cache not found")
			},
			saveFunc: func(ctx context.Context, cache *entity.BookRecommendationCache) (*entity.BookRecommendationCache, error) {
				cache.ID = 1
				return cache, nil
			},
		}

		mockArticleRepo := &mockArticleRepository{
			findAllFunc: func(ctx context.Context) ([]*entity.Article, error) {
				return []*entity.Article{
					{
						ID:      1,
						Title:   "Go言語の並行処理",
						URL:     "https://example.com/go-concurrent",
						Summary: "Goの並行処理パターン",
						Tags:    []string{"Go", "Concurrency"},
					},
					{
						ID:      2,
						Title:   "クリーンアーキテクチャ入門",
						URL:     "https://example.com/clean-arch",
						Summary: "アーキテクチャの設計",
						Tags:    []string{"Architecture"},
					},
				}, nil
			},
		}

		mockService := &mockBookRecommendationService{
			recommendBooksFunc: func(ctx context.Context, articles []*entity.Article) ([]entity.Book, error) {
				return []entity.Book{
					{
						Title:  "Goプログラミング実践入門",
						Author: "Shiju Varghese",
						ISBN:   "9784295004073",
						PurchaseLinks: entity.PurchaseLinks{
							Amazon: "https://www.amazon.co.jp/dp/4295004073",
						},
					},
					{
						Title:  "Clean Architecture",
						Author: "Robert C. Martin",
						ISBN:   "9780134494166",
						PurchaseLinks: entity.PurchaseLinks{
							Amazon: "https://www.amazon.co.jp/dp/0134494164",
						},
					},
					{
						Title:  "リーダブルコード",
						Author: "Dustin Boswell",
						ISBN:   "9784873115658",
						PurchaseLinks: entity.PurchaseLinks{
							Amazon:  "https://www.amazon.co.jp/dp/4873115655",
							Rakuten: "https://books.rakuten.co.jp/rb/11753651/",
						},
					},
				}, nil
			},
		}

		usecase := NewBookRecommendationUsecase(mockArticleRepo, mockBookRecommendationRepo, mockService)

		result, err := usecase.GetBookRecommendations(context.Background())

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(1), result.ID)
		assert.Equal(t, 3, len(result.Books))
		assert.Equal(t, "Goプログラミング実践入門", result.Books[0].Title)
		assert.Equal(t, "Clean Architecture", result.Books[1].Title)
		assert.Equal(t, "リーダブルコード", result.Books[2].Title)
		assert.True(t, result.IsValid())
	})

	t.Run("異常系：エンティティ作成時のバリデーションエラー", func(t *testing.T) {
		mockBookRecommendationRepo := &mockBookRecommendationRepository{
			findLatestValidFunc: func(ctx context.Context) (*entity.BookRecommendationCache, error) {
				return nil, errors.New("cache not found")
			},
		}

		mockArticleRepo := &mockArticleRepository{
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

		// タイトルが空の書籍を返す（バリデーションエラー）
		mockService := &mockBookRecommendationService{
			recommendBooksFunc: func(ctx context.Context, articles []*entity.Article) ([]entity.Book, error) {
				return []entity.Book{
					{
						Title:  "",
						Author: "著者名",
						ISBN:   "1234567890",
						PurchaseLinks: entity.PurchaseLinks{
							Amazon: "https://www.amazon.co.jp/dp/1234567890",
						},
					},
				}, nil
			},
		}

		usecase := NewBookRecommendationUsecase(mockArticleRepo, mockBookRecommendationRepo, mockService)

		result, err := usecase.GetBookRecommendations(context.Background())

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "book_recommendation_cache")
	})
}
