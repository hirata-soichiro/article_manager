package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"article-manager/internal/domain/entity"
	"article-manager/internal/domain/repository"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// tagsテーブルとのマッピング
type tagRow struct {
	ID        int64        `db:"id"`
	Name      string       `db:"name"`
	CreatedAt sql.NullTime `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

// TagRepositoryのMySQL実装
type mysqlTagRepository struct {
	db *sqlx.DB
}

// コンストラクタ
func NewMySQLTagRepository(db *sqlx.DB) repository.TagRepository {
	return &mysqlTagRepository{db: db}
}

// 新しいタグを保存
func (r *mysqlTagRepository) Create(ctx context.Context, tag *entity.Tag) (*entity.Tag, error) {
	if tag == nil {
		return nil, errors.New("tag is nil")
	}

	query := `INSERT INTO tags (name, created_at, updated_at) VALUES (?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query, tag.Name, tag.CreatedAt, tag.UpdatedAt)
	if err != nil {
		// UNIQUE制約違反のチェック
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 {
				return nil, fmt.Errorf("failed to insert tag: %w", err)
			}
		}
		return nil, fmt.Errorf("failed to insert tag: %w", err)
	}

	tagID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return r.FindByID(ctx, tagID)
}

// 指定されたIDのタグを取得
func (r *mysqlTagRepository) FindByID(ctx context.Context, id int64) (*entity.Tag, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}

	query := `SELECT id, name, created_at, updated_at FROM tags WHERE id = ?`

	var row tagRow
	err := r.db.GetContext(ctx, &row, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("tag not found: id=%d", id)
		}
		return nil, fmt.Errorf("failed to find tag: %w", err)
	}

	return tagRowToEntity(&row)
}

// 指定された名前のタグを取得
func (r *mysqlTagRepository) FindByName(ctx context.Context, name string) (*entity.Tag, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("name is required")
	}

	query := `SELECT id, name, created_at, updated_at FROM tags WHERE name = ?`

	var row tagRow
	err := r.db.GetContext(ctx, &row, query, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("tag not found: name=%s", name)
		}
		return nil, fmt.Errorf("failed to find tag: %w", err)
	}

	return tagRowToEntity(&row)
}

// 全てのタグを取得
func (r *mysqlTagRepository) FindAll(ctx context.Context) ([]*entity.Tag, error) {
	query := `SELECT id, name, created_at, updated_at FROM tags ORDER BY name ASC`

	var rows []tagRow
	err := r.db.SelectContext(ctx, &rows, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find all tags: %w", err)
	}

	tags := make([]*entity.Tag, 0, len(rows))
	for _, row := range rows {
		tag, err := tagRowToEntity(&row)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

// タグを更新
func (r *mysqlTagRepository) Update(ctx context.Context, tag *entity.Tag) (*entity.Tag, error) {
	if tag == nil {
		return nil, errors.New("tag is nil")
	}
	if tag.ID <= 0 {
		return nil, errors.New("invalid tag id")
	}

	query := `UPDATE tags SET name = ?, updated_at = ? WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, tag.Name, tag.UpdatedAt, tag.ID)
	if err != nil {
		// UNIQUE制約違反のチェック
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 {
				return nil, fmt.Errorf("tag with name '%s' already exists", tag.Name)
			}
		}
		return nil, fmt.Errorf("failed to update tag: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("tag not found: id=%d", tag.ID)
	}

	return r.FindByID(ctx, tag.ID)
}

// 指定されたIDのタグを削除
func (r *mysqlTagRepository) Delete(ctx context.Context, id int64) error {
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

	deleteArticleTagsQuery := `DELETE FROM article_tags WHERE tag_id = ?`
	_, err = tx.ExecContext(ctx, deleteArticleTagsQuery, id)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to delete article tags: %w", err)
	}

	deleteTagQuery := `DELETE FROM tags WHERE id = ?`
	result, err := tx.ExecContext(ctx, deleteTagQuery, id)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		_ = tx.Rollback()
		return fmt.Errorf("tag not found: id=%d", id)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// tagRowをentity.Tagに変換
func tagRowToEntity(row *tagRow) (*entity.Tag, error) {
	tag := &entity.Tag{
		ID:        row.ID,
		Name:      row.Name,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}

	return tag, nil
}
