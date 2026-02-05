import { useState, useEffect, useCallback, useRef } from 'react'
import { articleClient } from '@/lib/api/articleClient'
import type { Article, CreateArticleInput, UpdateArticleInput } from '@/types/article'
import { ApiError } from '@/lib/errors/ApiError'
import { CACHE_CONSTANTS } from '@/config/constants'

// カスタムフックの戻り値の型定義
interface UseArticlesReturn {
    articles: Article[]
    loading: boolean
    error: ApiError | Error | null
    createArticle: (input: CreateArticleInput) => Promise<void>
    updateArticle: (id: number, input: UpdateArticleInput) => Promise<void>
    deleteArticle: (id: number) => Promise<void>
    refetch: () => Promise<void>
}

// 簡易的なキャッシュ機構
const articlesCache = {
    data: null as Article[] | null,
    timestamp: 0,
    ttl: CACHE_CONSTANTS.CACHE_TTL,
}

// テスト用: キャッシュをリセットする関数
export const __resetArticlesCache = () => {
    articlesCache.data = null
    articlesCache.timestamp = 0
}

// 記事管理用カスタムフック(記事一覧、作成、更新、削除)
export function useArticles(): UseArticlesReturn {
    // 状態管理
    const [articles, setArticles] = useState<Article[]>([])
    const [loading, setLoading] = useState<boolean>(true)
    const [error, setError] = useState<ApiError | Error | null>(null)

    // 重複リクエスト防止
    const fetchingRef = useRef(false)

    // 記事一覧を取得
    const fetchArticles = useCallback(async (useCache = true) => {
        // 重複リクエスト防止
        if (fetchingRef.current) {
            return
        }

        // キャッシュチェック
        if (useCache && articlesCache.data && (Date.now() - articlesCache.timestamp) < articlesCache.ttl) {
            setArticles(articlesCache.data)
            setLoading(false)
            return
        }

        try {
            fetchingRef.current = true
            setLoading(true)
            setError(null)
            const data = await articleClient.getAll()

            // キャッシュに保存
            articlesCache.data = data
            articlesCache.timestamp = Date.now()

            setArticles(data)
        } catch (err) {
            if (err instanceof ApiError) {
                setError(err)
            } else {
                setError(err as Error)
            }
        } finally {
            setLoading(false)
            fetchingRef.current = false
        }
    }, [])

    // コンポーネントマウント時に記事一覧を取得
    useEffect(() => {
        fetchArticles()
    }, [fetchArticles])

    // 記事作成（楽観的UI更新 + キャッシュ無効化）
    const createArticle = useCallback(async (input: CreateArticleInput) => {
        try {
            setError(null)
            const newArticle = await articleClient.create(input)

            // 楽観的UI更新
            setArticles((prev) => [...prev, newArticle])

            // キャッシュ無効化
            articlesCache.data = null
        } catch (err) {
            if (err instanceof ApiError) {
                setError(err)
            } else {
                setError(err as Error)
            }
            throw err
        }
    }, [])

    // 記事更新（楽観的UI更新 + キャッシュ無効化）
    const updateArticle = useCallback(async (id: number, input: UpdateArticleInput) => {
        try {
            setError(null)
            const updatedArticle = await articleClient.update(id, input)

            // 楽観的UI更新
            setArticles((prev) =>
                prev.map((article) => (article.id === id ? updatedArticle : article))
            )

            // キャッシュ無効化
            articlesCache.data = null
        } catch (err) {
            if (err instanceof ApiError) {
                setError(err)
            } else {
                setError(err as Error)
            }
            throw err
        }
    }, [])

    // 記事削除（楽観的UI更新 + キャッシュ無効化）
    const deleteArticle = useCallback(async (id: number) => {
        try {
            setError(null)

            // 楽観的UI更新（削除前に即座にUIを更新）
            const previousArticles = articles
            setArticles((prev) => prev.filter((article) => article.id !== id))

            try {
                await articleClient.delete(id)
                // キャッシュ無効化
                articlesCache.data = null
            } catch (err) {
                // 削除失敗時は元に戻す
                setArticles(previousArticles)
                throw err
            }
        } catch (err) {
            if (err instanceof ApiError) {
                setError(err)
            } else {
                setError(err as Error)
            }
            throw err
        }
    }, [articles])

    // 記事一覧を再取得
    const refetch = useCallback(async () => {
        await fetchArticles(false) // キャッシュを使わずに再取得
    }, [fetchArticles])

     return {
        articles,
        loading,
        error,
        createArticle,
        updateArticle,
        deleteArticle,
        refetch,
     }
}
