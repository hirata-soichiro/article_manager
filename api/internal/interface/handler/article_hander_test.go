package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"article-manager/internal/infrastructure/repository"
	"article-manager/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// テスト用のハンドラのセットアップ
func setupHandler() *ArticleHandler {
	repo := repository.NewMemoryArticleRepository()
	uc := usecase.NewArticleUsecase(repo)
	return NewArticleHandler(uc)
}

// GET /api/articlesのテスト
func TestGetAllArticles(t *testing.T) {
	t.Run("正常系：全記事を取得できる", func(t *testing.T) {
		handler := setupHandler()

		// テストデータ作成
		ctx := context.Background()
		handler.usecase.CreateArticle(ctx, "記事1", "https://example.com/1", "要約1", []string{"Go"}, "メモ1")
		handler.usecase.CreateArticle(ctx, "記事2", "https://example.com/2", "要約2", []string{"Next.js"}, "メモ2")

		// リクエスト作成
		req := httptest.NewRequest(http.MethodGet, "/api/articles", nil)
		rec := httptest.NewRecorder()

		// ハンドラ実行
		handler.GetAllArticles(rec, req)

		// レスポンス検証
		require.Equal(t, http.StatusOK, rec.Code)

		var response []map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, 2, len(response))
		assert.Equal(t, "記事1", response[0]["title"])
		assert.Equal(t, "記事2", response[1]["title"])
	})

	t.Run("正常系：記事が0件の場合", func(t *testing.T) {
		handler := setupHandler()

		req := httptest.NewRequest(http.MethodGet, "/api/articles", nil)
		rec := httptest.NewRecorder()

		handler.GetAllArticles(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		var response []map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, 0, len(response))
	})
}

// GET /api/articles/{id}のテスト
func TestGetArticleByID(t *testing.T) {
	t.Run("正常系：IDで記事を取得できる", func(t *testing.T) {
		handler := setupHandler()

		// テストデータ作成
		ctx := context.Background()
		created, _ := handler.usecase.CreateArticle(ctx, "テスト記事", "https://example.com", "テスト要約", []string{"Go"}, "メモ")

		// リクエスト作成
		req := httptest.NewRequest(http.MethodGet, "/api/articles/1", nil)
		rec := httptest.NewRecorder()

		// ハンドラ実行
		handler.GetArticleByID(rec, req, created.ID)

		// レスポンス検証
		require.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "テスト記事", response["title"])
		assert.Equal(t, "https://example.com", response["url"])
	})

	t.Run("異常系：存在しないIDの場合", func(t *testing.T) {
		handler := setupHandler()

		req := httptest.NewRequest(http.MethodGet, "/api/articles/999", nil)
		rec := httptest.NewRecorder()

		handler.GetArticleByID(rec, req, 999)

		require.Equal(t, http.StatusNotFound, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "not found")
	})

	t.Run("異常系：不正なID(0以下)", func(t *testing.T) {
		handler := setupHandler()

		req := httptest.NewRequest(http.MethodGet, "/api/articles/0", nil)
		rec := httptest.NewRecorder()

		handler.GetArticleByID(rec, req, 0)

		require.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "invalid id")
	})
}

// PUT /api/articles/:idのテスト
func TestUpdateArticle(t *testing.T) {
	t.Run("正常系：記事を更新できる", func(t *testing.T) {
		handler := setupHandler()

		// テストデータ作成
		ctx := context.Background()
		created, _ := handler.usecase.CreateArticle(ctx, "旧タイトル", "https://example.com/old", "旧要約", []string{"Old"}, "旧メモ")

		requestBody := map[string]interface{}{
			"title":   "新タイトル",
			"url":     "https://example.com/new",
			"summary": "新要約",
			"tags":    []string{"New", "Go"},
			"memo":    "新メモ",
		}
		body, _ := json.Marshal(requestBody)

		// リクエスト作成
		req := httptest.NewRequest(http.MethodPut, "/api/articles/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		// ハンドラ実行
		handler.UpdateArticle(rec, req, created.ID)

		// レスポンス検証
		require.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "新タイトル", response["title"])
		assert.Equal(t, "https://example.com/new", response["url"])
	})

	t.Run("異常系：存在しないIDの場合", func(t *testing.T) {
		handler := setupHandler()

		requestBody := map[string]interface{}{
			"title":   "更新タイトル",
			"url":     "https://example.com",
			"summary": "更新要約",
			"tags":    []string{"Go"},
			"memo":    "",
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPut, "/api/articles/999", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.UpdateArticle(rec, req, 999)

		require.Equal(t, http.StatusNotFound, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "not found")
	})

	t.Run("異常系：バリデーションエラー", func(t *testing.T) {
		handler := setupHandler()

		ctx := context.Background()
		created, _ := handler.usecase.CreateArticle(ctx, "旧タイトル", "https://example.com/old", "旧要約", []string{"Old"}, "")

		// 不正なデータ(タイトルが空)
		requestBody := map[string]interface{}{
			"title":   "",
			"url":     "https://example.com",
			"summary": "要約",
			"tags":    []string{},
			"memo":    "",
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPut, "/api/articles/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.UpdateArticle(rec, req, created.ID)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

// DELETE /api/articles/:idのテスト
func TestDeleteArticle(t *testing.T) {
	t.Run("正常系：記事を削除できる", func(t *testing.T) {
		handler := setupHandler()

		// テストデータ作成
		ctx := context.Background()
		created, _ := handler.usecase.CreateArticle(ctx, "削除対象", "https://example.com", "要約", []string{"Go"}, "")

		// リクエスト作成
		req := httptest.NewRequest(http.MethodDelete, "/api/articles/1", nil)
		rec := httptest.NewRecorder()

		// ハンドラ実行
		handler.DeleteArticle(rec, req, created.ID)

		// レスポンス検証
		require.Equal(t, http.StatusNoContent, rec.Code)

		// 削除確認
		article, err := handler.usecase.GetArticleByID(ctx, created.ID)
		assert.Error(t, err)
		assert.Nil(t, article)
	})

	t.Run("異常系：存在しないIDの場合", func(t *testing.T) {
		handler := setupHandler()

		req := httptest.NewRequest(http.MethodDelete, "/api/articles/999", nil)
		rec := httptest.NewRecorder()

		handler.DeleteArticle(rec, req, 999)

		require.Equal(t, http.StatusNotFound, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "not found")
	})

	t.Run("異常系：不正なID", func(t *testing.T) {
		handler := setupHandler()

		req := httptest.NewRequest(http.MethodDelete, "/api/articles/0", nil)
		rec := httptest.NewRecorder()

		handler.DeleteArticle(rec, req, 0)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
