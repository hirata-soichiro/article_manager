package external

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"article-manager/internal/domain/service"
)

// Google Books API設定
type GoogleBooksConfig struct {
	APIKey        string
	BaseURL       string
	Timeout       time.Duration
	MaxRetries    int
	RetryWaitTime time.Duration
}

// デフォルトGoogle Books API設定
func DefaultGoogleBooksConfig(apiKey string) *GoogleBooksConfig {
	return &GoogleBooksConfig{
		APIKey:        apiKey,
		BaseURL:       "https://www.googleapis.com/books/v1",
		Timeout:       10 * time.Second,
		MaxRetries:    3,
		RetryWaitTime: 1 * time.Second,
	}
}

// Google Books APIクライアント
type GoogleBooksClient struct {
	config     *GoogleBooksConfig
	httpClient *http.Client
}

// 新しいクライアントを作成
func NewGoogleBooksClient(config *GoogleBooksConfig) *GoogleBooksClient {
	return &GoogleBooksClient{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// 書籍詳細情報
type BookDetail struct {
	Title         string
	Author        string
	ISBN          string
	ISBN13        string
	PurchaseLinks PurchaseLinks
}

// 購入リンク
type PurchaseLinks struct {
	Amazon  string
	Rakuten string
}

// Google Books APIレスポンス構造
type googleBooksResponse struct {
	Items []struct {
		VolumeInfo struct {
			Title               string   `json:"title"`
			Authors             []string `json:"authors"`
			IndustryIdentifiers []struct {
				Type       string `json:"type"`
				Identifier string `json:"identifier"`
			} `json:"industryIdentifiers"`
		} `json:"volumeInfo"`
	} `json:"items"`
}

// 書籍情報を検索
func (c *GoogleBooksClient) SearchBook(ctx context.Context, title, author string) (*BookDetail, error) {
	if title == "" {
		return nil, &service.BookRecommendationError{
			Code:    service.ErrCodeBooksAPIError,
			Message: "title is required",
		}
	}

	var lastErr error
	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-time.After(c.config.RetryWaitTime * time.Duration(attempt)):
			case <-ctx.Done():
				return nil, &service.BookRecommendationError{
					Code:    service.ErrCodeBooksAPIError,
					Message: "context cancelled",
					Err:     ctx.Err(),
				}
			}
		}

		book, err := c.searchBookRequest(ctx, title, author)
		if err == nil {
			return book, nil
		}

		lastErr = err
		// レート制限エラーやネットワークエラーの場合はリトライ
		if !c.isRetryable(err) {
			return nil, err
		}
	}

	return nil, &service.BookRecommendationError{
		Code:    service.ErrCodeBooksAPIError,
		Message: fmt.Sprintf("failed after %d retries", c.config.MaxRetries),
		Err:     lastErr,
	}
}

