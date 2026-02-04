package repository

import (
	"context"
	"database/sql"
	"strings"

	"article-manager/internal/domain/entity"
	domainerrors "article-manager/internal/domain/errors"
	"article-manager/internal/domain/repository"
	"article-manager/internal/infrastructure/logger"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
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
		logger.Error("Attempted to create nil tag")
		return nil, domainerrors.InvalidArgumentError("tag", "tag cannot be nil")
	}

	logger.Debug("Creating tag in database",
		zap.String("name", tag.Name),
	)

	query := `INSERT INTO tags (name, created_at, updated_at) VALUES (?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query, tag.Name, tag.CreatedAt, tag.UpdatedAt)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 {
				logger.Debug("Tag already exists",
					zap.String("name", tag.Name),
				)
				return nil, domainerrors.AlreadyExistsError("tag", tag.Name)
			}
		}
		logger.Error("Failed to insert tag",
			zap.Error(err),
			zap.String("name", tag.Name),
		)
		return nil, domainerrors.DatabaseError("insert tag", err)
	}

	tagID, err := result.LastInsertId()
	if err != nil {
		logger.Error("Failed to get last insert ID",
			zap.Error(err),
		)
		return nil, domainerrors.DatabaseError("get last insert id", err)
	}

	logger.Info("Successfully created tag in database",
		zap.Int64("id", tagID),
		zap.String("name", tag.Name),
	)

	return r.FindByID(ctx, tagID)
}

// 指定されたIDのタグを取得
func (r *mysqlTagRepository) FindByID(ctx context.Context, id int64) (*entity.Tag, error) {
	if id <= 0 {
		logger.Warn("Invalid tag ID",
			zap.Int64("id", id),
		)
		return nil, domainerrors.InvalidArgumentError("id", "id must be positive")
	}

	logger.Debug("Finding tag by ID",
		zap.Int64("id", id),
	)

	query := `SELECT id, name, created_at, updated_at FROM tags WHERE id = ?`

	var row tagRow
	err := r.db.GetContext(ctx, &row, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Debug("Tag not found",
				zap.Int64("id", id),
			)
			return nil, domainerrors.NotFoundError("tag", id)
		}
		logger.Error("Failed to find tag",
			zap.Error(err),
			zap.Int64("id", id),
		)
		return nil, domainerrors.DatabaseError("find tag", err)
	}

	logger.Debug("Successfully found tag",
		zap.Int64("id", id),
		zap.String("name", row.Name),
	)

	return tagRowToEntity(&row)
}

// 指定された名前のタグを取得
func (r *mysqlTagRepository) FindByName(ctx context.Context, name string) (*entity.Tag, error) {
	if strings.TrimSpace(name) == "" {
		logger.Warn("Invalid tag name",
			zap.String("name", name),
		)
		return nil, domainerrors.InvalidArgumentError("name", "name is required")
	}

	logger.Debug("Finding tag by name",
		zap.String("name", name),
	)

	query := `SELECT id, name, created_at, updated_at FROM tags WHERE name = ?`

	var row tagRow
	err := r.db.GetContext(ctx, &row, query, name)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Debug("Tag not found",
				zap.String("name", name),
			)
			return nil, domainerrors.NotFoundError("tag", name)
		}
		logger.Error("Failed to find tag",
			zap.Error(err),
			zap.String("name", name),
		)
		return nil, domainerrors.DatabaseError("find tag", err)
	}

	logger.Debug("Successfully found tag",
		zap.Int64("id", row.ID),
		zap.String("name", row.Name),
	)

	return tagRowToEntity(&row)
}

// 全てのタグを取得
func (r *mysqlTagRepository) FindAll(ctx context.Context) ([]*entity.Tag, error) {
	logger.Debug("Finding all tags")

	query := `SELECT id, name, created_at, updated_at FROM tags ORDER BY name ASC`

	var rows []tagRow
	err := r.db.SelectContext(ctx, &rows, query)
	if err != nil {
		logger.Error("Failed to find all tags",
			zap.Error(err),
		)
		return nil, domainerrors.DatabaseError("find all tags", err)
	}

	tags := make([]*entity.Tag, 0, len(rows))
	for _, row := range rows {
		tag, err := tagRowToEntity(&row)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	logger.Debug("Successfully found all tags",
		zap.Int("count", len(tags)),
	)

	return tags, nil
}

// タグを更新
func (r *mysqlTagRepository) Update(ctx context.Context, tag *entity.Tag) (*entity.Tag, error) {
	if tag == nil {
		logger.Error("Attempted to update nil tag")
		return nil, domainerrors.InvalidArgumentError("tag", "tag cannot be nil")
	}
	if tag.ID <= 0 {
		logger.Warn("Invalid tag ID for update",
			zap.Int64("id", tag.ID),
		)
		return nil, domainerrors.InvalidArgumentError("id", "id must be positive")
	}

	logger.Debug("Updating tag in database",
		zap.Int64("id", tag.ID),
		zap.String("name", tag.Name),
	)

	query := `UPDATE tags SET name = ?, updated_at = ? WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, tag.Name, tag.UpdatedAt, tag.ID)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 {
				logger.Debug("Tag name already exists",
					zap.String("name", tag.Name),
				)
				return nil, domainerrors.AlreadyExistsError("tag", tag.Name)
			}
		}
		logger.Error("Failed to update tag",
			zap.Error(err),
			zap.Int64("id", tag.ID),
		)
		return nil, domainerrors.DatabaseError("update tag", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Error("Failed to get rows affected",
			zap.Error(err),
		)
		return nil, domainerrors.DatabaseError("get rows affected", err)
	}
	if rowsAffected == 0 {
		logger.Debug("Tag not found for update",
			zap.Int64("id", tag.ID),
		)
		return nil, domainerrors.NotFoundError("tag", tag.ID)
	}

	logger.Info("Successfully updated tag in database",
		zap.Int64("id", tag.ID),
		zap.String("name", tag.Name),
	)

	return r.FindByID(ctx, tag.ID)
}

