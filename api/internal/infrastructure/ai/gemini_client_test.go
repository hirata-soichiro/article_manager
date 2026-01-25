package ai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"article-manager/internal/domain/service"
)

func TestGeminiClient_GenerateArticleFromURL_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		response := geminiResponse{
			Candidates: []struct {
				Content struct {
					Parts []struct {
						Text string `json:"text"`
					} `json:"parts"`
				} `json:"content"`
				FinishReason string `json:"finishReason"`
			}{
				{
					Content: struct {
						Parts []struct {
							Text string `json:"text"`
						} `json:"parts"`
					}{
						Parts: []struct {
							Text string `json:"text"`
						}{
							{
								Text: `{
									"title": "Test Article",
									"summary": "This is a test summary",
									"suggestedTags": ["test", "article", "sample"]
								}`,
							},
						},
					},
					FinishReason: "STOP",
				},
			},
			UsageMetadata: struct {
				TotalTokenCount int `json:"totalTokenCount"`
			}{
				TotalTokenCount: 100,
			},
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := DefaultGeminiConfig("test-api-key")
	config.BaseURL = server.URL
	config.MaxRetries = 0
	client := NewGeminiClient(config)

	req := service.ArticleGenerationRequest{
		URL: "https://example.com/article",
	}

	ctx := context.Background()
	result, err := client.GenerateArticleFromURL(ctx, req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Title != "Test Article" {
		t.Errorf("Expected title 'Test Article', got '%s'", result.Title)
	}

	if result.Summary != "This is a test summary" {
		t.Errorf("Expected summary 'This is a test summary', got '%s'", result.Summary)
	}

	if len(result.SuggestedTags) != 3 {
		t.Errorf("Expected 3 tags, got %d", len(result.SuggestedTags))
	}

	if result.SourceURL != "https://example.com/article" {
		t.Errorf("Expected source URL 'https://example.com/article', got '%s'", result.SourceURL)
	}

	if result.TokenUsed != 100 {
		t.Errorf("Expected 100 token used, got %d", result.TokenUsed)
	}
}

func TestGeminiClient_GenerateArticleFromURL_EmptyURL(t *testing.T) {
	config := DefaultGeminiConfig("test-api-key")
	client := NewGeminiClient(config)

	req := service.ArticleGenerationRequest{
		URL: "",
	}

	ctx := context.Background()
	_, err := client.GenerateArticleFromURL(ctx, req)

	if err == nil {
		t.Fatalf("Expected error for empty URL, got nil")
	}

	aiErr, ok := err.(*service.AIGeneratorError)
	if !ok {
		t.Fatalf("Expected AIGeneratorError, got %T", err)
	}

	if aiErr.Code != service.ErrCodeInvalidURL {
		t.Errorf("Expected error code %s, got %s", service.ErrCodeInvalidURL, aiErr.Code)
	}
}

func TestGeminiClient_GenerateArticleFromURL_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(geminiErrorResponse{
			Error: struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
				Status  string `json:"status"`
			}{
				Code:    401,
				Message: "Invalid API key",
				Status:  "UNAUTHORIZED",
			},
		})
	}))
	defer server.Close()

	config := DefaultGeminiConfig("invalid-key")
	config.BaseURL = server.URL
	config.MaxRetries = 0
	client := NewGeminiClient(config)

	req := service.ArticleGenerationRequest{
		URL: "https://example.com/article",
	}

	ctx := context.Background()
	_, err := client.GenerateArticleFromURL(ctx, req)

	if err == nil {
		t.Fatalf("Expected error for unauthorized, got nil")
	}

	aiErr, ok := err.(*service.AIGeneratorError)
	if !ok {
		t.Fatalf("Expected AIGeneratorError, got %T", err)
	}

	if aiErr.Code != service.ErrCodeUnauthorized {
		t.Errorf("Expected error code %s, got %s", service.ErrCodeUnauthorized, aiErr.Code)
	}
}

