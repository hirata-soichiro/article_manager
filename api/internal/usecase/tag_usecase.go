package usecase

import (
	"context"
	"errors"

	"article-manager/internal/domain/entity"
	"article-manager/internal/domain/repository"
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
	tag, err := entity.NewTag(name)
	if err != nil {
		return nil, err
	}

	savedTag, err := u.repo.Create(ctx, tag)
	if err != nil {
		return nil, err
	}

	return savedTag, nil
}

// 指定されたIDのタグを取得
func (u *TagUsecase) GetTagByID(ctx context.Context, id int64) (*entity.Tag, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}

	tag, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

// 指定された名前のタグを取得
func (u *TagUsecase) GetTagByName(ctx context.Context, name string) (*entity.Tag, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}

	tag, err := u.repo.FindByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

// 全てのタグを取得
func (u *TagUsecase) GetAllTags(ctx context.Context) ([]*entity.Tag, error) {
	tags, err := u.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

// タグを更新
func (u *TagUsecase) UpdateTag(ctx context.Context, id int64, name string) (*entity.Tag, error) {
	tag, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := tag.Update(name); err != nil {
		return nil, err
	}

	updatedTag, err := u.repo.Update(ctx, tag)
	if err != nil {
		return nil, err
	}

	return updatedTag, nil
}

// 指定されたIDのタグを削除
func (u *TagUsecase) DeleteTag(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("invalid id")
	}

	if err := u.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
