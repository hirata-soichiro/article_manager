package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"article-manager/internal/domain/entity"
	"article-manager/internal/domain/service"
	"article-manager/internal/usecase"
)

// 記事自動生成ハンドラー
type ArticleGeneratorHandler struct {
	generatorUsecase *usecase.ArticleGeneratorUsecase
}

// コンストラクタ
func NewArticleGeneratorHandler(uc *usecase.ArticleGeneratorUsecase) *ArticleGeneratorHandler {
	return &ArticleGeneratorHandler{
		generatorUsecase: uc,
	}
}

// 記事自動生成リクエストの構造体
type GenerateArticleRequest struct {
	URL  string `json:"url"`
	Memo string `json:"memo"`
}

// URLから記事を自動生成
func (h *ArticleGeneratorHandler) GenerateArticle(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var req GenerateArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	article, err := h.generatorUsecase.GenerateArticleFromURL(ctx, req.URL, req.Memo)
	if err != nil {
		h.handleGeneratorError(w, err)
		return
	}

	h.respondJSON(w, http.StatusOK, h.toArticleResponse(article))
}

// AI生成エラーのハンドリング
func (h *ArticleGeneratorHandler) handleGeneratorError(w http.ResponseWriter, err error) {
	if aiErr, ok := err.(*service.AIGeneratorError); ok {
		switch aiErr.Code {
		case service.ErrCodeAPILimit:
			h.respondError(w, http.StatusTooManyRequests, aiErr.Message)
		case service.ErrCodeTimeout:
			h.respondError(w, http.StatusGatewayTimeout, aiErr.Message)
		case service.ErrCodeInvalidResponse:
			h.respondError(w, http.StatusBadGateway, aiErr.Message)
		case service.ErrCodeUnauthorized:
			h.respondError(w, http.StatusUnauthorized, aiErr.Message)
		case service.ErrCodeContentBlocked:
			h.respondError(w, http.StatusForbidden, aiErr.Message)
		case service.ErrCodeNetworkError:
			h.respondError(w, http.StatusBadGateway, aiErr.Message)
		case service.ErrCodeInvalidURL:
			h.respondError(w, http.StatusBadRequest, aiErr.Message)
		default:
			h.respondError(w, http.StatusInternalServerError, aiErr.Message)
		}
		return
	}

	// バリデーションエラー
	if h.isValidationError(err) {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// その他のエラー
	h.respondError(w, http.StatusInternalServerError, "failed to generate article")
}

// JSON形式でレスポンスを返す
func (h *ArticleGeneratorHandler) respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// エラーレスポンスを返す
func (h *ArticleGeneratorHandler) respondError(w http.ResponseWriter, statusCode int, message string) {
	h.respondJSON(w, statusCode, ErrorResponse{Error: message})
}

// エンティティをレスポンス形式に変換する
func (h *ArticleGeneratorHandler) toArticleResponse(article *entity.Article) ArticleResponse {
	// JSTに変換
	jst, _ := time.LoadLocation("Asia/Tokyo")

	return ArticleResponse{
		ID:        article.ID,
		Title:     article.Title,
		URL:       article.URL,
		Summary:   article.Summary,
		Tags:      article.Tags,
		Memo:      article.Memo,
		CreatedAt: article.CreatedAt.In(jst).Format("2006-01-02 15:04:05"),
		UpdatedAt: article.UpdatedAt.In(jst).Format("2006-01-02 15:04:05"),
	}
}

// バリデーションエラーかどうかを判定する
func (h *ArticleGeneratorHandler) isValidationError(err error) bool {
	if err == nil {
		return false
	}

	message := err.Error()
	validationKeywords := []string{
		"required",
		"must be",
		"cannot be empty",
		"must start with",
		"invalid url format",
	}

	for _, keyword := range validationKeywords {
		if strings.Contains(message, keyword) {
			return true
		}
	}

	return false
}
