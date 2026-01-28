package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"article-manager/internal/domain/entity"
	"article-manager/internal/domain/repository"

	"github.com/jmoiron/sqlx"
)

// articlesテーブルとのマッピング
type articleRow struct {
	ID        int64          `db:"id"`
	Title     string         `db:"title"`
	URL       string         `db:"url"`
	Summary   string         `db:"summary"`
	Memo      sql.NullString `db:"memo"`
	CreatedAt sql.NullTime   `db:"created_at"`
	UpdatedAt sql.NullTime   `db:"updated_at"`
}

type articleWithTagRow struct {
	ID        int64          `db:"id"`
	Title     string         `db:"title"`
	URL       string         `db:"url"`
	Summary   string         `db:"summary"`
	Memo      sql.NullString `db:"memo"`
	CreatedAt sql.NullTime   `db:"created_at"`
	UpdatedAt sql.NullTime   `db:"updated_at"`
	TagName   sql.NullString `db:"tag_name"`
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

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	var memo sql.NullString
	if article.Memo != "" {
		memo = sql.NullString{String: article.Memo, Valid: true}
	}

	query := `INSERT INTO articles (title, url, summary, memo, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`

	result, err := tx.ExecContext(ctx, query, article.Title, article.URL, article.Summary, memo, article.CreatedAt, article.UpdatedAt)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("failed to insert article: %w", err)
	}

	articleID, err := result.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	// タグが指定されている場合、記事とタグの関連付けを保存
	if len(article.Tags) > 0 {
		if err := r.insertArticleTags(ctx, tx, articleID, article.Tags); err != nil {
			_ = tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return r.FindByID(ctx, articleID)
}

// 指定されたIDの記事を取得
func (r *mysqlArticleRepository) FindByID(ctx context.Context, id int64) (*entity.Article, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}

	query := `SELECT id, title, url, summary, memo, created_at, updated_at FROM articles WHERE id = ?`

	var row articleRow
	err := r.db.GetContext(ctx, &row, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("article not found: id=%d", id)
		}
		return nil, fmt.Errorf("failed to find article: %w", err)
	}

	tags, err := r.findTagsByArticleID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find tags: %w", err)
	}

	return rowToEntity(&row, tags)
}

// 全ての記事を取得
func (r *mysqlArticleRepository) FindAll(ctx context.Context) ([]*entity.Article, error) {
	query := `
			SELECT
				a.id,
				a.title,
				a.url,
				a.summary,
				a.memo,
				a.created_at,
				a.updated_at,
				t.name AS tag_name
			FROM articles a
			LEFT JOIN article_tags at ON a.id = at.article_id
			LEFT JOIN tags t ON at.tag_id = t.id
			ORDER BY a.created_at DESC, t.name ASC
	`

	var rows []articleWithTagRow
	err := r.db.SelectContext(ctx, &rows, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find all articles: %w", err)
	}

	articleMap := make(map[int64]*entity.Article)
	var articleOrder []int64

	for _, row := range rows {
		article, exists := articleMap[row.ID]
		if !exists {
			memo := ""
			if row.Memo.Valid {
				memo = row.Memo.String
			}

			article = &entity.Article{
				ID:        row.ID,
				Title:     row.Title,
				URL:       row.URL,
				Summary:   row.Summary,
				Tags:      []string{},
				Memo:      memo,
				CreatedAt: row.CreatedAt.Time,
				UpdatedAt: row.UpdatedAt.Time,
			}
			articleMap[row.ID] = article
			articleOrder = append(articleOrder, row.ID)
		}

		if row.TagName.Valid && row.TagName.String != "" {
			article.Tags = append(article.Tags, row.TagName.String)
		}
	}

	articles := make([]*entity.Article, 0, len(articleOrder))
	for _, id := range articleOrder {
		articles = append(articles, articleMap[id])
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

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	var memo sql.NullString
	if article.Memo != "" {
		memo = sql.NullString{String: article.Memo, Valid: true}
	}

	query := `UPDATE articles SET title = ?, url = ?, summary = ?, memo = ?, updated_at = ? WHERE id = ?`

	result, err := tx.ExecContext(ctx, query, article.Title, article.URL, article.Summary, memo, article.UpdatedAt, article.ID)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("failed to update article: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		_ = tx.Rollback()
		return nil, fmt.Errorf("article not found: id=%d", article.ID)
	}

	deleteQuery := `DELETE FROM article_tags WHERE article_id = ?`
	_, err = tx.ExecContext(ctx, deleteQuery, article.ID)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("failed to delete article tags: %w", err)
	}

	if len(article.Tags) > 0 {
		if err := r.insertArticleTags(ctx, tx, article.ID, article.Tags); err != nil {
			_ = tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return r.FindByID(ctx, article.ID)
}

// 指定されたIDの記事を削除
func (r *mysqlArticleRepository) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("invalid id")
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	deleteTagsQuery := `DELETE FROM article_tags WHERE article_id = ?`
	_, err = tx.ExecContext(ctx, deleteTagsQuery, id)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to delete article: %w", err)
	}

	deleteArticleQuery := `DELETE FROM articles WHERE id = ?`
	result, err := tx.ExecContext(ctx, deleteArticleQuery, id)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to delete article: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		_ = tx.Rollback()
		return fmt.Errorf("article not found: id=%d", id)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *mysqlArticleRepository) insertArticleTags(ctx context.Context, tx *sqlx.Tx, articleID int64, tagNames []string) error {
	for _, tagName := range tagNames {
		// タグ名からタグIDを取得
		var tagID int64
		query := `SELECT id FROM tags WHERE name = ?`
		err := tx.GetContext(ctx, &tagID, query, tagName)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// タグが存在しない場合は新規作成
				insertTagQuery := `INSERT INTO tags (name, created_at, updated_at) VALUES (?, NOW(), NOW())`
				result, err := tx.ExecContext(ctx, insertTagQuery, tagName)
				if err != nil {
					return fmt.Errorf("failed to create tag: %w", err)
				}
				newTagID, err := result.LastInsertId()
				if err != nil {
					return fmt.Errorf("failed to get new tag id: %w", err)
				}
				tagID = newTagID
			} else {
				return fmt.Errorf("failed to find tag: %w", err)
			}
		}

		// article_tagsテーブルに関連付けを保存
		insertQuery := `INSERT INTO article_tags (article_id, tag_id) VALUES (?, ?)`
		_, err = tx.ExecContext(ctx, insertQuery, articleID, tagID)
		if err != nil {
			return fmt.Errorf("failed to insert article tag: %w", err)
		}
	}

	return nil
}

func (r *mysqlArticleRepository) findTagsByArticleID(ctx context.Context, articleID int64) ([]string, error) {
	query := `
			SELECT t.name
			FROM tags t
			INNER JOIN article_tags at ON t.id = at.tag_id
			WHERE at.article_id = ?
			ORDER BY t.name ASC
	`

	var tags []string
	err := r.db.SelectContext(ctx, &tags, query, articleID)
	if err != nil {
		return nil, fmt.Errorf("failed to select tags: %w", err)
	}
	if tags == nil {
		tags = []string{}
	}

	return tags, nil
}

// artibleRowをentity.Articleに変換
func rowToEntity(row *articleRow, tags []string) (*entity.Article, error) {
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
