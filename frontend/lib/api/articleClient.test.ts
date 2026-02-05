import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { articleClient } from '@/lib/api/articleClient'
import { Article } from '@/types/article'

describe('ArticleClient', () => {
    // 各テストの前にfetchをモック化し、実際のHTTPリクエストを防ぐ
    beforeEach(() => {
        global.fetch = vi.fn()
    })

    // 各テストの後にすべてのモックを復元し、テスト間の影響を排除
    afterEach(() => {
        vi.restoreAllMocks()
    })

    // テスト用のダミーデータ
    const mockApiArticle = {
        id: 1,
        title: 'Test Article',
        url: 'https://example.com/test',
        summary: 'This is a test article',
        tags: ['test', 'example'],
        memo: 'Test memo',
        created_at: '2024-01-01 10:00:00',
        updated_at: '2024-01-01 10:00:00',
    }

    // フロントエンドで期待される形式
    const expectedArticle: Article = {
        id: 1,
        title: 'Test Article',
        url: 'https://example.com/test',
        summary: 'This is a test article',
        tags: ['test', 'example'],
        memo: 'Test memo',
        createdAt: '2024-01-01 10:00:00',
        updatedAt: '2024-01-01 10:00:00',
    }

    // 記事一覧取得のテスト
    describe('getAll', () => {
        it('全記事を取得できる', async () => {
            // 2件の記事データを準備
            const mockArticles = [mockApiArticle,{...mockApiArticle, id: 2, title: 'Article 2'}]

            // API成功レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: true,
                status: 200,
                json: async () => mockArticles,
            })

            // 実行
            const result = await articleClient.getAll()

            // 検証: 呼び出し回数、URL、データ件数、データ内容
            expect(global.fetch).toHaveBeenCalledTimes(1)
            expect(global.fetch).toHaveBeenCalledWith('http://localhost:8080/api/articles')
            expect(result).toHaveLength(2)
            expect(result[0]).toEqual(expectedArticle)
            expect(result[1].id).toBe(2)
            expect(result[1].title).toBe('Article 2')
        })

        it('空の配列が返ってきた場合も正常に処理できる', async () => {
            // 空配列のAPIレスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: true,
                status: 200,
                json: async () => [],
            })

            // 実行
            const result = await articleClient.getAll()

            // 検証: データ内容
            expect(result).toEqual([])
        })

        it('APIエラー時にエラーをスローする', async () => {
            // APIエラーレスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: false,
                status: 500,
                json: async () => ({ error: 'Failed to fetch articles'}),
            })

            // 実行と検証: エラーがスローされること
            await expect(articleClient.getAll()).rejects.toThrow('Failed to fetch articles')
        })

        it('ネットワークエラー時にエラーをスローする', async () => {
            // ネットワークエラーをモック化
            ;(global.fetch as any).mockRejectedValue(new Error('Network error'))

            // 実行と検証: エラーがスローされること
            await expect(articleClient.getAll()).rejects.toThrow('Network error')
        })
    })

    // 記事詳細取得のテスト
    describe('getByID', () => {
        it('指定IDの記事を取得できる', async () => {
            // API成功レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: true,
                status: 200,
                json: async () => mockApiArticle,
            })

            // 実行
            const result = await articleClient.getById(1)

            // 検証: URL、データ内容
            expect(global.fetch).toHaveBeenCalledWith('http://localhost:8080/api/articles/1')
            expect(result).toEqual(expectedArticle)
        })

        it('存在しないIDの場合は404エラーをスローする', async () => {
            // API404レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: false,
                status: 404,
                json: async () => ({ error: 'article not found' }),
            })

            // 実行と検証: エラーがスローされること
            await expect(articleClient.getById(999)).rejects.toThrow()

            const error = await articleClient.getById(999).catch(e => e)
            expect(error.message).toContain('article not found')
        })

        it('不正なIDの場合は400エラーをスローする', async () => {
            // API400レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: false,
                status: 400,
                json: async () => ({ error: 'invalid id' }),
            })

            // 実行と検証: エラーがスローされること
            const error = await articleClient.getById(-1).catch(e => e)
            expect(error).toBeDefined()
        })
    })

    // 記事作成のテスト
    describe('create', () => {
        it('新しい記事を作成できる', async () => {
            // 記事データの準備
            const createRequest = {
                title: 'New Article',
                url: 'https://example.com/new',
                summary: 'New article summary',
                tags: ['new', 'test'],
                memo: 'New memo',
            }

            const mockCreatedArticle = {
                ...mockApiArticle,
                ...createRequest,
                id: 10,
            }

            // API成功レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: true,
                status: 200,
                json: async () => mockCreatedArticle,
            })

            // 実行
            const result = await articleClient.create(createRequest)

            // 検証: 呼び出し回数、URL、リクエスト内容、レスポンス内容
            expect(global.fetch).toHaveBeenCalledWith(
                'http://localhost:8080/api/articles',
                expect.objectContaining({
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(createRequest),
                })
            )
            expect(result.id).toBe(10)
            expect(result.title).toBe('New Article')
        })

        it('必須フィールドが不足している場合は400エラーをスローする', async () => {
            // API400レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: false,
                status: 400,
                json: async () => ({ error: 'title is required' }),
            })

            const invalidRequest = {
                title: '',
                url: 'https://example.com',
                summary: 'summary',
                tags: [],
            }

            // 実行と検証: エラーがスローされること
            const error = await articleClient.create(invalidRequest).catch(e => e)
            expect(error).toBeDefined()
        })

        it('memoを省略して記事を作成できる', async () => {
            // 記事データの準備
            const createRequest = {
                title: 'No Memo Article',
                url: 'https://example.com/no-memo',
                summary: 'Article without memo',
                tags: ['test'],
            }

            // API成功レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: true,
                status: 200,
                json: async () => ({
                    ...mockApiArticle,
                    ...createRequest,
                    memo: '',
                }),
            })

            // 実行
            const result = await articleClient.create(createRequest)

            // 検証: memoが空文字であること
            expect(result.memo).toBe('')
        })
    })

    // 記事更新のテスト
    describe('update', () => {
        it('既存の記事を更新できる', async () => {
            // 更新データの準備
            const updateRequest = {
                title: 'Updated Title',
                url: 'https://example.com/updated',
                summary: 'Updated summary',
                tags: ['updated'],
                memo: 'Updated memo',
            }

            const mockUpdatedArticle = {
                ...mockApiArticle,
                ...updateRequest,
                updated_at: '2024-01-02 12:00:00',
            }

            // API成功レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: true,
                status: 200,
                json: async () => mockUpdatedArticle,
            })

            // 実行
            const result = await articleClient.update(1, updateRequest)

            // 検証: 呼び出し内容、レスポンス内容
            expect(global.fetch).toHaveBeenCalledWith(
                'http://localhost:8080/api/articles/1',
                expect.objectContaining({
                    method: 'PUT',
                    headers: {
                    'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(updateRequest),
                })
            )
            expect(result.title).toBe('Updated Title')
            expect(result.updatedAt).toBe('2024-01-02 12:00:00')
        })

        it('存在しないIDの場合は404エラーをスローする', async () => {
            // API404レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: false,
                status: 404,
                json: async () => ({ error: 'article not found' }),
            })

            const updateRequest = {
                title: 'Updated',
                url: 'https://example.com',
                summary: 'summary',
                tags: [],
            }

            // 実行と検証: エラーがスローされること
            const error = await articleClient.update(999, updateRequest).catch(e => e)
            expect(error).toBeDefined()
        })
    })

    // 記事削除のテスト
    describe('delete', () => {
        it('指定IDの記事を削除できる', async () => {
            // API成功レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: true,
                status: 204,
            })

            // 実行
            await articleClient.delete(1)

            // 検証: 呼び出し内容
            expect(global.fetch).toHaveBeenCalledWith(
                'http://localhost:8080/api/articles/1',
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
                json: async () => ({ error: 'article not found' }),
            })

            // 実行と検証: エラーがスローされること
            const error = await articleClient.delete(999).catch(e => e)
            expect(error).toBeDefined()
        })
    })

    // データ変換のテスト
    describe('Data Transformation', () => {
        it('snake_caseのAPIレスポンスをcamelCaseに変換できる', async () => {
            // APIレスポンスの準備
            const apiResponse = {
                id: 1,
                title: 'Test',
                url: 'https://example.com',
                summary: 'Summary',
                tags: ['tag1'],
                memo: 'memo',
                created_at: '2024-01-01 10:00:00',
                updated_at: '2024-01-01 11:00:00',
            }

            // API成功レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: true,
                status: 200,
                json: async () => apiResponse,
            })

            // 実行
            const result = await articleClient.getById(1)

            // 検証: フィールド名と値の確認
            expect(result.createdAt).toBe('2024-01-01 10:00:00')
            expect(result.updatedAt).toBe('2024-01-01 11:00:00')
            expect(result).not.toHaveProperty('created_at')
            expect(result).not.toHaveProperty('updated_at')
        })

        it('tagsがnullの場合は空配列に変換できる', async () => {
            // APIレスポンスの準備
            const apiResponse = {
                ...mockApiArticle,
                tags: null,
            }

            // API成功レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: true,
                status: 200,
                json: async () => apiResponse,
            })

            // 実行
            const result = await articleClient.getById(1)

            // 検証: tagsが空配列であること
            expect(result.tags).toEqual([])
        })

        it('memoがnullの場合は空文字列に変換できる', async () => {
            // APIレスポンスの準備
            const apiResponse = {
                ...mockApiArticle,
                memo: null,
            }

            // API成功レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: true,
                status: 200,
                json: async () => apiResponse,
            })

            // 実行
            const result = await articleClient.getById(1)

            // 検証: memoが空文字列であること
            expect(result.memo).toBe('')
        })
    })

    // AI自動生成APIのテスト
    describe('generate', () => {
        it('URLから記事を生成できる', async () => {
            // 生成リクエストデータの準備
            const generateRequest = {
                url: 'https://example.com/article',
                memo: 'テストメモ',
            }

            // 生成された記事データをモック化
            const mockGeneratedArticle = {
                id: 100,
                title: 'AI Generated Article',
                url: 'https://example.com/article',
                summary: 'This article was generated by AI from the provided URL',
                tags: ['ai', 'generated', 'test'],
                memo: 'テストメモ',
                created_at: '2024-01-15 14:30:00',
                updated_at: '2024-01-15 14:30:00',
            }

            // API成功レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: true,
                status: 200,
                json: async () => mockGeneratedArticle,
            })

            // 実行
            const result = await articleClient.generate(
                generateRequest.url,
                generateRequest.memo
            )

            // 検証: 呼び出し回数、URL、リクエスト内容、レスポンス内容
            expect(global.fetch).toHaveBeenCalledTimes(1)
            expect(global.fetch).toHaveBeenCalledWith(
                'http://localhost:8080/api/articles/generate',
                expect.objectContaining({
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(generateRequest),
                })
            )
            expect(result.id).toBe(100)
            expect(result.title).toBe('AI Generated Article')
            expect(result.url).toBe('https://example.com/article')
            expect(result.summary).toContain('AI')
            expect(result.tags).toEqual(['ai', 'generated', 'test'])
            expect(result.memo).toBe('テストメモ')
            expect(result.createdAt).toBe('2024-01-15 14:30:00')
            expect(result.updatedAt).toBe('2024-01-15 14:30:00')
        })

        it('memoなしで記事を生成できる', async () => {
            // memoなしのリクエストデータ
            const generateRequest = {
                url: 'https://example.com/article',
            }

            // 生成された記事データをモック化（memoなし）
            const mockGeneratedArticle = {
                id: 101,
                title: 'AI Generated Article Without Memo',
                url: 'https://example.com/article',
                summary: 'Generated summary',
                tags: ['ai', 'test'],
                memo: '',
                created_at: '2024-01-15 15:00:00',
                updated_at: '2024-01-15 15:00:00',
            }

            // API成功レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: true,
                status: 200,
                json: async () => mockGeneratedArticle,
            })

            // 実行（memoパラメータなし）
            const result = await articleClient.generate(generateRequest.url)

            // 検証: リクエスト内容とレスポンス
            expect(global.fetch).toHaveBeenCalledWith(
                'http://localhost:8080/api/articles/generate',
                expect.objectContaining({
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(generateRequest),
                })
            )
            expect(result.id).toBe(101)
            expect(result.memo).toBe('')
        })

        it('URLが空の場合はエラーをスローする', async () => {
            // API400レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: false,
                status: 400,
                json: async () => ({ error: 'url is required' }),
            })

            // 実行と検証: エラーがスローされること
            const error = await articleClient.generate('').catch(e => e)
            expect(error).toBeDefined()
            expect(error.message).toContain('url is required')
        })

        it('API制限エラー（429）の場合は適切なエラーメッセージをスローする', async () => {
            // API429レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: false,
                status: 429,
                json: async () => ({
                    error: 'API rate limit exceeded. Please try again later.'
                }),
            })

            // 実行と検証: エラーがスローされること
            const error = await articleClient.generate('https://example.com/article').catch(e => e)
            expect(error).toBeDefined()
            expect(error.message).toContain('rate limit')
        })

        it('タイムアウトエラー（504）の場合は適切なエラーメッセージをスローする', async () => {
            // API504レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: false,
                status: 504,
                json: async () => ({
                    error: 'Gateway timeout. The article generation took too long.'
                }),
            })

            // 実行と検証: エラーがスローされること
            const error = await articleClient.generate('https://example.com/article').catch(e => e)
            expect(error).toBeDefined()
            expect(error.message).toContain('timeout')
        })

        it('ネットワークエラー時にエラーをスローする', async () => {
            // ネットワークエラーをモック化
            ;(global.fetch as any).mockRejectedValue(new Error('Network connection failed'))

            // 実行と検証: エラーがスローされること
            await expect(
                articleClient.generate('https://example.com/article', 'メモ')
            ).rejects.toThrow('Network connection failed')
        })

        it('無効なURLフォーマットの場合はエラーをスローする', async () => {
            // API400レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: false,
                status: 400,
                json: async () => ({ error: 'invalid URL format' }),
            })

            // 実行と検証: エラーがスローされること
            const error = await articleClient.generate('not-a-valid-url').catch(e => e)
            expect(error).toBeDefined()
            expect(error.message).toContain('invalid URL')
        })
    })

    // 記事検索のテスト
    describe('searchArticles', () => {
        it('キーワードで記事を検索できる', async () => {
            // 検索結果データの準備
            const mockSearchResults = [
                mockApiArticle,
                { ...mockApiArticle, id: 2, title: 'Another Test Article' },
            ]

            // API成功レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: true,
                status: 200,
                json: async () => mockSearchResults,
            })

            // 実行
            const result = await articleClient.searchArticles('test')

            // 検証: 呼び出し回数、URL（クエリパラメータ付き）、データ件数、データ内容
            expect(global.fetch).toHaveBeenCalledTimes(1)
            expect(global.fetch).toHaveBeenCalledWith(
                'http://localhost:8080/api/articles/search?keyword=test'
            )
            expect(result).toHaveLength(2)
            expect(result[0]).toEqual(expectedArticle)
            expect(result[1].id).toBe(2)
            expect(result[1].title).toBe('Another Test Article')
        })

        it('キーワードに日本語が含まれる場合でも検索できる', async () => {
            // 日本語キーワードでの検索結果データの準備
            const mockSearchResults = [
                { ...mockApiArticle, title: 'テスト記事', summary: 'これはテスト記事です' },
            ]

            // API成功レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: true,
                status: 200,
                json: async () => mockSearchResults,
            })

            // 実行
            const result = await articleClient.searchArticles('テスト')

            // 検証: URLエンコードされたクエリパラメータ
            expect(global.fetch).toHaveBeenCalledWith(
                'http://localhost:8080/api/articles/search?keyword=%E3%83%86%E3%82%B9%E3%83%88'
            )
            expect(result).toHaveLength(1)
            expect(result[0].title).toBe('テスト記事')
        })

        it('キーワードにスペースが含まれる場合でも検索できる', async () => {
            // スペースを含むキーワードでの検索結果データの準備
            const mockSearchResults = [mockApiArticle]

            // API成功レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: true,
                status: 200,
                json: async () => mockSearchResults,
            })

            // 実行
            const result = await articleClient.searchArticles('test article')

            // 検証: スペースがエンコードされること
            expect(global.fetch).toHaveBeenCalledWith(
                'http://localhost:8080/api/articles/search?keyword=test%20article'
            )
            expect(result).toHaveLength(1)
        })

        it('検索結果が0件の場合は空配列を返す', async () => {
            // 空配列のAPIレスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: true,
                status: 200,
                json: async () => [],
            })

            // 実行
            const result = await articleClient.searchArticles('nonexistent')

            // 検証: 空配列が返ること
            expect(result).toEqual([])
            expect(result).toHaveLength(0)
        })

        it('空文字列で検索した場合は400エラーをスローする', async () => {
            // API400レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: false,
                status: 400,
                json: async () => ({ error: 'keyword is required' }),
            })

            // 実行と検証: エラーがスローされること
            const error = await articleClient.searchArticles('').catch(e => e)
            expect(error).toBeDefined()
            expect(error.message).toContain('keyword is required')
        })

        it('特殊文字を含むキーワードで検索できる', async () => {
            // 特殊文字を含むキーワードでの検索結果データの準備
            const mockSearchResults = [mockApiArticle]

            // API成功レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: true,
                status: 200,
                json: async () => mockSearchResults,
            })

            // 実行: 特殊文字を含むキーワード
            const result = await articleClient.searchArticles('test&example')

            // 検証: 特殊文字がエンコードされること
            expect(global.fetch).toHaveBeenCalledWith(
                'http://localhost:8080/api/articles/search?keyword=test%26example'
            )
            expect(result).toHaveLength(1)
        })

        it('検索結果のデータ変換が正しく行われる', async () => {
            // snake_caseのAPIレスポンスを準備
            const apiResponse = [
                {
                    id: 1,
                    title: 'Search Result',
                    url: 'https://example.com/search',
                    summary: 'This is a search result',
                    tags: ['search', 'test'],
                    memo: 'Search memo',
                    created_at: '2024-01-01 10:00:00',
                    updated_at: '2024-01-01 11:00:00',
                },
            ]

            // API成功レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: true,
                status: 200,
                json: async () => apiResponse,
            })

            // 実行
            const result = await articleClient.searchArticles('search')

            // 検証: camelCaseに変換されていること
            expect(result[0].createdAt).toBe('2024-01-01 10:00:00')
            expect(result[0].updatedAt).toBe('2024-01-01 11:00:00')
            expect(result[0]).not.toHaveProperty('created_at')
            expect(result[0]).not.toHaveProperty('updated_at')
        })

        it('検索結果のtagsがnullの場合は空配列に変換される', async () => {
            // tagsがnullのAPIレスポンスを準備
            const apiResponse = [
                {
                    ...mockApiArticle,
                    tags: null,
                },
            ]

            // API成功レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: true,
                status: 200,
                json: async () => apiResponse,
            })

            // 実行
            const result = await articleClient.searchArticles('test')

            // 検証: tagsが空配列であること
            expect(result[0].tags).toEqual([])
        })

        it('検索結果のmemoがnullの場合は空文字列に変換される', async () => {
            // memoがnullのAPIレスポンスを準備
            const apiResponse = [
                {
                    ...mockApiArticle,
                    memo: null,
                },
            ]

            // API成功レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: true,
                status: 200,
                json: async () => apiResponse,
            })

            // 実行
            const result = await articleClient.searchArticles('test')

            // 検証: memoが空文字列であること
            expect(result[0].memo).toBe('')
        })

        it('APIエラー時にエラーをスローする', async () => {
            // APIエラーレスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: false,
                status: 500,
                json: async () => ({ error: 'Internal Server Error' }),
            })

            // 実行と検証: エラーがスローされること
            await expect(articleClient.searchArticles('test')).rejects.toThrow(
                'Internal Server Error'
            )
        })

        it('ネットワークエラー時にエラーをスローする', async () => {
            // ネットワークエラーをモック化
            ;(global.fetch as any).mockRejectedValue(new Error('Network error'))

            // 実行と検証: エラーがスローされること
            await expect(articleClient.searchArticles('test')).rejects.toThrow('Network error')
        })

        it('非常に長いキーワードでも検索できる', async () => {
            // 長いキーワードでの検索結果データの準備
            const longKeyword = 'a'.repeat(1000)
            const mockSearchResults = [mockApiArticle]

            // API成功レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: true,
                status: 200,
                json: async () => mockSearchResults,
            })

            // 実行
            const result = await articleClient.searchArticles(longKeyword)

            // 検証: 正常に処理されること
            expect(global.fetch).toHaveBeenCalledTimes(1)
            expect(result).toHaveLength(1)
        })

        it('検索結果が多数の記事を含む場合でも正しく処理できる', async () => {
            // 多数の検索結果データの準備（100件）
            const mockSearchResults = Array.from({ length: 100 }, (_, i) => ({
                ...mockApiArticle,
                id: i + 1,
                title: `Article ${i + 1}`,
            }))

            // API成功レスポンスをモック化
            ;(global.fetch as any).mockResolvedValue({
                ok: true,
                status: 200,
                json: async () => mockSearchResults,
            })

            // 実行
            const result = await articleClient.searchArticles('article')

            // 検証: 全件が正しく変換されること
            expect(result).toHaveLength(100)
            expect(result[0].id).toBe(1)
            expect(result[99].id).toBe(100)
            expect(result[0].title).toBe('Article 1')
            expect(result[99].title).toBe('Article 100')
        })
    })
})
