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

// POST /api/articlesのテスト
func TestCreateArticle(t *testing.T) {
	t.Run("正常系：記事を作成できる", func(t *testing.T) {
		handler := setupHandler()

		requestBody := map[string]interface{}{
			"title":   "Go言語入門",
			"url":     "https://example.com/go-intro",
			"summary": "Go言語の基本的な使い方を解説します",
			"tags":    []string{"Go", "プログラミング"},
			"memo":    "初心者向けの記事",
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/api/articles", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.CreateArticle(rec, req)

		require.Equal(t, http.StatusCreated, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Go言語入門", response["title"])
		assert.Equal(t, "https://example.com/go-intro", response["url"])
		assert.NotNil(t, response["id"])
		assert.NotNil(t, response["created_at"])
	})

	t.Run("異常系：タイトルが空の場合", func(t *testing.T) {
		handler := setupHandler()

		requestBody := map[string]interface{}{
			"title":   "",
			"url":     "https://example.com/test",
			"summary": "テスト記事",
			"tags":    []string{"Test"},
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/api/articles", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.CreateArticle(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "validation failed")
	})

	t.Run("異常系：URLが空の場合", func(t *testing.T) {
		handler := setupHandler()

		requestBody := map[string]interface{}{
			"title":   "テスト記事",
			"url":     "",
			"summary": "テスト記事の要約",
			"tags":    []string{"Test"},
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/api/articles", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.CreateArticle(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "validation failed")
	})

	t.Run("異常系：URLが不正な形式の場合", func(t *testing.T) {
		handler := setupHandler()

		requestBody := map[string]interface{}{
			"title":   "テスト記事",
			"url":     "invalid-url",
			"summary": "テスト記事の要約",
			"tags":    []string{"Test"},
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/api/articles", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.CreateArticle(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "validation failed")
	})

	t.Run("異常系：要約が空の場合", func(t *testing.T) {
		handler := setupHandler()

		requestBody := map[string]interface{}{
			"title":   "テスト記事",
			"url":     "https://example.com/test",
			"summary": "",
			"tags":    []string{"Test"},
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/api/articles", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.CreateArticle(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "validation failed")
	})

	t.Run("異常系：不正なJSON", func(t *testing.T) {
		handler := setupHandler()

		invalidJSON := []byte(`{invalid json}`)

		req := httptest.NewRequest(http.MethodPost, "/api/articles", bytes.NewReader(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.CreateArticle(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})
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
		assert.Contains(t, response["error"], "invalid argument")
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

// GET /api/articles/searchのテスト
func TestSearchArticles(t *testing.T) {
	t.Run("正常系：キーワードで検索できる", func(t *testing.T) {
		handler := setupHandler()

		ctx := context.Background()
		handler.usecase.CreateArticle(ctx, "Go言語入門", "https://example.com/1", "Go言語の基本を解説", []string{"Go", "プログラミング"}, "初心者向け")
		handler.usecase.CreateArticle(ctx, "Next.js入門", "https://example.com/2", "Next.jsの使い方", []string{"Next.js", "React"}, "フロントエンド")
		handler.usecase.CreateArticle(ctx, "GoとNext.jsで作るアプリ", "https://example.com/3", "GoとNext.jsを組み合わせた開発", []string{"Go", "Next.js"}, "フルスタック")

		req := httptest.NewRequest(http.MethodGet, "/api/articles/search?keyword=Go", nil)
		rec := httptest.NewRecorder()

		handler.SearchArticles(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		var response []map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, 2, len(response))
		titles := []string{response[0]["title"].(string), response[1]["title"].(string)}
		assert.Contains(t, titles, "Go言語入門")
		assert.Contains(t, titles, "GoとNext.jsで作るアプリ")
	})

	t.Run("正常系：AND検索（複数キーワードをスペース区切り）", func(t *testing.T) {
		handler := setupHandler()

		ctx := context.Background()
		handler.usecase.CreateArticle(ctx, "Go言語入門", "https://example.com/1", "Go言語の基本を解説", []string{"Go"}, "初心者向け")
		handler.usecase.CreateArticle(ctx, "Next.js入門", "https://example.com/2", "Next.jsの使い方", []string{"Next.js"}, "フロントエンド")
		handler.usecase.CreateArticle(ctx, "GoとNext.jsで作るアプリ", "https://example.com/3", "GoとNext.jsを組み合わせた開発", []string{"Go", "Next.js"}, "フルスタック開発")

		req := httptest.NewRequest(http.MethodGet, "/api/articles/search?keyword=Go+Next.js", nil)
		rec := httptest.NewRecorder()

		handler.SearchArticles(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		var response []map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, 1, len(response))
		assert.Equal(t, "GoとNext.jsで作るアプリ", response[0]["title"])
	})

	t.Run("正常系：タイトルと要約から検索", func(t *testing.T) {
		handler := setupHandler()

		ctx := context.Background()
		handler.usecase.CreateArticle(ctx, "データベース設計", "https://example.com/1", "MySQL設計の基礎", []string{"Database"}, "")
		handler.usecase.CreateArticle(ctx, "API開発", "https://example.com/2", "RESTful APIの実装", []string{"Go", "API"}, "")

		// タイトルから検索
		req1 := httptest.NewRequest(http.MethodGet, "/api/articles/search?keyword=データベース", nil)
		rec1 := httptest.NewRecorder()
		handler.SearchArticles(rec1, req1)
		var response1 []map[string]interface{}
		json.Unmarshal(rec1.Body.Bytes(), &response1)
		assert.Equal(t, 1, len(response1))
		assert.Equal(t, "データベース設計", response1[0]["title"])

		// 要約から検索
		req2 := httptest.NewRequest(http.MethodGet, "/api/articles/search?keyword=RESTful", nil)
		rec2 := httptest.NewRecorder()
		handler.SearchArticles(rec2, req2)
		var response2 []map[string]interface{}
		json.Unmarshal(rec2.Body.Bytes(), &response2)
		assert.Equal(t, 1, len(response2))
		assert.Equal(t, "API開発", response2[0]["title"])
	})

	t.Run("正常系：検索結果が0件の場合", func(t *testing.T) {
		handler := setupHandler()

		ctx := context.Background()
		handler.usecase.CreateArticle(ctx, "Go言語入門", "https://example.com/1", "Go言語の基本", []string{"Go"}, "")

		req := httptest.NewRequest(http.MethodGet, "/api/articles/search?keyword=Python", nil)
		rec := httptest.NewRecorder()

		handler.SearchArticles(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		var response []map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, 0, len(response))
	})

	t.Run("正常系：記事が存在しない場合", func(t *testing.T) {
		handler := setupHandler()

		req := httptest.NewRequest(http.MethodGet, "/api/articles/search?keyword=Go", nil)
		rec := httptest.NewRecorder()

		handler.SearchArticles(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		var response []map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, 0, len(response))
	})

	t.Run("正常系：キーワードに余分なスペースがある場合", func(t *testing.T) {
		handler := setupHandler()

		ctx := context.Background()
		handler.usecase.CreateArticle(ctx, "Go言語入門", "https://example.com/1", "Go言語の基本", []string{"Go"}, "")

		req := httptest.NewRequest(http.MethodGet, "/api/articles/search?keyword=+Go+", nil)
		rec := httptest.NewRecorder()

		handler.SearchArticles(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		var response []map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, 1, len(response))
	})

	t.Run("正常系：大文字小文字を区別しない検索", func(t *testing.T) {
		handler := setupHandler()

		ctx := context.Background()
		handler.usecase.CreateArticle(ctx, "Go言語入門", "https://example.com/1", "go言語の基本", []string{"golang"}, "")

		// 小文字で検索
		req1 := httptest.NewRequest(http.MethodGet, "/api/articles/search?keyword=go", nil)
		rec1 := httptest.NewRecorder()
		handler.SearchArticles(rec1, req1)
		var response1 []map[string]interface{}
		json.Unmarshal(rec1.Body.Bytes(), &response1)
		assert.Equal(t, 1, len(response1))

		// 大文字で検索
		req2 := httptest.NewRequest(http.MethodGet, "/api/articles/search?keyword=GO", nil)
		rec2 := httptest.NewRecorder()
		handler.SearchArticles(rec2, req2)
		var response2 []map[string]interface{}
		json.Unmarshal(rec2.Body.Bytes(), &response2)
		assert.Equal(t, 1, len(response2))
	})

	t.Run("異常系：キーワードパラメータが空", func(t *testing.T) {
		handler := setupHandler()

		req := httptest.NewRequest(http.MethodGet, "/api/articles/search?keyword=", nil)
		rec := httptest.NewRecorder()

		handler.SearchArticles(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "keyword")
	})

	t.Run("異常系：キーワードパラメータが存在しない", func(t *testing.T) {
		handler := setupHandler()

		req := httptest.NewRequest(http.MethodGet, "/api/articles/search", nil)
		rec := httptest.NewRecorder()

		handler.SearchArticles(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "keyword")
	})

	t.Run("異常系：キーワードがスペースのみ", func(t *testing.T) {
		handler := setupHandler()

		req := httptest.NewRequest(http.MethodGet, "/api/articles/search?keyword=+++", nil)
		rec := httptest.NewRecorder()

		handler.SearchArticles(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "keyword")
	})
}
