import { useState, useEffect, useCallback, useRef } from 'react'
import { bookRecommendationClient } from '@/lib/api/bookRecommendationClient'
import type { BookRecommendation } from '@/types/book'
import { ApiError } from '@/lib/errors/ApiError'

// カスタムフックの戻り値の型定義
interface UseBookRecommendationsReturn {
    books: BookRecommendation[]
    loading: boolean
    error: ApiError | Error | null
    cached: boolean
    generatedAt: string | null
    expiresAt: string | null
    refetch: () => Promise<void>
}

// 簡易的なキャッシュ機構
const bookRecommendationsCache = {
    data: null as BookRecommendation[] | null,
    timestamp: 0,
    ttl: 24 * 60 * 60 * 1000, // 24時間（ミリ秒）
    cached: false,
    generatedAt: null as string | null,
    expiresAt: null as string | null,
}

// テスト用: キャッシュをリセットする関数
export const __resetBookRecommendationsCache = () => {
    bookRecommendationsCache.data = null
    bookRecommendationsCache.timestamp = 0
    bookRecommendationsCache.cached = false
    bookRecommendationsCache.generatedAt = null
    bookRecommendationsCache.expiresAt = null
}

// 書籍推薦用カスタムフック
export function useBookRecommendations(): UseBookRecommendationsReturn {
    // 状態管理
    const [books, setBooks] = useState<BookRecommendation[]>([])
    const [loading, setLoading] = useState<boolean>(true)
    const [error, setError] = useState<ApiError | Error | null>(null)
    const [cached, setCached] = useState<boolean>(false)
    const [generatedAt, setGeneratedAt] = useState<string | null>(null)
    const [expiresAt, setExpiresAt] = useState<string | null>(null)

    // 重複リクエスト防止
    const fetchingRef = useRef(false)

    // 書籍推薦を取得
    const fetchBookRecommendations = useCallback(async (useCache = true) => {
        // 重複リクエスト防止
        if (fetchingRef.current) {
            return
        }

        // キャッシュチェック
        if (
            useCache &&
            bookRecommendationsCache.data &&
            Date.now() - bookRecommendationsCache.timestamp < bookRecommendationsCache.ttl
        ) {
            setBooks(bookRecommendationsCache.data)
            setCached(bookRecommendationsCache.cached)
            setGeneratedAt(bookRecommendationsCache.generatedAt)
            setExpiresAt(bookRecommendationsCache.expiresAt)
            setLoading(false)
            return
        }

        try {
            fetchingRef.current = true
            setLoading(true)
            setError(null)
            const data = await bookRecommendationClient.getRecommendations()

            // キャッシュに保存
            bookRecommendationsCache.data = data.books
            bookRecommendationsCache.timestamp = Date.now()
            bookRecommendationsCache.cached = data.cached
            bookRecommendationsCache.generatedAt = data.generatedAt
            bookRecommendationsCache.expiresAt = data.expiresAt

            setBooks(data.books)
            setCached(data.cached)
            setGeneratedAt(data.generatedAt)
            setExpiresAt(data.expiresAt)
        } catch (err) {
            // エラー時は静かに失敗（エラーをログに記録するがUIには表示しない）
            console.error('Failed to fetch book recommendations:', err)
            if (err instanceof ApiError) {
                setError(err)
            } else {
                setError(err as Error)
            }
            // エラー時は空配列を設定
            setBooks([])
            setCached(false)
            setGeneratedAt(null)
            setExpiresAt(null)
        } finally {
            setLoading(false)
            fetchingRef.current = false
        }
    }, [])

    // コンポーネントマウント時に書籍推薦を取得
    useEffect(() => {
        fetchBookRecommendations()
    }, [fetchBookRecommendations])

    // 書籍推薦を再取得
    const refetch = useCallback(async () => {
        await fetchBookRecommendations(false) // キャッシュを使わずに再取得
    }, [fetchBookRecommendations])

    return {
        books,
        loading,
        error,
        cached,
        generatedAt,
        expiresAt,
        refetch,
    }
}