import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { renderHook, waitFor } from '@testing-library/react'
import {
    useBookRecommendations,
    __resetBookRecommendationsCache,
} from './useBookRecommendations'
import { bookRecommendationClient } from '@/lib/api/bookRecommendationClient'
import type { BookRecommendationsResponse } from '@/types/book'
import { ApiError } from '@/lib/errors/ApiError'

// モックデータ
const mockBooksResponse: BookRecommendationsResponse = {
    books: [
        {
            title: 'リーダブルコード',
            author: 'Dustin Boswell, Trevor Foucher',
            isbn: '9784873115658',
            purchaseLinks: {
                amazon: 'https://www.amazon.co.jp/dp/4873115655',
                rakuten: 'https://books.rakuten.co.jp/rb/11753651/',
            },
        },
        {
            title: 'Clean Code',
            author: 'Robert C. Martin',
            isbn: '9780132350884',
            purchaseLinks: {
                amazon: 'https://www.amazon.co.jp/dp/0132350884',
            },
        },
    ],
    cached: false,
    generatedAt: '2026-02-09T10:30:00+09:00',
    expiresAt: '2026-02-10T10:30:00+09:00',
}

const mockEmptyBooksResponse: BookRecommendationsResponse = {
    books: [],
    cached: false,
    generatedAt: null,
    expiresAt: null,
}

