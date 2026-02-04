package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"article-manager/internal/domain/entity"
	domainerrors "article-manager/internal/domain/errors"
	"article-manager/internal/domain/repository"
	"article-manager/internal/infrastructure/logger"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
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
		logger.Error("Attempted to create nil article")
		return nil, domainerrors.InvalidArgumentError("article", "article cannot be nil")
	}

	logger.Debug("Creating article in database",
		zap.String("title", article.Title),
		zap.String("url", article.URL),
	)

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		logger.Error("Failed to begin transaction",
			zap.Error(err),
			zap.String("operation", "Create"),
		)
		return nil, domainerrors.DatabaseError("begin transaction", err)
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
		logger.Error("Failed to insert article",
			zap.Error(err),
			zap.String("title", article.Title),
		)
		return nil, domainerrors.DatabaseError("insert article", err)
	}

	articleID, err := result.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		logger.Error("Failed to get last insert ID",
			zap.Error(err),
		)
		return nil, domainerrors.DatabaseError("get last insert id", err)
	}

	// タグが指定されている場合、記事とタグの関連付けを保存
	if len(article.Tags) > 0 {
		if err := r.insertArticleTags(ctx, tx, articleID, article.Tags); err != nil {
			_ = tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		logger.Error("Failed to commit transaction",
			zap.Error(err),
			zap.Int64("article_id", articleID),
		)
		return nil, domainerrors.DatabaseError("commit transaction", err)
	}

	logger.Info("Successfully created article in database",
		zap.Int64("id", articleID),
		zap.String("title", article.Title),
	)

	return r.FindByID(ctx, articleID)
}

// 指定されたIDの記事を取得
func (r *mysqlArticleRepository) FindByID(ctx context.Context, id int64) (*entity.Article, error) {
	if id <= 0 {
		logger.Warn("Invalid article ID",
			zap.Int64("id", id),
		)
		return nil, domainerrors.InvalidArgumentError("id", "id must be positive")
	}

	logger.Debug("Finding article by ID",
		zap.Int64("id", id),
	)

	query := `SELECT id, title, url, summary, memo, created_at, updated_at FROM articles WHERE id = ?`

	var row articleRow
	err := r.db.GetContext(ctx, &row, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Debug("Article not found",
				zap.Int64("id", id),
			)
			return nil, domainerrors.NotFoundError("article", id)
		}
		logger.Error("Failed to find article",
			zap.Error(err),
			zap.Int64("id", id),
		)
		return nil, domainerrors.DatabaseError("find article", err)
	}

	tags, err := r.findTagsByArticleID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find tags: %w", err)
	}

	logger.Debug("Successfully found article",
		zap.Int64("id", id),
		zap.String("title", row.Title),
	)

	return rowToEntity(&row, tags)
}

// 全ての記事を取得
func (r *mysqlArticleRepository) FindAll(ctx context.Context) ([]*entity.Article, error) {
	logger.Debug("Finding all articles")

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
		logger.Error("Failed to find all articles",
			zap.Error(err),
		)
		return nil, domainerrors.DatabaseError("find all articles", err)
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

	logger.Debug("Successfully found all articles",
		zap.Int("count", len(articles)),
	)

	return articles, nil
}

