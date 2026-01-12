import { useState, useEffect, useCallback } from 'react'
import { articleClient } from '@/lib/api/articleClient'
import type { Article, CreateArticleInput, UpdateArticleInput } from '@/types/article'

// カスタムフックの戻り値の型定義
interface UseArticlesReturn {
    articles: Article[]
    loading: boolean
    error: Error | null
    createArticle: (input: CreateArticleInput) => Promise<void>
    updateArticle: (id: number, input: UpdateArticleInput) => Promise<void>
    deleteArticle: (id: number) => Promise<void>
    refetch: () => Promise<void>
}

// 記事管理用カスタムフック(記事一覧、作成、更新、削除)
export function useArticles(): UseArticlesReturn {
    // 状態管理
    const [articles, setArticles] = useState<Article[]>([])
    const [loading, setLoading] = useState<boolean>(true)
    const [error, setError] = useState<Error | null>(null)

    // 記事一覧を取得
    const fetchArticles = useCallback(async () => {
        try {
            setLoading(true)
            setError(null)
            const data = await articleClient.getAll()
            setArticles(data)
        } catch (err) {
            setError(err as Error)
        } finally {
            setLoading(false)
        }
    }, [])

    // コンポーネントマウント時に記事一覧を取得
    useEffect(() => {
        fetchArticles()
    }, [fetchArticles])

    // 記事作成
    const createArticle = useCallback(async (input: CreateArticleInput) => {
        try {
            setError(null)
            const newArticle = await articleClient.create(input)
            setArticles((prev) => [...prev, newArticle])
        } catch (err) {
            setError(err as Error)
            throw err
        }
    }, [])

    // 記事更新
    const updateArticle = useCallback(async (id: number, input: UpdateArticleInput) => {
        try {
            setError(null)
            const updatedArticle = await articleClient.update(id, input)
            setArticles((prev) =>
                prev.map((article) => (article.id === id ? updatedArticle : article))
            )
        } catch (err) {
            setError(err as Error)
            throw err
        }
    }, [])

    // 記事削除
     const deleteArticle = useCallback(async (id: number) => {
        try {
            setError(null)
            await articleClient.delete(id)
            setArticles((prev) => prev.filter((article) => article.id !== id))
        } catch (err) {
            setError(err as Error)
            throw err
        }
     }, [])

     // 記事一覧を再取得
     const refetch = useCallback(async () => {
        await fetchArticles()
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
