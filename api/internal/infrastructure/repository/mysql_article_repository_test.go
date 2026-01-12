package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"testing"
	"time"

	"article-manager/internal/domain/entity"
	"article-manager/internal/infrastructure/database"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// テスト用のDB接続を取得
func setupTestDB(t *testing.T) *sqlx.DB {
	t.Helper()

	host := getEnv("TEST_DB_HOST", "localhost")
	port := getEnv("TEST_DB_PORT", "3307")
	user := getEnv("TEST_DB_USER", "test_user")
	password := getEnv("TEST_DB_PASSWORD", "test_password")
	dbname := getEnv("TEST_DB_NAME", "article_manager_test")

	db, err := database.NewMySQL(host, port, user, password, dbname)
	if err != nil {
		t.Skipf("テスト用データベースに接続できません： %v", err)
	}

	var tableExists int
	query := "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ? AND table_name = 'articles'"
	err = db.Get(&tableExists, query, dbname)
	if err != nil || tableExists == 0 {
		db.Close()
		t.Skip("articlesテーブルが存在しません")
	}

	return db
}

// 環境変数を取得
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// テーブルをクリーンアップ
func cleanupTable(t *testing.T, db *sqlx.DB) {
	t.Helper()
	_, err := db.Exec("TRUNCATE TABLE articles")
	require.NoError(t, err, "テーブルのクリーンアップに失敗")
}

// テスト用の記事データを作成
func createTestArticle(t *testing.T, title, url, summary string, tags []string, memo string) *entity.Article {
	t.Helper()
	article, err := entity.NewArticle(title, url, summary, tags, memo)
	require.NoError(t, err)
	return article
}