// 記事を更新
func (r *mysqlArticleRepository) Update(ctx context.Context, article *entity.Article) (*entity.Article, error) {
	if article == nil {
		logger.Error("Attempted to update nil article")
		return nil, domainerrors.InvalidArgumentError("article", "article cannot be nil")
	}
	if article.ID <= 0 {
		logger.Warn("Invalid article ID for update",
			zap.Int64("id", article.ID),
		)
		return nil, domainerrors.InvalidArgumentError("id", "id must be positive")
	}

	logger.Debug("Updating article in database",
		zap.Int64("id", article.ID),
		zap.String("title", article.Title),
	)

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		logger.Error("Failed to begin transaction",
			zap.Error(err),
			zap.String("operation", "Update"),
		)
		return nil, domainerrors.DatabaseError("begin transaction", err)
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
		logger.Error("Failed to update article",
			zap.Error(err),
			zap.Int64("id", article.ID),
		)
		return nil, domainerrors.DatabaseError("update article", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		logger.Error("Failed to get rows affected",
			zap.Error(err),
		)
		return nil, domainerrors.DatabaseError("get rows affected", err)
	}
	if rowsAffected == 0 {
		_ = tx.Rollback()
		logger.Debug("Article not found for update",
			zap.Int64("id", article.ID),
		)
		return nil, domainerrors.NotFoundError("article", article.ID)
	}

	deleteQuery := `DELETE FROM article_tags WHERE article_id = ?`
	_, err = tx.ExecContext(ctx, deleteQuery, article.ID)
	if err != nil {
		_ = tx.Rollback()
		logger.Error("Failed to delete article tags",
			zap.Error(err),
			zap.Int64("article_id", article.ID),
		)
		return nil, domainerrors.DatabaseError("delete article tags", err)
	}

	if len(article.Tags) > 0 {
		if err := r.insertArticleTags(ctx, tx, article.ID, article.Tags); err != nil {
			_ = tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		logger.Error("Failed to commit transaction",
			zap.Error(err),
			zap.Int64("article_id", article.ID),
		)
		return nil, domainerrors.DatabaseError("commit transaction", err)
	}

	logger.Info("Successfully updated article in database",
		zap.Int64("id", article.ID),
		zap.String("title", article.Title),
	)

	return r.FindByID(ctx, article.ID)
}

// 指定されたIDの記事を削除
func (r *mysqlArticleRepository) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		logger.Warn("Invalid article ID for deletion",
			zap.Int64("id", id),
		)
		return domainerrors.InvalidArgumentError("id", "id must be positive")
	}

	logger.Debug("Deleting article from database",
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

	deleteTagsQuery := `DELETE FROM article_tags WHERE article_id = ?`
	_, err = tx.ExecContext(ctx, deleteTagsQuery, id)
	if err != nil {
		_ = tx.Rollback()
		logger.Error("Failed to delete article tags",
			zap.Error(err),
			zap.Int64("article_id", id),
		)
		return domainerrors.DatabaseError("delete article tags", err)
	}

	deleteArticleQuery := `DELETE FROM articles WHERE id = ?`
	result, err := tx.ExecContext(ctx, deleteArticleQuery, id)
	if err != nil {
		_ = tx.Rollback()
		logger.Error("Failed to delete article",
			zap.Error(err),
			zap.Int64("id", id),
		)
		return domainerrors.DatabaseError("delete article", err)
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
		logger.Debug("Article not found for deletion",
			zap.Int64("id", id),
		)
		return domainerrors.NotFoundError("article", id)
	}

	if err := tx.Commit(); err != nil {
		logger.Error("Failed to commit transaction",
			zap.Error(err),
			zap.Int64("article_id", id),
		)
		return domainerrors.DatabaseError("commit transaction", err)
	}

	logger.Info("Successfully deleted article from database",
		zap.Int64("id", id),
	)

	return nil
}

