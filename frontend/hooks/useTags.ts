import { useState, useEffect, useCallback } from 'react'
import { tagClient } from '@/lib/api/tagClient'
import type { Tag, CreateTagInput, UpdateTagInput } from '@/types/tag'
import { ApiError } from '@/lib/errors/ApiError'

// カスタムフックの戻り値の型定義
  interface UseTagsReturn {
    tags: Tag[]
    loading: boolean
    error: ApiError | Error | null
    createTag: (input: CreateTagInput) => Promise<void>
    updateTag: (id: number, input: UpdateTagInput) => Promise<void>
    deleteTag: (id: number) => Promise<void>
    refetch: () => Promise<void>
  }

// タグ管理用カスタムフック(タグ一覧、作成、更新、削除)
export function useTags(): UseTagsReturn {
    // 状態管理
    const [tags, setTags] = useState<Tag[]>([])
    const [loading, setLoading] = useState<boolean>(true)
    const [error, setError] = useState<ApiError | Error | null>(null)

    // タグ一覧を取得
    const fetchTags = useCallback(async () => {
        try {
            setLoading(true)
            setError(null)
            const data = await tagClient.getAll()
            setTags(data)
        } catch (err) {
            if (err instanceof ApiError) {
                setError(err)
            } else {
                setError(err as Error)
            }
        } finally {
            setLoading(false)
        }
    }, [])

    // コンポーネントマウント時にタグ一覧を取得
    useEffect(() => {
        fetchTags()
    }, [fetchTags])

    // タグ作成
    const createTag = useCallback(async (input: CreateTagInput) => {
        try {
            setError(null)
            const newTag = await tagClient.create(input)
            setTags((prev) => [...prev, newTag])
        } catch (err) {
            if (err instanceof ApiError) {
                setError(err)
            } else {
                setError(err as Error)
            }
            throw err
        }
    }, [])

    // タグ更新
    const updateTag = useCallback(async (id: number, input: UpdateTagInput) => {
        try {
            setError(null)
            const updatedTag = await tagClient.update(id, input)
            setTags((prev) =>
                prev.map((tag) => (tag.id === id ? updatedTag : tag))
            )
        } catch (err) {
            if (err instanceof ApiError) {
                setError(err)
            } else {
                setError(err as Error)
            }
            throw err
        }
    }, [])

    // タグ削除
     const deleteTag = useCallback(async (id: number) => {
        try {
            setError(null)
            await tagClient.delete(id)
            setTags((prev) => prev.filter((tag) => tag.id !== id))
        } catch (err) {
            if (err instanceof ApiError) {
                setError(err)
            } else {
                setError(err as Error)
            }
            throw err
        }
     }, [])

     // タグ一覧を再取得
     const refetch = useCallback(async () => {
        await fetchTags()
     }, [fetchTags])

     return {
        tags,
        loading,
        error,
        createTag,
        updateTag,
        deleteTag,
        refetch,
     }
}
