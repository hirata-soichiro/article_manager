package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"article-manager/internal/domain/entity"
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
		Model:         "gemini-2.5-flash-lite",
		BaseURL:       "https://generativelanguage.googleapis.com/v1beta",
		Timeout:       30 * time.Second,
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

type BookRecommendationRequest struct {
	Articles []*entity.Article
}

type RecommendedBook struct {
	Title  string
	Author string
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
	return fmt.Sprintf("以下のURLの記事を分析し、記事管理用の情報をJSON形式で生成してください。\n\n"+
		"URL: %s\n\n"+
		"【重要な指示】\n"+
		"1. 出力は必ず以下のJSON形式のみとし、それ以外のテキスト（説明文、マークダウン、コードブロック記号など）は一切含めないでください\n"+
		"2. JSONオブジェクトのみを出力してください（前後に余分なテキストを含めない）\n"+
		"3. summaryは記事の核心を捉え、200文字以内で簡潔にまとめる\n"+
		"4. suggestedTagsは3-5個、検索しやすく具体的なものを選ぶ\n"+
		"5. 技術記事の場合は使用されている技術スタックをタグに含める\n"+
		"6. 日本語記事は日本語で、英語記事は日本語に翻訳して出力する\n\n"+
		"出力形式（このフォーマット通りに出力）:\n"+
		"{\n"+
		"  \"title\": \"記事のタイトル\",\n"+
		"  \"summary\": \"記事の要約（200文字以内）\",\n"+
		"  \"suggestedTags\": [\"タグ1\", \"タグ2\", \"タグ3\", \"タグ4\", \"タグ5\"]\n"+
		"}\n\n"+
		"必ず上記のJSON形式のみで回答してください。説明文は不要です。", url)
}

// 全記事から書籍を推薦
func (c *GeminiClient) RecommendBooks(ctx context.Context, articles []*entity.Article) ([]RecommendedBook, error) {
	if len(articles) == 0 {
		return []RecommendedBook{}, nil
	}

	prompt := c.buildBookRecommendationPrompt(articles)
	response, err := c.callAPI(ctx, prompt)
	if err != nil {
		return nil, err
	}

	books, err := c.parseBookRecommendationResponse(response)
	if err != nil {
		return nil, &service.AIGeneratorError{
			Code:    service.ErrCodeInvalidResponse,
			Message: "Failed to parse book recommendation response",
			Err:     err,
		}
	}

	return books, nil
}