func (r *mysqlArticleRepository) insertArticleTags(ctx context.Context, tx *sqlx.Tx, articleID int64, tagNames []string) error {
	if len(tagNames) == 0 {
		return nil
	}

	// 既存タグを一括取得
	query, args, err := sqlx.In(`SELECT id, name FORM tags WHERE name IN (?)`, tagNames)
	if err != nil {
		return domainerrors.DatabaseError("prepare tag query", err)
	}
	query = tx.Rebind(query)

	type tagResult struct {
		ID   int64  `db:"id"`
		Name string `db:"name"`
	}
	var existingTags []tagResult
	err = tx.SelectContext(ctx, &existingTags, query, args...)
	if err != nil {
		logger.Error("Failed to fetch existing tags",
			zap.Error(err),
		)
		return domainerrors.DatabaseError("fetch existing tags", err)
	}

	// タグ名からIDへのマッピングを作成
	tagIDMap := make(map[string]int64, len(existingTags))
	for _, tag := range existingTags {
		tagIDMap[tag.Name] = tag.ID
	}

	// 存在しないタグを抽出して一括作成
	var missingTags []string
	for _, tagName := range tagNames {
		if _, exists := tagIDMap[tagName]; !exists {
			missingTags = append(missingTags, tagName)
		}
	}

	if len(missingTags) > 0 {
		valueStrings := make([]string, 0, len(missingTags))
		valueArgs := make([]interface{}, 0, len(missingTags))
		now := "NOW()"

		for _, tagName := range missingTags {
			valueStrings = append(valueStrings, "(?, "+now+", "+now+")")
			valueArgs = append(valueArgs, tagName)
		}

		buldInsertQuery := fmt.Sprintf(
			"INSERT INTO tags (name, created_at, updated_at) VALUES %s",
			strings.Join(valueStrings, ","),
		)

		result, err := tx.ExecContext(ctx, buldInsertQuery, valueArgs...)
		if err != nil {
			logger.Error("Failed to bulk insert tags",
				zap.Error(err),
			)
			return domainerrors.DatabaseError("insert new tags", err)
		}

		// 新規作成されたタグのIDを取得
		lastID, err := result.LastInsertId()
		if err != nil {
			return domainerrors.DatabaseError("get last insert id", err)
		}

		// 新規タグのIDをマッピングに追加
		for i, tagName := range missingTags {
			tagIDMap[tagName] = lastID + int64(i)
		}

		logger.Debug("Inserted new tags",
			zap.Int("count", len(missingTags)),
		)
	}

	// article_tagsテーブルに一括挿入
	if len(tagNames) > 0 {
		valueStrings := make([]string, 0, len(tagNames))
		valueArgs := make([]interface{}, 0, len(tagNames)*2)

		for _, tagName := range tagNames {
			tagID := tagIDMap[tagName]
			valueStrings = append(valueStrings, "(?, ?)")
			valueArgs = append(valueArgs, articleID, tagID)
		}

		buldInsertQuery := fmt.Sprintf(
			"INSERT INTO article_tags (article_id, tag_id) VALUES %s",
			strings.Join(valueStrings, ","),
		)

		_, err = tx.ExecContext(ctx, buldInsertQuery, valueArgs...)
		if err != nil {
			logger.Error("Failed to bulk insert article tags",
				zap.Error(err),
				zap.Int64("article_id", articleID),
			)
			return domainerrors.DatabaseError("bulk insert article tags", err)
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
		logger.Error("Failed to select tags",
			zap.Error(err),
			zap.Int64("article_id", articleID),
		)
		return nil, domainerrors.DatabaseError("select tags", err)
	}
	if tags == nil {
		tags = []string{}
	}

	return tags, nil
}

// 曖昧検索機能
func (r *mysqlArticleRepository) Search(ctx context.Context, keyword string) ([]*entity.Article, error) {
	logger.Debug("Searching articles",
		zap.String("keyword", keyword),
	)

	trimmedKeyword := strings.TrimSpace(keyword)
	if trimmedKeyword == "" {
		return r.FindAll(ctx)
	}

	keywords := strings.Fields(trimmedKeyword)

	// FULLTEXTインデックスを活用した検索（ngramパーサー使用）
	var matchConditions []string
	var args []interface{}

	// MATCH AGAINST構文を使用（FULLTEXTインデックスを活用）
	for _, kw := range keywords {
		matchConditions = append(matchConditions, "MATCH(a.title, a.summary) AGAINST(? IN BOOLEAN MODE)")
		// BOOLEAN MODEで部分一致検索を可能にする
		args = append(args, "*"+kw+"*")
	}

	whereClause := strings.Join(matchConditions, " AND ")

	query := fmt.Sprintf(`
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
		WHERE %s
		ORDER BY a.created_at DESC, t.name ASC
	`, whereClause)

	var rows []articleWithTagRow
	err := r.db.SelectContext(ctx, &rows, query, args...)
	if err != nil {
		logger.Error("Failed to search articles",
			zap.Error(err),
			zap.String("keyword", trimmedKeyword),
		)
		return nil, domainerrors.DatabaseError("search articles", err)
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

	logger.Info("Successfully searched articles",
		zap.String("keyword", trimmedKeyword),
		zap.Int("count", len(articles)),
	)

	return articles, nil
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
