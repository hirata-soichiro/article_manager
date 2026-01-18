package repository

import (
	"context"

	"article-manager/internal/domain/entity"
)

// タグデータへのアクセス操作を定義
type TagRepository interface {
	// 新しいタグを保存
	Create(ctx context.Context, tag *entity.Tag) (*entity.Tag, error)

	// 指定されたIDのタグを取得
	FindByID(ctx context.Context, id int64) (*entity.Tag, error)

	// 指定されたIDのタグを取得
	FindByName(ctx context.Context, name string) (*entity.Tag, error)

	// すべてのタグを取得
	FindAll(ctx context.Context) ([]*entity.Tag, error)

	// タグを更新
	Update(ctx context.Context, tag *entity.Tag) (*entity.Tag, error)

	// 指定されたIDのタグを削除
	Delete(ctx context.Context, id int64) error
}
