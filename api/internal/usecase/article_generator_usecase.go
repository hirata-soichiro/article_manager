package usecase

import (
	"context"
	"strings"

	"article-manager/internal/domain/entity"
	domainerrors "article-manager/internal/domain/errors"
	"article-manager/internal/domain/repository"
	"article-manager/internal/domain/service"
	"article-manager/internal/infrastructure/logger"

	"go.uber.org/zap"
)

// 記事自動生成ユースケース
type ArticleGeneratorUsecase struct {
	aiGenerator service.AIGeneratorService
	articleRepo repository.ArticleRepository
	tagRepo     repository.TagRepository
}

func NewArticleGeneratorUsecase(
	aiGenerator service.AIGeneratorService,
	articleRepo repository.ArticleRepository,
	tagRepo repository.TagRepository,
) *ArticleGeneratorUsecase {
	return &ArticleGeneratorUsecase{
		aiGenerator: aiGenerator,
		articleRepo: articleRepo,
		tagRepo:     tagRepo,
	}
}

// URLから記事情報を自動生成してDBに保存
func (u *ArticleGeneratorUsecase) GenerateArticleFromURL(ctx context.Context, url string, memo string) (*entity.Article, error) {
	logger.Debug("Generating article from URL",
		zap.String("url", url),
		zap.String("memo", memo),
	)

	if url == "" {
		logger.Warn("URL is empty")
		return nil, domainerrors.InvalidArgumentError("url", "url is required")
	}
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		logger.Warn("Invalid URL format",
			zap.String("url", url),
		)
		return nil, domainerrors.InvalidArgumentError("url", "invalid url format")
	}

	logger.Info("Calling AI generator service",
		zap.String("url", url),
	)

	generated, err := u.aiGenerator.GenerateArticleFromURL(ctx, service.ArticleGenerationRequest{URL: url})
	if err != nil {
		logger.Error("AI generator service failed",
			zap.Error(err),
			zap.String("url", url),
		)
		return nil, err
	}

	if generated.Title == "" {
		logger.Warn("Generated title is empty",
			zap.String("url", url),
		)
		return nil, domainerrors.ValidationError("title", "title is required")
	}
	if generated.Summary == "" {
		logger.Warn("Generated summary is empty",
			zap.String("url", url),
		)
		return nil, domainerrors.ValidationError("summary", "summary is required")
	}

	logger.Info("AI successfully generated article content",
		zap.String("title", generated.Title),
		zap.Strings("suggested_tags", generated.SuggestedTags),
	)

	tags := []string{}
	for _, tagName := range generated.SuggestedTags {
		if tagName == "" {
			logger.Warn("Tag name is empty, skipping")
			continue
		}

		logger.Debug("Processing tag",
			zap.String("tag", tagName),
		)

		existingTag, err := u.tagRepo.FindByName(ctx, tagName)
		if err != nil {
			logger.Debug("Tag not found, creating new tag",
				zap.String("tag", tagName),
			)

			newTag, err := entity.NewTag(tagName)
			if err != nil {
				logger.Warn("Failed to create tag entity",
					zap.Error(err),
					zap.String("tag", tagName),
				)
				return nil, domainerrors.ValidationError("tag", err.Error())
			}
			createdTag, err := u.tagRepo.Create(ctx, newTag)
			if err != nil {
				logger.Error("Failed to create tag in repository",
					zap.Error(err),
					zap.String("tag", tagName),
				)
				return nil, err
			}
			tags = append(tags, createdTag.Name)
			logger.Debug("Created new tag",
				zap.Int64("tag_id", createdTag.ID),
				zap.String("tag", createdTag.Name),
			)
		} else {
			tags = append(tags, existingTag.Name)
			logger.Debug("Using existing tag",
				zap.Int64("tag_id", existingTag.ID),
				zap.String("tag", existingTag.Name),
			)
		}
	}

	article, err := entity.NewArticle(generated.Title, url, generated.Summary, tags, memo)
	if err != nil {
		logger.Warn("Failed to create article entity",
			zap.Error(err),
			zap.String("title", generated.Title),
		)
		return nil, domainerrors.ValidationError("article", err.Error())
	}

	savedArticle, err := u.articleRepo.Create(ctx, article)
	if err != nil {
		logger.Error("Failed to save article to repository",
			zap.Error(err),
			zap.String("title", generated.Title),
		)
		return nil, err
	}

	logger.Info("Successfully generated and saved article",
		zap.Int64("id", savedArticle.ID),
		zap.String("title", savedArticle.Title),
		zap.String("url", savedArticle.URL),
		zap.Strings("tags", savedArticle.Tags),
	)

	return savedArticle, nil
}
