// アプリケーション全体で使用する定数

// UI表示関連
export const UI_CONSTANTS = {
    TAG_DISPLAY_LIMIT: 8,           // タグ表示上限
    PAGINATION_VISIBLE_PAGES: 5,    // ページネーション表示数
    ARTICLES_PER_PAGE: 10,          // 1ページの記事数
} as const

// バリデーション関連
export const VALIDATION_CONSTANTS = {
    TAG_NAME_MAX_LENGTH: 50,        // タグ名最大長
    ARTICLE_TITLE_MIN_LENGTH: 1,    // 記事タイトル最小長
    ARTICLE_TITLE_MAX_LENGTH: 255,  // 記事タイトル最大長
    ARTICLE_URL_MIN_LENGTH: 1,      // 記事URL最小長
    ARTICLE_URL_MAX_LENGTH: 2048,   // 記事URL最大長
    ARTICLE_SUMMARY_MIN_LENGTH: 1,  // 記事要約最小長
    ARTICLE_SUMMARY_MAX_LENGTH: 500, // 記事要約最大長
} as const

// キャッシュ関連
export const CACHE_CONSTANTS = {
    CACHE_TTL: 60 * 1000,           // キャッシュTTL（1分）
} as const

// API関連
export const API_CONSTANTS = {
    BASE_URL: process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080',
} as const
