package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"article_manager/internal/domain/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// モックリポジトリ
type mockArticleRepository struct {
	createFunc   func(ctx context.Context, article *entity.Article) (*entity.Article, error)
	findByIDFunc func(ctx context.Context, id int64) (*entity.Article, error)
	findAllFunc  func(ctx context.Context) ([]*entity.Article, error)
	updateFunc   func(ctx context.Context, article *entity.Article) (*entity.Article, error)
	deleteFunc   func(ctx context.Context, id int64) error
}

func (m *mockArticleRepository) Create(ctx context.Context, article *entity.Article) (*entity.Article, error) {
	return m.createFunc(ctx, article)
}

func (m *mockArticleRepository) FindByID(ctx context.Context, id int64) (*entity.Article, error) {
	return m.findByIDFunc(ctx, id)
}

func (m *mockArticleRepository) FindAll(ctx context.Context) ([]*entity.Article, error) {
	return m.findAllFunc(ctx)
}

func (m *mockArticleRepository) Update(ctx context.Context, article *entity.Article) (*entity.Article, error) {
	return m.updateFunc(ctx, article)
}

func (m *mockArticleRepository) Delete(ctx context.Context, id int64) error {
	return m.deleteFunc(ctx, id)
}

// CreateArticleのテスト
func TestCreateArticle(t *testing.T) {
	t.Run("正常系：記事を作成できる", func(t *testing.T) {
		// モックリポジトリ準備
		mockRepo := &mockArticleRepository{
			createFunc: func(ctx context.Context, article *entity.Article) (*entity.Article, error) {
				article.ID = 1
				return article, nil
			},
		}

		// ユースケース作成
		usecase := NewArticleUsecase(mockRepo)

		// テスト実行
		result, err := usecase.CreateArticle(
			context.Background(),
			"テスト記事",
			"https://example.com",
			"これはテスト記事です",
			[]string{"Go", "テスト"},
			"テストメモ",
		)

		// 検証
		require.NoError(t, err)
		assert.Equal(t, int64(1), result.ID)
		assert.Equal(t, "テスト記事", result.Title)
		assert.Equal(t, "https://example.com", result.URL)
		assert.Equal(t, "これはテスト記事です", result.Description)
		assert.ElementsMatch(t, []string{"Go", "テスト"}, result.Tags)
		assert.Equal(t, "テストメモ", result.Memo)
	})

	t.Run("異常系：タイトルが空の場合エラー", func(t *testing.T) {
		mockRepo := &mockArticleRepository{}
		usecase := NewArticleUsecase(mockRepo)

		result, err := usecase.CreateArticle(
			context.Background(),
			"",
			"https://example.com",
			"これはテスト記事です",
			[]string{"Go"},
			"",
		)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "title is required")
	})

	t.Run("異常系：URLが空の場合エラー", func(t *testing.T) {
		mockRepo := &mockArticleRepository{}
		usecase := NewArticleUsecase(mockRepo)

		result, err := usecase.CreateArticle(
			context.Background(),
			"テスト記事",
			"",
			"これはテスト記事です",
			[]string{"Go"},
			"",
		)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "url is required")
	})

	t.Run("異常系：URLが不正な形式の場合エラー", func(t *testing.T) {
		mockRepo := &mockArticleRepository{}
		usecase := NewArticleUsecase(mockRepo)

		result, err := usecase.CreateArticle(
			context.Background(),
			"テスト記事",
			"invalid-url",
			"これはテスト記事です",
			[]string{"Go"},
			"",
		)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "url must start with http:// or https://")
	})

	t.Run("異常系：要約が空の場合エラー", func(t *testing.T) {
		mockRepo := &mockArticleRepository{}
		usecase := NewArticleUsecase(mockRepo)

		result, err := usecase.CreateArticle(
			context.Background(),
			"テスト記事",
			"https://example.com",
			"",
			[]string{"Go"},
			"",
		)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "summary is required")
	})

	t.Run("異常系：リポジトリがエラーを返す", func(t *testing.T) {
		mockRepo := &mockArticleRepository{
			createFunc: func(ctx context.Context, article *entity.Article) (*entity.Article, error) {
				return nil, errors.New("database error")
			},
		}
		usecase := NewArticleUsecase(mockRepo)

		result, err := usecase.CreateArticle(
			context.Background(),
			"テスト記事",
			"https://example.com",
			"これはテスト記事です",
			[]string{"Go"},
			"",
		)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database error")
	})
}

