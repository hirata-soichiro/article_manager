import { useState, useCallback } from 'react'
import { articleClient } from  '@/lib/api/articleClient'
import type { Article } from '@/types/article'
import { ApiError } from '@/lib/errors/ApiError'

// カスタムフックの戻り値の型定義
interface UseArticleSearchReturn {
    results: Article[]
    loading: boolean
    error: ApiError | Error | null
    keyword: string
    search: (keyword: string) => Promise<void>
    clearSearch: () => void
}

// 記事検索用カスタムフック
export function useArticleSearch(): UseArticleSearchReturn {
    const [results, setResults] = useState<Article[]>([])
    const [loading, setLoading] = useState<boolean>(false)
    const [error, setError] = useState<ApiError | Error | null>(null)
    const [keyword, setKeyword] = useState<string>('')

    // 記事を検索
    const search = useCallback(async (searchKeyword: string) => {
        try {
            setLoading(true)
            setError(null)
            setKeyword(searchKeyword)
            const data = await articleClient.searchArticles(searchKeyword)
            setResults(data)
        } catch (err) {
            if (err instanceof ApiError) {
                setError(err)
            } else {
                setError(err as Error)
            }
            setResults([])
        } finally {
            setLoading(false)
        }
    }, [])

    // 検索結果をクリア
    const clearSearch = useCallback(() => {
        setResults([])
        setKeyword('')
        setError(null)
        setLoading(false)
    }, [])

    return {
        results,
        loading,
        error,
        keyword,
        search,
        clearSearch,
    }
}
