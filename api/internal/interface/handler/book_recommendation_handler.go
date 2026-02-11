package handler

import (
	"net/http"

	"article-manager/internal/domain/entity"
	"article-manager/internal/infrastructure/logger"
	"article-manager/internal/infrastructure/timeutil"
	"article-manager/internal/usecase"

	"go.uber.org/zap"
)

// 書籍推薦に関するHTTPハンドラ
type BookRecommendationHandler struct {
	usecase *usecase.BookRecommendationUsecase
}

// BookRecommendationHandlerのコンストラクタ
func NewBookRecommendationHandler(uc *usecase.BookRecommendationUsecase) *BookRecommendationHandler {
	return &BookRecommendationHandler{
		usecase: uc,
	}
}

// 購入リンクのレスポンス構造体
type PurchaseLinksResponse struct {
	Amazon  string `json:"amazon,omitempty"`
	Rakuten string `json:"rakuten,omitempty"`
}

// 書籍のレスポンス構造体
type BookResponse struct {
	Title         string                `json:"title"`
	Author        string                `json:"author"`
	ISBN          string                `json:"isbn,omitempty"`
	PurchaseLinks PurchaseLinksResponse `json:"purchaseLinks"`
}

// 書籍推薦のレスポンス構造体
type BookRecommendationResponse struct {
	Books       []BookResponse `json:"books"`
	Cached      bool           `json:"cached"`
	GeneratedAt *string        `json:"generatedAt"`
	ExpiresAt   *string        `json:"expiresAt"`
}

// 書籍推薦を取得する
func (h *BookRecommendationHandler) GetBookRecommendations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger.Info("Getting book recommendations",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
	)

	cache, err := h.usecase.GetBookRecommendations(ctx)
	if err != nil {
		HandleError(w, err, "GetBookRecommendations")
		return
	}

	response := toBookRecommendationResponse(cache)

	logger.Info("Successfully retrieved book recommendations",
		zap.Int("book_count", len(cache.Books)),
		zap.Bool("cached", cache.ID > 0),
	)

	RespondSuccess(w, http.StatusOK, response)
}

// エンティティをレスポンス形式に変換する
func toBookRecommendationResponse(cache *entity.BookRecommendationCache) BookRecommendationResponse {
	// 書籍が0件の場合
	if len(cache.Books) == 0 {
		return BookRecommendationResponse{
			Books:       []BookResponse{},
			Cached:      false,
			GeneratedAt: nil,
			ExpiresAt:   nil,
		}
	}

	// 書籍リストの変換
	books := make([]BookResponse, 0, len(cache.Books))
	for _, book := range cache.Books {
		books = append(books, toBookResponse(book))
	}

	// 日時をISO 8601形式に変換
	generatedAt := timeutil.MustFormatInJST(cache.GeneratedAt)
	expiresAt := timeutil.MustFormatInJST(cache.ExpiresAt)

	// キャッシュされているかどうかを判定（IDが0より大きければキャッシュ済み）
	cached := cache.ID > 0

	return BookRecommendationResponse{
		Books:       books,
		Cached:      cached,
		GeneratedAt: &generatedAt,
		ExpiresAt:   &expiresAt,
	}
}

// 書籍エンティティをレスポンス形式に変換する
func toBookResponse(book entity.Book) BookResponse {
	return BookResponse{
		Title:  book.Title,
		Author: book.Author,
		ISBN:   book.ISBN,
		PurchaseLinks: PurchaseLinksResponse{
			Amazon:  book.PurchaseLinks.Amazon,
			Rakuten: book.PurchaseLinks.Rakuten,
		},
	}
}
