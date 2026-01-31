package usecase

import (
	"context"
	"errors"
	"strings"

	"article-manager/internal/domain/entity"
	"article-manager/internal/domain/repository"
	"article-manager/internal/domain/service"
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
	if url == "" {
		return nil, errors.New("url is required")
	}
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return nil, errors.New("invalid url format")
	}

	// AIで記事情報を生成
	generated, err := u.aiGenerator.GenerateArticleFromURL(ctx, service.ArticleGenerationRequest{URL: url})
	if err != nil {
		return nil, err
	}

	if generated.Title == "" {
		return nil, errors.New("title is required")
	}
	if generated.Summary == "" {
		return nil, errors.New("summary is required")
	}

	// タグの処理（既存タグを検索し、なければ作成）
	tags := []string{}
	for _, tagName := range generated.SuggestedTags {
		if tagName == "" {
			return nil, errors.New("tag cannot be empty")
		}

		// タグを検索
		existingTag, err := u.tagRepo.FindByName(ctx, tagName)
		if err != nil {
			// タグが存在しない場合は新規作成
			newTag, err := entity.NewTag(tagName)
			if err != nil {
				return nil, err
			}
			createdTag, err := u.tagRepo.Create(ctx, newTag)
			if err != nil {
				return nil, err
			}
			tags = append(tags, createdTag.Name)
		} else {
			// タグが既存の場合
			tags = append(tags, existingTag.Name)
		}
	}

	// 記事エンティティを作成
	article, err := entity.NewArticle(generated.Title, url, generated.Summary, tags, memo)
	if err != nil {
		return nil, err
	}

	// 記事をDBに保存
	savedArticle, err := u.articleRepo.Create(ctx, article)
	if err != nil {
		return nil, err
	}

	return savedArticle, nil
}
