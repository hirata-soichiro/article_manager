package usecase

import (
	"context"

	"article-manager/internal/domain/entity"
	domainerrors "article-manager/internal/domain/errors"
	"article-manager/internal/domain/repository"
	"article-manager/internal/infrastructure/logger"

	"go.uber.org/zap"
)

// タグに関するビジネスロジック
type TagUsecase struct {
	repo repository.TagRepository
}

func NewTagUsecase(repo repository.TagRepository) *TagUsecase {
	return &TagUsecase{repo: repo}
}

// 新しいタグを作成
func (u *TagUsecase) CreateTag(ctx context.Context, name string) (*entity.Tag, error) {
	logger.Debug("Creating tag",
		zap.String("name", name),
	)

	tag, err := entity.NewTag(name)
	if err != nil {
		logger.Warn("Failed to create tag entity",
			zap.Error(err),
			zap.String("name", name),
		)
		return nil, domainerrors.ValidationError("tag", err.Error())
	}

	savedTag, err := u.repo.Create(ctx, tag)
	if err != nil {
		logger.Error("Failed to save tag to repository",
			zap.Error(err),
			zap.String("name", name),
		)
		return nil, err
	}

	logger.Info("Successfully created tag",
		zap.Int64("id", savedTag.ID),
		zap.String("name", savedTag.Name),
	)

	return savedTag, nil
}

// 指定されたIDのタグを取得
func (u *TagUsecase) GetTagByID(ctx context.Context, id int64) (*entity.Tag, error) {
	logger.Debug("Getting tag by ID",
		zap.Int64("id", id),
	)

	if id <= 0 {
		logger.Warn("Invalid tag ID",
			zap.Int64("id", id),
		)
		return nil, domainerrors.InvalidArgumentError("id", "id must be positive")
	}

	tag, err := u.repo.FindByID(ctx, id)
	if err != nil {
		logger.Warn("Failed to find tag",
			zap.Error(err),
			zap.Int64("id", id),
		)
		return nil, err
	}

	logger.Debug("Successfully retrieved tag",
		zap.Int64("id", id),
		zap.String("name", tag.Name),
	)

	return tag, nil
}

// 指定された名前のタグを取得
func (u *TagUsecase) GetTagByName(ctx context.Context, name string) (*entity.Tag, error) {
	logger.Debug("Getting tag by name",
		zap.String("name", name),
	)

	if name == "" {
		logger.Warn("Tag name is empty")
		return nil, domainerrors.InvalidArgumentError("name", "name is required")
	}

	tag, err := u.repo.FindByName(ctx, name)
	if err != nil {
		logger.Warn("Failed to find tag by name",
			zap.Error(err),
			zap.String("name", name),
		)
		return nil, err
	}

	logger.Debug("Successfully retrieved tag by name",
		zap.Int64("id", tag.ID),
		zap.String("name", tag.Name),
	)

	return tag, nil
}

// 全てのタグを取得
func (u *TagUsecase) GetAllTags(ctx context.Context) ([]*entity.Tag, error) {
	logger.Debug("Getting all tags")

	tags, err := u.repo.FindAll(ctx)
	if err != nil {
		logger.Error("Failed to retrieve all tags",
			zap.Error(err),
		)
		return nil, err
	}

	logger.Debug("Successfully retrieved all tags",
		zap.Int("count", len(tags)),
	)

	return tags, nil
}

// タグを更新
func (u *TagUsecase) UpdateTag(ctx context.Context, id int64, name string) (*entity.Tag, error) {
	logger.Debug("Updating tag",
		zap.Int64("id", id),
		zap.String("name", name),
	)

	if id <= 0 {
		logger.Warn("Invalid tag ID",
			zap.Int64("id", id),
		)
		return nil, domainerrors.InvalidArgumentError("id", "id must be positive")
	}

	tag, err := u.repo.FindByID(ctx, id)
	if err != nil {
		logger.Warn("Failed to find tag",
			zap.Error(err),
			zap.Int64("id", id),
		)
		return nil, err
	}

	if err := tag.Update(name); err != nil {
		logger.Warn("Failed to update tag entity",
			zap.Error(err),
			zap.Int64("id", id),
		)
		return nil, domainerrors.ValidationError("tag", err.Error())
	}

	updatedTag, err := u.repo.Update(ctx, tag)
	if err != nil {
		logger.Error("Failed to update tag in repository",
			zap.Error(err),
			zap.Int64("id", id),
		)
		return nil, err
	}

	logger.Info("Successfully updated tag",
		zap.Int64("id", updatedTag.ID),
		zap.String("name", updatedTag.Name),
	)

	return updatedTag, nil
}

// 指定されたIDのタグを削除
func (u *TagUsecase) DeleteTag(ctx context.Context, id int64) error {
	logger.Debug("Deleting tag",
		zap.Int64("id", id),
	)

	if id <= 0 {
		logger.Warn("Invalid tag ID",
			zap.Int64("id", id),
		)
		return domainerrors.InvalidArgumentError("id", "id must be positive")
	}

	if err := u.repo.Delete(ctx, id); err != nil {
		logger.Error("Failed to delete tag from repository",
			zap.Error(err),
			zap.Int64("id", id),
		)
		return err
	}

	logger.Info("Successfully deleted tag",
		zap.Int64("id", id),
	)

	return nil
}
