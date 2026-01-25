package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"article-manager/internal/domain/service"
)

// Gemini API設定
type GeminiConfig struct {
	APIKey        string
	Model         string
	BaseURL       string
	Timeout       time.Duration
	MaxRetries    int
	RetryWaitTime time.Duration
}

// デフォルトGemini API設定
func DefaultGeminiConfig(apiKey string) *GeminiConfig {
	return &GeminiConfig{
		APIKey:        apiKey,
		Model:         "gemini-2.5-pro",
		BaseURL:       "https://api.example.com/v1/gemini/generate",
		Timeout:       15 * time.Second,
		MaxRetries:    3,
		RetryWaitTime: 2 * time.Second,
	}
}

// Gemini APIクライアント
type GeminiClient struct {
	config     *GeminiConfig
	httpClient *http.Client
}

// 新しいクライアントを作成
func NewGeminiClient(config *GeminiConfig) *GeminiClient {
	return &GeminiClient{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// Gemini APIリクエスト構造
type geminiRequest struct {
	Contents         []geminiContent         `json:"contents"`
	Tools            []geminiTool            `json:"tools,omitempty"`
	GenerationConfig *geminiGenerationConfig `json:"generationConfig"`
}

type geminiContent struct {
	Role  string       `json:"role"`
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiTool struct {
	URLContext *geminiURLContext `json:"url_context,omitempty"`
}

type geminiURLContext struct{}

type geminiGenerationConfig struct {
	Temperature     float32 `json:"temperature,omitempty"`
	MaxOutputTokens int     `json:"maxOutputTokens,omitempty"`
}

// Gemini APIレスポンス構造
type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
		FinishReason string `json:"finishReason"`
	} `json:"candidates"`
	UsageMetadata struct {
		TotalTokenCount int `json:"totalTokenCount"`
	} `json:"usageMetadata"`
}

type geminiErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

// URLから記事を生成
func (c *GeminiClient) GenerateArticleFromURL(ctx context.Context, req service.ArticleGenerationRequest) (*service.GeneratedArticle, error) {
	if req.URL == "" {
		return nil, &service.AIGeneratorError{
			Code:    service.ErrCodeInvalidURL,
			Message: "URL is empty",
		}
	}

	prompt := c.buildPrompt(req.URL)
	response, err := c.callAPI(ctx, prompt)
	if err != nil {
		return nil, err
	}

	article, err := c.parseResponse(response, req.URL)
	if err != nil {
		return nil, &service.AIGeneratorError{
			Code:    service.ErrCodeInvalidResponse,
			Message: "Failed to parse response",
			Err:     err,
		}
	}

	article.GeneratedAt = time.Now()
	return article, nil
}

// プロンプト構築
func (c *GeminiClient) buildPrompt(url string) string {
	return fmt.Sprintf(`以下のURLの記事を分析し、記事管理用の情報を生成してください。

URL: %s

要件:
- summaryは記事の核心を捉え、200文字以内にまとめる
- suggestedTagsは3-5個、検索しやすく具体的なものを選ぶ
- 技術記事の場合は技術スタックをタグに含める
- 日本語記事は日本語で、英語記事は翻訳し日本語で出力する

以下のJSON形式のみで出力してください。説明文やマークダウンのコードブロックは不要です。
{
  "title": "記事のタイトル",
  "summary": "記事の要約",
  "suggestedTags": ["タグ1", "タグ2", "タグ3"]
}`, url)
}

// Gemini API呼び出し
func (c *GeminiClient) callAPI(ctx context.Context, prompt string) (*geminiResponse, error) {
	var lastErr error

	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-time.After(c.config.RetryWaitTime * time.Duration(attempt)):
			case <-ctx.Done():
				return nil, &service.AIGeneratorError{
					Code:    service.ErrCodeTimeout,
					Message: "Context cancelled",
					Err:     ctx.Err(),
				}
			}
		}

		response, err := c.makeRequest(ctx, prompt)
		if err == nil {
			return response, nil
		}

		lastErr = err
		if !c.isRetryable(err) {
			return nil, err
		}
	}

	// リトライ可能なエラーで全て失敗した場合
	if lastErr != nil {
		return nil, lastErr
	}

	return nil, &service.AIGeneratorError{
		Code:    service.ErrCodeNetworkError,
		Message: "Unknown error",
	}
}

