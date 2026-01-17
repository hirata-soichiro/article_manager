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
func setupTestDBForTag(t *testing.T) *sqlx.DB {
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
func getEnvForTag(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// テーブルをクリーンアップ
func cleanupTableForTag(t *testing.T, db *sqlx.DB) {
	t.Helper()

	_, err := db.Exec("DELETE FROM article_tags")
	require.NoError(t, err, "article_tagsテーブルのクリーンアップに失敗")

	_, err = db.Exec("DELETE FROM articles")
	require.NoError(t, err, "articlesテーブルのクリーンアップに失敗")

	_, err = db.Exec("DELETE FROM tags")
	require.NoError(t, err, "tagsテーブルのクリーンアップに失敗")
}

// テスト用の記事をデータベースに挿入
func ensureArticlesExist(t *testing.T, db *sqlx.DB, count int) []int64 {
	t.Helper()

	articleIDs := make([]int64, count)
	now := time.Now()

	for i := 0; i < count; i++ {
		query := `INSERT INTO articles (title, url, summary, memo, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`
		result, err := db.Exec(query,
			"テスト記事"+string(rune('A'+i)),
			"https://example.com/"+string(rune('a'+i)),
			"テスト要約"+string(rune('A'+i)),
			sql.NullString{},
			now,
			now,
		)
		require.NoError(t, err, "記事の挿入に失敗")

		articleID, err := result.LastInsertId()
		require.NoError(t, err, "記事IDの取得に失敗")
		articleIDs[i] = articleID
	}

	return articleIDs
}

// テスト用のタグデータを作成
func createTestTag(t *testing.T, name string) *entity.Tag {
	t.Helper()
	tag, err := entity.NewTag(name)
	require.NoError(t, err)
	return tag
}

// データベースに直接タグを挿入
func insertTagDirectly(t *testing.T, db *sqlx.DB, tag *entity.Tag) int64 {
	t.Helper()

	query := `INSERT INTO tags (name, created_at, updated_at) VALUES (?, ?, ?)`
	result, err := db.Exec(query, tag.Name, tag.CreatedAt, tag.UpdatedAt)
	require.NoError(t, err, "タグの挿入に失敗")

	id, err := result.LastInsertId()
	require.NoError(t, err, "タグIDの取得に失敗")

	return id
}

func TestMySQLTagRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	t.Run("正常系：タグを作成できる", func(t *testing.T) {
		cleanupTableForTag(t, db)
		repo := NewMySQLTagRepository(db)

		tag := createTestTag(t, "Go")

		ctx := context.Background()
		created, err := repo.Create(ctx, tag)

		require.NoError(t, err)
		require.NotNil(t, created)
		assert.Greater(t, created.ID, int64(0))
		assert.Equal(t, tag.Name, created.Name)
		assert.False(t, created.CreatedAt.IsZero())
		assert.False(t, created.UpdatedAt.IsZero())
	})

	t.Run("正常系：複数のタグを作成できる", func(t *testing.T) {
		cleanupTableForTag(t, db)
		repo := NewMySQLTagRepository(db)

		tag1 := createTestTag(t, "Go")
		tag2 := createTestTag(t, "Python")

		ctx := context.Background()
		created1, err := repo.Create(ctx, tag1)
		require.NoError(t, err)

		created2, err := repo.Create(ctx, tag2)
		require.NoError(t, err)

		assert.NotEqual(t, created1.ID, created2.ID)
	})

	t.Run("異常系：nilのタグを作成しようとするとエラー", func(t *testing.T) {
		cleanupTableForTag(t, db)
		repo := NewMySQLTagRepository(db)

		ctx := context.Background()
		created, err := repo.Create(ctx, nil)

		require.Error(t, err)
		assert.Nil(t, created)
	})

	t.Run("異常系：キャンセルされたコンテキストではエラー", func(t *testing.T) {
		cleanupTableForTag(t, db)
		repo := NewMySQLTagRepository(db)

		tag := createTestTag(t, "Go")

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		created, err := repo.Create(ctx, tag)

		require.Error(t, err)
		assert.Nil(t, created)
	})

	t.Run("異常系：重複した名前のタグを作成しようとするとエラー", func(t *testing.T) {
		cleanupTableForTag(t, db)
		repo := NewMySQLTagRepository(db)

		tag1 := createTestTag(t, "Go")
		tag2 := createTestTag(t, "Go")

		ctx := context.Background()
		_, err := repo.Create(ctx, tag1)
		require.NoError(t, err)

		created, err := repo.Create(ctx, tag2)
		require.Error(t, err)
		assert.Nil(t, created)
	})
}

func TestMySQLTageRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	t.Run("正常系：IDでタグを取得できる", func(t *testing.T) {
		cleanupTableForTag(t, db)
		repo := NewMySQLTagRepository(db)

		tag := createTestTag(t, "Go")
		id := insertTagDirectly(t, db, tag)

		ctx := context.Background()
		found, err := repo.FindByID(ctx, id)

		require.NoError(t, err)
		require.NotNil(t, found)
		assert.Equal(t, id, found.ID)
		assert.Equal(t, tag.Name, found.Name)
	})

	t.Run("異常系：存在しないIDで取得するとエラー", func(t *testing.T) {
		cleanupTableForTag(t, db)
		repo := NewMySQLTagRepository(db)

		ctx := context.Background()
		found, err := repo.FindByID(ctx, 99999)

		require.Error(t, err)
		assert.Nil(t, found)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("異常系：負のIDで取得するとエラー", func(t *testing.T) {
		cleanupTableForTag(t, db)
		repo := NewMySQLTagRepository(db)

		ctx := context.Background()
		found, err := repo.FindByID(ctx, -1)

		require.Error(t, err)
		assert.Nil(t, found)
	})

	t.Run("異常系：キャンセルされたコンテキストではエラー", func(t *testing.T) {
		cleanupTableForTag(t, db)
		repo := NewMySQLTagRepository(db)

		tag := createTestTag(t, "Go")
		id := insertTagDirectly(t, db, tag)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		found, err := repo.FindByID(ctx, id)

		require.Error(t, err)
		assert.Nil(t, found)
	})
}

func TestMySQLTagRepository_FindByName(t *testing.T) {
	db := setupTestDBForTag(t)
	defer db.Close()

	t.Run("正常系：名前でタグを取得できる", func(t *testing.T) {
		cleanupTableForTag(t, db)
		repo := NewMySQLTagRepository(db)

		tag := createTestTag(t, "Go")
		id := insertTagDirectly(t, db, tag)

		ctx := context.Background()
		found, err := repo.FindByName(ctx, "Go")

		require.NoError(t, err)
		require.NotNil(t, found)
		assert.Equal(t, id, found.ID)
		assert.Equal(t, "Go", found.Name)
	})

	t.Run("異常系：存在しない名前で取得するとエラー", func(t *testing.T) {
		cleanupTableForTag(t, db)
		repo := NewMySQLTagRepository(db)

		ctx := context.Background()
		found, err := repo.FindByName(ctx, "NonExistent")

		require.Error(t, err)
		assert.Nil(t, found)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("異常系：空文字列で取得するとエラー", func(t *testing.T) {
		cleanupTableForTag(t, db)
		repo := NewMySQLTagRepository(db)

		ctx := context.Background()
		found, err := repo.FindByName(ctx, "")

		require.Error(t, err)
		assert.Nil(t, found)
	})

	t.Run("異常系：キャンセルされたコンテキストではエラー", func(t *testing.T) {
		cleanupTableForTag(t, db)
		repo := NewMySQLTagRepository(db)

		tag := createTestTag(t, "Go")
		insertTagDirectly(t, db, tag)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		found, err := repo.FindByName(ctx, "Go")

		require.Error(t, err)
		assert.Nil(t, found)
	})
}

func TestMySQLTagRepository_FindAll(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	t.Run("正常系：複数のタグを取得できる（アルファベット順）", func(t *testing.T) {
		cleanupTableForTag(t, db)
		repo := NewMySQLTagRepository(db)

		tag1 := createTestTag(t, "Python")
		tag2 := createTestTag(t, "Go")
		tag3 := createTestTag(t, "Rust")

		insertTagDirectly(t, db, tag1)
		insertTagDirectly(t, db, tag2)
		insertTagDirectly(t, db, tag3)

		ctx := context.Background()
		tags, err := repo.FindAll(ctx)

		require.NoError(t, err)
		require.NotNil(t, tags)
		assert.Len(t, tags, 3)

		// アルファベット順で並んでいることを確認
		assert.Equal(t, "Go", tags[0].Name)
		assert.Equal(t, "Python", tags[1].Name)
		assert.Equal(t, "Rust", tags[2].Name)
	})

	t.Run("正常系：タグが0件の場合は空配列を返す", func(t *testing.T) {
		cleanupTableForTag(t, db)
		repo := NewMySQLTagRepository(db)

		ctx := context.Background()
		tags, err := repo.FindAll(ctx)

		require.NoError(t, err)
		require.NotNil(t, tags)
		assert.Empty(t, tags)
	})

	t.Run("正常系：1件のタグを取得できる", func(t *testing.T) {
		cleanupTableForTag(t, db)
		repo := NewMySQLTagRepository(db)

		tag := createTestTag(t, "Go")
		insertTagDirectly(t, db, tag)

		ctx := context.Background()
		tags, err := repo.FindAll(ctx)

		require.NoError(t, err)
		require.NotNil(t, tags)
		assert.Len(t, tags, 1)
		assert.Equal(t, "Go", tags[0].Name)
	})

	t.Run("異常系：キャンセルされたコンテキストではエラー", func(t *testing.T) {
		cleanupTableForTag(t, db)
		repo := NewMySQLTagRepository(db)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		tags, err := repo.FindAll(ctx)

		require.Error(t, err)
		assert.Nil(t, tags)
	})
}

func TestMySQLTagRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	t.Run("正常系：タグを更新できる", func(t *testing.T) {
		cleanupTableForTag(t, db)
		repo := NewMySQLTagRepository(db)

		tag := createTestTag(t, "Go")
		id := insertTagDirectly(t, db, tag)

		tag.ID = id
		time.Sleep(10 * time.Millisecond)
		err := tag.Update("Golang")
		require.NoError(t, err)

		ctx := context.Background()
		updated, err := repo.Update(ctx, tag)

		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, id, updated.ID)
		assert.Equal(t, "Golang", updated.Name)
	})

	t.Run("正常系：UpdatedAtが更新される", func(t *testing.T) {
		cleanupTableForTag(t, db)
		repo := NewMySQLTagRepository(db)

		tag := createTestTag(t, "Go")
		id := insertTagDirectly(t, db, tag)

		time.Sleep(100 * time.Millisecond)

		tag.ID = id
		originalUpdatedAt := tag.UpdatedAt
		err := tag.Update("Golang")
		require.NoError(t, err)

		ctx := context.Background()
		updated, err := repo.Update(ctx, tag)

		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.True(t, updated.UpdatedAt.After(originalUpdatedAt))
	})

	t.Run("異常系：存在しないIDのタグを更新しようとするとエラー", func(t *testing.T) {
		cleanupTableForTag(t, db)
		repo := NewMySQLTagRepository(db)

		tag := createTestTag(t, "Go")
		tag.ID = 99999

		ctx := context.Background()
		updated, err := repo.Update(ctx, tag)

		require.Error(t, err)
		assert.Nil(t, updated)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("異常系：IDが0のタグを更新しようとするとエラー", func(t *testing.T) {
		cleanupTableForTag(t, db)
		repo := NewMySQLTagRepository(db)

		tag := createTestTag(t, "Go")
		tag.ID = 0

		ctx := context.Background()
		updated, err := repo.Update(ctx, tag)

		require.Error(t, err)
		assert.Nil(t, updated)
	})

	t.Run("異常系：nilのタグを更新しようとするとエラー", func(t *testing.T) {
		cleanupTableForTag(t, db)
		repo := NewMySQLTagRepository(db)

		ctx := context.Background()
		updated, err := repo.Update(ctx, nil)

		require.Error(t, err)
		assert.Nil(t, updated)
	})

	t.Run("異常系：キャンセルされたコンテキストではエラー", func(t *testing.T) {
		cleanupTableForTag(t, db)
		repo := NewMySQLTagRepository(db)

		tag := createTestTag(t, "Go")
		id := insertTagDirectly(t, db, tag)
		tag.ID = id

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		updated, err := repo.Update(ctx, tag)

		require.Error(t, err)
		assert.Nil(t, updated)
	})

	t.Run("異常系：重複した名前に更新しようとするとエラー", func(t *testing.T) {
		cleanupTableForTag(t, db)
		repo := NewMySQLTagRepository(db)

		tag1 := createTestTag(t, "Go")
		tag2 := createTestTag(t, "Python")

		id1 := insertTagDirectly(t, db, tag1)
		insertTagDirectly(t, db, tag2)

		tag1.ID = id1
		err := tag1.Update("Python")
		require.NoError(t, err)

		ctx := context.Background()
		updated, err := repo.Update(ctx, tag1)

		require.Error(t, err)
		assert.Nil(t, updated)
	})
}

func TestMySQLTagRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	t.Run("正常系：タグを削除できる", func(t *testing.T) {
		cleanupTableForTag(t, db)
		repo := NewMySQLTagRepository(db)

		tag := createTestTag(t, "Go")
		id := insertTagDirectly(t, db, tag)

		ctx := context.Background()
		err := repo.Delete(ctx, id)

		require.NoError(t, err)

		found, err := repo.FindByID(ctx, id)
		require.Error(t, err)
		assert.Nil(t, found)
	})

	t.Run("正常系：タグを削除するとarticle_tagsも削除される", func(t *testing.T) {
		cleanupTableForTag(t, db)
		repo := NewMySQLTagRepository(db)

		articleIDs := ensureArticlesExist(t, db, 1)

		tag := createTestTag(t, "Go")
		tagID := insertTagDirectly(t, db, tag)

		_, err := db.Exec("INSERT INTO article_tags (article_id, tag_id) VALUES (?, ?)", articleIDs[0], tagID)
		require.NoError(t, err)

		ctx := context.Background()
		err = repo.Delete(ctx, tagID)

		require.NoError(t, err)

		// article_tagsテーブルからも削除されていることを確認
		var count int
		err = db.Get(&count, "SELECT COUNT(*) FROM article_tags WHERE tag_id = ?", tagID)
		require.NoError(t, err)
		assert.Equal(t, 0, count, "article_tagsが削除されていません")
	})

	t.Run("正常系：複数のタグのうち1つを削除できる", func(t *testing.T) {
		cleanupTableForTag(t, db)
		repo := NewMySQLTagRepository(db)

		tag1 := createTestTag(t, "Go")
		tag2 := createTestTag(t, "Python")

		id1 := insertTagDirectly(t, db, tag1)
		id2 := insertTagDirectly(t, db, tag2)

		ctx := context.Background()
		err := repo.Delete(ctx, id1)
		require.NoError(t, err)

		// id2のタグは残っていることを確認
		found, err := repo.FindByID(ctx, id2)
		require.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, id2, found.ID)
	})

	t.Run("異常系：存在しないIDのタグを削除しようとするとエラー", func(t *testing.T) {
		cleanupTableForTag(t, db)
		repo := NewMySQLTagRepository(db)

		ctx := context.Background()
		err := repo.Delete(ctx, 99999)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("異常系：IDが0のタグを削除しようとするとエラー", func(t *testing.T) {
		cleanupTableForTag(t, db)
		repo := NewMySQLTagRepository(db)

		ctx := context.Background()
		err := repo.Delete(ctx, 0)

		require.Error(t, err)
	})

	t.Run("異常系：負のIDのタグを削除しようとするとエラー", func(t *testing.T) {
		cleanupTableForTag(t, db)
		repo := NewMySQLTagRepository(db)

		ctx := context.Background()
		err := repo.Delete(ctx, -1)

		require.Error(t, err)
	})

	t.Run("異常系：キャンセルされたコンテキストではエラー", func(t *testing.T) {
		cleanupTableForTag(t, db)
		repo := NewMySQLTagRepository(db)

		tag := createTestTag(t, "Go")
		id := insertTagDirectly(t, db, tag)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := repo.Delete(ctx, id)

		require.Error(t, err)
	})

	t.Run("異常系：同じIDを2回削除しようとすると2回目はエラー", func(t *testing.T) {
		cleanupTableForTag(t, db)
		repo := NewMySQLTagRepository(db)

		tag := createTestTag(t, "Go")
		id := insertTagDirectly(t, db, tag)

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