describe('useBookRecommendations', () => {
    beforeEach(() => {
        // 各テストの前にキャッシュをリセット
        __resetBookRecommendationsCache()
        vi.clearAllMocks()
    })

    afterEach(() => {
        vi.restoreAllMocks()
    })

    describe('初期状態', () => {
        it('ローディング状態から始まる', () => {
            // APIをモック（永遠に解決しないPromise）
            vi.spyOn(bookRecommendationClient, 'getRecommendations').mockImplementation(
                () => new Promise(() => {})
            )

            const { result } = renderHook(() => useBookRecommendations())

            expect(result.current.loading).toBe(true)
            expect(result.current.books).toEqual([])
            expect(result.current.error).toBeNull()
            expect(result.current.cached).toBe(false)
        })
    })

    describe('API呼び出し成功', () => {
        it('書籍推薦データを正しく取得する', async () => {
            // APIをモック
            vi.spyOn(bookRecommendationClient, 'getRecommendations').mockResolvedValue(
                mockBooksResponse
            )

            const { result } = renderHook(() => useBookRecommendations())

            // 初期状態はローディング中
            expect(result.current.loading).toBe(true)

            // データ取得完了まで待機
            await waitFor(() => {
                expect(result.current.loading).toBe(false)
            })

            // データが正しく設定されている
            expect(result.current.books).toEqual(mockBooksResponse.books)
            expect(result.current.error).toBeNull()
            expect(result.current.cached).toBe(false)
            expect(result.current.generatedAt).toBe('2026-02-09T10:30:00+09:00')
            expect(result.current.expiresAt).toBe('2026-02-10T10:30:00+09:00')
        })

        it('記事が0件の場合、空配列を返す', async () => {
            // APIをモック
            vi.spyOn(bookRecommendationClient, 'getRecommendations').mockResolvedValue(
                mockEmptyBooksResponse
            )

            const { result } = renderHook(() => useBookRecommendations())

            // データ取得完了まで待機
            await waitFor(() => {
                expect(result.current.loading).toBe(false)
            })

            // 空配列が返る
            expect(result.current.books).toEqual([])
            expect(result.current.error).toBeNull()
            expect(result.current.generatedAt).toBeNull()
            expect(result.current.expiresAt).toBeNull()
        })
    })

    describe('API呼び出し失敗', () => {
        it('エラー時に空配列を返し、エラーを記録する', async () => {
            // コンソールエラーをモック
            const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {})

            // APIエラーをモック
            const apiError = new ApiError(
                'Internal Server Error',
                500,
                '/api/book-recommendations',
                'GET'
            )
            vi.spyOn(bookRecommendationClient, 'getRecommendations').mockRejectedValue(
                apiError
            )

            const { result } = renderHook(() => useBookRecommendations())

            // データ取得完了まで待機
            await waitFor(() => {
                expect(result.current.loading).toBe(false)
            })

            // エラー時は空配列
            expect(result.current.books).toEqual([])
            expect(result.current.error).toBeInstanceOf(ApiError)
            expect(result.current.error?.message).toBe('Internal Server Error')
            expect(result.current.generatedAt).toBeNull()
            expect(result.current.expiresAt).toBeNull()

            // コンソールエラーが呼ばれたことを確認
            expect(consoleErrorSpy).toHaveBeenCalledWith(
                'Failed to fetch book recommendations:',
                apiError
            )

            consoleErrorSpy.mockRestore()
        })

        it('一般的なエラーも正しく処理する', async () => {
            // コンソールエラーをモック
            const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {})

            // 一般的なエラーをモック
            const genericError = new Error('Network error')
            vi.spyOn(bookRecommendationClient, 'getRecommendations').mockRejectedValue(
                genericError
            )

            const { result } = renderHook(() => useBookRecommendations())

            // データ取得完了まで待機
            await waitFor(() => {
                expect(result.current.loading).toBe(false)
            })

            // エラー時は空配列
            expect(result.current.books).toEqual([])
            expect(result.current.error).toBeInstanceOf(Error)
            expect(result.current.error?.message).toBe('Network error')

            consoleErrorSpy.mockRestore()
        })
    })

    describe('キャッシュ機能', () => {
        it('キャッシュが有効な場合、APIを再呼び出ししない', async () => {
            // 1回目のAPI呼び出しをモック
            const getRecommendationsSpy = vi
                .spyOn(bookRecommendationClient, 'getRecommendations')
                .mockResolvedValue(mockBooksResponse)

            // 1回目のレンダリング
            const { result: result1, unmount: unmount1 } = renderHook(() =>
                useBookRecommendations()
            )

            await waitFor(() => {
                expect(result1.current.loading).toBe(false)
            })

            expect(getRecommendationsSpy).toHaveBeenCalledTimes(1)
            expect(result1.current.books).toEqual(mockBooksResponse.books)

            unmount1()

            // 2回目のレンダリング（キャッシュが有効）
            const { result: result2 } = renderHook(() => useBookRecommendations())

            await waitFor(() => {
                expect(result2.current.loading).toBe(false)
            })

            // API呼び出しは1回のまま（キャッシュが使われた）
            expect(getRecommendationsSpy).toHaveBeenCalledTimes(1)
            expect(result2.current.books).toEqual(mockBooksResponse.books)
        })

        it('キャッシュが期限切れの場合、APIを再呼び出しする', async () => {
            // 1回目のAPI呼び出しをモック
            const getRecommendationsSpy = vi
                .spyOn(bookRecommendationClient, 'getRecommendations')
                .mockResolvedValue(mockBooksResponse)

            // 1回目のレンダリング
            const { result: result1, unmount: unmount1 } = renderHook(() =>
                useBookRecommendations()
            )

            await waitFor(() => {
                expect(result1.current.loading).toBe(false)
            })

            expect(getRecommendationsSpy).toHaveBeenCalledTimes(1)

            unmount1()

            // キャッシュをリセット（期限切れをシミュレート）
            __resetBookRecommendationsCache()

            // 2回目のレンダリング（キャッシュが期限切れ）
            const { result: result2 } = renderHook(() => useBookRecommendations())

            await waitFor(() => {
                expect(result2.current.loading).toBe(false)
            })

            // API呼び出しが再度実行された
            expect(getRecommendationsSpy).toHaveBeenCalledTimes(2)
            expect(result2.current.books).toEqual(mockBooksResponse.books)
        })
    })

    describe('refetch機能', () => {
        it('refetchを呼び出すとキャッシュを無視して再取得する', async () => {
            // 1回目のAPI呼び出しをモック
            const getRecommendationsSpy = vi
                .spyOn(bookRecommendationClient, 'getRecommendations')
                .mockResolvedValue(mockBooksResponse)

            const { result } = renderHook(() => useBookRecommendations())

            // 初回データ取得完了まで待機
            await waitFor(() => {
                expect(result.current.loading).toBe(false)
            })

            expect(getRecommendationsSpy).toHaveBeenCalledTimes(1)

            // refetchを呼び出す
            await result.current.refetch()

            // データ取得完了まで待機
            await waitFor(() => {
                expect(result.current.loading).toBe(false)
            })

            // API呼び出しが再度実行された（キャッシュを無視）
            expect(getRecommendationsSpy).toHaveBeenCalledTimes(2)
            expect(result.current.books).toEqual(mockBooksResponse.books)
        })

        it('refetch中もloadingがtrueになる', async () => {
            // API呼び出しをモック（遅延を追加）
            let resolveFirst: ((value: BookRecommendationsResponse) => void) | undefined
            let resolveSecond: ((value: BookRecommendationsResponse) => void) | undefined
            let callCount = 0

            const getRecommendationsSpy = vi
                .spyOn(bookRecommendationClient, 'getRecommendations')
                .mockImplementation(() => {
                    callCount++
                    if (callCount === 1) {
                        return new Promise((resolve) => {
                            resolveFirst = resolve
                        })
                    } else {
                        return new Promise((resolve) => {
                            resolveSecond = resolve
                        })
                    }
                })

            const { result } = renderHook(() => useBookRecommendations())

            // 初回ローディング中
            expect(result.current.loading).toBe(true)

            // 初回データを解決
            resolveFirst!(mockBooksResponse)

            // 初回データ取得完了まで待機
            await waitFor(() => {
                expect(result.current.loading).toBe(false)
            })

            expect(getRecommendationsSpy).toHaveBeenCalledTimes(1)

            // refetchを呼び出す
            result.current.refetch()

            // ローディング状態になる
            await waitFor(() => {
                expect(result.current.loading).toBe(true)
            })

            // refetchデータを解決
            resolveSecond!(mockBooksResponse)

            // データ取得完了
            await waitFor(() => {
                expect(result.current.loading).toBe(false)
            })

            expect(getRecommendationsSpy).toHaveBeenCalledTimes(2)
        })
    })

    describe('重複リクエスト防止', () => {
        it('fetchingRef により同一インスタンス内での重複リクエストが防止される', async () => {
            // API呼び出しをモック（遅延を追加）
            let resolvePromise: ((value: BookRecommendationsResponse) => void) | undefined
            const getRecommendationsSpy = vi
                .spyOn(bookRecommendationClient, 'getRecommendations')
                .mockImplementation(
                    () =>
                        new Promise((resolve) => {
                            resolvePromise = resolve
                        })
                )

            const { result } = renderHook(() => useBookRecommendations())

            // 初回ローディング中
            expect(result.current.loading).toBe(true)

            // データを解決
            resolvePromise!(mockBooksResponse)

            // データ取得完了まで待機
            await waitFor(() => {
                expect(result.current.loading).toBe(false)
            })

            // API呼び出しは1回
            expect(getRecommendationsSpy).toHaveBeenCalledTimes(1)
            expect(result.current.books).toEqual(mockBooksResponse.books)
        })

        it('複数回マウントされた場合、キャッシュにより追加のAPI呼び出しが防止される', async () => {
            // API呼び出しをモック
            const getRecommendationsSpy = vi
                .spyOn(bookRecommendationClient, 'getRecommendations')
                .mockResolvedValue(mockBooksResponse)

            // 1回目のレンダリング
            const { result: result1, unmount: unmount1 } = renderHook(() =>
                useBookRecommendations()
            )

            await waitFor(() => {
                expect(result1.current.loading).toBe(false)
            })

            expect(getRecommendationsSpy).toHaveBeenCalledTimes(1)

            unmount1()

            // 2回目のレンダリング（キャッシュが使われる）
            const { result: result2 } = renderHook(() => useBookRecommendations())

            await waitFor(() => {
                expect(result2.current.loading).toBe(false)
            })

            // API呼び出しは1回のまま（キャッシュが使われた）
            expect(getRecommendationsSpy).toHaveBeenCalledTimes(1)
            expect(result2.current.books).toEqual(mockBooksResponse.books)
        })
    })
})