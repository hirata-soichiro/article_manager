package repository

import (
	"context"
	"database/sql"
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

	// 必要なテーブルが全て存在するかチェック
	requiredTables := []string{"articles", "tags", "article_tags"}
	for _, tableName := range requiredTables {
		var tableExists int
		query := "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ? AND table_name = ?"
		err = db.Get(&tableExists, query, dbname, tableName)
		if err != nil || tableExists == 0 {
			db.Close()
			t.Skipf("%sテーブルが存在しません", tableName)
		}
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

	_, err := db.Exec("DELETE FROM article_tags")
	require.NoError(t, err, "article_tagsテーブルのクリーンアップに失敗")

	_, err = db.Exec("DELETE FROM articles")
	require.NoError(t, err, "articlesテーブルのクリーンアップに失敗")

	_, err = db.Exec("DELETE FROM tags")
	require.NoError(t, err, "tagsテーブルのクリーンアップに失敗")
}

// テスト用のタグをデータベースに挿入
func ensureTagsExist(t *testing.T, db *sqlx.DB, tagNames []string) map[string]int64 {
	t.Helper()

	tagIDMap := make(map[string]int64)

	for _, tagName := range tagNames {
		var tagID int64
		query := `SELECT id FROM tags WHERE name = ?`
		err := db.Get(&tagID, query, tagName)

		if err == sql.ErrNoRows {
			insertQuery := `INSERT INTO tags (name) VALUES (?)`
			result, err := db.Exec(insertQuery, tagName)
			require.NoError(t, err, "タグの挿入に失敗: %s", tagName)

			tagID, err = result.LastInsertId()
			require.NoError(t, err, "タグIDの取得に失敗")
		} else if err != nil {
			require.NoError(t, err, "タグの検索に失敗")
		}

		tagIDMap[tagName] = tagID
	}

	return tagIDMap
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

	var memo sql.NullString
	if article.Memo != "" {
		memo = sql.NullString{String: article.Memo, Valid: true}
	}

	query := `INSERT INTO articles (title, url, summary, memo, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := db.Exec(query, article.Title, article.URL, article.Summary, memo, article.CreatedAt, article.UpdatedAt)
	require.NoError(t, err, "記事の挿入に失敗")

	articleID, err := result.LastInsertId()
	require.NoError(t, err, "記事IDの取得に失敗")

	if len(article.Tags) > 0 {
		tagIDMap := ensureTagsExist(t, db, article.Tags)

		for _, tagName := range article.Tags {
			tagID := tagIDMap[tagName]
			insertQuery := `INSERT INTO article_tags (article_id, tag_id) VALUES (?, ?)`
			_, err := db.Exec(insertQuery, articleID, tagID)
			require.NoError(t, err, "article_tagsの挿入に失敗")
		}
	}

	return articleID
}

func TestMySQLArticleRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	t.Run("正常系：記事を作成できる", func(t *testing.T) {
		cleanupTable(t, db)

		ensureTagsExist(t, db, []string{"Go", "プログラミング"})

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
		assert.ElementsMatch(t, article.Tags, created.Tags)
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

		ensureTagsExist(t, db, []string{"Go"})

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

		ensureTagsExist(t, db, []string{"tag1", "tag2"})

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

		ensureTagsExist(t, db, []string{"Go"})

		repo := NewMySQLArticleRepository(db)
		article := createTestArticle(t, "Go言語入門", "https://example.com/go", "Go言語の基本", []string{"Go"}, "")

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		created, err := repo.Create(ctx, article)

		require.Error(t, err)
		assert.Nil(t, created)
	})

	t.Run("異常系：存在しないタグを指定すると作成に失敗する", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		// タグを事前に作成しない
		article := createTestArticle(t, "Go言語入門", "https://example.com/go", "Go言語の基本", []string{"NonExistentTag"}, "")

		ctx := context.Background()
		created, err := repo.Create(ctx, article)

		require.Error(t, err)
		assert.Nil(t, created)
		assert.Contains(t, err.Error(), "tag not found")
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
		assert.ElementsMatch(t, article.Tags, found.Tags)
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

	t.Run("正常系：複数のタグを持つ記事を取得できる", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article := createTestArticle(t, "Go言語入門", "https://example.com/go", "Go言語の基本", []string{"Go", "プログラミング", "バックエンド"}, "")
		insertArticleDirectly(t, db, article)

		ctx := context.Background()
		articles, err := repo.FindAll(ctx)

		require.NoError(t, err)
		require.NotNil(t, articles)
		assert.Len(t, articles, 1)
		assert.ElementsMatch(t, []string{"Go", "プログラミング", "バックエンド"}, articles[0].Tags)
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

		ensureTagsExist(t, db, []string{"Go", "完全ガイド"})

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
		assert.ElementsMatch(t, []string{"Go", "完全ガイド"}, updated.Tags)
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

	t.Run("異常系：存在しないタグで更新しようとするとエラー", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article := createTestArticle(t, "Go言語入門", "https://example.com/go", "Go言語の基本", []string{"Go"}, "")
		id := insertArticleDirectly(t, db, article)

		article.ID = id
		err := article.Update("Go言語完全ガイド", "https://example.com/go-guide", "Go言語の完全版", []string{"NonExistentTag"}, "")
		require.NoError(t, err)

		ctx := context.Background()
		updated, err := repo.Update(ctx, article)

		require.Error(t, err)
		assert.Nil(t, updated)
		assert.Contains(t, err.Error(), "tag not found")
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

	t.Run("正常系：記事を削除するとarticle_tagsも削除される", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article := createTestArticle(t, "Go言語入門", "https://example.com/go", "Go言語の基本", []string{"Go", "プログラミング"}, "")
		id := insertArticleDirectly(t, db, article)

		ctx := context.Background()
		err := repo.Delete(ctx, id)

		require.NoError(t, err)

		var count int
		err = db.Get(&count, "SELECT COUNT(*) FROM article_tags WHERE article_id = ?", id)
		require.NoError(t, err)
		assert.Equal(t, 0, count, "article_tagsが削除されていません")
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

func TestMySQLArticleRepository_Search(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	t.Run("正常系：単一キーワードでタイトルを検索できる", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article1 := createTestArticle(t, "Go言語入門", "https://example.com/go", "Go言語の基本を学ぶ", []string{"Go"}, "")
		article2 := createTestArticle(t, "Python入門", "https://example.com/python", "Pythonの基本を学ぶ", []string{"Python"}, "")
		article3 := createTestArticle(t, "Go言語完全ガイド", "https://example.com/go-guide", "Go言語の完全版", []string{"Go"}, "")

		insertArticleDirectly(t, db, article1)
		insertArticleDirectly(t, db, article2)
		insertArticleDirectly(t, db, article3)

		ctx := context.Background()
		results, err := repo.Search(ctx, "Go言語")

		require.NoError(t, err)
		require.NotNil(t, results)
		assert.Len(t, results, 2)
		// 新しい順（created_at DESC）で並ぶことを確認
		assert.Equal(t, "Go言語完全ガイド", results[0].Title)
		assert.Equal(t, "Go言語入門", results[1].Title)
	})

	t.Run("正常系：単一キーワードで要約（Summary）を検索できる", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article1 := createTestArticle(t, "プログラミング入門", "https://example.com/prog1", "Go言語とPythonの比較", []string{"Go", "Python"}, "")
		article2 := createTestArticle(t, "Web開発入門", "https://example.com/web", "JavaScriptでWeb開発", []string{"JavaScript"}, "")
		article3 := createTestArticle(t, "データベース入門", "https://example.com/db", "MySQLとPostgreSQLの比較", []string{"MySQL"}, "")

		insertArticleDirectly(t, db, article1)
		insertArticleDirectly(t, db, article2)
		insertArticleDirectly(t, db, article3)

		ctx := context.Background()
		results, err := repo.Search(ctx, "比較")

		require.NoError(t, err)
		require.NotNil(t, results)
		assert.Len(t, results, 2)
		assert.Equal(t, "データベース入門", results[0].Title)
		assert.Equal(t, "プログラミング入門", results[1].Title)
	})

	t.Run("正常系：複数キーワード（2単語）でAND検索できる", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article1 := createTestArticle(t, "Go言語入門", "https://example.com/go", "Go言語の基本を学ぶ", []string{"Go"}, "")
		article2 := createTestArticle(t, "Go言語完全ガイド", "https://example.com/go-guide", "Go言語の完全版", []string{"Go"}, "")
		article3 := createTestArticle(t, "Go言語による設計パターン", "https://example.com/go-pattern", "設計パターンをGoで実装", []string{"Go"}, "")
		article4 := createTestArticle(t, "Python入門", "https://example.com/python", "Pythonの基本を学ぶ", []string{"Python"}, "")

		insertArticleDirectly(t, db, article1)
		insertArticleDirectly(t, db, article2)
		insertArticleDirectly(t, db, article3)
		insertArticleDirectly(t, db, article4)

		ctx := context.Background()
		// "Go言語"と"完全"の両方を含む記事を検索
		results, err := repo.Search(ctx, "Go言語 完全")

		require.NoError(t, err)
		require.NotNil(t, results)
		assert.Len(t, results, 1)
		assert.Equal(t, "Go言語完全ガイド", results[0].Title)
	})

	t.Run("正常系：複数キーワード（3単語）でAND検索できる", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article1 := createTestArticle(t, "Go言語入門", "https://example.com/go", "Go言語の基本を学ぶ", []string{"Go"}, "")
		article2 := createTestArticle(t, "Go言語完全ガイド", "https://example.com/go-guide", "Go言語の基本から応用まで完全網羅", []string{"Go"}, "")
		article3 := createTestArticle(t, "Go言語による設計パターン", "https://example.com/go-pattern", "設計パターンをGoで実装", []string{"Go"}, "")

		insertArticleDirectly(t, db, article1)
		insertArticleDirectly(t, db, article2)
		insertArticleDirectly(t, db, article3)

		ctx := context.Background()
		// "Go言語"、"基本"、"完全"の全てを含む記事を検索
		results, err := repo.Search(ctx, "Go言語 基本 完全")

		require.NoError(t, err)
		require.NotNil(t, results)
		assert.Len(t, results, 1)
		assert.Equal(t, "Go言語完全ガイド", results[0].Title)
	})

	t.Run("正常系：タイトルと要約の両方を検索対象とする", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		// タイトルにのみ"Go言語"を含む記事
		article1 := createTestArticle(t, "Go言語入門", "https://example.com/go1", "プログラミングの基礎", []string{"Go"}, "")
		// 要約にのみ"Go言語"を含む記事
		article2 := createTestArticle(t, "プログラミング入門", "https://example.com/go2", "Go言語を使った開発", []string{"Go"}, "")
		// タイトルと要約の両方に"Go言語"を含む記事
		article3 := createTestArticle(t, "Go言語完全ガイド", "https://example.com/go3", "Go言語の完全版", []string{"Go"}, "")
		// どちらにも含まない記事
		article4 := createTestArticle(t, "Python入門", "https://example.com/python", "Pythonの基本", []string{"Python"}, "")

		insertArticleDirectly(t, db, article1)
		insertArticleDirectly(t, db, article2)
		insertArticleDirectly(t, db, article3)
		insertArticleDirectly(t, db, article4)

		ctx := context.Background()
		results, err := repo.Search(ctx, "Go言語")

		require.NoError(t, err)
		require.NotNil(t, results)
		assert.Len(t, results, 3)
	})

	t.Run("正常系：複数キーワードがタイトルと要約に分散している場合もマッチする", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		// "Go言語"はタイトル、"パフォーマンス"は要約に存在
		article1 := createTestArticle(t, "Go言語入門", "https://example.com/go1", "パフォーマンスの最適化", []string{"Go"}, "")
		// 両方ともタイトルに存在
		article2 := createTestArticle(t, "Go言語パフォーマンス最適化", "https://example.com/go2", "実践的な最適化手法", []string{"Go"}, "")
		// "Go言語"は要約、"パフォーマンス"はタイトルに存在
		article3 := createTestArticle(t, "パフォーマンスチューニング", "https://example.com/go3", "Go言語での最適化", []string{"Go"}, "")
		// 片方のキーワードのみ
		article4 := createTestArticle(t, "Go言語完全ガイド", "https://example.com/go4", "基本から応用まで", []string{"Go"}, "")

		insertArticleDirectly(t, db, article1)
		insertArticleDirectly(t, db, article2)
		insertArticleDirectly(t, db, article3)
		insertArticleDirectly(t, db, article4)

		ctx := context.Background()
		results, err := repo.Search(ctx, "Go言語 パフォーマンス")

		require.NoError(t, err)
		require.NotNil(t, results)
		assert.Len(t, results, 3)
		// article4は"Go言語"のみでマッチしない
	})

	t.Run("正常系：検索キーワードが空文字列の場合は全件取得と同じ", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article1 := createTestArticle(t, "Go言語入門", "https://example.com/go", "Go言語の基本", []string{"Go"}, "")
		article2 := createTestArticle(t, "Python入門", "https://example.com/python", "Pythonの基本", []string{"Python"}, "")

		insertArticleDirectly(t, db, article1)
		time.Sleep(10 * time.Millisecond)
		insertArticleDirectly(t, db, article2)

		ctx := context.Background()
		results, err := repo.Search(ctx, "")

		require.NoError(t, err)
		require.NotNil(t, results)
		assert.Len(t, results, 2)
	})

	t.Run("正常系：検索キーワードがスペースのみの場合は全件取得と同じ", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article1 := createTestArticle(t, "Go言語入門", "https://example.com/go", "Go言語の基本", []string{"Go"}, "")
		article2 := createTestArticle(t, "Python入門", "https://example.com/python", "Pythonの基本", []string{"Python"}, "")

		insertArticleDirectly(t, db, article1)
		time.Sleep(10 * time.Millisecond)
		insertArticleDirectly(t, db, article2)

		ctx := context.Background()
		results, err := repo.Search(ctx, "   ")

		require.NoError(t, err)
		require.NotNil(t, results)
		assert.Len(t, results, 2)
	})

	t.Run("正常系：マッチする記事が0件の場合は空配列を返す", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article1 := createTestArticle(t, "Go言語入門", "https://example.com/go", "Go言語の基本", []string{"Go"}, "")
		article2 := createTestArticle(t, "Python入門", "https://example.com/python", "Pythonの基本", []string{"Python"}, "")

		insertArticleDirectly(t, db, article1)
		insertArticleDirectly(t, db, article2)

		ctx := context.Background()
		results, err := repo.Search(ctx, "存在しないキーワード")

		require.NoError(t, err)
		require.NotNil(t, results)
		assert.Empty(t, results)
	})

	t.Run("正常系：記事が0件の状態で検索しても空配列を返す", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		ctx := context.Background()
		results, err := repo.Search(ctx, "Go言語")

		require.NoError(t, err)
		require.NotNil(t, results)
		assert.Empty(t, results)
	})

	t.Run("正常系：大文字小文字を区別せず検索する", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article1 := createTestArticle(t, "Go Language Tutorial", "https://example.com/go", "Learn GO programming", []string{"Go"}, "")
		article2 := createTestArticle(t, "python tutorial", "https://example.com/python", "Learn python", []string{"Python"}, "")

		insertArticleDirectly(t, db, article1)
		insertArticleDirectly(t, db, article2)

		ctx := context.Background()
		results, err := repo.Search(ctx, "go")

		require.NoError(t, err)
		require.NotNil(t, results)
		assert.Len(t, results, 1)
		assert.Equal(t, "Go Language Tutorial", results[0].Title)
	})

	t.Run("正常系：前後の空白を無視して検索する", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article1 := createTestArticle(t, "Go言語入門", "https://example.com/go", "Go言語の基本", []string{"Go"}, "")
		article2 := createTestArticle(t, "Python入門", "https://example.com/python", "Pythonの基本", []string{"Python"}, "")

		insertArticleDirectly(t, db, article1)
		insertArticleDirectly(t, db, article2)

		ctx := context.Background()
		results, err := repo.Search(ctx, "  Go言語  ")

		require.NoError(t, err)
		require.NotNil(t, results)
		assert.Len(t, results, 1)
		assert.Equal(t, "Go言語入門", results[0].Title)
	})

	t.Run("正常系：複数の連続するスペースも1つの区切りとして扱う", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article1 := createTestArticle(t, "Go言語完全ガイド", "https://example.com/go", "Go言語の基本から応用まで", []string{"Go"}, "")
		article2 := createTestArticle(t, "Go言語入門", "https://example.com/go2", "初心者向けの内容", []string{"Go"}, "")

		insertArticleDirectly(t, db, article1)
		insertArticleDirectly(t, db, article2)

		ctx := context.Background()
		// 複数スペースで区切っても正しく検索できる
		results, err := repo.Search(ctx, "Go言語    完全")

		require.NoError(t, err)
		require.NotNil(t, results)
		assert.Len(t, results, 1)
		assert.Equal(t, "Go言語完全ガイド", results[0].Title)
	})

	t.Run("正常系：部分一致で検索できる", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article1 := createTestArticle(t, "Go言語入門", "https://example.com/go", "Go言語の基本", []string{"Go"}, "")
		article2 := createTestArticle(t, "Python入門", "https://example.com/python", "Pythonの基本", []string{"Python"}, "")

		insertArticleDirectly(t, db, article1)
		insertArticleDirectly(t, db, article2)

		ctx := context.Background()
		// "Go"で検索して"Go言語"にマッチする
		results, err := repo.Search(ctx, "Go")

		require.NoError(t, err)
		require.NotNil(t, results)
		assert.Len(t, results, 1)
		assert.Equal(t, "Go言語入門", results[0].Title)
	})

	t.Run("正常系：結果はcreated_atの降順でソートされる", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article1 := createTestArticle(t, "Go言語入門", "https://example.com/go1", "Go言語の基本", []string{"Go"}, "")
		article2 := createTestArticle(t, "Go言語応用", "https://example.com/go2", "Go言語の応用", []string{"Go"}, "")
		article3 := createTestArticle(t, "Go言語完全ガイド", "https://example.com/go3", "Go言語の完全版", []string{"Go"}, "")

		insertArticleDirectly(t, db, article1)
		time.Sleep(10 * time.Millisecond)
		insertArticleDirectly(t, db, article2)
		time.Sleep(10 * time.Millisecond)
		insertArticleDirectly(t, db, article3)

		ctx := context.Background()
		results, err := repo.Search(ctx, "Go言語")

		require.NoError(t, err)
		require.NotNil(t, results)
		assert.Len(t, results, 3)
		// 新しい順に並ぶ
		assert.Equal(t, "Go言語完全ガイド", results[0].Title)
		assert.Equal(t, "Go言語応用", results[1].Title)
		assert.Equal(t, "Go言語入門", results[2].Title)
	})

	t.Run("正常系：検索結果にタグも含まれる", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article1 := createTestArticle(t, "Go言語入門", "https://example.com/go", "Go言語の基本", []string{"Go", "プログラミング", "入門"}, "")
		insertArticleDirectly(t, db, article1)

		ctx := context.Background()
		results, err := repo.Search(ctx, "Go言語")

		require.NoError(t, err)
		require.NotNil(t, results)
		assert.Len(t, results, 1)
		assert.ElementsMatch(t, []string{"Go", "プログラミング", "入門"}, results[0].Tags)
	})

	t.Run("正常系：検索結果にメモも含まれる", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article1 := createTestArticle(t, "Go言語入門", "https://example.com/go", "Go言語の基本", []string{"Go"}, "後で読む")
		insertArticleDirectly(t, db, article1)

		ctx := context.Background()
		results, err := repo.Search(ctx, "Go言語")

		require.NoError(t, err)
		require.NotNil(t, results)
		assert.Len(t, results, 1)
		assert.Equal(t, "後で読む", results[0].Memo)
	})

	t.Run("異常系：キャンセルされたコンテキストではエラー", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article := createTestArticle(t, "Go言語入門", "https://example.com/go", "Go言語の基本", []string{"Go"}, "")
		insertArticleDirectly(t, db, article)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		results, err := repo.Search(ctx, "Go言語")

		require.Error(t, err)
		assert.Nil(t, results)
	})

	t.Run("正常系：特殊文字を含むキーワードでも検索できる", func(t *testing.T) {
		cleanupTable(t, db)
		repo := NewMySQLArticleRepository(db)

		article1 := createTestArticle(t, "C++プログラミング", "https://example.com/cpp", "C++の基本", []string{"C++"}, "")
		article2 := createTestArticle(t, "C#プログラミング", "https://example.com/csharp", "C#の基本", []string{"C#"}, "")

		insertArticleDirectly(t, db, article1)
		insertArticleDirectly(t, db, article2)

		ctx := context.Background()
		results, err := repo.Search(ctx, "C++")

		require.NoError(t, err)
		require.NotNil(t, results)
		assert.Len(t, results, 1)
		assert.Equal(t, "C++プログラミング", results[0].Title)
	})
}
