package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"article-manager/internal/domain/entity"
	"article-manager/internal/domain/repository"

	"github.com/jmoiron/sqlx"
)

// articlesテーブルとのマッピング
type articleRow struct {
	ID        int64           `db:"id"`
	Title     string          `db:"title"`
	URL       string          `db:"url"`
	Summary   string          `db:"summary"`
	Tags      json.RawMessage `db:"tags"`
	Memo      sql.NullString  `db:"memo"`
	CreatedAt sql.NullTime    `db:"created_at"`
	UpdatedAt sql.NullTime    `db:"updated_at"`
}

// ArticeleRepositoryのMySQL実装
type mysqlArticleRepository struct {
	db *sqlx.DB
}

// コンストラクタ
func NewMySQLArticleRepository(db *sqlx.DB) repository.ArticleRepository {
	return &mysqlArticleRepository{db: db}
}

// 新しい記事を保存
func (r *mysqlArticleRepository) Create(ctx context.Context, article *entity.Article) (*entity.Article, error) {
	if article == nil {
		return nil, errors.New("article is nil")
	}

	tagsJson, err := json.Marshal(article.Tags)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tags: %w", err)
	}

	var memo sql.NullString
	if article.Memo != "" {
		memo = sql.NullString{String: article.Memo, Valid: true}
	}

	query := `INSERT INTO articles (title, url, summary, tags, memo, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query, article.Title, article.URL, article.Summary, tagsJson, memo, article.CreatedAt, article.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert article: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return r.FindByID(ctx, id)
}

// 指定されたIDの記事を取得
func (r *mysqlArticleRepository) FindByID(ctx context.Context, id int64) (*entity.Article, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}

	query := `SELECT id, title, url, summary, tags, memo, created_at, updated_at FROM articles WHERE id = ?`

	var row articleRow
	err := r.db.GetContext(ctx, &row, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("article not found: id=%d", id)
		}
		return nil, fmt.Errorf("failed to find article: %w", err)
	}

	return rowToEntity(&row)
}

// 全ての記事を取得
func (r *mysqlArticleRepository) FindAll(ctx context.Context) ([]*entity.Article, error) {
	query := `SELECT id, title, url, summary, tags, memo, created_at, updated_at FROM articles ORDER BY created_at DESC`

	var rows []articleRow
	err := r.db.SelectContext(ctx, &rows, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find all articles: %w", err)
	}

	articles := make([]*entity.Article, 0, len(rows))
	for _, row := range rows {
		article, err := rowToEntity(&row)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	return articles, nil
}

// 記事を更新
func (r *mysqlArticleRepository) Update(ctx context.Context, article *entity.Article) (*entity.Article, error) {
	if article == nil {
		return nil, errors.New("article is nil")
	}
	if article.ID <= 0 {
		return nil, errors.New("invalid article id")
	}

	tagsJson, err := json.Marshal(article.Tags)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tags: %w", err)
	}

	var memo sql.NullString
	if article.Memo != "" {
		memo = sql.NullString{String: article.Memo, Valid: true}
	}

	query := `UPDATE articles SET title = ?, url = ?, summary = ?, tags = ?, memo = ?, updated_at = ? WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, article.Title, article.URL, article.Summary, tagsJson, memo, article.UpdatedAt, article.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update article: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("article not found: id=%d", article.ID)
	}

	return r.FindByID(ctx, article.ID)
}

// 指定されたIDの記事を削除
func (r *mysqlArticleRepository) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("invalid id")
	}

	query := `DELETE FROM articles WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete article: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("article not found: id=%d", id)
	}

	return nil
}

// artibleRowをentity.Articleに変換
func rowToEntity(row *articleRow) (*entity.Article, error) {
	var tags []string
	if err := json.Unmarshal(row.Tags, &tags); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
	}

	memo := ""
	if row.Memo.Valid {
		memo = row.Memo.String
	}

	article := &entity.Article{
		ID:        row.ID,
		Title:     row.Title,
		URL:       row.URL,
		Summary:   row.Summary,
		Tags:      tags,
		Memo:      memo,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}

	return article, nil
}
