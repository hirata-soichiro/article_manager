package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"article-manager/internal/infrastructure/repository"
	"article-manager/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// テスト用のハンドラのセットアップ
func setupTagHandler() *TagHandler {
	repo := repository.NewMemoryTagRepository()
	uc := usecase.NewTagUsecase(repo)
	return NewTagHandler(uc)
}

// POST /api/tagsのテスト
func TestCreateTag(t *testing.T) {
	t.Run("正常系：タグを作成できる", func(t *testing.T) {
		handler := setupTagHandler()

		requestBody := map[string]interface{}{
			"name": "Go",
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/api/tags", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.CreateTag(rec, req)

		require.Equal(t, http.StatusCreated, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Go", response["name"])
		assert.NotNil(t, response["id"])
		assert.NotNil(t, response["created_at"])
	})

	t.Run("異常系：名前が空の場合", func(t *testing.T) {
		handler := setupTagHandler()

		requestBody := map[string]interface{}{
			"name": "",
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/api/tags", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.CreateTag(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "name is required")
	})

	t.Run("異常系：名前が長すぎる場合", func(t *testing.T) {
		handler := setupTagHandler()

		longName := strings.Repeat("あ", 51)
		requestBody := map[string]interface{}{
			"name": longName,
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/api/tags", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.CreateTag(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "50 characters or less")
	})

	t.Run("異常系：不正なJSON", func(t *testing.T) {
		handler := setupTagHandler()

		invalidJSON := []byte(`{invalid json}`)

		req := httptest.NewRequest(http.MethodPost, "/api/tags", bytes.NewReader(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.CreateTag(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

// GET /api/tagsのテスト
func TestGetAllTags(t *testing.T) {
	t.Run("正常系：全タグを取得できる", func(t *testing.T) {
		handler := setupTagHandler()

		// テストデータ作成
		ctx := context.Background()
		handler.usecase.CreateTag(ctx, "Go")
		handler.usecase.CreateTag(ctx, "Next.js")

		// リクエスト作成
		req := httptest.NewRequest(http.MethodGet, "/api/tags", nil)
		rec := httptest.NewRecorder()

		// ハンドラ実行
		handler.GetAllTags(rec, req)

		// レスポンス検証
		require.Equal(t, http.StatusOK, rec.Code)

		var response []map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, 2, len(response))
		assert.Equal(t, "Go", response[0]["name"])
		assert.Equal(t, "Next.js", response[1]["name"])
	})

	t.Run("正常系：タグが0件の場合", func(t *testing.T) {
		handler := setupTagHandler()

		req := httptest.NewRequest(http.MethodGet, "/api/tags", nil)
		rec := httptest.NewRecorder()

		handler.GetAllTags(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		var response []map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, 0, len(response))
	})
}

// GET /api/tags/{id}のテスト
func TestGetTagByID(t *testing.T) {
	t.Run("正常系：IDでタグを取得できる", func(t *testing.T) {
		handler := setupTagHandler()

		// テストデータ作成
		ctx := context.Background()
		created, _ := handler.usecase.CreateTag(ctx, "Go")

		// リクエスト作成
		req := httptest.NewRequest(http.MethodGet, "/api/tags/1", nil)
		rec := httptest.NewRecorder()

		// ハンドラ実行
		handler.GetTagByID(rec, req, created.ID)

		// レスポンス検証
		require.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Go", response["name"])
		assert.NotNil(t, response["id"])
	})

	t.Run("異常系：存在しないIDの場合", func(t *testing.T) {
		handler := setupTagHandler()

		req := httptest.NewRequest(http.MethodGet, "/api/tags/999", nil)
		rec := httptest.NewRecorder()

		handler.GetTagByID(rec, req, 999)

		require.Equal(t, http.StatusNotFound, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "not found")
	})

	t.Run("異常系：不正なID(0以下)", func(t *testing.T) {
		handler := setupTagHandler()

		req := httptest.NewRequest(http.MethodGet, "/api/tags/0", nil)
		rec := httptest.NewRecorder()

		handler.GetTagByID(rec, req, 0)

		require.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "invalid id")
	})
}

// PUT /api/tags/:idのテスト
func TestUpdateTag(t *testing.T) {
	t.Run("正常系：タグを更新できる", func(t *testing.T) {
		handler := setupTagHandler()

		// テストデータ作成
		ctx := context.Background()
		created, _ := handler.usecase.CreateTag(ctx, "旧タグ名")

		requestBody := map[string]interface{}{
			"name": "NewName",
		}
		body, _ := json.Marshal(requestBody)

		// リクエスト作成
		req := httptest.NewRequest(http.MethodPut, "/api/tags/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		// ハンドラ実行
		handler.UpdateTag(rec, req, created.ID)

		// レスポンス検証
		require.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "NewName", response["name"])
	})

	t.Run("異常系：存在しないIDの場合", func(t *testing.T) {
		handler := setupTagHandler()

		requestBody := map[string]interface{}{
			"name": "UpdatedName",
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPut, "/api/tags/999", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.UpdateTag(rec, req, 999)

		require.Equal(t, http.StatusNotFound, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "not found")
	})

	t.Run("異常系：バリデーションエラー", func(t *testing.T) {
		handler := setupTagHandler()

		ctx := context.Background()
		created, _ := handler.usecase.CreateTag(ctx, "旧タイトル")

		// 不正なデータ(タイトルが空)
		requestBody := map[string]interface{}{
			"name": "",
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPut, "/api/tags/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.UpdateTag(rec, req, created.ID)

		require.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "name is required")
	})

	t.Run("異常系：不正なID", func(t *testing.T) {
		handler := setupTagHandler()

		requestBody := map[string]interface{}{
			"name": "NewName",
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPut, "/api/tags/0", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.UpdateTag(rec, req, 0)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

// DELETE /api/tags/:idのテスト
func TestDeleteTag(t *testing.T) {
	t.Run("正常系：記事を削除できる", func(t *testing.T) {
		handler := setupTagHandler()

		// テストデータ作成
		ctx := context.Background()
		created, _ := handler.usecase.CreateTag(ctx, "DeleteTarget")

		// リクエスト作成
		req := httptest.NewRequest(http.MethodDelete, "/api/tags/1", nil)
		rec := httptest.NewRecorder()

		// ハンドラ実行
		handler.DeleteTag(rec, req, created.ID)

		// レスポンス検証
		require.Equal(t, http.StatusNoContent, rec.Code)

		// 削除確認
		tag, err := handler.usecase.GetTagByID(ctx, created.ID)
		assert.Error(t, err)
		assert.Nil(t, tag)
	})

	t.Run("異常系：存在しないIDの場合", func(t *testing.T) {
		handler := setupTagHandler()

		req := httptest.NewRequest(http.MethodDelete, "/api/tags/999", nil)
		rec := httptest.NewRecorder()

		handler.DeleteTag(rec, req, 999)

		require.Equal(t, http.StatusNotFound, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "not found")
	})

	t.Run("異常系：不正なID", func(t *testing.T) {
		handler := setupTagHandler()

		req := httptest.NewRequest(http.MethodDelete, "/api/tags/0", nil)
		rec := httptest.NewRecorder()

		handler.DeleteTag(rec, req, 0)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