// 指定されたIDのタグを削除
func (r *mysqlTagRepository) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		logger.Warn("Invalid tag ID for deletion",
			zap.Int64("id", id),
		)
		return domainerrors.InvalidArgumentError("id", "id must be positive")
	}

	logger.Debug("Deleting tag from database",
		zap.Int64("id", id),
	)

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		logger.Error("Failed to begin transaction",
			zap.Error(err),
			zap.String("operation", "Delete"),
		)
		return domainerrors.DatabaseError("begin transaction", err)
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
		logger.Error("Failed to delete article tags",
			zap.Error(err),
			zap.Int64("tag_id", id),
		)
		return domainerrors.DatabaseError("delete article tags", err)
	}

	deleteTagQuery := `DELETE FROM tags WHERE id = ?`
	result, err := tx.ExecContext(ctx, deleteTagQuery, id)
	if err != nil {
		_ = tx.Rollback()
		logger.Error("Failed to delete tag",
			zap.Error(err),
			zap.Int64("id", id),
		)
		return domainerrors.DatabaseError("delete tag", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		logger.Error("Failed to get rows affected",
			zap.Error(err),
		)
		return domainerrors.DatabaseError("get rows affected", err)
	}
	if rowsAffected == 0 {
		_ = tx.Rollback()
		logger.Debug("Tag not found for deletion",
			zap.Int64("id", id),
		)
		return domainerrors.NotFoundError("tag", id)
	}

	if err := tx.Commit(); err != nil {
		logger.Error("Failed to commit transaction",
			zap.Error(err),
			zap.Int64("tag_id", id),
		)
		return domainerrors.DatabaseError("commit transaction", err)
	}

	logger.Info("Successfully deleted tag from database",
		zap.Int64("id", id),
	)

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
