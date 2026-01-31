package repository

import (
	"context"

	"article-manager/internal/domain/entity"
)

// 記事データへのアクセス操作を定義
type ArticleRepository interface {
	// 新しい記事を保存
	Create(ctx context.Context, article *entity.Article) (*entity.Article, error)

	// 指定されたIDの記事を取得
	FindByID(ctx context.Context, id int64) (*entity.Article, error)

	// すべての記事を取得
	FindAll(ctx context.Context) ([]*entity.Article, error)

	// 記事を更新
	Update(ctx context.Context, article *entity.Article) (*entity.Article, error)

	// 指定されたIDの記事を削除
	Delete(ctx context.Context, id int64) error

	// 曖昧検索でタイトルまたは要約から記事を検索
	Search(ctx context.Context, keyword string) ([]*entity.Article, error)
}
