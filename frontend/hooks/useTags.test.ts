import { renderHook, waitFor } from '@testing-library/react'
import { describe, it, expect, beforeEach, vi } from 'vitest'
import { useTags } from './useTags'
import { tagClient } from '@/lib/api/tagClient'
import type { Tag, CreateTagInput, UpdateTagInput } from '@/types/tag'

// tagClientをモック化
vi.mock('@/lib/api/tagClient')

// モック関数の型定義
const mockGetAll = tagClient.getAll as ReturnType<typeof vi.fn>
const mockCreate = tagClient.create as ReturnType<typeof vi.fn>
const mockUpdate = tagClient.update as ReturnType<typeof vi.fn>
const mockDelete = tagClient.delete as ReturnType<typeof vi.fn>

describe('useTags', () => {
    // テスト用データ
    const mockTags: Tag[] = [
        {
            id: 1,
            name: 'Go',
            createdAt: '2024-01-01 10:00:00',
            updatedAt: '2024-01-01 10:00:00',
        },
        {
            id: 2,
            name: 'React',
            createdAt: '2024-01-02 10:00:00',
            updatedAt: '2024-01-02 10:00:00',
        },
        {
            id: 3,
            name: 'TypeScript',
            createdAt: '2024-01-03 10:00:00',
            updatedAt: '2024-01-03 10:00:00',
        },
    ]

    // 各テスト前にモックをクリア
    beforeEach(() => {
        vi.clearAllMocks()
    })

    describe('タグ一覧の取得', () => {
        it('初期状態ではloading=true, tags=[]であること', () => {
            mockGetAll.mockResolvedValue(mockTags)

            const { result } = renderHook(() => useTags())

            // 初期状態を確認
            expect(result.current.loading).toBe(true)
            expect(result.current.articles).toEqual([])
            expect(result.current.error).toBeNull()
        })

        it('タグ一覧取得成功時、loading=false, tagsにタグが格納されること', async () => {
            mockGetAll.mockResolvedValue(mockTags)
            const { result } = renderHook(() => useTags())

            // 非同期処理完了を待つ
            await waitFor(() => {
                expect(result.current.loading).toBe(false)
            })

            expect(result.current.tags).toEqual(mockTags)
            expect(result.current.error).toBeNull()
            expect(mockGetAll).toHaveBeenCalledTimes(1)
        })

        it('タグ一覧取得失敗時、loading=false、errorにエラー情報が格納されること', async () => {
            const mockError = new Error('Failed to fetch tags')
            mockGetAll.mockRejectedValue(mockError)

            const { result } = renderHook(() => useTags())

            await waitFor(() => {
                expect(result.current.loading).toBe(false)
            })

            expect(result.current.tags).toEqual([])
            expect(result.current.error).toBe(mockError)
            expect(mockGetAll).toHaveBeenCalledTimes(1)
        })

        it('タグが名前順にソートされて取得されること', async () => {
            const sortedTags: Tag[] = [
                { id: 2, name: 'Go', createdAt: '2024-01-01 10:00:00', updatedAt: '2024-01-01 10:00:00' },
                { id: 3, name: 'Python', createdAt: '2024-01-02 10:00:00', updatedAt: '2024-01-02 10:00:00' },
                { id: 1, name: 'React', createdAt: '2024-01-03 10:00:00', updatedAt: '2024-01-03 10:00:00' },
            ]
            mockGetAll.mockResolvedValue(sortedTags)

            const { result } = renderHook(() => useTags())

            await waitFor(() => {
                expect(result.current.loading).toBe(false)
            })

            expect(result.current.tags).toEqual(sortedTags)
            expect(result.current.tags[0].name).toBe('Go')
            expect(result.current.tags[1].name).toBe('Python')
            expect(result.current.tags[2].name).toBe('React')
        })
    })

    describe('タグの作成', () => {
        it('タグ作成成功時、新しいタグが一覧に追加されること', async () => {
            const newTag: Tag = {
                id: 4,
                name: 'Rust',
                createdAt: '2024-01-04 10:00:00',
                updatedAt: '2024-01-04 10:00:00',
            }

            const createInput: CreateTagInput = {
                name: 'Rust',
            }

            mockGetAll.mockResolvedValue(mockTags)
            mockCreate.mockResolvedValue(newTag)

            const { result } = renderHook(() => useTags())

            await waitFor(() => {
                expect(result.current.loading).toBe(false)
            })

            // タグ作成を実行
            await result.current.createTag(createInput)

            // 新しいタグが追加されたか確認（状態更新完了を待つ）
            await waitFor(() => {
                expect(result.current.tags).toHaveLength(3)
            })
            expect(result.current.tags[3]).toEqual(newTag)
            expect(mockCreate).toHaveBeenCalledWith(createInput)
        })

        it('タグ作成失敗時、errorに情報が格納され一覧は変更されないこと', async () => {
            const createInput: CreateTagInput = {
                name: 'Rust',
            }

            const mockError = new Error('Failed to create tag')

            mockGetAll.mockResolvedValue(mockTags)
            mockCreate.mockRejectedValue(mockError)

            const { result } = renderHook(() => useTags())

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

        it('空のタグ名での作成が失敗すること', async () => {
            const createInput: CreateTagInput = {
                name: '',
            }

            const mockError = new Error('name is required')

            mockGetAll.mockResolvedValue(mockTags)
            mockCreate.mockRejectedValue(mockError)

            const { result } = renderHook(() => useTags())

            await waitFor(() => {
                expect(result.current.loading).toBe(false)
            })

            try {
                await result.current.createTag(createInput)
            } catch (err) {
                // エラーが投げられることを期待
            }

            await waitFor(() => {
                expect(result.current.error).toBe(mockError)
            })
        })

        it('重複したタグ名での作成が失敗すること', async () => {
            const createInput: CreateTagInput = {
                name: 'Go',
            }

            const mockError = new Error('tag already exists')

            mockGetAll.mockResolvedValue(mockTags)
            mockCreate.mockRejectedValue(mockError)

            const { result } = renderHook(() => useTags())

            await waitFor(() => {
                expect(result.current.loading).toBe(false)
            })

            try {
                await result.current.createTag(createInput)
            } catch (err) {
                // エラーが投げられることを期待
            }

            await waitFor(() => {
                expect(result.current.error).toBe(mockError)
            })
        })
    })

    describe('タグの更新', () => {
        it('タグ更新成功時、一覧内の該当タグが更新されること', async () => {
            const updateInput: UpdateTagInput = {
                name: 'Golang',
            }

            const updatedTag: Tag = {
                ...mockTags[0],
                name: 'Golang',
                updatedAt: '2024-01-05 12:00:00',
            }

            mockGetAll.mockResolvedValue(mockTags)
            mockUpdate.mockResolvedValue(updatedTag)

            const { result } = renderHook(() => useTags())

            await waitFor(() => {
                expect(result.current.loading).toBe(false)
            })

            // タグ更新を実行
            await result.current.updateTag(1, updateInput)

            // 一覧内の該当タグが更新されたか確認（状態更新完了を待つ）
            await waitFor(() => {
                const updatedInList = result.current.tags.find((t: Tag) => t.id === 1)
                expect(updatedInList).toEqual(updatedTag)
            })
            expect(mockUpdate).toHaveBeenCalledWith(1, updateInput)
        })

        it('タグ更新失敗時、errorに情報が格納され一覧は変更されないこと', async () => {
            const updateInput: UpdateTagInput = {
                name: 'Golang',
            }

            const mockError = new Error('Failed to update tag')

            mockGetAll.mockResolvedValue(mockTags)
            mockUpdate.mockRejectedValue(mockError)

            const { result } = renderHook(() => useTags())

            await waitFor(() => {
                expect(result.current.loading).toBe(false)
            })

            const tagsBeforeUpdate = [...result.current.tags]

            // 更新失敗を実行（エラーがthrowされる）
            try {
                await result.current.updateTag(1, updateInput)
            } catch (err) {
                // エラーが投げられることを期待
            }

            // エラーが格納され、一覧は変更されていないか（状態更新完了を待つ）
            await waitFor(() => {
                expect(result.current.error).toBe(mockError)
            })
            expect(result.current.tags).toEqual(tagsBeforeUpdate)
        })

        it('存在しないIDのタグ更新が失敗すること', async () => {
            const updateInput: UpdateTagInput = {
                name: 'Updated',
            }

            const mockError = new Error('tag not found')

            mockGetAll.mockResolvedValue(mockTags)
            mockUpdate.mockRejectedValue(mockError)

            const { result } = renderHook(() => useTags())

            await waitFor(() => {
                expect(result.current.loading).toBe(false)
            })

            try {
                await result.current.updateTag(999, updateInput)
            } catch (err) {
                // エラーが投げられることを期待
            }

            await waitFor(() => {
                expect(result.current.error).toBe(mockError)
            })
        })
    })

    describe('タグの削除', () => {
        it('タグ削除成功時、一覧から該当タグが削除されること', async () => {
            mockGetAll.mockResolvedValue(mockTags)
            mockDelete.mockResolvedValue(undefined)

            const { result } = renderHook(() => useTags())

            await waitFor(() => {
                expect(result.current.loading).toBe(false)
            })

            expect(result.current.tags).toHaveLength(2)

            // タグ削除を実行
            await result.current.deleteTag(1)

            // 該当タグが削除されたか確認（状態更新完了を待つ）
            await waitFor(() => {
                expect(result.current.tags).toHaveLength(1)
            })
            expect(result.current.tags.find((t: Tag) => t.id === 1)).toBeUndefined()
            expect(mockDelete).toHaveBeenCalledWith(1)
        })

        it('タグ削除失敗時、errorに情報が格納され一覧は変更されないこと', async () => {
            const mockError = new Error('Failed to delete tag')

            mockGetAll.mockResolvedValue(mockTags)
            mockDelete.mockRejectedValue(mockError)

            const { result } = renderHook(() => useTags())

            await waitFor(() => {
                expect(result.current.loading).toBe(false)
            })

            const tagsBeforeDelete = [...result.current.tags]

            // 削除失敗を実行（エラーがthrowされる）
            try {
                await result.current.deleteTag(1)
            } catch (err) {
                // エラーが投げられることを期待
            }

            // エラーが格納され、一覧は変更されていないか（状態更新完了を待つ）
            await waitFor(() => {
                expect(result.current.error).toBe(mockError)
            })
            expect(result.current.tags).toEqual(tagsBeforeDelete)
        })

        it('存在しないIDのタグ削除が失敗すること', async () => {
            const mockError = new Error('tag not found')

            mockGetAll.mockResolvedValue(mockTags)
            mockDelete.mockRejectedValue(mockError)

            const { result } = renderHook(() => useTags())

            await waitFor(() => {
                expect(result.current.loading).toBe(false)
            })

            try {
                await result.current.deleteTag(999)
            } catch (err) {
                // エラーが投げられることを期待
            }

            await waitFor(() => {
                expect(result.current.error).toBe(mockError)
            })
        })
    })

    describe('再取得(refetch)機能', () => {
        it('refetch呼び出し時、記事一覧が再取得されること', async () => {
            const updatedTags: Tag[] = [
                ...mockTags,
                {
                    id: 4,
                    name: 'Rust',
                    createdAt: '2024-01-05 10:00:00',
                    updatedAt: '2024-01-05 10:00:00',
                },
            ]

            // 初回3件、再取得時4件
            mockGetAll.mockResolvedValueOnce(mockTags).mockResolvedValueOnce(updatedTags)

            const { result } = renderHook(() => useTags())

            await waitFor(() => {
                expect(result.current.loading).toBe(false)
            })

            expect(result.current.articles).toHaveLength(3)

            // 再取得を実行
            await result.current.refetch()

            await waitFor(() => {
                expect(result.current.loading).toBe(false)
            })

            // 再取得後のデータが反映されているか
            expect(result.current.tags).toHaveLength(4)
            expect(result.current.tags).toEqual(updatedTags)
            expect(mockGetAll).toHaveBeenCalledTimes(2)
        })
    })
})
