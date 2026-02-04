package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"article-manager/internal/domain/entity"
	domainerrors "article-manager/internal/domain/errors"
	"article-manager/internal/infrastructure/logger"
	"article-manager/internal/usecase"

	"go.uber.org/zap"
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
	ctx := r.Context()

	logger.Info("Getting all tags",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
	)

	tags, err := h.usecase.GetAllTags(ctx)
	if err != nil {
		HandleError(w, err, "GetAllTags")
		return
	}

	response := make([]TagResponse, 0, len(tags))
	for _, tag := range tags {
		response = append(response, toTagResponse(tag))
	}

	logger.Info("Successfully retrieved all tags",
		zap.Int("count", len(tags)),
	)

	RespondSuccess(w, http.StatusOK, response)
}

// 指定されたIDのタグを取得
func (h *TagHandler) GetTagByID(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := r.Context()

	logger.Info("Getting tag by ID",
		zap.Int64("id", id),
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
	)

	tag, err := h.usecase.GetTagByID(ctx, id)
	if err != nil {
		HandleError(w, err, "GetTagByID")
		return
	}

	logger.Info("Successfully retrieved tag",
		zap.Int64("id", id),
		zap.String("name", tag.Name),
	)

	RespondSuccess(w, http.StatusOK, toTagResponse(tag))
}

// 新しいタグを作成する
func (h *TagHandler) CreateTag(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateTagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Warn("Failed to decode request body",
			zap.Error(err),
			zap.String("operation", "CreateTag"),
		)
		HandleError(w, domainerrors.InvalidArgumentError("request body", "invalid JSON format"), "CreateTag")
		return
	}

	logger.Info("Creating tag",
		zap.String("name", req.Name),
	)

	tag, err := h.usecase.CreateTag(ctx, req.Name)
	if err != nil {
		HandleError(w, err, "CreateTag")
		return
	}

	logger.Info("Successfully created tag",
		zap.Int64("id", tag.ID),
		zap.String("name", tag.Name),
	)

	RespondSuccess(w, http.StatusCreated, toTagResponse(tag))
}

// タグを更新する
func (h *TagHandler) UpdateTag(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := r.Context()

	var req UpdateTagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Warn("Failed to decode request body",
			zap.Error(err),
			zap.String("operation", "UpdateTag"),
			zap.Int64("id", id),
		)
		HandleError(w, domainerrors.InvalidArgumentError("request body", "invalid JSON format"), "UpdateTag")
		return
	}

	logger.Info("Updating tag",
		zap.Int64("id", id),
		zap.String("name", req.Name),
	)

	tag, err := h.usecase.UpdateTag(ctx, id, req.Name)
	if err != nil {
		HandleError(w, err, "UpdateTag")
		return
	}

	logger.Info("Successfully updated tag",
		zap.Int64("id", tag.ID),
		zap.String("name", tag.Name),
	)

	RespondSuccess(w, http.StatusOK, toTagResponse(tag))
}

// タグを削除する
func (h *TagHandler) DeleteTag(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := r.Context()

	logger.Info("Deleting tag",
		zap.Int64("id", id),
	)

	err := h.usecase.DeleteTag(ctx, id)
	if err != nil {
		HandleError(w, err, "DeleteTag")
		return
	}

	logger.Info("Successfully deleted tag",
		zap.Int64("id", id),
	)

	w.WriteHeader(http.StatusNoContent)
}

// エンティティをレスポンス形式に変換する
func toTagResponse(tag *entity.Tag) TagResponse {
	// JSTに変換
	jst, _ := time.LoadLocation("Asia/Tokyo")

	return TagResponse{
		ID:        tag.ID,
		Name:      tag.Name,
		CreatedAt: tag.CreatedAt.In(jst).Format("2006-01-02 15:04:05"),
		UpdatedAt: tag.UpdatedAt.In(jst).Format("2006-01-02 15:04:05"),
	}
}
