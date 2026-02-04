package usecase

import (
	"context"
	"strings"

	"article-manager/internal/domain/entity"
	domainerrors "article-manager/internal/domain/errors"
	"article-manager/internal/domain/repository"
	"article-manager/internal/infrastructure/logger"

	"go.uber.org/zap"
)

// 記事に関するユースケース
type ArticleUsecase struct {
	repo repository.ArticleRepository
}

// コンストラクタ
func NewArticleUsecase(repo repository.ArticleRepository) *ArticleUsecase {
	return &ArticleUsecase{repo: repo}
}

// 新しい記事を作成
func (u *ArticleUsecase) CreateArticle(ctx context.Context, title, url, summary string, tags []string, memo string) (*entity.Article, error) {
	logger.Debug("Creating new article",
		zap.String("title", title),
		zap.String("url", url),
		zap.Strings("tags", tags),
	)

	article, err := entity.NewArticle(title, url, summary, tags, memo)
	if err != nil {
		logger.Warn("Failed to create article entity",
			zap.Error(err),
			zap.String("title", title),
		)
		return nil, domainerrors.ValidationError("article", err.Error())
	}

	savedArticle, err := u.repo.Create(ctx, article)
	if err != nil {
		logger.Error("Failed to save article to repository",
			zap.Error(err),
			zap.String("title", title),
		)
		return nil, err
	}

	logger.Info("Successfully created article",
		zap.Int64("id", savedArticle.ID),
		zap.String("title", savedArticle.Title),
	)

	return savedArticle, nil
}

// 指定されたIDの記事を取得
func (u *ArticleUsecase) GetArticleByID(ctx context.Context, id int64) (*entity.Article, error) {
	logger.Debug("Getting article by ID",
		zap.Int64("id", id),
	)

	if id <= 0 {
		logger.Warn("Invalid article ID",
			zap.Int64("id", id),
		)
		return nil, domainerrors.InvalidArgumentError("id", "id must be positive")
	}

	article, err := u.repo.FindByID(ctx, id)
	if err != nil {
		logger.Warn("Failed to find article",
			zap.Error(err),
			zap.Int64("id", id),
		)
		return nil, err
	}

	logger.Debug("Successfully retrieved article",
		zap.Int64("id", id),
		zap.String("title", article.Title),
	)

	return article, nil
}

// 全ての記事を取得
func (u *ArticleUsecase) GetAllArticles(ctx context.Context) ([]*entity.Article, error) {
	logger.Debug("Getting all articles")

	articles, err := u.repo.FindAll(ctx)
	if err != nil {
		logger.Error("Failed to retrive all articles",
			zap.Error(err),
		)
		return nil, err
	}

	logger.Debug("Successfully retrieved all articles",
		zap.Int("count", len(articles)),
	)

	return articles, nil
}

// 記事を更新
func (u *ArticleUsecase) UpdateArticle(ctx context.Context, id int64, title, url, summary string, tags []string, memo string) (*entity.Article, error) {
	logger.Debug("Updating article",
		zap.Int64("id", id),
		zap.String("title", title),
	)

	article, err := u.repo.FindByID(ctx, id)
	if err != nil {
		logger.Warn("Failed to find article for update",
			zap.Error(err),
			zap.Int64("id", id),
		)
		return nil, err
	}

	if err := article.Update(title, url, summary, tags, memo); err != nil {
		logger.Warn("Failed to update article entity",
			zap.Error(err),
			zap.Int64("id", id),
		)
		return nil, domainerrors.ValidationError("article", err.Error())
	}

	updatedArticle, err := u.repo.Update(ctx, article)
	if err != nil {
		logger.Error("Failed to update article in repository",
			zap.Error(err),
			zap.Int64("id", id),
		)
		return nil, err
	}

	logger.Info("Successfully updated article",
		zap.Int64("id", updatedArticle.ID),
		zap.String("title", updatedArticle.Title),
	)

	return updatedArticle, nil
}

// 指定されたIDの記事を削除
func (u *ArticleUsecase) DeleteArticle(ctx context.Context, id int64) error {
	logger.Debug("Deleting article",
		zap.Int64("id", id),
	)

	if id <= 0 {
		logger.Warn("Invalid article ID for deletion",
			zap.Int64("id", id),
		)
		return domainerrors.InvalidArgumentError("id", "id must be positive")
	}

	if err := u.repo.Delete(ctx, id); err != nil {
		logger.Error("Failed to delete article",
			zap.Error(err),
			zap.Int64("id", id),
		)
		return err
	}

	logger.Info("Successfully deleted article",
		zap.Int64("id", id),
	)

	return nil
}

// キーワードで記事を検索
func (u *ArticleUsecase) SearchArticles(ctx context.Context, keyword string) ([]*entity.Article, error) {
	logger.Debug("Searching articles",
		zap.String("keyword", keyword),
	)

	trimmedKeyword := strings.TrimSpace(keyword)
	if trimmedKeyword == "" {
		logger.Warn("Empty search keyword")
		return nil, domainerrors.ValidationError("keyword", "keyword cannot be empty")
	}

	articles, err := u.repo.Search(ctx, trimmedKeyword)
	if err != nil {
		logger.Error("Failed to search articles",
			zap.Error(err),
			zap.String("keyword", trimmedKeyword),
		)
		return nil, err
	}

	logger.Info("Successfully searched articles",
		zap.String("keyword", trimmedKeyword),
		zap.Int("count", len(articles)),
	)

	return articles, nil
}