// 単一の書籍検索リクエスト
func (c *GoogleBooksClient) searchBookRequest(ctx context.Context, title, author string) (*BookDetail, error) {
	// クエリ構築
	query := title
	if author != "" {
		query = fmt.Sprintf("%s %s", title, author)
	}

	// URLエンコード
	encodedQuery := url.QueryEscape(query)
	apiURL := fmt.Sprintf("%s/volumes?q=%s&maxResults=1", c.config.BaseURL, encodedQuery)

	if c.config.APIKey != "" {
		apiURL = fmt.Sprintf("%s&key=%s", apiURL, c.config.APIKey)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, &service.BookRecommendationError{
			Code:    service.ErrCodeBooksAPIError,
			Message: "failed to create request",
			Err:     err,
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &service.BookRecommendationError{
			Code:    service.ErrCodeBooksAPIError,
			Message: "request failed",
			Err:     err,
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &service.BookRecommendationError{
			Code:    service.ErrCodeBooksAPIError,
			Message: "failed to read response",
			Err:     err,
		}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleError(resp.StatusCode, body)
	}

	var booksResp googleBooksResponse
	if err := json.Unmarshal(body, &booksResp); err != nil {
		return nil, &service.BookRecommendationError{
			Code:    service.ErrCodeBooksAPIError,
			Message: "failed to parse response",
			Err:     err,
		}
	}

	// 検索結果が0件の場合(エラーではなく基本情報のみ返す)
	if len(booksResp.Items) == 0 {
		return &BookDetail{
			Title:  title,
			Author: author,
		}, nil
	}

	// 最初の結果を使用
	item := booksResp.Items[0]
	volumeInfo := item.VolumeInfo

	// ISBNを抽出
	isbn10 := ""
	isbn13 := ""
	for _, identifier := range volumeInfo.IndustryIdentifiers {
		switch identifier.Type {
		case "ISBN_10":
			isbn10 = identifier.Identifier
		case "ISBN_13":
			isbn13 = identifier.Identifier
		}
	}

	primaryISBN := isbn10
	if primaryISBN == "" {
		primaryISBN = isbn13
	}

	// 著署名を結合
	authorName := author
	if len(volumeInfo.Authors) > 0 {
		authorName = strings.Join(volumeInfo.Authors, ", ")
	}

	// 購入リンクを生成
	purchaseLinks := c.generatePurchaseLinks(primaryISBN, isbn13)

	return &BookDetail{
		Title:         volumeInfo.Title,
		Author:        authorName,
		ISBN:          primaryISBN,
		ISBN13:        isbn13,
		PurchaseLinks: purchaseLinks,
	}, nil
}

// 購入リンクを生成
func (c *GoogleBooksClient) generatePurchaseLinks(isbn10, isbn13 string) PurchaseLinks {
	links := PurchaseLinks{}

	if isbn10 != "" {
		cleanISBN := strings.ReplaceAll(isbn10, "-", "")
		links.Amazon = fmt.Sprintf("https://www.amazon.co.jp/dp/%s", cleanISBN)
	} else if isbn13 != "" {
		cleanISBN := strings.ReplaceAll(isbn13, "-", "")
		links.Amazon = fmt.Sprintf("https://www.amazon.co.jp/dp/%s", cleanISBN)
	}

	if isbn13 != "" {
		cleanISBN := strings.ReplaceAll(isbn13, "-", "")
		links.Rakuten = fmt.Sprintf("https://books.rakuten.co.jp/search?sitem=%s", cleanISBN)
	}

	return links
}

// APIエラーの処理
func (c *GoogleBooksClient) handleError(statusCode int, body []byte) error {
	var errResp struct {
		Error struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	json.Unmarshal(body, &errResp)

	switch statusCode {
	case http.StatusBadRequest:
		return &service.BookRecommendationError{
			Code:    service.ErrCodeBooksAPIError,
			Message: fmt.Sprintf("bad request: %s", errResp.Error.Message),
		}
	case http.StatusUnauthorized, http.StatusForbidden:
		return &service.BookRecommendationError{
			Code:    service.ErrCodeBooksAPIError,
			Message: fmt.Sprintf("invalid API key"),
		}
	case http.StatusTooManyRequests:
		return &service.BookRecommendationError{
			Code:    service.ErrCodeBooksAPIError,
			Message: fmt.Sprintf("rate limit exceeded"),
		}
	case http.StatusInternalServerError, http.StatusServiceUnavailable:
		return &service.BookRecommendationError{
			Code:    service.ErrCodeBooksAPIError,
			Message: fmt.Sprintf("API error: status=%d", statusCode),
		}
	default:
		return &service.BookRecommendationError{
			Code:    service.ErrCodeBooksAPIError,
			Message: fmt.Sprintf("unexpected error: status=%d", statusCode),
		}
	}
}

// リトライ可能なエラーか判定
func (c *GoogleBooksClient) isRetryable(err error) bool {
	bookErr, ok := err.(*service.BookRecommendationError)
	if !ok {
		return false
	}

	if bookErr.Code == service.ErrCodeBooksAPIError {
		msg := bookErr.Message
		return strings.Contains(msg, "rate limit") ||
			strings.Contains(msg, "status=500") ||
			strings.Contains(msg, "status=503") ||
			strings.Contains(msg, "request failed")
	}

	return false
}
