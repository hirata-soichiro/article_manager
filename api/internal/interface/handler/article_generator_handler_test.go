package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"article-manager/internal/domain/service"
	"article-manager/internal/infrastructure/repository"
	"article-manager/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// モックAIGeneratorService
type mockAIGeneratorService struct {
	generateFunc func(ctx context.Context, req service.ArticleGenerationRequest) (*service.GeneratedArticle, error)
}

func (m *mockAIGeneratorService) GenerateArticleFromURL(ctx context.Context, req service.ArticleGenerationRequest) (*service.GeneratedArticle, error) {
	if m.generateFunc != nil {
		return m.generateFunc(ctx, req)
	}
	return nil, nil
}

// テスト用ハンドラのセットアップ
func setupGeneratorHandler(aiService service.AIGeneratorService) *ArticleGeneratorHandler {
	articleRepo := repository.NewMemoryArticleRepository()
	tagRepo := repository.NewMemoryTagRepository()
	generatorUsecase := usecase.NewArticleGeneratorUsecase(aiService, articleRepo, tagRepo)
	return NewArticleGeneratorHandler(generatorUsecase)
}

// POST /api/articles/generateのテスト
func TestGenerateArticle(t *testing.T) {
	t.Run("正常系：URLから記事を生成して保存成功", func(t *testing.T) {
		mockAI := &mockAIGeneratorService{
			generateFunc: func(ctx context.Context, req service.ArticleGenerationRequest) (*service.GeneratedArticle, error) {
				return &service.GeneratedArticle{
					Title:         "AI生成記事タイトル",
					Summary:       "AI生成記事の要約です。",
					SuggestedTags: []string{"Go", "AI"},
					SourceURL:     req.URL,
					TokenUsed:     100,
					GeneratedAt:   time.Now(),
				}, nil
			},
		}
		handler := setupGeneratorHandler(mockAI)

		requestBody := map[string]interface{}{
			"url":  "https://example.com/article",
			"memo": "テストメモ",
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/api/articles/generate", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.GenerateArticle(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "AI生成記事タイトル", response["title"])
		assert.Equal(t, "https://example.com/article", response["url"])
		assert.Equal(t, "AI生成記事の要約です。", response["summary"])
		assert.Equal(t, "テストメモ", response["memo"])
		assert.NotNil(t, response["id"])
		assert.NotNil(t, response["created_at"])
	})

	t.Run("正常系：memoなしで記事生成成功", func(t *testing.T) {
		mockAI := &mockAIGeneratorService{
			generateFunc: func(ctx context.Context, req service.ArticleGenerationRequest) (*service.GeneratedArticle, error) {
				return &service.GeneratedArticle{
					Title:         "技術記事",
					Summary:       "技術的な内容の解説です。",
					SuggestedTags: []string{"Tech"},
					SourceURL:     req.URL,
					TokenUsed:     80,
					GeneratedAt:   time.Now(),
				}, nil
			},
		}
		handler := setupGeneratorHandler(mockAI)

		requestBody := map[string]interface{}{
			"url": "https://example.com/tech",
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/api/articles/generate", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.GenerateArticle(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "技術記事", response["title"])
		assert.Equal(t, "", response["memo"])
	})

	t.Run("異常系：URLが空の場合", func(t *testing.T) {
		mockAI := &mockAIGeneratorService{}
		handler := setupGeneratorHandler(mockAI)

		requestBody := map[string]interface{}{
			"url":  "",
			"memo": "メモ",
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/api/articles/generate", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.GenerateArticle(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "url")
	})

	t.Run("異常系：URLが不正な形式の場合", func(t *testing.T) {
		mockAI := &mockAIGeneratorService{}
		handler := setupGeneratorHandler(mockAI)

		requestBody := map[string]interface{}{
			"url":  "invalid-url",
			"memo": "",
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/api/articles/generate", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.GenerateArticle(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "url")
	})

	t.Run("異常系：AI生成サービスがエラーを返す（API制限）", func(t *testing.T) {
		mockAI := &mockAIGeneratorService{
			generateFunc: func(ctx context.Context, req service.ArticleGenerationRequest) (*service.GeneratedArticle, error) {
				return nil, &service.AIGeneratorError{
					Code:    service.ErrCodeAPILimit,
					Message: "API rate limit exceeded",
				}
			},
		}
		handler := setupGeneratorHandler(mockAI)

		requestBody := map[string]interface{}{
			"url": "https://example.com/article",
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/api/articles/generate", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.GenerateArticle(rec, req)

		require.Equal(t, http.StatusTooManyRequests, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "rate limit")
	})

	t.Run("異常系：AI生成サービスがエラーを返す（タイムアウト）", func(t *testing.T) {
		mockAI := &mockAIGeneratorService{
			generateFunc: func(ctx context.Context, req service.ArticleGenerationRequest) (*service.GeneratedArticle, error) {
				return nil, &service.AIGeneratorError{
					Code:    service.ErrCodeTimeout,
					Message: "request timeout",
				}
			},
		}
		handler := setupGeneratorHandler(mockAI)

		requestBody := map[string]interface{}{
			"url": "https://example.com/article",
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/api/articles/generate", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.GenerateArticle(rec, req)

		require.Equal(t, http.StatusGatewayTimeout, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "timeout")
	})

	t.Run("異常系：AI生成サービスがエラーを返す（不正なレスポンス）", func(t *testing.T) {
		mockAI := &mockAIGeneratorService{
			generateFunc: func(ctx context.Context, req service.ArticleGenerationRequest) (*service.GeneratedArticle, error) {
				return nil, &service.AIGeneratorError{
					Code:    service.ErrCodeInvalidResponse,
					Message: "invalid response format",
				}
			},
		}
		handler := setupGeneratorHandler(mockAI)

		requestBody := map[string]interface{}{
			"url": "https://example.com/article",
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/api/articles/generate", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.GenerateArticle(rec, req)

		require.Equal(t, http.StatusBadGateway, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "invalid response")
	})

	t.Run("異常系：AI生成サービスがエラーを返す（認証エラー）", func(t *testing.T) {
		mockAI := &mockAIGeneratorService{
			generateFunc: func(ctx context.Context, req service.ArticleGenerationRequest) (*service.GeneratedArticle, error) {
				return nil, &service.AIGeneratorError{
					Code:    service.ErrCodeUnauthorized,
					Message: "Invalid API key",
				}
			},
		}
		handler := setupGeneratorHandler(mockAI)

		requestBody := map[string]interface{}{
			"url": "https://example.com/article",
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/api/articles/generate", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.GenerateArticle(rec, req)

		require.Equal(t, http.StatusUnauthorized, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "API key")
	})

	t.Run("異常系：AI生成サービスがエラーを返す（コンテンツブロック）", func(t *testing.T) {
		mockAI := &mockAIGeneratorService{
			generateFunc: func(ctx context.Context, req service.ArticleGenerationRequest) (*service.GeneratedArticle, error) {
				return nil, &service.AIGeneratorError{
					Code:    service.ErrCodeContentBlocked,
					Message: "Content blocked by safety filters",
				}
			},
		}
		handler := setupGeneratorHandler(mockAI)

		requestBody := map[string]interface{}{
			"url": "https://example.com/article",
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/api/articles/generate", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.GenerateArticle(rec, req)

		require.Equal(t, http.StatusForbidden, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "blocked")
	})

	t.Run("異常系：不正なJSON", func(t *testing.T) {
		mockAI := &mockAIGeneratorService{}
		handler := setupGeneratorHandler(mockAI)

		invalidJSON := []byte(`{invalid json}`)

		req := httptest.NewRequest(http.MethodPost, "/api/articles/generate", bytes.NewReader(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.GenerateArticle(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "request")
	})

	t.Run("異常系：ネットワークエラー", func(t *testing.T) {
		mockAI := &mockAIGeneratorService{
			generateFunc: func(ctx context.Context, req service.ArticleGenerationRequest) (*service.GeneratedArticle, error) {
				return nil, &service.AIGeneratorError{
					Code:    service.ErrCodeNetworkError,
					Message: "Network connection failed",
				}
			},
		}
		handler := setupGeneratorHandler(mockAI)

		requestBody := map[string]interface{}{
			"url": "https://example.com/article",
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/api/articles/generate", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.GenerateArticle(rec, req)

		require.Equal(t, http.StatusBadGateway, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "Network")
	})
}
