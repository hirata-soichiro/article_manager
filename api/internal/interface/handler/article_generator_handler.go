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
	ctx := r.Context()

	var req GenerateArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Warn("Failed to decode request body",
			zap.Error(err),
			zap.String("operation", "GenerateArticle"),
		)
		HandleError(w, domainerrors.InvalidArgumentError("request body", "invalid JSON format"), "GenerateArticle")
		return
	}

	logger.Info("Generating article from URL",
		zap.String("url", req.URL),
		zap.String("memo", req.Memo),
	)

	article, err := h.generatorUsecase.GenerateArticleFromURL(ctx, req.URL, req.Memo)
	if err != nil {
		HandleError(w, err, "GenerateArticle")
		return
	}

	logger.Info("Successfully generated article",
		zap.Int64("id", article.ID),
		zap.String("title", article.Title),
		zap.String("url", article.URL),
	)

	RespondSuccess(w, http.StatusOK, h.toArticleResponse(article))
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
