import { BookRecommendation, BookRecommendationsResponse } from '@/types/book'
import { BaseApiClient } from './baseClient'

// APIから返却される書籍推薦データの型（スネークケース）
interface ApiPurchaseLinks {
    amazon?: string
    rakuten?: string
}

interface ApiBookRecommendation {
    title: string
    author: string
    isbn?: string
    purchaseLinks: ApiPurchaseLinks
}

interface ApiBookRecommendationsResponse {
    books: ApiBookRecommendation[]
    cached: boolean
    generatedAt: string | null
    expiresAt: string | null
}

// バックエンドのAPIと通信するクライアント
class BookRecommendationClient extends BaseApiClient {
    // 書籍推薦を取得
    async getRecommendations(): Promise<BookRecommendationsResponse> {
        const data = await this.fetchWithErrorHandling<ApiBookRecommendationsResponse>(
            '/api/book-recommendations'
        )
        return this.convertToCamelCase(data)
    }

    // APIレスポンスをフロントエンド用に変換
    private convertToCamelCase(
        apiResponse: ApiBookRecommendationsResponse
    ): BookRecommendationsResponse {
        return {
            books: apiResponse.books.map((book) => ({
                title: book.title,
                author: book.author,
                isbn: book.isbn,
                purchaseLinks: {
                    amazon: book.purchaseLinks.amazon,
                    rakuten: book.purchaseLinks.rakuten,
                },
            })),
            cached: apiResponse.cached,
            generatedAt: apiResponse.generatedAt,
            expiresAt: apiResponse.expiresAt,
        }
    }
}

// シングルトンインスタンスをエクスポート
export const bookRecommendationClient = new BookRecommendationClient()
