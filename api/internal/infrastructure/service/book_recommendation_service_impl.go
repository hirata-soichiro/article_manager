package service

import (
	"context"

	"article-manager/internal/domain/entity"
	domainservice "article-manager/internal/domain/service"
	"article-manager/internal/infrastructure/ai"
	"article-manager/internal/infrastructure/logger"

	"go.uber.org/zap"
)

// BookRecommendationServiceの実装
type bookRecommendationServiceImpl struct {
	geminiClient *ai.GeminiClient
}

// コンストラクタ
func NewBookRecommendationService(
	geminiClient *ai.GeminiClient,
) domainservice.BookRecommendationService {
	return &bookRecommendationServiceImpl{
		geminiClient: geminiClient,
	}
}

// 書籍を推薦する
func (s *bookRecommendationServiceImpl) RecommendBooks(ctx context.Context, articles []*entity.Article) ([]entity.Book, error) {
	logger.Debug("Recommending books from articles",
		zap.Int("article_count", len(articles)),
	)

	// 記事が0件の場合は空の書籍リストを返す
	if len(articles) == 0 {
		logger.Debug("No articles provided, returning empty book list")
		return []entity.Book{}, nil
	}

	// GeminiでAI推薦を取得
	recommendedBooks, err := s.geminiClient.RecommendBooks(ctx, articles)
	if err != nil {
		logger.Error("Failed to get book recommendations from Gemini",
			zap.Error(err),
		)
		return nil, &domainservice.BookRecommendationError{
			Code:    domainservice.ErrCodeAIError,
			Message: "Failed to generate book recommendations from AI",
			Err:     err,
		}
	}

	logger.Debug("Received book recommendations from Gemini",
		zap.Int("book_count", len(recommendedBooks)),
	)

	// Geminiのレスポンスを直接entity.Bookに変換
	books := make([]entity.Book, 0, len(recommendedBooks))
	for i, rec := range recommendedBooks {
		logger.Debug("Processing book recommendation",
			zap.Int("index", i),
			zap.String("title", rec.Title),
		)

		// entity.Bookに変換
		book := entity.Book{
			Title: rec.Title,
			PurchaseLinks: entity.PurchaseLinks{
				Amazon:  rec.AmazonURL,
				Rakuten: rec.RakutenURL,
			},
		}

		// バリデーション
		if book.Title == "" {
			logger.Warn("Skipping book with missing title",
				zap.String("title", book.Title),
			)
			continue
		}

		books = append(books, book)

		logger.Debug("Successfully added book",
			zap.Int("index", i),
			zap.String("title", book.Title),
		)
	}

	if len(books) == 0 {
		logger.Error("No valid books after processing recommendations")
		return nil, &domainservice.BookRecommendationError{
			Code:    domainservice.ErrCodeBooksAPIError,
			Message: "Failed to process any book recommendations",
		}
	}

	logger.Info("Successfully recommended books",
		zap.Int("book_count", len(books)),
	)

	return books, nil
}
