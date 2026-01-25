package service

import (
	"context"
	"time"
)

// 記事生成リクエスト
type ArticleGenerationRequest struct {
	URL string
}

// 生成された記事の結果
type GeneratedArticle struct {
	Title         string
	Summary       string
	SuggestedTags []string
	SourceURL     string
	TokenUsed     int
	GeneratedAt   time.Time
}

// AIを使用した記事生成サービスのインターフェース
type AIGeneratorService interface {
	GenerateArticleFromURL(ctx context.Context, req ArticleGenerationRequest) (*GeneratedArticle, error)
}

// AI生成処理固有のエラー
type AIGeneratorError struct {
	Code    string
	Message string
	Err     error
}

func (e *AIGeneratorError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *AIGeneratorError) Unwrap() error {
	return e.Err
}

// エラーコード定数
const (
	ErrCodeInvalidURL      = "INVALID_URL"
	ErrCodeAPILimit        = "API_LIMIT_EXCEEDED"
	ErrCodeTimeout         = "TIMEOUT"
	ErrCodeContentBlocked  = "CONTENT_BLOCKED"
	ErrCodeInvalidResponse = "INVALID_RESPONSE"
	ErrCodeNetworkError    = "NETWORK_ERROR"
	ErrCodeUnauthorized    = "UNAUTHORIZED"
)
