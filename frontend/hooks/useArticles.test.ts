import { renderHook, waitFor } from '@testing-library/react'
import { describe, it, expect, beforeEach, vi } from 'vitest'
import { useArticles, __resetArticlesCache } from './useArticles'
import { articleClient } from '@/lib/api/articleClient'
import type { Article, CreateArticleInput, UpdateArticleInput } from '@/types/article'

// articleClientをモック化
vi.mock('@/lib/api/articleClient')

// モック関数の型定義
const mockGetAll = articleClient.getAll as ReturnType<typeof vi.fn>
const mockCreate = articleClient.create as ReturnType<typeof vi.fn>
const mockUpdate = articleClient.update as ReturnType<typeof vi.fn>
const mockDelete = articleClient.delete as ReturnType<typeof vi.fn>

describe('useArticles', () => {
    // テスト用データ
    const mockArticles: Article[] = [
        {
            id: 1,
            title: 'Test Article 1',
            url: 'https://example.com/1',
            summary: 'Test summary 1',
            tags: ['test', 'article'],
            memo: 'Test memo 1',
            createdAt: '2024-01-01 10:00:00',
            updatedAt: '2024-01-01 10:00:00',
        },
        {
            id: 2,
            title: 'Test Article 2',
            url: 'https://example.com/2',
            summary: 'Test summary 2',
            tags: ['test'],
            memo: 'Test memo 2',
            createdAt: '2024-01-02 10:00:00',
            updatedAt: '2024-01-02 10:00:00',
        },
    ]

    // 各テスト前にモックとキャッシュをクリア
    beforeEach(() => {
        vi.clearAllMocks()
        __resetArticlesCache()
    })

    describe('記事一覧の取得', () => {
        it('初期状態ではloading=true, articles=[]であること', () => {
            mockGetAll.mockResolvedValue(mockArticles)

            const { result } = renderHook(() => useArticles())

            // 初期状態を確認
            expect(result.current.loading).toBe(true)
            expect(result.current.articles).toEqual([])
            expect(result.current.error).toBeNull()
        })

        it('記事一覧取得成功時、loading=false, articlesに記事が格納されること', async () => {
            mockGetAll.mockResolvedValue(mockArticles)
            const { result } = renderHook(() => useArticles())

            // 非同期処理完了を待つ
            await waitFor(() => {
                expect(result.current.loading).toBe(false)
            })

            expect(result.current.articles).toEqual(mockArticles)
            expect(result.current.error).toBeNull()
            expect(mockGetAll).toHaveBeenCalledTimes(1)
        })

        it('記事一覧取得失敗時、loading=false、errorにエラー情報が格納されること', async () => {
            const mockError = new Error('Failed to fetch articles')
            mockGetAll.mockRejectedValue(mockError)

            const { result } = renderHook(() => useArticles())

            await waitFor(() => {
                expect(result.current.loading).toBe(false)
            })

            expect(result.current.articles).toEqual([])
            expect(result.current.error).toBe(mockError)
            expect(mockGetAll).toHaveBeenCalledTimes(1)
        })
    })

    describe('記事の作成', () => {
        it('記事作成成功時、新しい記事が一覧に追加されること', async () => {
            const newArticle: Article = {
                id: 3,
                title: 'New Article',
                url: 'https://example.com/new',
                summary: 'New summary',
                tags: ['new'],
                memo: 'New memo',
                createdAt: '2024-01-03 10:00:00',
                updatedAt: '2024-01-03 10:00:00',
            }

            const createInput: CreateArticleInput = {
                title: 'New Article',
                url: 'https://example.com/new',
                summary: 'New summary',
                tags: ['new'],
                memo: 'New memo',
            }

            mockGetAll.mockResolvedValue(mockArticles)
            mockCreate.mockResolvedValue(newArticle)

            const { result } = renderHook(() => useArticles())

            await waitFor(() => {
                expect(result.current.loading).toBe(false)
            })

            // 記事作成を実行
            await result.current.createArticle(createInput)

            // 新しい記事が追加されたか確認（状態更新完了を待つ）
            await waitFor(() => {
                expect(result.current.articles).toHaveLength(3)
            })
            expect(result.current.articles[2]).toEqual(newArticle)
            expect(mockCreate).toHaveBeenCalledWith(createInput)
        })

        it('記事作成失敗時、errorに情報が格納され一覧は変更されないこと', async () => {
            const createInput: CreateArticleInput = {
                title: 'New Article',
                url: 'https://example.com/new',
                summary: 'New summary',
                tags: ['new'],
            }

            const mockError = new Error('Failed to create article')

            mockGetAll.mockResolvedValue(mockArticles)
            mockCreate.mockRejectedValue(mockError)

            const { result } = renderHook(() => useArticles())

            await waitFor(() => {
                expect(result.current.loading).toBe(false)
            })

            const articlesBeforeCreate = [...result.current.articles]

            // 作成失敗を実行（エラーがthrowされる）
            try {
                await result.current.createArticle(createInput)
            } catch (err) {
                // エラーが投げられることを期待
            }

            // エラーが格納され、一覧は変更されていないか（状態更新完了を待つ）
            await waitFor(() => {
                expect(result.current.error).toBe(mockError)
            })
            expect(result.current.articles).toEqual(articlesBeforeCreate)
        })
    })

    describe('記事の更新', () => {
        it('記事更新成功時、一覧内の該当記事が更新されること', async () => {
            const updateInput: UpdateArticleInput = {
                title: 'Updated Title',
                url: 'https://example.com/updated',
                summary: 'Updated summary',
                tags: ['updated'],
                memo: 'Updated memo',
            }

            const updatedArticle: Article = {
                ...mockArticles[0],
                ...updateInput,
                updatedAt: '2024-01-04 10:00:00',
            }

            mockGetAll.mockResolvedValue(mockArticles)
            mockUpdate.mockResolvedValue(updatedArticle)

            const { result } = renderHook(() => useArticles())

            await waitFor(() => {
                expect(result.current.loading).toBe(false)
            })

            // 記事更新を実行
            await result.current.updateArticle(1, updateInput)

            // 該当記事が更新されたか確認（状態更新完了を待つ）
            await waitFor(() => {
                const updatedInList = result.current.articles.find((a: Article) => a.id === 1)
                expect(updatedInList).toEqual(updatedArticle)
            })
            expect(mockUpdate).toHaveBeenCalledWith(1, updateInput)
        })

        it('記事更新失敗時、errorに情報が格納され一覧は変更されないこと', async () => {
            const updateInput: UpdateArticleInput = {
                title: 'Updated Title',
                url: 'https://example.com/updated',
                summary: 'Updated summary',
                tags: ['updated'],
            }

            const mockError = new Error('Failed to update article')

            mockGetAll.mockResolvedValue(mockArticles)
            mockUpdate.mockRejectedValue(mockError)

            const { result } = renderHook(() => useArticles())

            await waitFor(() => {
                expect(result.current.loading).toBe(false)
            })

            const articlesBeforeUpdate = [...result.current.articles]

            // 更新失敗を実行（エラーがthrowされる）
            try {
                await result.current.updateArticle(1, updateInput)
            } catch (err) {
                // エラーが投げられることを期待
            }

            // エラーが格納され、一覧は変更されていないか（状態更新完了を待つ）
            await waitFor(() => {
                expect(result.current.error).toBe(mockError)
            })
            expect(result.current.articles).toEqual(articlesBeforeUpdate)
        })
    })

    describe('記事の削除', () => {
        it('記事削除成功時、一覧から該当記事が削除されること', async () => {
            mockGetAll.mockResolvedValue(mockArticles)
            mockDelete.mockResolvedValue(undefined)

            const { result } = renderHook(() => useArticles())

            await waitFor(() => {
                expect(result.current.loading).toBe(false)
            })

            expect(result.current.articles).toHaveLength(2)

            // 記事削除を実行
            await result.current.deleteArticle(1)

            // 該当記事が削除されたか確認（状態更新完了を待つ）
            await waitFor(() => {
                expect(result.current.articles).toHaveLength(1)
            })
            expect(result.current.articles.find((a: Article) => a.id === 1)).toBeUndefined()
            expect(mockDelete).toHaveBeenCalledWith(1)
        })

        it('記事削除失敗時、errorに情報が格納され一覧は変更されないこと', async () => {
            const mockError = new Error('Failed to delete article')

            mockGetAll.mockResolvedValue(mockArticles)
            mockDelete.mockRejectedValue(mockError)

            const { result } = renderHook(() => useArticles())

            await waitFor(() => {
                expect(result.current.loading).toBe(false)
            })

            const articlesBeforeDelete = [...result.current.articles]

            // 削除失敗を実行（エラーがthrowされる）
            try {
                await result.current.deleteArticle(1)
            } catch (err) {
                // エラーが投げられることを期待
            }

            // エラーが格納され、一覧は変更されていないか（状態更新完了を待つ）
            await waitFor(() => {
                expect(result.current.error).toBe(mockError)
            })
            expect(result.current.articles).toEqual(articlesBeforeDelete)
        })
    })

    describe('再取得(refetch)機能', () => {
        it('refetch呼び出し時、記事一覧が再取得されること', async () => {
            const updatedArticles: Article[] = [
                ...mockArticles,
                {
                    id: 3,
                    title: 'Refetched Article',
                    url: 'https://example.com/refetch',
                    summary: 'Refetched summary',
                    tags: ['refetch'],
                    memo: 'Refetched memo',
                    createdAt: '2024-01-05 10:00:00',
                    updatedAt: '2024-01-05 10:00:00',
                },
            ]

            // 初回2件、再取得時3件
            mockGetAll.mockResolvedValueOnce(mockArticles).mockResolvedValueOnce(updatedArticles)

            const { result } = renderHook(() => useArticles())

            await waitFor(() => {
                expect(result.current.loading).toBe(false)
            })

            expect(result.current.articles).toHaveLength(2)

            // 再取得を実行
            await result.current.refetch()

            await waitFor(() => {
                expect(result.current.loading).toBe(false)
            })

            // 再取得後のデータが反映されているか
            expect(result.current.articles).toHaveLength(3)
            expect(result.current.articles).toEqual(updatedArticles)
            expect(mockGetAll).toHaveBeenCalledTimes(2)
        })
    })
})
