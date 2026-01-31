package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"article-manager/internal/domain/entity"
	"article-manager/internal/usecase"
)

// 記事に関するHTTPハンドラ
type ArticleHandler struct {
	usecase *usecase.ArticleUsecase
}

// ArticleHandlerのコンストラクタ
func NewArticleHandler(uc *usecase.ArticleUsecase) *ArticleHandler {
	return &ArticleHandler{
		usecase: uc,
	}
}

// 記事作成リクエストの構造体
type CreateArticleRequest struct {
	Title   string   `json:"title"`
	URL     string   `json:"url"`
	Summary string   `json:"summary"`
	Tags    []string `json:"tags"`
	Memo    string   `json:"memo"`
}

// 記事更新リクエストの構造体
type UpdateArticleRequest struct {
	Title   string   `json:"title"`
	URL     string   `json:"url"`
	Summary string   `json:"summary"`
	Tags    []string `json:"tags"`
	Memo    string   `json:"memo"`
}

// エラーレスポンスの構造体
type ErrorResponse struct {
	Error string `json:"error"`
}

// 記事レスポンスの構造体
type ArticleResponse struct {
	ID        int64    `json:"id"`
	Title     string   `json:"title"`
	URL       string   `json:"url"`
	Summary   string   `json:"summary"`
	Tags      []string `json:"tags"`
	Memo      string   `json:"memo"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

// 全記事の取得
func (h *ArticleHandler) GetAllArticles(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	articles, err := h.usecase.GetAllArticles(ctx)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "failed to get articles")
		return
	}

	response := make([]ArticleResponse, 0, len(articles))
	for _, article := range articles {
		response = append(response, toArticleResponse(article))
	}

	h.respondJSON(w, http.StatusOK, response)
}

// 指定されたIDの記事を取得
func (h *ArticleHandler) GetArticleByID(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := context.Background()

	article, err := h.usecase.GetArticleByID(ctx, id)
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

	h.respondJSON(w, http.StatusOK, toArticleResponse(article))
}

// 新しい記事を作成する
func (h *ArticleHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var req CreateArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	article, err := h.usecase.CreateArticle(ctx, req.Title, req.URL, req.Summary, req.Tags, req.Memo)
	if err != nil {
		if isValidationError(err) {
			h.respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.respondError(w, http.StatusInternalServerError, "failed to create article")
		return
	}

	h.respondJSON(w, http.StatusCreated, toArticleResponse(article))
}

// 記事を更新する
func (h *ArticleHandler) UpdateArticle(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := context.Background()

	var req UpdateArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	article, err := h.usecase.UpdateArticle(ctx, id, req.Title, req.URL, req.Summary, req.Tags, req.Memo)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			h.respondError(w, http.StatusNotFound, "article not found")
			return
		}
		if isValidationError(err) {
			h.respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.respondError(w, http.StatusInternalServerError, "failed to update article")
		return
	}

	h.respondJSON(w, http.StatusOK, toArticleResponse(article))
}

// 記事を削除する
func (h *ArticleHandler) DeleteArticle(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := context.Background()

	err := h.usecase.DeleteArticle(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "invalid id") {
			h.respondError(w, http.StatusBadRequest, "invalid id")
			return
		}
		if strings.Contains(err.Error(), "not found") {
			h.respondError(w, http.StatusNotFound, "article not found")
			return
		}
		h.respondError(w, http.StatusInternalServerError, "failed to delete articles")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// 記事を検索する
func (h *ArticleHandler) SearchArticles(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	keyword := r.URL.Query().Get("keyword")

	trimmedKeyword := strings.TrimSpace(keyword)
	if trimmedKeyword == "" {
		h.respondError(w, http.StatusBadRequest, "keyword parameter is required and cannot be empty")
		return
	}

	articles, err := h.usecase.SearchArticles(ctx, trimmedKeyword)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "failed to search articles")
		return
	}

	response := make([]ArticleResponse, 0, len(articles))
	for _, article := range articles {
		response = append(response, toArticleResponse(article))
	}

	h.respondJSON(w, http.StatusOK, response)
}

// JSON形式でレスポンスを返す
func (h *ArticleHandler) respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// エラーレスポンスを返す
func (h *ArticleHandler) respondError(w http.ResponseWriter, statusCode int, message string) {
	h.respondJSON(w, statusCode, ErrorResponse{Error: message})
}

// エンティティをレスポンス形式に変換する
func toArticleResponse(article *entity.Article) ArticleResponse {
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
func isValidationError(err error) bool {
	if err == nil {
		return false
	}

	message := err.Error()
	validationKeywords := []string{
		"required",
		"must be",
		"cannot be empty",
		"must start with",
		"already exists",
	}

	for _, keyword := range validationKeywords {
		if strings.Contains(message, keyword) {
			return true
		}
	}

	return false
}
