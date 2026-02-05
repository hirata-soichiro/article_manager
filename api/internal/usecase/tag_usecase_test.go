package usecase

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"article-manager/internal/domain/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// モックリポジトリ
type mockTagRepository struct {
	createFunc     func(ctx context.Context, tag *entity.Tag) (*entity.Tag, error)
	findByIDFunc   func(ctx context.Context, id int64) (*entity.Tag, error)
	findByNameFunc func(ctx context.Context, name string) (*entity.Tag, error)
	findAllFunc    func(ctx context.Context) ([]*entity.Tag, error)
	updateFunc     func(ctx context.Context, tag *entity.Tag) (*entity.Tag, error)
	deleteFunc     func(ctx context.Context, id int64) error
}

func (m *mockTagRepository) Create(ctx context.Context, tag *entity.Tag) (*entity.Tag, error) {
	return m.createFunc(ctx, tag)
}

func (m *mockTagRepository) FindByID(ctx context.Context, id int64) (*entity.Tag, error) {
	return m.findByIDFunc(ctx, id)
}

func (m *mockTagRepository) FindByName(ctx context.Context, name string) (*entity.Tag, error) {
	return m.findByNameFunc(ctx, name)
}

func (m *mockTagRepository) FindAll(ctx context.Context) ([]*entity.Tag, error) {
	return m.findAllFunc(ctx)
}

func (m *mockTagRepository) Update(ctx context.Context, tag *entity.Tag) (*entity.Tag, error) {
	return m.updateFunc(ctx, tag)
}

func (m *mockTagRepository) Delete(ctx context.Context, id int64) error {
	return m.deleteFunc(ctx, id)
}

// CreateTagのテスト
func TestCreateTag(t *testing.T) {
	t.Run("正常系：タグを作成できる", func(t *testing.T) {
		// モックリポジトリ準備
		mockRepo := &mockTagRepository{
			createFunc: func(ctx context.Context, tag *entity.Tag) (*entity.Tag, error) {
				tag.ID = 1
				return tag, nil
			},
		}

		// ユースケース作成
		usecase := NewTagUsecase(mockRepo)

		// テスト実行
		result, err := usecase.CreateTag(context.Background(), "Go")

		// 検証
		require.NoError(t, err)
		assert.Equal(t, int64(1), result.ID)
		assert.Equal(t, "Go", result.Name)
	})

	t.Run("異常系：タグ名が空の場合エラー", func(t *testing.T) {
		mockRepo := &mockTagRepository{}
		usecase := NewTagUsecase(mockRepo)

		result, err := usecase.CreateTag(context.Background(), "")

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "name is required")
	})

	t.Run("異常系：タグ名が空白のみの場合エラー", func(t *testing.T) {
		mockRepo := &mockTagRepository{}
		usecase := NewTagUsecase(mockRepo)

		result, err := usecase.CreateTag(context.Background(), "   ")

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "name cannot be only whitespace")
	})

	t.Run("異常系：タグ名が50文字を超える場合エラー", func(t *testing.T) {
		mockRepo := &mockTagRepository{}
		usecase := NewTagUsecase(mockRepo)

		longName := strings.Repeat("あ", 51)
		result, err := usecase.CreateTag(context.Background(), longName)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "name must be 50 characters or less")
	})

	t.Run("異常系：リポジトリがエラーを返す", func(t *testing.T) {
		mockRepo := &mockTagRepository{
			createFunc: func(ctx context.Context, tag *entity.Tag) (*entity.Tag, error) {
				return nil, errors.New("database error")
			},
		}
		usecase := NewTagUsecase(mockRepo)

		result, err := usecase.CreateTag(context.Background(), "Go")

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database error")
	})
}