// 書籍推薦用プロンプト構築
func (c *GeminiClient) buildBookRecommendationPrompt(articles []*entity.Article) string {
	var articlesInfo strings.Builder

	articlesInfo.WriteString("【登録されている記事一覧】\n\n")
	for i, article := range articles {
		articlesInfo.WriteString(fmt.Sprintf("%d. タイトル: %s\n", i+1, article.Title))
		articlesInfo.WriteString(fmt.Sprintf("   要約: %s\n", article.Summary))
		if len(article.Tags) > 0 {
			articlesInfo.WriteString(fmt.Sprintf("   タグ: %s\n", strings.Join(article.Tags, ", ")))
		}
		if article.Memo != "" {
			articlesInfo.WriteString(fmt.Sprintf("   メモ: %s\n", article.Memo))
		}
		articlesInfo.WriteString("\n")
	}

	return fmt.Sprintf("上記のユーザーが登録している記事の内容を総合的に分析し、ユーザーの興味・関心領域に基づいておすすめの書籍を5冊推薦してください。\n\n"+
		"%s\n\n"+
		"【重要な指示】\n"+
		"1. 出力は必ず以下のJSON形式のみとし、それ以外のテキスト（説明文、マークダウン、コードブロック記号など）は一切含めないでください\n"+
		"2. JSONオブジェクトのみを出力してください（前後に余分なテキストを含めない）\n"+
		"3. 推薦書籍は正確な書籍タイトルと著者名を記載してください\n"+
		"4. 実在する書籍のみを推薦してください（架空の書籍は不可）\n"+
		"5. 記事の内容から推測されるユーザーの専門性や興味に合った書籍を選んでください\n"+
		"6. 技術書、ビジネス書、専門書など、実用的な書籍を優先してください\n"+
		"7. 必ず5冊推薦してください\n"+
		"8. 著者名は正式名称（フルネーム）で記載してください\n"+
		"9. 日本語の書籍は日本語で、英語の書籍は原題のまま記載してください\n\n"+
		"出力形式（このフォーマット通りに出力）:\n"+
		"{\n"+
		"  \"books\": [\n"+
		"    {\n"+
		"      \"title\": \"書籍タイトル1\",\n"+
		"      \"author\": \"著者名1\"\n"+
		"    },\n"+
		"    {\n"+
		"      \"title\": \"書籍タイトル2\",\n"+
		"      \"author\": \"著者名2\"\n"+
		"    },\n"+
		"    {\n"+
		"      \"title\": \"書籍タイトル3\",\n"+
		"      \"author\": \"著者名3\"\n"+
		"    },\n"+
		"    {\n"+
		"      \"title\": \"書籍タイトル4\",\n"+
		"      \"author\": \"著者名4\"\n"+
		"    },\n"+
		"    {\n"+
		"      \"title\": \"書籍タイトル5\",\n"+
		"      \"author\": \"著者名5\"\n"+
		"    }\n"+
		"  ]\n"+
		"}\n\n"+
		"必ず上記のJSON形式のみで回答してください。説明文や推薦理由は不要です。", articlesInfo.String())
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
			Temperature:     0.3, // より決定論的な出力のため低く設定
			MaxOutputTokens: 4096,
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

// マークダウンのコードブロックからJSONを抽出（フェイルセーフ用）
func (c *GeminiClient) extractJSON(text string) string {
	text = strings.TrimSpace(text)

	// ```json ... ``` または ``` ... ``` で囲まれている場合は中身を抽出
	if strings.HasPrefix(text, "```") {
		lines := strings.Split(text, "\n")
		if len(lines) > 2 {
			// 最初の行（```json または ```）を除去
			content := strings.Join(lines[1:], "\n")
			// 最後の ``` を除去
			if idx := strings.LastIndex(content, "```"); idx != -1 {
				content = content[:idx]
			}
			text = strings.TrimSpace(content)
		}
	}

	// 念のため、最初の { から最後の } までを抽出
	startIdx := strings.Index(text, "{")
	endIdx := strings.LastIndex(text, "}")

	if startIdx != -1 && endIdx != -1 && endIdx > startIdx {
		text = text[startIdx : endIdx+1]
	}

	return strings.TrimSpace(text)
}

// レスポンスをパース
func (c *GeminiClient) parseResponse(resp *geminiResponse, sourceURL string) (*service.GeneratedArticle, error) {
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("empty response from API")
	}

	originalText := resp.Candidates[0].Content.Parts[0].Text

	// マークダウンのコードブロックを除去
	extractedText := c.extractJSON(originalText)

	var data struct {
		Title         string   `json:"title"`
		Summary       string   `json:"summary"`
		SuggestedTags []string `json:"suggestedTags"`
	}

	if err := json.Unmarshal([]byte(extractedText), &data); err != nil {
		// デバッグ情報を含めたエラーメッセージ
		return nil, fmt.Errorf("failed to parse JSON: %w\nOriginal text: %s\nExtracted text: %s", err, originalText, extractedText)
	}

	if data.Title == "" || data.Summary == "" {
		return nil, fmt.Errorf("missing required fields (title or summary is empty)\nParsed data: title=%s, summary=%s, tags=%v", data.Title, data.Summary, data.SuggestedTags)
	}

	// タグが空の場合はデフォルト値を設定
	if len(data.SuggestedTags) == 0 {
		data.SuggestedTags = []string{}
	}

	return &service.GeneratedArticle{
		Title:         data.Title,
		Summary:       data.Summary,
		SuggestedTags: data.SuggestedTags,
		SourceURL:     sourceURL,
		TokenUsed:     resp.UsageMetadata.TotalTokenCount,
	}, nil
}

// 書籍推薦レスポンスをパース
func (c *GeminiClient) parseBookRecommendationResponse(resp *geminiResponse) ([]RecommendedBook, error) {
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("empty response from API")
	}

	originalText := resp.Candidates[0].Content.Parts[0].Text

	// マークダウンのコードブロックを除去（既存のextractJSONメソッドを利用）
	extractedText := c.extractJSON(originalText)

	var data struct {
		Books []struct {
			Title  string `json:"title"`
			Author string `json:"author"`
		} `json:"books"`
	}

	if err := json.Unmarshal([]byte(extractedText), &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w\nOriginal text: %s\nExtracted text: %s", err, originalText, extractedText)
	}

	if len(data.Books) == 0 {
		return nil, fmt.Errorf("no books found in response")
	}

	books := make([]RecommendedBook, 0, len(data.Books))
	for i, book := range data.Books {
		if book.Title == "" || book.Author == "" {
			return nil, fmt.Errorf("book at index %d is missing title or author (title=%s, author=%s)", i, book.Title, book.Author)
		}
		books = append(books, RecommendedBook{
			Title:  book.Title,
			Author: book.Author,
		})
	}

	return books, nil
}
