import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { describe, it, expect, beforeEach, vi } from 'vitest'
import ArticleForm from '@/components/ArticleForm'
import { articleClient } from '@/lib/api/articleClient'
import { tagClient } from '@/lib/api/tagClient'

// apiクライアントをモック化
vi.mock('@/lib/api/articleClient')
vi.mock('@/lib/api/tagClient')

// next.jsのuseRouterをモック化
const mockPush = vi.fn()
vi.mock('next/navigation', () => ({
    useRouter: () => ({
        push: mockPush,
    }),
}))

describe('ArticleForm', () => {
    // テスト用のタグデータ
    const mockTags = [
        { id: 1, name: 'Go', createdAt: '2024-01-01 10:00:00', updatedAt: '2024-01-01 10:00:00' },
        { id: 2, name: 'React', createdAt: '2024-01-02 10:00:00', updatedAt: '2024-01-02 10:00:00' },
        { id: 3, name: 'TypeScript', createdAt: '2024-01-03 10:00:00', updatedAt: '2024-01-03 10:00:00' },
    ]

    beforeEach(() => {
        // 各テスト前にモックをクリア
        vi.clearAllMocks()

        // タグ一覧取得のモック（デフォルト）
        vi.mocked(tagClient.getAll).mockResolvedValue(mockTags)
    })

    describe('初期表示', () => {
        it('フォームが正しく表示される', async () => {
            render(<ArticleForm />)

            // フォーム要素が存在することを確認
            expect(screen.getByLabelText(/タイトル/i)).toBeInTheDocument()
            expect(screen.getByLabelText(/URL/i)).toBeInTheDocument()
            expect(screen.getByLabelText(/要約/i)).toBeInTheDocument()
            expect(screen.getByLabelText(/メモ/i)).toBeInTheDocument()
            expect(screen.getByRole('button', { name: /登録/i })).toBeInTheDocument()
        })

        it('タグ一覧が表示される', async () => {
            render(<ArticleForm />)

            // タグが読み込まれるまで待機
            await waitFor(() => {
                expect(screen.getByText('Go')).toBeInTheDocument()
            })

            expect(screen.getByText('React')).toBeInTheDocument()
            expect(screen.getByText('TypeScript')).toBeInTheDocument()
        })

        it('初期状態では送信ボタンが無効', () => {
            render(<ArticleForm />)

            const submitButton = screen.getByRole('button', { name: /登録/i })
            expect(submitButton).toBeDisabled()
        })
    })

    describe('入力フィールドの動作', () => {
        it('タイトルを入力できる', () => {
            render(<ArticleForm />)

            const titleInput = screen.getByLabelText(/タイトル/i) as HTMLInputElement
            fireEvent.change(titleInput, { target: { value: 'Go言語入門' } })

            expect(titleInput.value).toBe('Go言語入門')
        })

        it('URLを入力できる', () => {
            render(<ArticleForm />)

            const urlInput = screen.getByLabelText(/URL/i) as HTMLInputElement
            fireEvent.change(urlInput, { target: { value: 'https://example.com' } })

            expect(urlInput.value).toBe('https://example.com')
        })

        it('要約を入力できる', () => {
            render(<ArticleForm />)

            const summaryInput = screen.getByLabelText(/要約/i) as HTMLTextAreaElement
            fireEvent.change(summaryInput, { target: { value: 'Go言語の基礎を学びます' } })

            expect(summaryInput.value).toBe('Go言語の基礎を学びます')
        })

        it('メモを入力できる', () => {
            render(<ArticleForm />)

            const memoInput = screen.getByLabelText(/メモ/i) as HTMLTextAreaElement
            fireEvent.change(memoInput, { target: { value: '後で読む' } })

            expect(memoInput.value).toBe('後で読む')
        })
    })

    describe('バリデーション', () => {
        it('タイトルが空の場合、エラーメッセージが表示される', async () => {
            render(<ArticleForm />)

            const titleInput = screen.getByLabelText(/タイトル/i)
            const urlInput = screen.getByLabelText(/URL/i)
            const summaryInput = screen.getByLabelText(/要約/i)

            // タイトルを空のまま、他を入力
            fireEvent.change(urlInput, { target: { value: 'https://example.com' } })
            fireEvent.change(summaryInput, { target: { value: 'テスト要約' } })

            // タイトルにフォーカスして離れる（blur）
            fireEvent.focus(titleInput)
            fireEvent.blur(titleInput)

            await waitFor(() => {
                expect(screen.getByText(/タイトルは必須です/i)).toBeInTheDocument()
            })
        })

        it('URLが空の場合、エラーメッセージが表示される', async () => {
            render(<ArticleForm />)

            const urlInput = screen.getByLabelText(/URL/i)

            // URLにフォーカスして離れる（blur）
            fireEvent.focus(urlInput)
            fireEvent.blur(urlInput)

            await waitFor(() => {
                expect(screen.getByText(/URLは必須です/i)).toBeInTheDocument()
            })
        })

        it('URLの形式が不正な場合、エラーメッセージが表示される', async () => {
            render(<ArticleForm />)

            const urlInput = screen.getByLabelText(/URL/i)
            fireEvent.change(urlInput, { target: { value: 'invalid-url' } })
            fireEvent.blur(urlInput)

            await waitFor(() => {
                expect(screen.getByText(/正しいURL形式で入力してください/i)).toBeInTheDocument()
            })
        })

        it('要約が空の場合、エラーメッセージが表示される', async () => {
            render(<ArticleForm />)

            const summaryInput = screen.getByLabelText(/要約/i)

            // 要約にフォーカスして離れる（blur）
            fireEvent.focus(summaryInput)
            fireEvent.blur(summaryInput)

            await waitFor(() => {
                expect(screen.getByText(/要約は必須です/i)).toBeInTheDocument()
            })
        })

        it('すべての必須項目が入力されている場合、送信ボタンが有効になる', async () => {
            render(<ArticleForm />)

            const titleInput = screen.getByLabelText(/タイトル/i)
            const urlInput = screen.getByLabelText(/URL/i)
            const summaryInput = screen.getByLabelText(/要約/i)

            fireEvent.change(titleInput, { target: { value: 'Go言語入門' } })
            fireEvent.change(urlInput, { target: { value: 'https://example.com' } })
            fireEvent.change(summaryInput, { target: { value: 'Go言語の基礎を学びます' } })

            await waitFor(() => {
                const submitButton = screen.getByRole('button', { name: /登録/i })
                expect(submitButton).not.toBeDisabled()
            })
        })
    })

    describe('タグ選択機能', () => {
        it('タグをクリックして選択できる', async () => {
            render(<ArticleForm />)

            // タグが読み込まれるまで待機
            await waitFor(() => {
                expect(screen.getByText('Go')).toBeInTheDocument()
            })

            const goTag = screen.getByText('Go')
            fireEvent.click(goTag)

            // 選択状態のスタイルが適用される（例：背景色が変わる）
            await waitFor(() => {
                expect(goTag.closest('button')).toHaveClass('bg-blue-600')
            })
        })

        it('複数のタグを選択できる', async () => {
            render(<ArticleForm />)

            await waitFor(() => {
                expect(screen.getByText('Go')).toBeInTheDocument()
            })

            const goTag = screen.getByText('Go')
            const reactTag = screen.getByText('React')

            fireEvent.click(goTag)
            fireEvent.click(reactTag)

            await waitFor(() => {
                expect(goTag.closest('button')).toHaveClass('bg-blue-600')
                expect(reactTag.closest('button')).toHaveClass('bg-blue-600')
            })
        })

        it('選択したタグを再度クリックして選択解除できる', async () => {
            render(<ArticleForm />)

            await waitFor(() => {
                expect(screen.getByText('Go')).toBeInTheDocument()
            })

            const goTag = screen.getByText('Go')

            // 選択
            fireEvent.click(goTag)
            await waitFor(() => {
                expect(goTag.closest('button')).toHaveClass('bg-blue-600')
            })

            // 選択解除
            fireEvent.click(goTag)
            await waitFor(() => {
                expect(goTag.closest('button')).not.toHaveClass('bg-blue-600')
            })
        })
    })

    describe('フォーム送信処理', () => {
        it('正常系：記事が正常に作成される', async () => {
            const mockCreatedArticle = {
                id: 1,
                title: 'Go言語入門',
                url: 'https://example.com',
                summary: 'Go言語の基礎を学びます',
                tags: ['Go'],
                memo: '後で読む',
                createdAt: '2024-01-01 10:00:00',
                updatedAt: '2024-01-01 10:00:00',
            }

            vi.mocked(articleClient.create).mockResolvedValue(mockCreatedArticle)

            render(<ArticleForm />)

            // タグが読み込まれるまで待機
            await waitFor(() => {
                expect(screen.getByText('Go')).toBeInTheDocument()
            })

            // フォーム入力
            const titleInput = screen.getByLabelText(/タイトル/i)
            const urlInput = screen.getByLabelText(/URL/i)
            const summaryInput = screen.getByLabelText(/要約/i)
            const memoInput = screen.getByLabelText(/メモ/i)

            fireEvent.change(titleInput, { target: { value: 'Go言語入門' } })
            fireEvent.change(urlInput, { target: { value: 'https://example.com' } })
            fireEvent.change(summaryInput, { target: { value: 'Go言語の基礎を学びます' } })
            fireEvent.change(memoInput, { target: { value: '後で読む' } })

            // タグを選択
            const goTag = screen.getByText('Go')
            fireEvent.click(goTag)

            // フォーム送信
            const submitButton = screen.getByRole('button', { name: /登録/i })
            fireEvent.click(submitButton)

            // API が呼ばれたことを確認
            await waitFor(() => {
                expect(articleClient.create).toHaveBeenCalledWith({
                    title: 'Go言語入門',
                    url: 'https://example.com',
                    summary: 'Go言語の基礎を学びます',
                    tags: ['Go'],
                    memo: '後で読む',
                })
            })

            // リダイレクトされることを確認
            expect(mockPush).toHaveBeenCalledWith('/articles')
        })

        it('正常系：メモとタグなしでも記事を作成できる', async () => {
            const mockCreatedArticle = {
                id: 2,
                title: 'React入門',
                url: 'https://example.com/react',
                summary: 'Reactの基礎を学びます',
                tags: [],
                memo: '',
                createdAt: '2024-01-02 10:00:00',
                updatedAt: '2024-01-02 10:00:00',
            }

            vi.mocked(articleClient.create).mockResolvedValue(mockCreatedArticle)

            render(<ArticleForm />)

            // フォーム入力（必須項目のみ）
            const titleInput = screen.getByLabelText(/タイトル/i)
            const urlInput = screen.getByLabelText(/URL/i)
            const summaryInput = screen.getByLabelText(/要約/i)

            fireEvent.change(titleInput, { target: { value: 'React入門' } })
            fireEvent.change(urlInput, { target: { value: 'https://example.com/react' } })
            fireEvent.change(summaryInput, { target: { value: 'Reactの基礎を学びます' } })

            // フォーム送信
            const submitButton = screen.getByRole('button', { name: /登録/i })
            fireEvent.click(submitButton)

            // API が呼ばれたことを確認
            await waitFor(() => {
                expect(articleClient.create).toHaveBeenCalledWith({
                    title: 'React入門',
                    url: 'https://example.com/react',
                    summary: 'Reactの基礎を学びます',
                    tags: [],
                    memo: '',
                })
            })

            // リダイレクトされることを確認
            expect(mockPush).toHaveBeenCalledWith('/articles')
        })

        it('異常系：API エラーが発生した場合、エラーメッセージが表示される', async () => {
            vi.mocked(articleClient.create).mockRejectedValue(new Error('Failed to create article'))

            render(<ArticleForm />)

            // フォーム入力
            const titleInput = screen.getByLabelText(/タイトル/i)
            const urlInput = screen.getByLabelText(/URL/i)
            const summaryInput = screen.getByLabelText(/要約/i)

            fireEvent.change(titleInput, { target: { value: 'テスト記事' } })
            fireEvent.change(urlInput, { target: { value: 'https://example.com' } })
            fireEvent.change(summaryInput, { target: { value: 'テスト要約' } })

            // フォーム送信
            const submitButton = screen.getByRole('button', { name: /登録/i })
            fireEvent.click(submitButton)

            // エラーメッセージが表示されることを確認
            await waitFor(() => {
                expect(screen.getByText(/記事の作成に失敗しました/i)).toBeInTheDocument()
            })

            // リダイレクトされないことを確認
            expect(mockPush).not.toHaveBeenCalled()
        })

        it('送信中は送信ボタンが無効になる', async () => {
            // API呼び出しを遅延させる
            vi.mocked(articleClient.create).mockImplementation(() =>
                new Promise(resolve => setTimeout(() => resolve({
                    id: 1,
                    title: 'テスト',
                    url: 'https://example.com',
                    summary: 'テスト',
                    tags: [],
                    memo: '',
                    createdAt: '2024-01-01 10:00:00',
                    updatedAt: '2024-01-01 10:00:00',
                }), 100))
            )

            render(<ArticleForm />)

            // フォーム入力
            const titleInput = screen.getByLabelText(/タイトル/i)
            const urlInput = screen.getByLabelText(/URL/i)
            const summaryInput = screen.getByLabelText(/要約/i)

            fireEvent.change(titleInput, { target: { value: 'テスト記事' } })
            fireEvent.change(urlInput, { target: { value: 'https://example.com' } })
            fireEvent.change(summaryInput, { target: { value: 'テスト要約' } })

            // フォーム送信
            const submitButton = screen.getByRole('button', { name: /登録/i })
            fireEvent.click(submitButton)

            // 送信中はボタンが無効になることを確認
            await waitFor(() => {
                expect(submitButton).toBeDisabled()
                expect(screen.getByText(/登録中.../i)).toBeInTheDocument()
            })
        })
    })

    describe('タグ読み込みエラー', () => {
        it('タグ読み込みに失敗した場合、エラーメッセージが表示される', async () => {
            vi.mocked(tagClient.getAll).mockRejectedValue(new Error('Failed to fetch tags'))

            render(<ArticleForm />)

            await waitFor(() => {
                expect(screen.getByText(/タグの読み込みに失敗しました/i)).toBeInTheDocument()
            })
        })
    })

    describe('キャンセル機能', () => {
        it('キャンセルボタンをクリックすると記事一覧に戻る', () => {
            render(<ArticleForm />)

            const cancelButton = screen.getByRole('button', { name: /キャンセル/i })
            fireEvent.click(cancelButton)

            expect(mockPush).toHaveBeenCalledWith('/articles')
        })
    })
})