// 単一のAPIリクエストを送信
func (c *GeminiClient) makeRequest(ctx context.Context, prompt string) (*geminiResponse, error) {
	reqBody := geminiRequest{
		Contents: []geminiContent{
			{
				Role:  "user",
				Parts: []geminiPart{{Text: prompt}},
			},
		},
		Tools: []geminiTool{
			{URLContext: &geminiURLContext{}},
		},
		GenerationConfig: &geminiGenerationConfig{
			Temperature:     0.7,
			MaxOutputTokens: 8192,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, &service.AIGeneratorError{
			Code:    service.ErrCodeInvalidURL,
			Message: "Failed to marshal request",
			Err:     err,
		}
	}

	url := fmt.Sprintf("%s/models/%s:generateContent?key=%s",
		c.config.BaseURL,
		c.config.Model,
		c.config.APIKey,
	)

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, &service.AIGeneratorError{
			Code:    service.ErrCodeNetworkError,
			Message: "Failed to create request",
			Err:     err,
		}
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &service.AIGeneratorError{
			Code:    service.ErrCodeNetworkError,
			Message: "Request failed",
			Err:     err,
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &service.AIGeneratorError{
			Code:    service.ErrCodeNetworkError,
			Message: "Failed to read response",
			Err:     err,
		}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleError(resp.StatusCode, body)
	}

	var geminiResp geminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return nil, &service.AIGeneratorError{
			Code:    service.ErrCodeInvalidResponse,
			Message: "Failed to parse response",
			Err:     err,
		}
	}

	return &geminiResp, nil
}

// APIエラーの処理
func (c *GeminiClient) handleError(statusCode int, body []byte) error {
	var errResp geminiErrorResponse
	json.Unmarshal(body, &errResp)

	switch statusCode {
	case http.StatusUnauthorized, http.StatusForbidden:
		return &service.AIGeneratorError{
			Code:    service.ErrCodeUnauthorized,
			Message: "Invalid API key",
		}
	case http.StatusTooManyRequests:
		return &service.AIGeneratorError{
			Code:    service.ErrCodeAPILimit,
			Message: "Rate limit exceeded",
		}
	case http.StatusBadRequest:
		if strings.Contains(errResp.Error.Message, "SAFETY") {
			return &service.AIGeneratorError{
				Code:    service.ErrCodeContentBlocked,
				Message: "Content blocked by safety filters",
			}
		}
		return &service.AIGeneratorError{
			Code:    service.ErrCodeInvalidURL,
			Message: errResp.Error.Message,
		}
	default:
		return &service.AIGeneratorError{
			Code:    service.ErrCodeNetworkError,
			Message: fmt.Sprintf("API error: %s", errResp.Error.Message),
		}
	}
}

// リトライ可能なエラーか判定
func (c *GeminiClient) isRetryable(err error) bool {
	aiErr, ok := err.(*service.AIGeneratorError)
	if !ok {
		return false
	}

	return aiErr.Code == service.ErrCodeNetworkError ||
		aiErr.Code == service.ErrCodeAPILimit ||
		aiErr.Code == service.ErrCodeTimeout
}

// レスポンスをパース
func (c *GeminiClient) parseResponse(resp *geminiResponse, sourceURL string) (*service.GeneratedArticle, error) {
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("empty response")
	}

	text := resp.Candidates[0].Content.Parts[0].Text

	var data struct {
		Title         string   `json:"title"`
		Summary       string   `json:"summary"`
		SuggestedTags []string `json:"suggestedTags"`
	}

	if err := json.Unmarshal([]byte(text), &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if data.Title == "" || data.Summary == "" {
		return nil, fmt.Errorf("missing required fields")
	}

	return &service.GeneratedArticle{
		Title:         data.Title,
		Summary:       data.Summary,
		SuggestedTags: data.SuggestedTags,
		SourceURL:     sourceURL,
		TokenUsed:     resp.UsageMetadata.TotalTokenCount,
	}, nil
}
