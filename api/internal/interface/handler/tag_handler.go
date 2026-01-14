package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"article-manager/internal/domain/entity"
	"article-manager/internal/usecase"
)

// タグに関するHTTPハンドラ
type TagHandler struct {
	usecase *usecase.TagUsecase
}

// TagHandlerのコンストラクタ
func NewTagHandler(uc *usecase.TagUsecase) *TagHandler {
	return &TagHandler{
		usecase: uc,
	}
}

// タグ作成リクエストの構造体
type CreateTagRequest struct {
	Name string `json:"name"`
}

// タグ更新リクエストの構造体
type UpdateTagRequest struct {
	Name string `json:"name"`
}

// タグレスポンスの構造体
type TagResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// 全タグの取得
func (h *TagHandler) GetAllTags(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	tags, err := h.usecase.GetAllTags(ctx)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "failed to get tags")
		return
	}

	response := make([]TagResponse, 0, len(tags))
	for _, tag := range tags {
		response = append(response, toTagResponse(tag))
	}

	h.respondJSON(w, http.StatusOK, response)
}

// 指定されたIDのタグを取得
func (h *TagHandler) GetTagByID(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := context.Background()

	tag, err := h.usecase.GetTagByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "invalid id") {
			h.respondError(w, http.StatusBadRequest, "invalid id")
			return
		}
		if strings.Contains(err.Error(), "not found") {
			h.respondError(w, http.StatusNotFound, "article not found")
			return
		}
		h.respondError(w, http.StatusInternalServerError, "failed to get articles")
		return
	}

	h.respondJSON(w, http.StatusOK, toTagResponse(tag))
}

// 新しいタグを作成する
func (h *TagHandler) CreateTag(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var req CreateTagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	tag, err := h.usecase.CreateTag(ctx, req.Name)
	if err != nil {
		if isValidationError(err) {
			h.respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.respondError(w, http.StatusInternalServerError, "failed to create tag")
		return
	}

	h.respondJSON(w, http.StatusCreated, toTagResponse(tag))
}

// タグを更新する
func (h *TagHandler) UpdateTag(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := context.Background()

	// IDの事前チェック
	if id <= 0 {
		h.respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req UpdateTagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	tag, err := h.usecase.UpdateTag(ctx, id, req.Name)
	if err != nil {
		if strings.Contains(err.Error(), "invalid id") {
			h.respondError(w, http.StatusBadRequest, "invalid id")
			return
		}
		if strings.Contains(err.Error(), "not found") {
			h.respondError(w, http.StatusNotFound, "tag not found")
			return
		}
		if isValidationError(err) {
			h.respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.respondError(w, http.StatusInternalServerError, "failed to update article")
		return
	}

	h.respondJSON(w, http.StatusOK, toTagResponse(tag))
}

// タグを削除する
func (h *TagHandler) DeleteTag(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := context.Background()

	err := h.usecase.DeleteTag(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "invalid id") {
			h.respondError(w, http.StatusBadRequest, "invalid id")
			return
		}
		if strings.Contains(err.Error(), "not found") {
			h.respondError(w, http.StatusNotFound, "tag not found")
			return
		}
		h.respondError(w, http.StatusInternalServerError, "failed to delete tag")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// JSON形式でレスポンスを返す
func (h *TagHandler) respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// エラーレスポンスを返す
func (h *TagHandler) respondError(w http.ResponseWriter, statusCode int, message string) {
	h.respondJSON(w, statusCode, ErrorResponse{Error: message})
}

// エンティティをレスポンス形式に変換する
func toTagResponse(tag *entity.Tag) TagResponse {
	return TagResponse{
		ID:        tag.ID,
		Name:      tag.Name,
		CreatedAt: tag.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: tag.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