func TestGeminiClient_GenerateArticleFromURL_RateLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(geminiErrorResponse{
			Error: struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
				Status  string `json:"status"`
			}{
				Code:    429,
				Message: "Rate limit exceeded",
				Status:  "RESOURCE_EXHAUSTED",
			},
		})
	}))
	defer server.Close()

	config := DefaultGeminiConfig("test-api-key")
	config.BaseURL = server.URL
	config.MaxRetries = 0
	client := NewGeminiClient(config)

	req := service.ArticleGenerationRequest{
		URL: "https://example.com/article",
	}

	ctx := context.Background()
	_, err := client.GenerateArticleFromURL(ctx, req)

	if err == nil {
		t.Fatalf("Expected error for rate limit, got nil")
	}

	aiErr, ok := err.(*service.AIGeneratorError)
	if !ok {
		t.Fatalf("Expected AIGeneratorError, got %T", err)
	}

	if aiErr.Code != service.ErrCodeAPILimit {
		t.Errorf("Expected error code %s, got %s", service.ErrCodeAPILimit, aiErr.Code)
	}
}

func TestGeminiClient_GenerateArticleFromURL_Retry(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(geminiErrorResponse{
				Error: struct {
					Code    int    `json:"code"`
					Message string `json:"message"`
					Status  string `json:"status"`
				}{
					Code:    429,
					Message: "Rate limit",
					Status:  "RESOURCE_EXHAUSTED",
				},
			})
			return
		}

		response := geminiResponse{
			Candidates: []struct {
				Content struct {
					Parts []struct {
						Text string `json:"text"`
					} `json:"parts"`
				} `json:"content"`
				FinishReason string `json:"finishReason"`
			}{
				{
					Content: struct {
						Parts []struct {
							Text string `json:"text"`
						} `json:"parts"`
					}{
						Parts: []struct {
							Text string `json:"text"`
						}{
							{
								Text: `{"title": "Success", "summary": "After retry", "suggestedTags": ["test"]}`,
							},
						},
					},
				},
			},
			UsageMetadata: struct {
				TotalTokenCount int `json:"totalTokenCount"`
			}{
				TotalTokenCount: 50,
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := DefaultGeminiConfig("test-api-key")
	config.BaseURL = server.URL
	config.MaxRetries = 3
	config.RetryWaitTime = 10 * time.Millisecond
	client := NewGeminiClient(config)

	req := service.ArticleGenerationRequest{
		URL: "https://example.com/article",
	}

	ctx := context.Background()
	result, err := client.GenerateArticleFromURL(ctx, req)

	if err != nil {
		t.Fatalf("Expected success after retry, got error: %v", err)
	}

	if result.Title != "Success" {
		t.Errorf("Expected title 'Success', got '%s'", result.Title)
	}

	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestGeminiClient_GenerateArticleFromURL_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Microsecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := DefaultGeminiConfig("invalid-key")
	config.BaseURL = server.URL
	config.Timeout = 50 * time.Millisecond
	config.MaxRetries = 0
	client := NewGeminiClient(config)

	req := service.ArticleGenerationRequest{
		URL: "https://example.com/article",
	}

	ctx := context.Background()
	_, err := client.GenerateArticleFromURL(ctx, req)

	if err == nil {
		t.Fatalf("Expected timeout error, got nil")
	}

	aiErr, ok := err.(*service.AIGeneratorError)
	if !ok {
		t.Fatalf("Expected AIGeneratorError, got %T", err)
	}

	if aiErr.Code != service.ErrCodeNetworkError {
		t.Errorf("Expected error code %s, got %s", service.ErrCodeNetworkError, aiErr.Code)
	}
}

func TestGeminiClient_GenerateArticleFromURL_ContentBlocked(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(geminiErrorResponse{
			Error: struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
				Status  string `json:"status"`
			}{
				Code:    400,
				Message: "Content blocked by SAFETY filters",
				Status:  "INVALID_ARGUMENT",
			},
		})
	}))
	defer server.Close()

	config := DefaultGeminiConfig("test-api-key")
	config.BaseURL = server.URL
	config.MaxRetries = 0
	client := NewGeminiClient(config)

	req := service.ArticleGenerationRequest{
		URL: "https://example.com/inappropriate",
	}

	ctx := context.Background()
	_, err := client.GenerateArticleFromURL(ctx, req)

	if err == nil {
		t.Fatalf("Expected content blocked error, got nil")
	}

	aiErr, ok := err.(*service.AIGeneratorError)
	if !ok {
		t.Fatalf("Expected AIGeneratorError, got %T", err)
	}

	if aiErr.Code != service.ErrCodeContentBlocked {
		t.Errorf("Expected error code %s, got %s", service.ErrCodeContentBlocked, aiErr.Code)
	}
}