// データベースに直接記事を挿入
func insertArticleDirectly(t *testing.T, db *sqlx.DB, article *entity.Article) int64 {
	t.Helper()

	tagsJSON, err := json.Marshal(article.Tags)
	require.NoError(t, err)

	var memo sql.NullString
	if article.Memo != "" {
		memo = sql.NullString{String: article.Memo, Valid: true}
	}

	query := `INSERT INTO articles (title, url, summary, tags, memo, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	result, err := db.Exec(query, article.Title, article.URL, article.Summary, tagsJSON, memo, article.CreatedAt, article.UpdatedAt)
	require.NoError(t, err)

	id, err := result.LastInsertId()
	require.NoError(t, err)

	return id
}

func TestMySQLArticleRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	t.Run("正常系：記事を作成できる", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article := createTestArticle(t, "Go言語入門", "https://example.com/go", "Go言語の基本", []string{"Go", "プログラミング"}, "後で読む")

		ctx := context.Background()
		created, err := repo.Create(ctx, article)

		require.NoError(t, err)
		require.NotNil(t, created)
		assert.Greater(t, created.ID, int64(0))
		assert.Equal(t, article.Title, created.Title)
		assert.Equal(t, article.URL, created.URL)
		assert.Equal(t, article.Summary, created.Summary)
		assert.Equal(t, article.Tags, created.Tags)
		assert.Equal(t, article.Memo, created.Memo)
		assert.False(t, created.CreatedAt.IsZero())
		assert.False(t, created.UpdatedAt.IsZero())
	})

	t.Run("正常系：タグが空配列の記事を作成できる", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article := createTestArticle(t, "Go言語入門", "https://example.com/go", "Go言語の基本", []string{}, "")

		ctx := context.Background()
		created, err := repo.Create(ctx, article)

		require.NoError(t, err)
		require.NotNil(t, created)
		assert.Greater(t, created.ID, int64(0))
		assert.Empty(t, created.Tags)
	})

	t.Run("正常系：メモが空文字列の記事を作成できる", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article := createTestArticle(t, "Go言語入門", "https://example.com/go", "Go言語の基本", []string{"Go"}, "")

		ctx := context.Background()
		created, err := repo.Create(ctx, article)

		require.NoError(t, err)
		require.NotNil(t, created)
		assert.Empty(t, created.Memo)
	})

	t.Run("正常系：複数の記事を作成できる", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article1 := createTestArticle(t, "記事1", "https://example.com/1", "要約1", []string{"tag1"}, "")
		article2 := createTestArticle(t, "記事2", "https://example.com/2", "要約2", []string{"tag2"}, "")

		ctx := context.Background()
		created1, err := repo.Create(ctx, article1)
		require.NoError(t, err)

		created2, err := repo.Create(ctx, article2)
		require.NoError(t, err)

		assert.NotEqual(t, created1.ID, created2.ID)
	})

	t.Run("異常系：nilの記事を作成しようとするとエラー", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		ctx := context.Background()
		created, err := repo.Create(ctx, nil)

		require.Error(t, err)
		assert.Nil(t, created)
	})

	t.Run("異常系：キャンセルされたコンテキストではエラー", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article := createTestArticle(t, "Go言語入門", "https://example.com/go", "Go言語の基本", []string{"Go"}, "")

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		created, err := repo.Create(ctx, article)

		require.Error(t, err)
		assert.Nil(t, created)
	})
}

func TestMySQLArticleRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	t.Run("正常系：IDで記事を取得できる", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article := createTestArticle(t, "Go言語入門", "https://example.com/go", "Go言語の基本", []string{"Go", "プログラミング"}, "後で読む")
		id := insertArticleDirectly(t, db, article)

		ctx := context.Background()
		found, err := repo.FindByID(ctx, id)

		require.NoError(t, err)
		require.NotNil(t, found)
		assert.Equal(t, id, found.ID)
		assert.Equal(t, article.Title, found.Title)
		assert.Equal(t, article.URL, found.URL)
		assert.Equal(t, article.Summary, found.Summary)
		assert.Equal(t, article.Tags, found.Tags)
		assert.Equal(t, article.Memo, found.Memo)
	})

	t.Run("正常系：タグが空配列の記事を取得できる", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article := createTestArticle(t, "Go言語入門", "https://example.com/go", "Go言語の基本", []string{}, "")
		id := insertArticleDirectly(t, db, article)

		ctx := context.Background()
		found, err := repo.FindByID(ctx, id)

		require.NoError(t, err)
		require.NotNil(t, found)
		assert.Empty(t, found.Tags)
	})

	t.Run("正常系：メモが空の記事を取得できる", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article := createTestArticle(t, "Go言語入門", "https://example.com/go", "Go言語の基本", []string{"Go"}, "")
		id := insertArticleDirectly(t, db, article)

		ctx := context.Background()
		found, err := repo.FindByID(ctx, id)

		require.NoError(t, err)
		require.NotNil(t, found)
		assert.Empty(t, found.Memo)
	})

	t.Run("異常系：存在しないIDで記事を取得するとエラー", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		ctx := context.Background()
		found, err := repo.FindByID(ctx, 99999)

		require.Error(t, err)
		assert.Nil(t, found)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("異常系：負のIDで記事を取得するとエラー", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		ctx := context.Background()
		found, err := repo.FindByID(ctx, -1)

		require.Error(t, err)
		assert.Nil(t, found)
	})

	t.Run("異常系：キャンセルされたコンテキストではエラー", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article := createTestArticle(t, "Go言語入門", "https://example.com/go", "Go言語の基本", []string{"Go"}, "")
		id := insertArticleDirectly(t, db, article)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		found, err := repo.FindByID(ctx, id)

		require.Error(t, err)
		assert.Nil(t, found)
	})
}

func TestMySQLArticleRepository_FindAll(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	t.Run("正常系：複数の記事を取得できる", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article1 := createTestArticle(t, "記事1", "https://example.com/1", "要約1", []string{"tag1"}, "メモ1")
		article2 := createTestArticle(t, "記事2", "https://example.com/2", "要約2", []string{"tag2"}, "メモ2")
		article3 := createTestArticle(t, "記事3", "https://example.com/3", "要約3", []string{"tag3"}, "")

		insertArticleDirectly(t, db, article1)
		time.Sleep(10 * time.Millisecond)
		insertArticleDirectly(t, db, article2)
		time.Sleep(10 * time.Millisecond)
		insertArticleDirectly(t, db, article3)

		ctx := context.Background()
		articles, err := repo.FindAll(ctx)

		require.NoError(t, err)
		require.NotNil(t, articles)
		assert.Len(t, articles, 3)

		assert.Equal(t, "記事3", articles[0].Title)
		assert.Equal(t, "記事2", articles[1].Title)
		assert.Equal(t, "記事1", articles[2].Title)
	})

	t.Run("正常系：記事が0件の場合は空配列を返す", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		ctx := context.Background()
		articles, err := repo.FindAll(ctx)

		require.NoError(t, err)
		require.NotNil(t, articles)
		assert.Empty(t, articles)
	})

	t.Run("正常系：1件の記事を取得できる", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article := createTestArticle(t, "Go言語入門", "https://example.com/go", "Go言語の基本", []string{"Go"}, "")
		insertArticleDirectly(t, db, article)

		ctx := context.Background()
		articles, err := repo.FindAll(ctx)

		require.NoError(t, err)
		require.NotNil(t, articles)
		assert.Len(t, articles, 1)
		assert.Equal(t, article.Title, articles[0].Title)
	})

	t.Run("正常系：タグが空配列の記事も取得できる", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article := createTestArticle(t, "Go言語入門", "https://example.com/go", "Go言語の基本", []string{}, "")
		insertArticleDirectly(t, db, article)

		ctx := context.Background()
		articles, err := repo.FindAll(ctx)

		require.NoError(t, err)
		require.NotNil(t, articles)
		assert.Len(t, articles, 1)
		assert.Empty(t, articles[0].Tags)
	})

	t.Run("異常系：キャンセルされたコンテキストではエラー", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		articles, err := repo.FindAll(ctx)

		require.Error(t, err)
		assert.Nil(t, articles)
	})
}

func TestMySQLArticleRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	t.Run("正常系：記事を更新できる", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article := createTestArticle(t, "Go言語入門", "https://example.com/go", "Go言語の基本", []string{"Go"}, "後で読む")
		id := insertArticleDirectly(t, db, article)

		article.ID = id
		time.Sleep(10 * time.Millisecond)
		err := article.Update("Go言語完全ガイド", "https://example.com/go-guide", "Go言語の完全版", []string{"Go", "完全ガイド"}, "重要")
		require.NoError(t, err)

		ctx := context.Background()
		updated, err := repo.Update(ctx, article)

		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, id, updated.ID)
		assert.Equal(t, "Go言語完全ガイド", updated.Title)
		assert.Equal(t, "https://example.com/go-guide", updated.URL)
		assert.Equal(t, "Go言語の完全版", updated.Summary)
		assert.Equal(t, []string{"Go", "完全ガイド"}, updated.Tags)
		assert.Equal(t, "重要", updated.Memo)
	})

	t.Run("正常系：タグを空配列に更新できる", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article := createTestArticle(t, "Go言語入門", "https://example.com/go", "Go言語の基本", []string{"Go", "プログラミング"}, "")
		id := insertArticleDirectly(t, db, article)

		article.ID = id
		err := article.Update("Go言語入門", "https://example.com/go", "Go言語の基本", []string{}, "")
		require.NoError(t, err)

		ctx := context.Background()
		updated, err := repo.Update(ctx, article)

		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Empty(t, updated.Tags)
	})

	t.Run("正常系：メモを空文字列に更新できる", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article := createTestArticle(t, "Go言語入門", "https://example.com/go", "Go言語の基本", []string{"Go"}, "後で読む")
		id := insertArticleDirectly(t, db, article)

		article.ID = id
		err := article.Update("Go言語入門", "https://example.com/go", "Go言語の基本", []string{"Go"}, "")
		require.NoError(t, err)

		ctx := context.Background()
		updated, err := repo.Update(ctx, article)

		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Empty(t, updated.Memo)
	})

	t.Run("正常系：UpdatedAtが更新される", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article := createTestArticle(t, "Go言語入門", "https://example.com/go", "Go言語の基本", []string{"Go"}, "")
		id := insertArticleDirectly(t, db, article)

		// 少し待機してから更新
		time.Sleep(100 * time.Millisecond)

		article.ID = id
		originalUpdatedAt := article.UpdatedAt
		err := article.Update("Go言語完全ガイド", "https://example.com/go-guide", "Go言語の完全版", []string{"Go"}, "")
		require.NoError(t, err)

		ctx := context.Background()
		updated, err := repo.Update(ctx, article)

		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.True(t, updated.UpdatedAt.After(originalUpdatedAt))
	})

	t.Run("異常系：存在しないIDの記事を更新しようとするとエラー", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article := createTestArticle(t, "Go言語入門", "https://example.com/go", "Go言語の基本", []string{"Go"}, "")
		article.ID = 99999

		ctx := context.Background()
		updated, err := repo.Update(ctx, article)

		require.Error(t, err)
		assert.Nil(t, updated)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("異常系：IDが0の記事を更新しようとするとエラー", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article := createTestArticle(t, "Go言語入門", "https://example.com/go", "Go言語の基本", []string{"Go"}, "")
		article.ID = 0

		ctx := context.Background()
		updated, err := repo.Update(ctx, article)

		require.Error(t, err)
		assert.Nil(t, updated)
	})

	t.Run("異常系：nilの記事を更新しようとするとエラー", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		ctx := context.Background()
		updated, err := repo.Update(ctx, nil)

		require.Error(t, err)
		assert.Nil(t, updated)
	})

	t.Run("異常系：キャンセルされたコンテキストではエラー", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article := createTestArticle(t, "Go言語入門", "https://example.com/go", "Go言語の基本", []string{"Go"}, "")
		id := insertArticleDirectly(t, db, article)
		article.ID = id

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		updated, err := repo.Update(ctx, article)

		require.Error(t, err)
		assert.Nil(t, updated)
	})
}

func TestMySQLArticleRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	t.Run("正常系：記事を削除できる", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article := createTestArticle(t, "Go言語入門", "https://example.com/go", "Go言語の基本", []string{"Go"}, "")
		id := insertArticleDirectly(t, db, article)

		ctx := context.Background()
		err := repo.Delete(ctx, id)

		require.NoError(t, err)

		// 削除後に取得できないことを確認
		found, err := repo.FindByID(ctx, id)
		require.Error(t, err)
		assert.Nil(t, found)
	})

	t.Run("正常系：複数の記事のうち1つを削除できる", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article1 := createTestArticle(t, "記事1", "https://example.com/1", "要約1", []string{"tag1"}, "")
		article2 := createTestArticle(t, "記事2", "https://example.com/2", "要約2", []string{"tag2"}, "")

		id1 := insertArticleDirectly(t, db, article1)
		id2 := insertArticleDirectly(t, db, article2)

		ctx := context.Background()
		err := repo.Delete(ctx, id1)
		require.NoError(t, err)

		// id2の記事は残っていることを確認
		found, err := repo.FindByID(ctx, id2)
		require.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, id2, found.ID)
	})

	t.Run("異常系：存在しないIDの記事を削除しようとするとエラー", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		ctx := context.Background()
		err := repo.Delete(ctx, 99999)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("異常系：IDが0の記事を削除しようとするとエラー", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		ctx := context.Background()
		err := repo.Delete(ctx, 0)

		require.Error(t, err)
	})

	t.Run("異常系：負のIDの記事を削除しようとするとエラー", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		ctx := context.Background()
		err := repo.Delete(ctx, -1)

		require.Error(t, err)
	})

	t.Run("異常系：キャンセルされたコンテキストではエラー", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article := createTestArticle(t, "Go言語入門", "https://example.com/go", "Go言語の基本", []string{"Go"}, "")
		id := insertArticleDirectly(t, db, article)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := repo.Delete(ctx, id)

		require.Error(t, err)
	})

	t.Run("異常系：同じIDを2回削除しようとすると2回目はエラー", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article := createTestArticle(t, "Go言語入門", "https://example.com/go", "Go言語の基本", []string{"Go"}, "")
		id := insertArticleDirectly(t, db, article)

		ctx := context.Background()

		// 1回目の削除は成功
		err := repo.Delete(ctx, id)
		require.NoError(t, err)

		// 2回目の削除は失敗
		err = repo.Delete(ctx, id)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}
