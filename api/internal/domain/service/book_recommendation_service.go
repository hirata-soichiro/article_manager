package service

import (
	"context"

	"article-manager/internal/domain/entity"
)

// 書籍推薦サービスのインターフェース
type BookRecomendationService interface {
	RecommendBooks(ctx context.Context, articles []*entity.Article) ([]entity.Book, error)
}

// 書籍推薦処理固有のエラー
type BookRecommendationError struct {
	Code    string
	Message string
	Err     error
}

func (e *BookRecommendationError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *BookRecommendationError) Unwrap() error {
	return e.Err
}

const (
	ErrCodeNoArticles    = "NO_ARTICLES"
	ErrCodeAIError       = "AI_ERROR"
	ErrCodeBooksAPIError = "BOOKS_API_ERROR"
)
