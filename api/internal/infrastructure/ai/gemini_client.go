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
	Title      string
	AmazonURL  string
	RakutenURL string
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
	// 書籍推薦ではURLコンテキストツールは不要なので、falseを指定
	response, err := c.callAPIWithTools(ctx, prompt, false)
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

	return fmt.Sprintf("%s\n\n"+
		"【タスク】上記のユーザーが登録している記事の内容を総合的に分析し、ユーザーの興味・関心領域に基づいておすすめの書籍を5冊推薦してください。\n\n"+
		"【最重要事項 - 必ず遵守してください】\n"+
		"- あなたは書籍情報をJSON形式で返す専用BOTです\n"+
		"- 説明文、理由、コメント、挨拶、マークダウン記号は絶対に出力しないでください\n"+
		"- 最初の文字は必ず「{」で、最後の文字は必ず「}」で終わってください\n"+
		"- JSON以外の出力は厳禁です\n\n"+
		"【重要な指示】\n"+
		"1. 出力は必ず以下のJSON形式のみ\n"+
		"2. JSONオブジェクトのみを出力（前後に余分なテキストを含めない）\n"+
		"3. 推薦書籍は正確な書籍タイトルと著者名を記載してください\n"+
		"4. 実在する書籍のみを推薦してください（架空の書籍は不可）\n"+
		"5. 記事の内容から推測されるユーザーの専門性や興味に合った書籍を選んでください\n"+
		"6. 技術書、ビジネス書、専門書など、実用的な書籍を優先してください\n"+
		"7. 必ず5冊推薦してください\n"+
		"8. 著者名は正式名称（フルネーム）で記載してください\n"+
		"9. 【重要】日本の出版社から日本語で出版されている書籍のみを推薦してください（翻訳書を含む）\n"+
		"10. 【重要】洋書（原書が英語で海外出版社から出版されている書籍）は絶対に推薦しないでください\n"+
		"11. 【重要】書籍タイトルは必ず日本語で記載してください（ローマ字表記は不可）\n"+
		"12. 【超重要】各書籍のAmazon ASIN（10桁の商品コード）を正確に記載してください\n"+
		"    - ASINはAmazon.co.jpで実際に使用されている10桁のコードです\n"+
		"    - ISBN-10が存在する場合、多くの場合ISBN-10がASINと一致します\n"+
		"    - 例：「実践Rustプログラミング入門」のASINは「4798061700」です\n"+
		"    - 例：「Kubernetes完全ガイド 第2版」のASINは「4295009792」です\n"+
		"13. amazonUrlには、ASINを使用したAmazon.co.jp直接リンク（https://www.amazon.co.jp/dp/ASIN）を記載してください\n"+
		"    - 正しい形式：https://www.amazon.co.jp/dp/4798061700\n"+
		"    - 間違った形式：13桁のISBN-13を使用しないでください\n"+
		"14. rakutenUrlには、書籍タイトルをURLエンコードした楽天ブックス検索URL（https://books.rakuten.co.jp/search?g=001&sitem=書籍タイトル）を記載してください\n\n"+
		"出力形式（このフォーマット通りに出力）:\n"+
		"{\n"+
		"  \"books\": [\n"+
		"    {\n"+
		"      \"title\": \"実践Rustプログラミング入門\",\n"+
		"      \"amazonUrl\": \"https://www.amazon.co.jp/dp/4798061700\",\n"+
		"      \"rakutenUrl\": \"https://books.rakuten.co.jp/search?g=001&sitem=実践Rustプログラミング入門\"\n"+
		"    },\n"+
		"    {\n"+
		"      \"title\": \"Kubernetes完全ガイド 第2版\",\n"+
		"      \"amazonUrl\": \"https://www.amazon.co.jp/dp/4295009792\",\n"+
		"      \"rakutenUrl\": \"https://books.rakuten.co.jp/search?g=001&sitem=Kubernetes完全ガイド 第2版\"\n"+
		"    },\n"+
		"    {\n"+
		"      \"title\": \"リーダブルコード\",\n"+
		"      \"amazonUrl\": \"https://www.amazon.co.jp/dp/4873115655\",\n"+
		"      \"rakutenUrl\": \"https://books.rakuten.co.jp/search?g=001&sitem=リーダブルコード\"\n"+
		"    },\n"+
		"    {\n"+
		"      \"title\": \"入門 監視\",\n"+
		"      \"amazonUrl\": \"https://www.amazon.co.jp/dp/4873118646\",\n"+
		"      \"rakutenUrl\": \"https://books.rakuten.co.jp/search?g=001&sitem=入門 監視\"\n"+
		"    },\n"+
		"    {\n"+
		"      \"title\": \"プログラミング言語Go\",\n"+
		"      \"amazonUrl\": \"https://www.amazon.co.jp/dp/4621300253\",\n"+
		"      \"rakutenUrl\": \"https://books.rakuten.co.jp/search?g=001&sitem=プログラミング言語Go\"\n"+
		"    }\n"+
		"  ]\n"+
		"}\n\n"+
		"必ず上記のJSON形式のみで回答してください。説明文や推薦理由は不要です。", articlesInfo.String())
}

// Gemini API呼び出し
func (c *GeminiClient) callAPI(ctx context.Context, prompt string) (*geminiResponse, error) {
	return c.callAPIWithTools(ctx, prompt, true)
}

// ツールの有無を指定してGemini API呼び出し
func (c *GeminiClient) callAPIWithTools(ctx context.Context, prompt string, includeURLContext bool) (*geminiResponse, error) {
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

		response, err := c.makeRequestWithTools(ctx, prompt, includeURLContext)
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
	return c.makeRequestWithTools(ctx, prompt, true)
}

// ツールの有無を指定してAPIリクエストを送信
func (c *GeminiClient) makeRequestWithTools(ctx context.Context, prompt string, includeURLContext bool) (*geminiResponse, error) {
	reqBody := geminiRequest{
		Contents: []geminiContent{
			{
				Role:  "user",
				Parts: []geminiPart{{Text: prompt}},
			},
		},
		GenerationConfig: &geminiGenerationConfig{
			Temperature:     0.3, // より決定論的な出力のため低く設定
			MaxOutputTokens: 4096,
		},
	}

	// URL Contextツールを条件付きで追加
	if includeURLContext {
		reqBody.Tools = []geminiTool{
			{URLContext: &geminiURLContext{}},
		}
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
			Title      string `json:"title"`
			AmazonURL  string `json:"amazonUrl"`
			RakutenURL string `json:"rakutenUrl"`
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
		if book.Title == "" {
			return nil, fmt.Errorf("book at index %d is missing title (title=%s)", i, book.Title)
		}
		books = append(books, RecommendedBook{
			Title:      book.Title,
			AmazonURL:  book.AmazonURL,
			RakutenURL: book.RakutenURL,
		})
	}

	return books, nil
}
