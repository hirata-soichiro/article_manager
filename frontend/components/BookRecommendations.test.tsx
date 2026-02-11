import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import BookRecommendations from './BookRecommendations'
import * as useBookRecommendationsHook from '@/hooks/useBookRecommendations'
import type { BookRecommendation } from '@/types/book'
import { ApiError } from '@/lib/errors/ApiError'

// モックデータ
const mockBooks: BookRecommendation[] = [
    {
        title: 'リーダブルコード',
        purchaseLinks: {
            amazon: 'https://www.amazon.co.jp/dp/4873115655',
            rakuten: 'https://books.rakuten.co.jp/rb/11753651/',
        },
    },
    {
        title: 'Clean Code',
        purchaseLinks: {
            amazon: 'https://www.amazon.co.jp/dp/0132350884',
        },
    },
    {
        title: 'Test Book Without Purchase Link',
        purchaseLinks: {},
    },
]

describe('BookRecommendations', () => {
    beforeEach(() => {
        vi.clearAllMocks()
    })

    afterEach(() => {
        vi.restoreAllMocks()
    })

    describe('ローディング状態', () => {
        it('ローディング中はスケルトンUIを表示する', () => {
            // useBookRecommendationsフックをモック
            vi.spyOn(useBookRecommendationsHook, 'useBookRecommendations').mockReturnValue({
                books: [],
                loading: true,
                error: null,
                cached: false,
                generatedAt: null,
                expiresAt: null,
                refetch: vi.fn(),
            })

            render(<BookRecommendations />)

            // セクションタイトルが表示されている
            expect(screen.getByText('あなたにおすすめの書籍')).toBeInTheDocument()

            // スケルトンカードが5個表示されている
            const skeletonCards = screen.getAllByRole('generic').filter((element) => {
                return element.classList.contains('animate-pulse')
            })
            expect(skeletonCards.length).toBeGreaterThan(0)

            // aria-labelが設定されている
            expect(screen.getByLabelText('書籍推薦を読み込み中')).toBeInTheDocument()
        })
    })

    describe('書籍推薦の表示', () => {
        it('書籍データを正しく表示する', async () => {
            // useBookRecommendationsフックをモック
            vi.spyOn(useBookRecommendationsHook, 'useBookRecommendations').mockReturnValue({
                books: mockBooks,
                loading: false,
                error: null,
                cached: false,
                generatedAt: '2026-02-09T10:30:00+09:00',
                expiresAt: '2026-02-10T10:30:00+09:00',
                refetch: vi.fn(),
            })

            render(<BookRecommendations />)

            await waitFor(() => {
                // セクションタイトルが表示されている
                expect(screen.getByText('あなたにおすすめの書籍')).toBeInTheDocument()

                // 各書籍のタイトルが表示されている
                expect(screen.getByText('リーダブルコード')).toBeInTheDocument()
                expect(screen.getByText('Clean Code')).toBeInTheDocument()
                expect(screen.getByText('Test Book Without Purchase Link')).toBeInTheDocument()
            })
        })

        it('Amazonの購入リンクを優先的に表示する', async () => {
            // useBookRecommendationsフックをモック
            vi.spyOn(useBookRecommendationsHook, 'useBookRecommendations').mockReturnValue({
                books: [mockBooks[0]], // リーダブルコード（Amazon + 楽天）
                loading: false,
                error: null,
                cached: false,
                generatedAt: '2026-02-09T10:30:00+09:00',
                expiresAt: '2026-02-10T10:30:00+09:00',
                refetch: vi.fn(),
            })

            render(<BookRecommendations />)

            await waitFor(() => {
                // Amazonの購入リンクが表示されている
                const amazonLink = screen.getByRole('link', {
                    name: 'リーダブルコードを購入',
                })
                expect(amazonLink).toHaveAttribute(
                    'href',
                    'https://www.amazon.co.jp/dp/4873115655'
                )
                expect(amazonLink).toHaveAttribute('target', '_blank')
                expect(amazonLink).toHaveAttribute('rel', 'noopener noreferrer')

                // ボタンのテキストが「Amazonで購入」
                expect(screen.getByText('Amazonで購入')).toBeInTheDocument()
            })
        })

        it('Amazonがない場合は楽天の購入リンクを表示する', async () => {
            // 楽天のみのモックデータ
            const rakutenOnlyBook: BookRecommendation = {
                title: '楽天限定書籍',
                purchaseLinks: {
                    rakuten: 'https://books.rakuten.co.jp/rb/12345678/',
                },
            }

            // useBookRecommendationsフックをモック
            vi.spyOn(useBookRecommendationsHook, 'useBookRecommendations').mockReturnValue({
                books: [rakutenOnlyBook],
                loading: false,
                error: null,
                cached: false,
                generatedAt: '2026-02-09T10:30:00+09:00',
                expiresAt: '2026-02-10T10:30:00+09:00',
                refetch: vi.fn(),
            })

            render(<BookRecommendations />)

            await waitFor(() => {
                // 楽天の購入リンクが表示されている
                const rakutenLink = screen.getByRole('link', {
                    name: '楽天限定書籍を購入',
                })
                expect(rakutenLink).toHaveAttribute(
                    'href',
                    'https://books.rakuten.co.jp/rb/12345678/'
                )

                // ボタンのテキストが「楽天で購入」
                expect(screen.getByText('楽天で購入')).toBeInTheDocument()
            })
        })

        it('購入リンクがない場合は無効ボタンを表示する', async () => {
            // useBookRecommendationsフックをモック
            vi.spyOn(useBookRecommendationsHook, 'useBookRecommendations').mockReturnValue({
                books: [mockBooks[2]], // 購入リンクなし
                loading: false,
                error: null,
                cached: false,
                generatedAt: '2026-02-09T10:30:00+09:00',
                expiresAt: '2026-02-10T10:30:00+09:00',
                refetch: vi.fn(),
            })

            render(<BookRecommendations />)

            await waitFor(() => {
                // 購入リンクなしのテキストが表示されている
                const disabledButton = screen.getByText('購入リンクなし')
                expect(disabledButton).toBeInTheDocument()

                // ボタンはdivとして表示される
                expect(disabledButton.tagName).toBe('DIV')

                // クラスが適用されている
                expect(disabledButton.className).toContain('cursor-not-allowed')
                expect(disabledButton.className).toContain('bg-gray-200')
            })
        })

        it('aria-labelが正しく設定されている', async () => {
            // useBookRecommendationsフックをモック
            vi.spyOn(useBookRecommendationsHook, 'useBookRecommendations').mockReturnValue({
                books: mockBooks,
                loading: false,
                error: null,
                cached: false,
                generatedAt: '2026-02-09T10:30:00+09:00',
                expiresAt: '2026-02-10T10:30:00+09:00',
                refetch: vi.fn(),
            })

            render(<BookRecommendations />)

            await waitFor(() => {
                // セクションのaria-labelが設定されている
                expect(
                    screen.getByLabelText('あなたにおすすめの書籍')
                ).toBeInTheDocument()
            })
        })
    })

    describe('エラー状態', () => {
        it('エラー時はセクション全体を表示しない（静かに失敗）', () => {
            // useBookRecommendationsフックをモック
            const apiError = new ApiError(
                'Internal Server Error',
                500,
                '/api/book-recommendations',
                'GET'
            )
            vi.spyOn(useBookRecommendationsHook, 'useBookRecommendations').mockReturnValue({
                books: [],
                loading: false,
                error: apiError,
                cached: false,
                generatedAt: null,
                expiresAt: null,
                refetch: vi.fn(),
            })

            const { container } = render(<BookRecommendations />)

            // セクション全体が表示されない
            expect(container.firstChild).toBeNull()
        })
    })

    describe('空状態', () => {
        it('書籍が0件の場合はセクション全体を表示しない', () => {
            // useBookRecommendationsフックをモック
            vi.spyOn(useBookRecommendationsHook, 'useBookRecommendations').mockReturnValue({
                books: [],
                loading: false,
                error: null,
                cached: false,
                generatedAt: null,
                expiresAt: null,
                refetch: vi.fn(),
            })

            const { container } = render(<BookRecommendations />)

            // セクション全体が表示されない
            expect(container.firstChild).toBeNull()
        })
    })

    describe('レスポンシブデザイン', () => {
        it('グリッドレイアウトが正しく適用されている', async () => {
            // useBookRecommendationsフックをモック
            vi.spyOn(useBookRecommendationsHook, 'useBookRecommendations').mockReturnValue({
                books: mockBooks,
                loading: false,
                error: null,
                cached: false,
                generatedAt: '2026-02-09T10:30:00+09:00',
                expiresAt: '2026-02-10T10:30:00+09:00',
                refetch: vi.fn(),
            })

            render(<BookRecommendations />)

            await waitFor(() => {
                // グリッドコンテナを取得
                const gridContainer = screen
                    .getByText('リーダブルコード')
                    .closest('.grid')

                // グリッドクラスが適用されている
                expect(gridContainer).toHaveClass('grid')
                expect(gridContainer).toHaveClass('grid-cols-1')
                expect(gridContainer).toHaveClass('sm:grid-cols-2')
                expect(gridContainer).toHaveClass('md:grid-cols-3')
                expect(gridContainer).toHaveClass('lg:grid-cols-5')
                expect(gridContainer).toHaveClass('gap-4')
            })
        })
    })

    describe('スタイリング', () => {
        it('書籍カードに正しいスタイルが適用されている', async () => {
            // useBookRecommendationsフックをモック
            vi.spyOn(useBookRecommendationsHook, 'useBookRecommendations').mockReturnValue({
                books: [mockBooks[0]],
                loading: false,
                error: null,
                cached: false,
                generatedAt: '2026-02-09T10:30:00+09:00',
                expiresAt: '2026-02-10T10:30:00+09:00',
                refetch: vi.fn(),
            })

            render(<BookRecommendations />)

            await waitFor(() => {
                // 書籍カードを取得
                const bookCard = screen.getByText('リーダブルコード').closest('div')

                // カードスタイルが適用されている
                expect(bookCard).toHaveClass('bg-white')
                expect(bookCard).toHaveClass('rounded-lg')
                expect(bookCard).toHaveClass('shadow-sm')
                expect(bookCard).toHaveClass('border')
                expect(bookCard).toHaveClass('border-gray-200')
            })
        })

        it('購入ボタンに正しいスタイルが適用されている', async () => {
            // useBookRecommendationsフックをモック
            vi.spyOn(useBookRecommendationsHook, 'useBookRecommendations').mockReturnValue({
                books: [mockBooks[0]],
                loading: false,
                error: null,
                cached: false,
                generatedAt: '2026-02-09T10:30:00+09:00',
                expiresAt: '2026-02-10T10:30:00+09:00',
                refetch: vi.fn(),
            })

            render(<BookRecommendations />)

            await waitFor(() => {
                // 購入ボタンを取得
                const purchaseButton = screen.getByRole('link', {
                    name: 'リーダブルコードを購入',
                })

                // ボタンスタイルが適用されている
                expect(purchaseButton).toHaveClass('bg-gradient-to-r')
                expect(purchaseButton).toHaveClass('from-indigo-600')
                expect(purchaseButton).toHaveClass('to-purple-600')
                expect(purchaseButton).toHaveClass('text-white')
                expect(purchaseButton).toHaveClass('rounded-lg')
            })
        })

        it('セクションタイトルに正しいスタイルが適用されている', async () => {
            // useBookRecommendationsフックをモック
            vi.spyOn(useBookRecommendationsHook, 'useBookRecommendations').mockReturnValue({
                books: mockBooks,
                loading: false,
                error: null,
                cached: false,
                generatedAt: '2026-02-09T10:30:00+09:00',
                expiresAt: '2026-02-10T10:30:00+09:00',
                refetch: vi.fn(),
            })

            render(<BookRecommendations />)

            await waitFor(() => {
                // タイトルを取得
                const title = screen.getByText('あなたにおすすめの書籍')

                // タイトルスタイルが適用されている
                expect(title).toHaveClass('text-lg')
                expect(title).toHaveClass('font-semibold')
                expect(title).toHaveClass('text-gray-800')
            })
        })
    })

    describe('タイトルの行数制限', () => {
        it('タイトルが2行に制限されている', async () => {
            // useBookRecommendationsフックをモック
            vi.spyOn(useBookRecommendationsHook, 'useBookRecommendations').mockReturnValue({
                books: [mockBooks[0]],
                loading: false,
                error: null,
                cached: false,
                generatedAt: '2026-02-09T10:30:00+09:00',
                expiresAt: '2026-02-10T10:30:00+09:00',
                refetch: vi.fn(),
            })

            render(<BookRecommendations />)

            await waitFor(() => {
                // タイトル要素を取得
                const titleElement = screen.getByText('リーダブルコード')

                // line-clamp-2が適用されている
                expect(titleElement).toHaveClass('line-clamp-2')
            })
        })
    })
})