// GetArticleByIDのテスト
func TestGetArticleByID(t *testing.T) {
	t.Run("正常系：IDで記事を取得できる", func(t *testing.T) {
		expected := &entity.Article{
			ID:        1,
			Title:     "テスト記事",
			URL:       "https://example.com",
			Summary:   "これはテスト記事です",
			Tags:      []string{"Go", "テスト"},
			Memo:      "テストメモ",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRepo := &mockArticleRepository{
			findByIDFunc: func(ctx context.Context, id int64) (*entity.Article, error) {
				return expected, nil
			},
		}

		usecase := NewArticleUsecase(mockRepo)
		result, err := usecase.GetArticleByID(context.Background(), 1)

		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("異常系：記事が見つからない", func(t *testing.T) {
		mockRepo := &mockArticleRepository{
			findByIDFunc: func(ctx context.Context, id int64) (*entity.Article, error) {
				return nil, errors.New("article not found")
			},
		}

		usecase := NewArticleUsecase(mockRepo)
		result, err := usecase.GetArticleByID(context.Background(), 999)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "article not found")
	})

	t.Run("異常系：不正なID（0以下）", func(t *testing.T) {
		mockRepo := &mockArticleRepository{}
		usecase := NewArticleUsecase(mockRepo)

		result, err := usecase.GetArticleByID(context.Background(), 0)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid id")
	})
}

// GetAllArticlesのテスト
func TestGetAllArticles(t *testing.T) {
	t.Run("正常系：全記事を取得できる", func(t *testing.T) {
		expected := []*entity.Article{
			{
				ID:      1,
				Title:   "記事1",
				URL:     "https://example.com/1",
				Summary: "記事1の要約",
				Tags:    []string{"Go"},
				Memo:    "",
			},
			{
				ID:      2,
				Title:   "記事2",
				URL:     "https://example.com/2",
				Summary: "記事2の要約",
				Tags:    []string{"Next.js"},
				Memo:    "",
			},
		}

		mockRepo := &mockArticleRepository{
			findAllFunc: func(ctx context.Context) ([]*entity.Article, error) {
				return expected, nil
			},
		}

		usecase := NewArticleUsecase(mockRepo)
		result, err := usecase.GetAllArticles(context.Background())

		require.NoError(t, err)
		assert.Equal(t, 2, len(result))
		assert.Equal(t, expected, result)
	})

	t.Run("正常系：記事が0件の場合", func(t *testing.T) {
		mockRepo := &mockArticleRepository{
			findAllFunc: func(ctx context.Context) ([]*entity.Article, error) {
				return []*entity.Article{}, nil
			},
		}

		usecase := NewArticleUsecase(mockRepo)
		result, err := usecase.GetAllArticles(context.Background())

		require.NoError(t, err)
		assert.Equal(t, 0, len(result))
		assert.NotNil(t, result)
	})

	t.Run("異常系：リポジトリがエラーを返す", func(t *testing.T) {
		mockRepo := &mockArticleRepository{
			findAllFunc: func(ctx context.Context) ([]*entity.Article, error) {
				return nil, errors.New("database error")
			},
		}

		usecase := NewArticleUsecase(mockRepo)
		result, err := usecase.GetAllArticles(context.Background())

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database error")
	})
}

// UpdateArticleのテスト
func TestUpdateArticle(t *testing.T) {
	t.Run("正常系：記事を更新できる", func(t *testing.T) {
		existingArticle := &entity.Article{
			ID:        1,
			Title:     "旧タイトル",
			URL:       "https://example.com/old",
			Summary:   "旧要約",
			Tags:      []string{"Old"},
			Memo:      "旧メモ",
			CreatedAt: time.Now().Add(-24 * time.Hour),
			UpdatedAt: time.Now().Add(-24 * time.Hour),
		}

		mockRepo := &mockArticleRepository{
			findByIDFunc: func(ctx context.Context, id int64) (*entity.Article, error) {
				return existingArticle, nil
			},
			updateFunc: func(ctx context.Context, article *entity.Article) (*entity.Article, error) {
				return article, nil
			},
		}

		usecase := NewArticleUsecase(mockRepo)
		result, err := usecase.UpdateArticle(
			context.Background(),
			1,
			"新タイトル",
			"https://example.com/new",
			"新要約",
			[]string{"New", "Go"},
			"新メモ",
		)

		require.NoError(t, err)
		assert.Equal(t, "新タイトル", result.Title)
		assert.Equal(t, "https://example.com/new", result.URL)
		assert.Equal(t, "新要約", result.Summary)
		assert.Equal(t, []string{"New", "Go"}, result.Tags)
		assert.Equal(t, "新メモ", result.Memo)
	})

	t.Run("異常系：記事が存在しない", func(t *testing.T) {
		mockRepo := &mockArticleRepository{
			findByIDFunc: func(ctx context.Context, id int64) (*entity.Article, error) {
				return nil, errors.New("article not found")
			},
		}

		usecase := NewArticleUsecase(mockRepo)
		result, err := usecase.UpdateArticle(
			context.Background(),
			999,
			"新タイトル",
			"https://example.com/new",
			"新要約",
			[]string{"New"},
			"",
		)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "article not found")
	})

	t.Run("異常系：更新データのバリデーションエラー", func(t *testing.T) {
		existingArticle := &entity.Article{
			ID:        1,
			Title:     "旧タイトル",
			URL:       "https://example.com/old",
			Summary:   "旧要約",
			Tags:      []string{"Old"},
			Memo:      "",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRepo := &mockArticleRepository{
			findByIDFunc: func(ctx context.Context, id int64) (*entity.Article, error) {
				return existingArticle, nil
			},
		}

		usecase := NewArticleUsecase(mockRepo)
		result, err := usecase.UpdateArticle(
			context.Background(),
			1,
			"", // 空のタイトル
			"https://example.com/new",
			"新要約",
			[]string{"New"},
			"",
		)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "title is required")
	})

	t.Run("異常系：リポジトリの更新でエラー", func(t *testing.T) {
		existingArticle := &entity.Article{
			ID:        1,
			Title:     "旧タイトル",
			URL:       "https://example.com/old",
			Summary:   "旧要約",
			Tags:      []string{"Old"},
			Memo:      "",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRepo := &mockArticleRepository{
			findByIDFunc: func(ctx context.Context, id int64) (*entity.Article, error) {
				return existingArticle, nil
			},
			updateFunc: func(ctx context.Context, article *entity.Article) (*entity.Article, error) {
				return nil, errors.New("database error")
			},
		}

		usecase := NewArticleUsecase(mockRepo)
		result, err := usecase.UpdateArticle(
			context.Background(),
			1,
			"新タイトル",
			"https://example.com/new",
			"新要約",
			[]string{"New"},
			"",
		)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database error")
	})
}

// DeleteArticleのテスト
func TestDeleteArticle(t *testing.T) {
	t.Run("正常系：記事を削除できる", func(t *testing.T) {
		mockRepo := &mockArticleRepository{
			deleteFunc: func(ctx context.Context, id int64) error {
				return nil
			},
		}

		usecase := NewArticleUsecase(mockRepo)
		err := usecase.DeleteArticle(context.Background(), 1)

		require.NoError(t, err)
	})

	t.Run("異常系：不正なID（0以下）", func(t *testing.T) {
		mockRepo := &mockArticleRepository{}
		usecase := NewArticleUsecase(mockRepo)

		err := usecase.DeleteArticle(context.Background(), 0)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid id")
	})

	t.Run("異常系：記事が存在しない", func(t *testing.T) {
		mockRepo := &mockArticleRepository{
			deleteFunc: func(ctx context.Context, id int64) error {
				return errors.New("article not found")
			},
		}

		usecase := NewArticleUsecase(mockRepo)
		err := usecase.DeleteArticle(context.Background(), 999)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "article not found")
	})

	t.Run("異常系：リポジトリがエラーを返す", func(t *testing.T) {
		mockRepo := &mockArticleRepository{
			deleteFunc: func(ctx context.Context, id int64) error {
				return errors.New("database error")
			},
		}

		usecase := NewArticleUsecase(mockRepo)
		err := usecase.DeleteArticle(context.Background(), 1)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
	})
}
