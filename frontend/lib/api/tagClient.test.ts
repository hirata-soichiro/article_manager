import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { tagClient } from '@/lib/api/tagClient'
import { Tag } from '@/types/tag'

describe('TagClient', () => {
    // 各テストの前にfetchをモック化し、実際のHTTPリクエストを防ぐ
    beforeEach(() => {
        global.fetch = vi.fn()
    })

    // 各テストの後にすべてのモックを復元し、テスト間の影響を排除
    afterEach(() => {
        vi.restoreAllMocks()
    })

    // テスト用のダミーデータ
    const mockApiTag = {
        id: 1,
        name: 'Go',
        created_at: '2024-01-01T10:00:00Z',
        updated_at: '2024-01-01T10:00:00Z',
    }

    // フロントエンドで期待される形式
    const expectedTag: Tag = {
        id: 1,
        name: 'Go',
        createdAt: '2024-01-01T10:00:00Z',
        updatedAt: '2024-01-01T10:00:00Z',
    }

    // タグ一覧取得のテスト
    describe('getAll', () => {
        it('全タグを取得できる', async () => {
            // 3件のタグデータを準備
            const mockTags = [
                mockApiTag,
                { ...mockApiTag, id: 2, name: 'React' },
                { ...mockApiTag, id: 3, name: 'TypeScript' },
            ]

            // API成功レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: true,
                status: 200,
                json: async () => mockTags,
            })

            // 実行
            const result = await tagClient.getAll()

            // 検証: 呼び出し回数、URL、データ件数、データ内容
            expect(global.fetch).toHaveBeenCalledTimes(1)
            expect(global.fetch).toHaveBeenCalledWith('http://localhost:8080/api/tags')
            expect(result).toHaveLength(3)
            expect(result[0]).toEqual(expectedTag)
            expect(result[1].id).toBe(2)
            expect(result[1].name).toBe('React')
            expect(result[2].id).toBe(3)
            expect(result[2].name).toBe('TypeScript')
        })

        it('空の配列が返ってきた場合も正常に処理できる', async () => {
            // 空配列のAPIレスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: true,
                status: 200,
                json: async () => [],
            })

            // 実行
            const result = await tagClient.getAll()

            // 検証: データ内容
            expect(result).toEqual([])
        })

        it('APIエラー時にエラーをスローする', async () => {
            // APIエラーレスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: false,
                status: 500,
                json: async () => ({ error: 'Internal Server Error'}),
            })

            // 実行と検証: エラーがスローされること
            await expect(tagClient.getAll()).rejects.toThrow('Failed to fetch tags')
        })

        it('ネットワークエラー時にエラーをスローする', async () => {
            // ネットワークエラーをモック化
            ;(global.fetch as any).mockRejectedValue(new Error('Network error'))

            // 実行と検証: エラーがスローされること
            await expect(tagClient.getAll()).rejects.toThrow('Network error')
        })

        it('タグが名前順にソートされて取得できる', async () => {
            // ソート済みのタグデータを準備
            const mockTags = [
                { ...mockApiTag, id: 2, name: 'Go' },
                { ...mockApiTag, id: 3, name: 'Python' },
                { ...mockApiTag, id: 1, name: 'React' },
            ]

            // API成功レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: true,
                status: 200,
                json: async () => mockTags,
            })

            // 実行
            const result = await tagClient.getAll()

            // 検証: 名前順になっていること
            expect(result[0].name).toBe('Go')
            expect(result[1].name).toBe('Python')
            expect(result[2].name).toBe('React')
        })
    })

    // 記事詳細取得のテスト
    describe('getByID', () => {
        it('指定IDのタグを取得できる', async () => {
            // API成功レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: true,
                status: 200,
                json: async () => mockApiTag,
            })

            // 実行
            const result = await tagClient.getById(1)

            // 検証: URL、データ内容
            expect(global.fetch).toHaveBeenCalledWith('http://localhost:8080/api/tags/1')
            expect(result).toEqual(expectedTag)
        })

        it('存在しないIDの場合は404エラーをスローする', async () => {
            // API404レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: false,
                status: 404,
                json: async () => ({ error: 'tag not found' }),
            })

            // 実行と検証: エラーがスローされること
            await expect(tagClient.getById(999)).rejects.toThrow()

            const error = await tagClient.getById(999).catch(e => e)
            expect(error.message).toContain('tag not found')
        })

        it('不正なIDの場合は400エラーをスローする', async () => {
            // API400レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: false,
                status: 400,
                json: async () => ({ error: 'invalid id' }),
            })

            // 実行と検証: エラーがスローされること
            const error = await tagClient.getById(-1).catch(e => e)
            expect(error).toBeDefined()
        })
    })

    // タグ作成のテスト
    describe('create', () => {
        it('新しいタグを作成できる', async () => {
            // タグデータの準備
            const createRequest = {
                name: 'Rust',
            }

            const mockCreatedTag = {
                id: 10,
                name: 'Rust',
                created_at: '2024-01-02T10:00:00Z',
                updated_at: '2024-01-02T10:00:00Z',
            }

            // API成功レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: true,
                status: 200,
                json: async () => mockCreatedTag,
            })

            // 実行
            const result = await tagClient.create(createRequest)

            // 検証: 呼び出し回数、URL、リクエスト内容、レスポンス内容
            expect(global.fetch).toHaveBeenCalledWith(
                'http://localhost:8080/api/tags',
                expect.objectContaining({
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(createRequest),
                })
            )
            expect(result.id).toBe(10)
            expect(result.name).toBe('Rust')
        })

        it('空のタグ名の場合は400エラーをスローする', async () => {
            // API400レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: false,
                status: 400,
                json: async () => ({ error: 'name is required' }),
            })

            const invalidRequest = {
                name: '',
            }

            // 実行と検証: エラーがスローされること
            const error = await tagClient.create(invalidRequest).catch(e => e)
            expect(error).toBeDefined()
            expect(error.message).toContain('name is required')
        })

        it('重複したタグ名の場合はエラーをスローする', async () => {
            // API409レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: false,
                status: 409,
                json: async () => ({ error: 'tag already exists' }),
            })

            const duplicateRequest = {
                name: 'Go',
            }

            // 実行と検証: エラーがスローされること
            const error = await tagClient.create(duplicateRequest).catch(e => e)
            expect(error).toBeDefined()
        })

        it('50文字を超えるタグ名の場合はエラーをスローする', async () => {
            // API400レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: false,
                status: 400,
                json: async () => ({ error: 'name must be 50 characters or less' }),
            })

            const longNameRequest = {
                name: 'a'.repeat(51),
            }

            // 実行と検証: エラーがスローされること
            const error = await tagClient.create(longNameRequest).catch(e => e)
            expect(error).toBeDefined()
        })
    })

    // タグ更新のテスト
    describe('update', () => {
        it('既存のタグを更新できる', async () => {
            // 更新データの準備
            const updateRequest = {
                name: 'Golang',
            }

            const mockUpdatedTag = {
                id: 1,
                name: 'Golang',
                created_at: '2024-01-01T10:00:00Z',
                updated_at: '2024-01-02T12:00:00Z',
            }

            // API成功レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: true,
                status: 200,
                json: async () => mockUpdatedTag,
            })

            // 実行
            const result = await tagClient.update(1, updateRequest)

            // 検証: 呼び出し内容、レスポンス内容
            expect(global.fetch).toHaveBeenCalledWith(
                'http://localhost:8080/api/tags/1',
                expect.objectContaining({
                    method: 'PUT',
                    headers: {
                    'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(updateRequest),
                })
            )
            expect(result.name).toBe('Golang')
            expect(result.updatedAt).toBe('2024-01-02 12:00:00')
        })

        it('存在しないIDの場合は404エラーをスローする', async () => {
            // API404レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: false,
                status: 404,
                json: async () => ({ error: 'tag not found' }),
            })

            const updateRequest = {
                name: 'Updated',
            }

            // 実行と検証: エラーがスローされること
            const error = await tagClient.update(999, updateRequest).catch(e => e)
            expect(error).toBeDefined()
            expect(error.message).toContain('tag not found')
        })

        it('重複した名前に更新しようとするとエラーをスローする', async () => {
            // API409レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: false,
                status: 409,
                json: async () => ({ error: 'tag with this name already exists' }),
            })

            const updateRequest = {
                name: 'React',
            }

            // 実行と検証: エラーがスローされること
            const error = await tagClient.update(1, updateRequest).catch(e => e)
            expect(error).toBeDefined()
        })

        it('空のタグ名に更新しようとするとエラーをスローする', async () => {
            // API400レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: false,
                status: 400,
                json: async () => ({ error: 'name is required' }),
            })

            const updateRequest = {
                name: '',
            }

            // 実行と検証: エラーがスローされること
            const error = await tagClient.update(1, updateRequest).catch(e => e)
            expect(error).toBeDefined()
        })
    })

    // タグ削除のテスト
    describe('delete', () => {
        it('指定IDのタグを削除できる', async () => {
            // API成功レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: true,
                status: 204,
            })

            // 実行
            await tagClient.delete(1)

            // 検証: 呼び出し内容
            expect(global.fetch).toHaveBeenCalledWith(
                'http://localhost:8080/api/tags/1',
                expect.objectContaining({
                    method: 'DELETE',
                })
            )
        })

        it('存在しないIDの場合は404エラーをスローする', async () => {
            // API404レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: false,
                status: 404,
                json: async () => ({ error: 'tag not found' }),
            })

            // 実行と検証: エラーがスローされること
            const error = await tagClient.delete(999).catch(e => e)
            expect(error).toBeDefined()
            expect(error.message).toContain('tag not found')
        })

        it('不正なID（0以下）の場合はエラーをスローする', async () => {
            // API400レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: false,
                status: 400,
                json: async () => ({ error: 'invalid id' }),
            })

            // 実行と検証: エラーがスローされること
            const error = await tagClient.delete(0).catch(e => e)
            expect(error).toBeDefined()
        })
    })

    // データ変換のテスト
    describe('Data Transformation', () => {
        it('snake_caseのAPIレスポンスをcamelCaseに変換できる', async () => {
            // APIレスポンスの準備
            const apiResponse = {
                id: 1,
                name: 'Go',
                created_at: '2024-01-01T10:00:00Z',
                updated_at: '2024-01-01T11:00:00Z',
            }

            // API成功レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: true,
                status: 200,
                json: async () => apiResponse,
            })

            // 実行
            const result = await tagClient.getById(1)

            // 検証: フィールド名と値の確認
            expect(result.createdAt).toBe('2024-01-01 10:00:00')
            expect(result.updatedAt).toBe('2024-01-01 11:00:00')
            expect(result).not.toHaveProperty('created_at')
            expect(result).not.toHaveProperty('updated_at')
        })

        it('複数のタグの変換も正しく行える', async () => {
            // 複数のAPIレスポンスを準備
            const apiResponses = [
                {
                    id: 1,
                    name: 'Go',
                    created_at: '2024-01-01T10:00:00Z',
                    updated_at: '2024-01-01T10:00:00Z',
                },
                {
                    id: 2,
                    name: 'React',
                    created_at: '2024-01-02T10:00:00Z',
                    updated_at: '2024-01-02T10:00:00Z',
                },
            ]

            // API成功レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: true,
                status: 200,
                json: async () => apiResponses,
            })

            // 実行
            const result = await tagClient.getAll()

            // 検証: 全てのタグが正しく変換されていること
            result.forEach((tag) => {
                expect(tag).toHaveProperty('createdAt')
                expect(tag).toHaveProperty('updatedAt')
                expect(tag).not.toHaveProperty('created_at')
                expect(tag).not.toHaveProperty('updated_at')
            })
        })
    })
})
