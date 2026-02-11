package repository

import (
	"context"

	"article-manager/internal/domain/entity"
)

// 記事推薦キャッシュデータへのアクセス操作を定義
type BookRecommendationRepository interface {
	// 最新の有効なキャッシュを取得
	FindLatestValid(ctx context.Context) (*entity.BookRecommendationCache, error)
	// キャッシュを保存
	Save(ctx context.Context, cache *entity.BookRecommendationCache) (*entity.BookRecommendationCache, error)
}
