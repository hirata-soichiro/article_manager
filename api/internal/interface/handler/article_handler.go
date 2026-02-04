package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"article-manager/internal/domain/entity"
	domainerrors "article-manager/internal/domain/errors"
	"article-manager/internal/infrastructure/logger"
	"article-manager/internal/infrastructure/timeutil"
	"article-manager/internal/usecase"

	"go.uber.org/zap"
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
	ctx := r.Context()

	logger.Info("Getting all articles",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
	)

	articles, err := h.usecase.GetAllArticles(ctx)
	if err != nil {
		HandleError(w, err, "GetAllArticles")
		return
	}

	response := make([]ArticleResponse, 0, len(articles))
	for _, article := range articles {
		response = append(response, toArticleResponse(article))
	}

	logger.Info("Successfully retrieved all articles",
		zap.Int("count", len(articles)),
	)

	RespondSuccess(w, http.StatusOK, response)
}

// 指定されたIDの記事を取得
func (h *ArticleHandler) GetArticleByID(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := r.Context()

	logger.Info("Getting article by ID",
		zap.Int64("id", id),
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
	)

	article, err := h.usecase.GetArticleByID(ctx, id)
	if err != nil {
		HandleError(w, err, "GetArticleByID")
		return
	}

	logger.Info("Successfully retrieved article",
		zap.Int64("id", article.ID),
		zap.String("title", article.Title),
	)

	RespondSuccess(w, http.StatusOK, toArticleResponse(article))
}

// 新しい記事を作成する
func (h *ArticleHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Warn("Failed to decode request body",
			zap.Error(err),
			zap.String("operation", "CreateArticle"),
		)
		HandleError(w, domainerrors.InvalidArgumentError("request body", "invalid JSON format"), "CreateArticle")
		return
	}

	logger.Info("Creating article",
		zap.String("title", req.Title),
		zap.String("url", req.URL),
		zap.Strings("tags", req.Tags),
	)

	article, err := h.usecase.CreateArticle(ctx, req.Title, req.URL, req.Summary, req.Tags, req.Memo)
	if err != nil {
		HandleError(w, err, "CreateArticle")
		return
	}

	logger.Info("Successfully created article",
		zap.Int64("id", article.ID),
		zap.String("title", article.Title),
	)

	RespondSuccess(w, http.StatusCreated, toArticleResponse(article))
}

// 記事を更新する
func (h *ArticleHandler) UpdateArticle(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := r.Context()

	var req UpdateArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Warn("Failed to decode request body",
			zap.Error(err),
			zap.String("operation", "UpdateArticle"),
			zap.Int64("id", id),
		)
		HandleError(w, domainerrors.InvalidArgumentError("request body", "invalid JSON format"), "UpdateArticle")
		return
	}

	logger.Info("Updating article",
		zap.Int64("id", id),
		zap.String("title", req.Title),
		zap.String("url", req.URL),
	)

	article, err := h.usecase.UpdateArticle(ctx, id, req.Title, req.URL, req.Summary, req.Tags, req.Memo)
	if err != nil {
		HandleError(w, err, "UpdateArticle")
		return
	}

	logger.Info("Successfully updated article",
		zap.Int64("id", article.ID),
		zap.String("title", article.Title),
	)

	RespondSuccess(w, http.StatusOK, toArticleResponse(article))
}

// 記事を削除する
func (h *ArticleHandler) DeleteArticle(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := r.Context()

	logger.Info("Deleting article",
		zap.Int64("id", id),
	)

	err := h.usecase.DeleteArticle(ctx, id)
	if err != nil {
		HandleError(w, err, "DeleteArticle")
		return
	}

	logger.Info("Successfully deleted article",
		zap.Int64("id", id),
	)

	w.WriteHeader(http.StatusNoContent)
}

// 記事を検索する
func (h *ArticleHandler) SearchArticles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	keyword := r.URL.Query().Get("keyword")

	logger.Info("Searching articles",
		zap.String("keyword", keyword),
	)

	trimmedKeyword := strings.TrimSpace(keyword)
	if trimmedKeyword == "" {
		HandleError(w, domainerrors.ValidationError("keyword", "keyword parameter is required and cannot be empty"), "SearchArticles")
		return
	}

	articles, err := h.usecase.SearchArticles(ctx, trimmedKeyword)
	if err != nil {
		HandleError(w, err, "SearchArticles")
		return
	}

	response := make([]ArticleResponse, 0, len(articles))
	for _, article := range articles {
		response = append(response, toArticleResponse(article))
	}

	logger.Info("Successfully searched articles",
		zap.String("keyword", trimmedKeyword),
		zap.Int("count", len(articles)),
	)

	RespondSuccess(w, http.StatusOK, response)
}

// エンティティをレスポンス形式に変換する
func toArticleResponse(article *entity.Article) ArticleResponse {
	return ArticleResponse{
		ID:        article.ID,
		Title:     article.Title,
		URL:       article.URL,
		Summary:   article.Summary,
		Tags:      article.Tags,
		Memo:      article.Memo,
		CreatedAt: timeutil.MustFormatInJST(article.CreatedAt),
		UpdatedAt: timeutil.MustFormatInJST(article.UpdatedAt),
	}
}
