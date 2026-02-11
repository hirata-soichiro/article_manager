// 購入リンクオブジェクト
export interface PurchaseLinks {
    amazon?: string
    rakuten?: string
}

// 書籍推薦情報
export interface BookRecommendation {
    title: string
    purchaseLinks: PurchaseLinks
}

// 書籍推薦APIレスポンス
export interface BookRecommendationsResponse {
    books: BookRecommendation[]
    cached: boolean
    generatedAt: string | null
    expiresAt: string | null
}