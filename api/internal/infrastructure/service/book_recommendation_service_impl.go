package service

import (
	"context"

	"article-manager/internal/domain/entity"
	domainservice "article-manager/internal/domain/service"
	"article-manager/internal/infrastructure/ai"
	"article-manager/internal/infrastructure/external"
	"article-manager/internal/infrastructure/logger"

	"go.uber.org/zap"
)

// BookRecommendationServiceの実装
type bookRecommendationServiceImpl struct {
	geminiClient      *ai.GeminiClient
	googleBooksClient *external.GoogleBooksClient
}

// コンストラクタ
func NewBookRecommendationService(
	geminiClient *ai.GeminiClient,
	googleBooksClient *external.GoogleBooksClient,
) domainservice.BookRecommendationService {
	return &bookRecommendationServiceImpl{
		geminiClient:      geminiClient,
		googleBooksClient: googleBooksClient,
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

	// 各書籍の詳細情報をGoogle Books APIから取得
	books := make([]entity.Book, 0, len(recommendedBooks))
	for i, rec := range recommendedBooks {
		logger.Debug("Fetching book details from Google Books",
			zap.Int("index", i),
			zap.String("title", rec.Title),
			zap.String("author", rec.Author),
		)

		// Google Books APIで書籍詳細を検索
		bookDetail, err := s.googleBooksClient.SearchBook(ctx, rec.Title, rec.Author)
		if err != nil {
			// Google Books APIのエラーはログに記録するが、推薦は続行
			logger.Warn("Failed to fetch book details from Google Books API, using basic info",
				zap.String("title", rec.Title),
				zap.String("author", rec.Author),
				zap.Error(err),
			)
			// 基本情報のみで書籍を作成
			bookDetail = &external.BookDetail{
				Title:  rec.Title,
				Author: rec.Author,
				PurchaseLinks: external.PurchaseLinks{
					Amazon:  "",
					Rakuten: "",
				},
			}
		}

		// entity.Bookに変換
		book := entity.Book{
			Title:  bookDetail.Title,
			Author: bookDetail.Author,
			ISBN:   bookDetail.ISBN,
			PurchaseLinks: entity.PurchaseLinks{
				Amazon:  bookDetail.PurchaseLinks.Amazon,
				Rakuten: bookDetail.PurchaseLinks.Rakuten,
			},
		}

		// バリデーション（NewBookを使わずに直接作成）
		if book.Title == "" || book.Author == "" {
			logger.Warn("Skipping book with missing title or author",
				zap.String("title", book.Title),
				zap.String("author", book.Author),
			)
			continue
		}

		books = append(books, book)

		logger.Debug("Successfully added book",
			zap.Int("index", i),
			zap.String("title", book.Title),
			zap.String("author", book.Author),
			zap.String("isbn", book.ISBN),
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
