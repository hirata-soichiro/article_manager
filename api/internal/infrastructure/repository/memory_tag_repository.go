package repository

import (
	"context"
	"errors"
	"sort"
	"sync"

	"article-manager/internal/domain/entity"
	"article-manager/internal/domain/repository"
)

// メモリ上でタグを管理するリポジトリ
type MemoryTagRepository struct {
	tags   map[int64]*entity.Tag
	nextID int64
	mu     sync.RWMutex
}

// 新しいインメモリリポジトリの作成
func NewMemoryTagRepository() repository.TagRepository {
	return &MemoryTagRepository{
		tags:   make(map[int64]*entity.Tag),
		nextID: 1,
	}
}

// 新しいタグを保存
func (r *MemoryTagRepository) Create(ctx context.Context, tag *entity.Tag) (*entity.Tag, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// IDを自動裁判
	tag.ID = r.nextID
	r.nextID++

	saved := *tag
	r.tags[saved.ID] = &saved

	return &saved, nil
}

// 指定されたIDのタグを取得
func (r *MemoryTagRepository) FindByID(ctx context.Context, id int64) (*entity.Tag, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tag, exists := r.tags[id]
	if !exists {
		return nil, errors.New("tag not found")
	}

	result := *tag
	return &result, nil
}

// 指定された名前のタグを取得
func (r *MemoryTagRepository) FindByName(ctx context.Context, name string) (*entity.Tag, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, tag := range r.tags {
		if tag.Name == name {
			result := *tag
			return &result, nil
		}
	}

	return nil, errors.New("tag not found")
}

// すべてのタグを取得
func (r *MemoryTagRepository) FindAll(ctx context.Context) ([]*entity.Tag, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	result := make([]*entity.Tag, 0, len(r.tags))
	for _, tag := range r.tags {
		copied := *tag
		result = append(result, &copied)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})

	return result, nil
}

// タグを更新
func (r *MemoryTagRepository) Update(ctx context.Context, tag *entity.Tag) (*entity.Tag, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tags[tag.ID]; !exists {
		return nil, errors.New("tag not found")
	}

	updated := *tag
	r.tags[updated.ID] = &updated

	return &updated, nil
}

// 指定されたIDのタグを削除
func (r *MemoryTagRepository) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tags[id]; !exists {
		return errors.New("tag not found")
	}

	delete(r.tags, id)
	return nil
}
