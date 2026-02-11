package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"article-manager/internal/domain/entity"
	domainerrors "article-manager/internal/domain/errors"
	"article-manager/internal/domain/repository"
	"article-manager/internal/infrastructure/logger"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// book_recommendationsテーブルとのマッピング
type bookRecommendationRow struct {
	ID                  int64        `db:"id"`
	RecommendationsJSON string       `db:"recommendations_json"`
	CreatedAt           sql.NullTime `db:"created_at"`
	ExpiresAt           sql.NullTime `db:"expires_at"`
}

// BookRecommendationRepositoryのMySQL実装
type mysqlBookRecommendationRepository struct {
	db *sqlx.DB
}

// コンストラクタ
func NewMySQLBookRecommendationRepository(db *sqlx.DB) repository.BookRecommendationRepository {
	return &mysqlBookRecommendationRepository{db: db}
}

// 最新の有効なキャッシュを取得
func (r *mysqlBookRecommendationRepository) FindLatestValid(ctx context.Context) (*entity.BookRecommendationCache, error) {
	logger.Debug("Finding latest valid book recommendation cache")

	// 現在時刻より後に期限切れとなるキャッシュ(有効なキャッシュ)を取得
	query := `
		SELECT id, recommendations_json, created_at, expires_at
		FROM book_recommendations
		WHERE expires_at > NOW()
		ORDER BY created_at DESC
		LIMIT 1
    `

	var row bookRecommendationRow
	err := r.db.GetContext(ctx, &row, query)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Debug("No valid book recommendation cache found")
			return nil, domainerrors.NotFoundError("book_recommendation_cache", "valid cache")
		}
		logger.Error("Failed to find book recommendation cache",
			zap.Error(err),
		)
		return nil, domainerrors.DatabaseError("find book recommendation cache", err)
	}

	cache, err := rowToBookRecommendationCache(&row)
	if err != nil {
		return nil, err
	}

	logger.Debug("Successfully found valid book recommendation cache",
		zap.Int64("id", cache.ID),
		zap.Time("expires_at", cache.ExpiresAt),
	)

	return cache, nil
}

// キャッシュを保存
func (r *mysqlBookRecommendationRepository) Save(ctx context.Context, cache *entity.BookRecommendationCache) (*entity.BookRecommendationCache, error) {
	if cache == nil {
		logger.Error("Attempted to save nil book recommendation cache")
		return nil, domainerrors.InvalidArgumentError("cache", "cache cannnot be nil")
	}

	logger.Debug("Saving book recommendation cache",
		zap.Int("book_count", len(cache.Books)),
		zap.Time("expires_at", cache.ExpiresAt),
	)

	// BooksをJSONに変換
	booksJSON, err := json.Marshal(cache.Books)
	if err != nil {
		logger.Error("Failed to marshal books to JSON",
			zap.Error(err),
		)
		return nil, domainerrors.DatabaseError("marshal books", err)
	}

	// トランザクション開始
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		logger.Error("Failed to begin transaction",
			zap.Error(err),
			zap.String("operation", "Save"),
		)
		return nil, domainerrors.DatabaseError("begin transaction", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	// 既存のキャッシュを全て削除
	deleteQuery := `DELETE FROM book_recommendations`
	_, err = tx.ExecContext(ctx, deleteQuery)
	if err != nil {
		_ = tx.Rollback()
		logger.Error("Failed to delete old cache",
			zap.Error(err),
		)
		return nil, domainerrors.DatabaseError("delete old cache", err)
	}

	// 新しいキャッシュを挿入
	insertQuery := `
		INSERT INTO book_recommendations (recommendations_json, created_at, expires_at)
		VALUES (?, ?, ?)
    `

	result, err := tx.ExecContext(ctx, insertQuery, string(booksJSON), cache.GeneratedAt, cache.ExpiresAt)
	if err != nil {
		_ = tx.Rollback()
		logger.Error("Failed to insert book recommendation cache",
			zap.Error(err),
		)
		return nil, domainerrors.DatabaseError("insert book recommendation cache", err)
	}

	cacheID, err := result.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		logger.Error("Failed to get last insert ID",
			zap.Error(err),
		)
		return nil, domainerrors.DatabaseError("get last insert id", err)
	}

	// トランザクションコミット
	if err := tx.Commit(); err != nil {
		logger.Error("Failed to commit transaction",
			zap.Error(err),
			zap.Int64("cache_id", cacheID),
		)
		return nil, domainerrors.DatabaseError("commit transaction", err)
	}

	logger.Info("Successfully saved book recommendation cache",
		zap.Int64("id", cacheID),
		zap.Int("book_count", len(cache.Books)),
	)

	cache.ID = cacheID
	return cache, nil
}

// bookRecommendationRowをentity.BookRecommendationCacheに変換
func rowToBookRecommendationCache(row *bookRecommendationRow) (*entity.BookRecommendationCache, error) {
	var books []entity.Book
	if err := json.Unmarshal([]byte(row.RecommendationsJSON), &books); err != nil {
		logger.Error("Failed to unmarshal books JSON",
			zap.Error(err),
			zap.String("json", row.RecommendationsJSON),
		)
		return nil, domainerrors.DatabaseError("unmarshal books", err)
	}

	cache := &entity.BookRecommendationCache{
		ID:          row.ID,
		Books:       books,
		GeneratedAt: row.CreatedAt.Time,
		ExpiresAt:   row.ExpiresAt.Time,
	}

	return cache, nil
}
