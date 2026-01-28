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

// URLから記事情報を自動生成（DB保存なし）
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

	// タグ名のバリデーションのみ実行（DB作成はしない）
	tags := []string{}
	for _, tagName := range generated.SuggestedTags {
		if tagName == "" {
			return nil, errors.New("tag cannot be empty")
		}
		tags = append(tags, tagName)
	}

	// 生成された情報のみを返却（DBには保存しない）
	article := &entity.Article{
		ID:        0, // IDは0（まだ保存されていない）
		Title:     generated.Title,
		URL:       url,
		Summary:   generated.Summary,
		Tags:      tags,
		Memo:      memo,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return article, nil
}
