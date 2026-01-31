import { renderHook, waitFor, act } from '@testing-library/react'
import { describe, it, expect, beforeEach, vi } from 'vitest'
import { useArticleSearch } from './useArticleSearch'
import { articleClient } from '@/lib/api/articleClient'
import type { Article } from '@/types/article'

// articleClientをモック化
vi.mock('@/lib/api/articleClient')

// モック関数の型定義
const mockSearchArticles = articleClient.searchArticles as ReturnType<typeof vi.fn>

describe('useArticleSearch', () => {
    // テスト用データ
    const mockSearchResults: Article[] = [
        {
            id: 1,
            title: 'React Hooks Guide',
            url: 'https://example.com/react-hooks',
            summary: 'A comprehensive guide to React Hooks',
            tags: ['react', 'hooks', 'javascript'],
            memo: 'Great resource for learning hooks',
            createdAt: '2024-01-01 10:00:00',
            updatedAt: '2024-01-01 10:00:00',
        },
        {
            id: 2,
            title: 'Advanced React Patterns',
            url: 'https://example.com/react-patterns',
            summary: 'Learn advanced React patterns and best practices',
            tags: ['react', 'patterns', 'advanced'],
            memo: 'Contains useful patterns',
            createdAt: '2024-01-02 10:00:00',
            updatedAt: '2024-01-02 10:00:00',
        },
    ]

    // 各テスト前にモックをクリア
    beforeEach(() => {
        vi.clearAllMocks()
    })

    describe('初期状態', () => {
        it('初期状態ではloading=false, results=[], keyword=""であること', () => {
            const { result } = renderHook(() => useArticleSearch())

            // 初期状態を確認
            expect(result.current.loading).toBe(false)
            expect(result.current.results).toEqual([])
            expect(result.current.keyword).toBe('')
            expect(result.current.error).toBeNull()
        })

        it('検索が実行されていない初期状態ではAPIは呼ばれないこと', () => {
            renderHook(() => useArticleSearch())

            // 初期状態ではAPIは呼ばれない
            expect(mockSearchArticles).not.toHaveBeenCalled()
        })
    })

    describe('検索実行', () => {
        it('検索実行中はloading=trueになること', async () => {
            // 検索API呼び出しを遅延させるPromiseを使用
            let resolveSearch: (value: Article[]) => void
            const searchPromise = new Promise<Article[]>((resolve) => {
                resolveSearch = resolve
            })
            mockSearchArticles.mockReturnValue(searchPromise)

            const { result } = renderHook(() => useArticleSearch())

            // 検索を実行
            act(() => {
                result.current.search('react')
            })

            // loading状態を確認
            expect(result.current.loading).toBe(true)
            expect(result.current.keyword).toBe('react')

            // 検索を完了
            act(() => {
                resolveSearch!(mockSearchResults)
            })

            await waitFor(() => {
                expect(result.current.loading).toBe(false)
            })
        })

        it('検索成功時、resultsに検索結果が格納されること', async () => {
            mockSearchArticles.mockResolvedValue(mockSearchResults)

            const { result } = renderHook(() => useArticleSearch())

            // 検索を実行
            await act(async () => {
                await result.current.search('react')
            })

            // 検索結果が格納されているか確認
            await waitFor(() => {
                expect(result.current.results).toEqual(mockSearchResults)
            })
            expect(result.current.keyword).toBe('react')
            expect(result.current.loading).toBe(false)
            expect(result.current.error).toBeNull()
            expect(mockSearchArticles).toHaveBeenCalledWith('react')
            expect(mockSearchArticles).toHaveBeenCalledTimes(1)
        })

        it('検索失敗時、errorにエラー情報が格納されること', async () => {
            const mockError = new Error('Failed to search articles')
            mockSearchArticles.mockRejectedValue(mockError)

            const { result } = renderHook(() => useArticleSearch())

            // 検索を実行
            await act(async () => {
                await result.current.search('react')
            })

            // エラーが格納されているか確認
            await waitFor(() => {
                expect(result.current.error).toBe(mockError)
            })
            expect(result.current.loading).toBe(false)
            expect(result.current.results).toEqual([])
            expect(mockSearchArticles).toHaveBeenCalledWith('react')
        })

        it('空文字列で検索を実行した場合、APIが呼ばれること', async () => {
            mockSearchArticles.mockResolvedValue([])

            const { result } = renderHook(() => useArticleSearch())

            // 空文字列で検索を実行
            await act(async () => {
                await result.current.search('')
            })

            await waitFor(() => {
                expect(result.current.loading).toBe(false)
            })

            expect(mockSearchArticles).toHaveBeenCalledWith('')
            expect(result.current.results).toEqual([])
            expect(result.current.keyword).toBe('')
        })

        it('検索結果が0件の場合、resultsが空配列になること', async () => {
            mockSearchArticles.mockResolvedValue([])

            const { result } = renderHook(() => useArticleSearch())

            await act(async () => {
                await result.current.search('nonexistent')
            })

            await waitFor(() => {
                expect(result.current.results).toEqual([])
            })
            expect(result.current.keyword).toBe('nonexistent')
            expect(result.current.loading).toBe(false)
            expect(result.current.error).toBeNull()
        })

        it('連続して検索を実行した場合、最新の検索結果が反映されること', async () => {
            const firstResults: Article[] = [mockSearchResults[0]]
            const secondResults: Article[] = [mockSearchResults[1]]

            mockSearchArticles
                .mockResolvedValueOnce(firstResults)
                .mockResolvedValueOnce(secondResults)

            const { result } = renderHook(() => useArticleSearch())

            // 1回目の検索
            await act(async () => {
                await result.current.search('hooks')
            })

            await waitFor(() => {
                expect(result.current.results).toEqual(firstResults)
            })
            expect(result.current.keyword).toBe('hooks')

            // 2回目の検索
            await act(async () => {
                await result.current.search('patterns')
            })

            await waitFor(() => {
                expect(result.current.results).toEqual(secondResults)
            })
            expect(result.current.keyword).toBe('patterns')
            expect(mockSearchArticles).toHaveBeenCalledTimes(2)
        })

        it('検索エラー後に再検索した場合、エラーがクリアされること', async () => {
            const mockError = new Error('Failed to search')
            mockSearchArticles
                .mockRejectedValueOnce(mockError)
                .mockResolvedValueOnce(mockSearchResults)

            const { result } = renderHook(() => useArticleSearch())

            // 1回目の検索（エラー）
            await act(async () => {
                await result.current.search('error')
            })

            await waitFor(() => {
                expect(result.current.error).toBe(mockError)
            })

            // 2回目の検索（成功）
            await act(async () => {
                await result.current.search('success')
            })

            await waitFor(() => {
                expect(result.current.error).toBeNull()
            })
            expect(result.current.results).toEqual(mockSearchResults)
        })
    })

    describe('検索クリア', () => {
        it('clearSearch呼び出し時、検索結果とキーワードがクリアされること', async () => {
            mockSearchArticles.mockResolvedValue(mockSearchResults)

            const { result } = renderHook(() => useArticleSearch())

            // 検索を実行
            await act(async () => {
                await result.current.search('react')
            })

            await waitFor(() => {
                expect(result.current.results).toEqual(mockSearchResults)
            })
            expect(result.current.keyword).toBe('react')

            // 検索をクリア
            act(() => {
                result.current.clearSearch()
            })

            // クリアされているか確認
            expect(result.current.results).toEqual([])
            expect(result.current.keyword).toBe('')
            expect(result.current.error).toBeNull()
            expect(result.current.loading).toBe(false)
        })

        it('検索実行前にclearSearchを呼んでも問題なく動作すること', () => {
            const { result } = renderHook(() => useArticleSearch())

            // 初期状態でクリアを実行
            act(() => {
                result.current.clearSearch()
            })

            // 状態が維持されているか確認
            expect(result.current.results).toEqual([])
            expect(result.current.keyword).toBe('')
            expect(result.current.error).toBeNull()
            expect(result.current.loading).toBe(false)
        })

        it('エラー状態でclearSearchを呼んだ場合、エラーもクリアされること', async () => {
            const mockError = new Error('Search failed')
            mockSearchArticles.mockRejectedValue(mockError)

            const { result } = renderHook(() => useArticleSearch())

            // エラーが発生する検索を実行
            await act(async () => {
                await result.current.search('error')
            })

            await waitFor(() => {
                expect(result.current.error).toBe(mockError)
            })

            // クリアを実行
            act(() => {
                result.current.clearSearch()
            })

            // エラーもクリアされているか確認
            expect(result.current.error).toBeNull()
            expect(result.current.results).toEqual([])
            expect(result.current.keyword).toBe('')
        })
    })

    describe('エッジケース', () => {
        it('検索実行中にclearSearchを呼んだ場合、適切に処理されること', async () => {
            // 検索を遅延させる
            let resolveSearch: (value: Article[]) => void
            const searchPromise = new Promise<Article[]>((resolve) => {
                resolveSearch = resolve
            })
            mockSearchArticles.mockReturnValue(searchPromise)

            const { result } = renderHook(() => useArticleSearch())

            // 検索を開始
            act(() => {
                result.current.search('react')
            })

            expect(result.current.loading).toBe(true)

            // 検索中にクリアを実行
            act(() => {
                result.current.clearSearch()
            })

            // クリアが実行されているか確認
            expect(result.current.results).toEqual([])
            expect(result.current.keyword).toBe('')
            expect(result.current.loading).toBe(false)

            // 検索を完了（クリア後なので状態は変わらないはず）
            act(() => {
                resolveSearch!(mockSearchResults)
            })

            // 状態がクリアされたまま維持されているか確認
            await waitFor(() => {
                expect(result.current.results).toEqual([])
            })
        })

        it('特殊文字を含むキーワードで検索できること', async () => {
            const specialKeyword = 'React & Vue.js #framework'
            mockSearchArticles.mockResolvedValue(mockSearchResults)

            const { result } = renderHook(() => useArticleSearch())

            await act(async () => {
                await result.current.search(specialKeyword)
            })

            await waitFor(() => {
                expect(result.current.keyword).toBe(specialKeyword)
            })
            expect(mockSearchArticles).toHaveBeenCalledWith(specialKeyword)
        })

        it('非常に長いキーワードで検索できること', async () => {
            const longKeyword = 'a'.repeat(1000)
            mockSearchArticles.mockResolvedValue([])

            const { result } = renderHook(() => useArticleSearch())

            await act(async () => {
                await result.current.search(longKeyword)
            })

            await waitFor(() => {
                expect(result.current.keyword).toBe(longKeyword)
            })
            expect(mockSearchArticles).toHaveBeenCalledWith(longKeyword)
        })
    })

    describe('型の整合性', () => {
        it('返却される型が正しいこと', () => {
            const { result } = renderHook(() => useArticleSearch())

            // 戻り値の型を確認
            expect(typeof result.current.loading).toBe('boolean')
            expect(Array.isArray(result.current.results)).toBe(true)
            expect(typeof result.current.keyword).toBe('string')
            expect(typeof result.current.search).toBe('function')
            expect(typeof result.current.clearSearch).toBe('function')
            expect(result.current.error === null || result.current.error instanceof Error).toBe(true)
        })
    })
})
