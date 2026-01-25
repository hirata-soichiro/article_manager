package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

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

// URLから記事を自動生成して保存
func (u *ArticleGeneratorUsecase) GenerateArticleFromURL(ctx context.Context, url string, memo string) (*entity.Article, error) {
	if url == "" {
		return nil, errors.New("url is required")
	}
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return nil, errors.New("invalid url format")
	}

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

	tags := []string{}
	for _, tagName := range generated.SuggestedTags {
		if tagName == "" {
			return nil, errors.New("tag cannot be empty")
		}

		exsitingTag, err := u.tagRepo.FindByName(ctx, tagName)
		if err != nil {
			if !strings.Contains(err.Error(), "tag not found") {
				return nil, err
			}
			newTag := &entity.Tag{
				Name:      tagName,
				CreatedAt: time.Now(),
			}
			createdTag, err := u.tagRepo.Create(ctx, newTag)
			if err != nil {
				return nil, err
			}
			tags = append(tags, createdTag.Name)
		} else {
			tags = append(tags, exsitingTag.Name)
		}
	}

	article := &entity.Article{
		Title:     generated.Title,
		URL:       url,
		Summary:   generated.Summary,
		Tags:      tags,
		Memo:      memo,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	savedArticle, err := u.articleRepo.Create(ctx, article)
	if err != nil {
		return nil, err
	}

	return savedArticle, nil
}