// GetAllTagsのテスト
func TestGetAllTags(t *testing.T) {
	t.Run("正常系：全タグを取得できる", func(t *testing.T) {
		expected := []*entity.Tag{
			{
				ID:        1,
				Name:      "Go",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        2,
				Name:      "Next.js",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		mockRepo := &mockTagRepository{
			findAllFunc: func(ctx context.Context) ([]*entity.Tag, error) {
				return expected, nil
			},
		}

		usecase := NewTagUsecase(mockRepo)
		result, err := usecase.GetAllTags(context.Background())

		require.NoError(t, err)
		assert.Equal(t, 2, len(result))
		assert.Equal(t, expected, result)
	})

	t.Run("正常系：タグが0件の場合", func(t *testing.T) {
		mockRepo := &mockTagRepository{
			findAllFunc: func(ctx context.Context) ([]*entity.Tag, error) {
				return []*entity.Tag{}, nil
			},
		}

		usecase := NewTagUsecase(mockRepo)
		result, err := usecase.GetAllTags(context.Background())

		require.NoError(t, err)
		assert.Equal(t, 0, len(result))
		assert.NotNil(t, result)
	})

	t.Run("異常系：リポジトリがエラーを返す", func(t *testing.T) {
		mockRepo := &mockTagRepository{
			findAllFunc: func(ctx context.Context) ([]*entity.Tag, error) {
				return nil, errors.New("database error")
			},
		}

		usecase := NewTagUsecase(mockRepo)
		result, err := usecase.GetAllTags(context.Background())

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database error")
	})
}

// GetTagByIDのテスト
func TestGetTagByID(t *testing.T) {
	t.Run("正常系：IDでタグを取得できる", func(t *testing.T) {
		expected := &entity.Tag{
			ID:        1,
			Name:      "Go",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRepo := &mockTagRepository{
			findByIDFunc: func(ctx context.Context, id int64) (*entity.Tag, error) {
				return expected, nil
			},
		}

		usecase := NewTagUsecase(mockRepo)
		result, err := usecase.GetTagByID(context.Background(), 1)

		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("異常系：タグが見つからない", func(t *testing.T) {
		mockRepo := &mockTagRepository{
			findByIDFunc: func(ctx context.Context, id int64) (*entity.Tag, error) {
				return nil, errors.New("tag not found")
			},
		}

		usecase := NewTagUsecase(mockRepo)
		result, err := usecase.GetTagByID(context.Background(), 999)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "tag not found")
	})

	t.Run("異常系：不正なID（0以下）", func(t *testing.T) {
		mockRepo := &mockTagRepository{}
		usecase := NewTagUsecase(mockRepo)

		result, err := usecase.GetTagByID(context.Background(), 0)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "id must be positive")
	})
}

// GetTagByNameのテスト
func TestGetTagByName(t *testing.T) {
	t.Run("正常系：名前でタグを取得できる", func(t *testing.T) {
		expected := &entity.Tag{
			ID:        1,
			Name:      "Go",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRepo := &mockTagRepository{
			findByNameFunc: func(ctx context.Context, name string) (*entity.Tag, error) {
				return expected, nil
			},
		}

		usecase := NewTagUsecase(mockRepo)
		result, err := usecase.GetTagByName(context.Background(), "Go")

		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("異常系：タグが見つからない", func(t *testing.T) {
		mockRepo := &mockTagRepository{
			findByNameFunc: func(ctx context.Context, name string) (*entity.Tag, error) {
				return nil, errors.New("tag not found")
			},
		}

		usecase := NewTagUsecase(mockRepo)
		result, err := usecase.GetTagByName(context.Background(), "NotExist")

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "tag not found")
	})

	t.Run("異常系：タグ名が空の場合エラー", func(t *testing.T) {
		mockRepo := &mockTagRepository{}
		usecase := NewTagUsecase(mockRepo)

		result, err := usecase.GetTagByName(context.Background(), "")

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "name is required")
	})
}

// UpdateTagのテスト
func TestUpdateTag(t *testing.T) {
	t.Run("正常系：タグを更新できる", func(t *testing.T) {
		existingTag := &entity.Tag{
			ID:        1,
			Name:      "旧タグ名",
			CreatedAt: time.Now().Add(-24 * time.Hour),
			UpdatedAt: time.Now().Add(-24 * time.Hour),
		}

		mockRepo := &mockTagRepository{
			findByIDFunc: func(ctx context.Context, id int64) (*entity.Tag, error) {
				return existingTag, nil
			},
			updateFunc: func(ctx context.Context, tag *entity.Tag) (*entity.Tag, error) {
				return tag, nil
			},
		}

		usecase := NewTagUsecase(mockRepo)
		result, err := usecase.UpdateTag(context.Background(), 1, "新タグ名")

		require.NoError(t, err)
		assert.Equal(t, "新タグ名", result.Name)
	})

	t.Run("異常系：タグが存在しない", func(t *testing.T) {
		mockRepo := &mockTagRepository{
			findByIDFunc: func(ctx context.Context, id int64) (*entity.Tag, error) {
				return nil, errors.New("tag not found")
			},
		}

		usecase := NewTagUsecase(mockRepo)
		result, err := usecase.UpdateTag(context.Background(), 999, "新タグ名")

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "tag not found")
	})

	t.Run("異常系：更新データのバリデーションエラー", func(t *testing.T) {
		existingTag := &entity.Tag{
			ID:        1,
			Name:      "旧タグ名",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRepo := &mockTagRepository{
			findByIDFunc: func(ctx context.Context, id int64) (*entity.Tag, error) {
				return existingTag, nil
			},
		}

		usecase := NewTagUsecase(mockRepo)
		result, err := usecase.UpdateTag(context.Background(), 1, "")

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "name is required")
	})

	t.Run("異常系：リポジトリの更新でエラー", func(t *testing.T) {
		existingTag := &entity.Tag{
			ID:        1,
			Name:      "旧タグ名",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRepo := &mockTagRepository{
			findByIDFunc: func(ctx context.Context, id int64) (*entity.Tag, error) {
				return existingTag, nil
			},
			updateFunc: func(ctx context.Context, tag *entity.Tag) (*entity.Tag, error) {
				return nil, errors.New("database error")
			},
		}

		usecase := NewTagUsecase(mockRepo)
		result, err := usecase.UpdateTag(context.Background(), 1, "新タグ名")

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database error")
	})
}

// DeleteTagのテスト
func TestDeleteTag(t *testing.T) {
	t.Run("正常系：タグを削除できる", func(t *testing.T) {
		mockRepo := &mockTagRepository{
			deleteFunc: func(ctx context.Context, id int64) error {
				return nil
			},
		}

		usecase := NewTagUsecase(mockRepo)
		err := usecase.DeleteTag(context.Background(), 1)

		require.NoError(t, err)
	})

	t.Run("異常系：不正なID（0以下）", func(t *testing.T) {
		mockRepo := &mockTagRepository{}
		usecase := NewTagUsecase(mockRepo)

		err := usecase.DeleteTag(context.Background(), 0)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "id must be positive")
	})

	t.Run("異常系：タグが存在しない", func(t *testing.T) {
		mockRepo := &mockTagRepository{
			deleteFunc: func(ctx context.Context, id int64) error {
				return errors.New("tag not found")
			},
		}

		usecase := NewTagUsecase(mockRepo)
		err := usecase.DeleteTag(context.Background(), 999)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "tag not found")
	})

	t.Run("異常系：リポジトリがエラーを返す", func(t *testing.T) {
		mockRepo := &mockTagRepository{
			deleteFunc: func(ctx context.Context, id int64) error {
				return errors.New("database error")
			},
		}

		usecase := NewTagUsecase(mockRepo)
		err := usecase.DeleteTag(context.Background(), 1)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
	})
}
