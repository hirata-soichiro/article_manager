package repository

import (
	"context"
	"errors"
	"sort"
	"strings"
	"sync"

	"article-manager/internal/domain/entity"
	"article-manager/internal/domain/repository"
)

// メモリ上で記事を管理するリポジトリ
type MemoryArticleRepository struct {
	articles map[int64]*entity.Article
	nextID   int64
	mu       sync.RWMutex
}

// 新しいインメモリリポジトリの作成
func NewMemoryArticleRepository() repository.ArticleRepository {
	return &MemoryArticleRepository{
		articles: make(map[int64]*entity.Article),
		nextID:   1,
	}
}

// 新しい記事を保存
func (r *MemoryArticleRepository) Create(ctx context.Context, article *entity.Article) (*entity.Article, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// IDを自動裁判
	article.ID = r.nextID
	r.nextID++

	saved := *article
	r.articles[saved.ID] = &saved

	return &saved, nil
}

// 指定されたIDの記事を取得
func (r *MemoryArticleRepository) FindByID(ctx context.Context, id int64) (*entity.Article, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	article, exists := r.articles[id]
	if !exists {
		return nil, errors.New("article not found")
	}

	result := *article
	return &result, nil
}

// すべての記事を取得
func (r *MemoryArticleRepository) FindAll(ctx context.Context) ([]*entity.Article, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	result := make([]*entity.Article, 0, len(r.articles))
	for _, article := range r.articles {
		copied := *article
		result = append(result, &copied)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})

	return result, nil
}

// 記事を更新
func (r *MemoryArticleRepository) Update(ctx context.Context, article *entity.Article) (*entity.Article, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.articles[article.ID]; !exists {
		return nil, errors.New("article not found")
	}

	updated := *article
	r.articles[updated.ID] = &updated

	return &updated, nil
}

// 指定されたIDの記事を削除
func (r *MemoryArticleRepository) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.articles[id]; !exists {
		return errors.New("article not found")
	}

	delete(r.articles, id)
	return nil
}

// 曖昧検索でタイトルまたは要約から記事を検索
func (r *MemoryArticleRepository) Search(ctx context.Context, keyword string) ([]*entity.Article, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	trimmedKeyword := strings.TrimSpace(keyword)

	// 空文字の場合は全件取得
	if trimmedKeyword == "" {
		result := make([]*entity.Article, 0, len(r.articles))
		for _, article := range r.articles {
			copied := *article
			result = append(result, &copied)
		}

		sort.Slice(result, func(i, j int) bool {
			return result[i].CreatedAt.After(result[j].CreatedAt)
		})

		return result, nil
	}

	keywords := strings.Fields(trimmedKeyword)

	result := make([]*entity.Article, 0)
	for _, article := range r.articles {
		// 全てのキーワードがタイトルまたは要約に含まれているかチェック
		allMatch := true
		for _, kw := range keywords {
			kwLower := strings.ToLower(kw)
			titleLower := strings.ToLower(article.Title)
			summaryLower := strings.ToLower(article.Summary)

			if !strings.Contains(titleLower, kwLower) && !strings.Contains(summaryLower, kwLower) {
				allMatch = false
				break
			}
		}

		if allMatch {
			copied := *article
			result = append(result, &copied)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.After(result[j].CreatedAt)
	})

	return result, nil
}
