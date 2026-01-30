package usecase

import (
	"context"
	"errors"
	"strings"

	"article-manager/internal/domain/entity"
	"article-manager/internal/domain/repository"
)

// 記事に関するビジネスロジック
type ArticleUsecase struct {
	repo repository.ArticleRepository
}

func NewArticleUsecase(repo repository.ArticleRepository) *ArticleUsecase {
	return &ArticleUsecase{repo: repo}
}

// 新しい記事を作成
func (u *ArticleUsecase) CreateArticle(ctx context.Context, title, url, summary string, tags []string, memo string) (*entity.Article, error) {
	article, err := entity.NewArticle(title, url, summary, tags, memo)
	if err != nil {
		return nil, err
	}

	savedArticle, err := u.repo.Create(ctx, article)
	if err != nil {
		return nil, err
	}

	return savedArticle, nil
}

// 指定されたIDの記事を取得
func (u *ArticleUsecase) GetArticleByID(ctx context.Context, id int64) (*entity.Article, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}

	article, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return article, nil
}

// 全ての記事を取得
func (u *ArticleUsecase) GetAllArticles(ctx context.Context) ([]*entity.Article, error) {
	articles, err := u.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

// 記事を更新
func (u *ArticleUsecase) UpdateArticle(ctx context.Context, id int64, title, url, summary string, tags []string, memo string) (*entity.Article, error) {
	article, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := article.Update(title, url, summary, tags, memo); err != nil {
		return nil, err
	}

	updatedArticle, err := u.repo.Update(ctx, article)
	if err != nil {
		return nil, err
	}

	return updatedArticle, nil
}

// 指定されたIDの記事を削除
func (u *ArticleUsecase) DeleteArticle(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("invalid id")
	}

	if err := u.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

// キーワードで記事を検索
func (u *ArticleUsecase) SearchArticles(ctx context.Context, keyword string) ([]*entity.Article, error) {
	trimmedKeyword := strings.TrimSpace(keyword)
	if trimmedKeyword == "" {
		return nil, errors.New("keyword cannot be empty")
	}

	articles, err := u.repo.Search(ctx, trimmedKeyword)
	if err != nil {
		return nil, err
	}

	return articles, nil
}
