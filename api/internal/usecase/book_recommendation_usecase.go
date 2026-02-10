package usecase

import (
	"context"

	"article-manager/internal/domain/entity"
	domainerrors "article-manager/internal/domain/errors"
	"article-manager/internal/domain/repository"
	"article-manager/internal/domain/service"
	"article-manager/internal/infrastructure/logger"

	"go.uber.org/zap"
)

// 書籍推薦ユースケース
type BookRecommendationUsecase struct {
	articleRepo               repository.ArticleRepository
	bookRecommendationRepo    repository.BookRecommendationRepository
	bookRecommendationService service.BookRecommendationService
}

// コンストラクタ
func NewBookRecommendationUsecase(
	articleRepo repository.ArticleRepository,
	bookRecommendationRepo repository.BookRecommendationRepository,
	bookRecommendationService service.BookRecommendationService,
) *BookRecommendationUsecase {
	return &BookRecommendationUsecase{
		articleRepo:               articleRepo,
		bookRecommendationRepo:    bookRecommendationRepo,
		bookRecommendationService: bookRecommendationService,
	}
}

// 書籍推薦を取得
func (u *BookRecommendationUsecase) GetBookRecommendations(ctx context.Context) (*entity.BookRecommendationCache, error) {
	logger.Debug("Getting book recommendations")

	// キャッシュ確認
	cache, err := u.bookRecommendationRepo.FindLatestValid(ctx)
	if err == nil && cache != nil && cache.IsValid() {
		logger.Info("Returning cached book recommendations",
			zap.Time("generated_at", cache.GeneratedAt),
			zap.Time("expires_at", cache.ExpiresAt),
			zap.Int("count", len(cache.Books)),
		)
		return cache, nil
	}

	logger.Debug("Cache not found or expired, generating new recommendations")

	// 全記事を取得
	articles, err := u.articleRepo.FindAll(ctx)
	if err != nil {
		logger.Error("Failed to retrieve articles for recommendations",
			zap.Error(err),
		)
		return nil, err
	}

	// 記事が0件の場合は空の推薦を返す
	if len(articles) == 0 {
		logger.Info("No articles found, returning empty recommendations")
		emptyCache, err := entity.NewBookRecommendationCache([]entity.Book{})
		if err != nil {
			return nil, err
		}
		return emptyCache, nil
	}

	logger.Debug("Retrieved articles for recommendation",
		zap.Int("article_count", len(articles)),
	)

	// AIで書籍を推薦
	books, err := u.bookRecommendationService.RecommendBooks(ctx, articles)
	if err != nil {
		logger.Error("Failed to generate book recommendations from AI",
			zap.Error(err),
		)
		return nil, err
	}

	logger.Debug("AI generated book recommendation",
		zap.Int("book_count", len(books)),
	)

	// キャッシュを作成して保存
	cache, err = entity.NewBookRecommendationCache(books)
	if err != nil {
		logger.Error("Failed to create book recommendations cache",
			zap.Error(err),
		)
		return nil, domainerrors.ValidationError("book_recommendation_cache", err.Error())
	}

	savedCache, err := u.bookRecommendationRepo.Save(ctx, cache)
	if err != nil {
		logger.Error("Failed to save book recommendations cache",
			zap.Error(err),
		)
		return nil, err
	}

	logger.Info("Successfully generated and cached book recommendations",
		zap.Int64("cache_id", savedCache.ID),
		zap.Int("book_count", len(savedCache.Books)),
		zap.Time("expires_at", savedCache.ExpiresAt),
	)

	return savedCache, nil
}
