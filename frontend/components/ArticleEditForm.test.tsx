import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { describe, it, expect, beforeEach, vi } from 'vitest'
import ArticleEditForm from '@/components/ArticleEditForm'
import { articleClient } from '@/lib/api/articleClient'
import { tagClient } from '@/lib/api/tagClient'

// API クライアントをモック化
vi.mock('@/lib/api/articleClient')
vi.mock('@/lib/api/tagClient')

// Next.jsのuseRouterとuseParamsをモック化
const mockPush = vi.fn()
const mockParams = { id: '1' }
vi.mock('next/navigation', () => ({
    useRouter: () => ({
        push: mockPush,
    }),
    useParams: () => mockParams,
}))

describe('ArticleEditForm', () => {
    const mockTags = [
        { id: 1, name: 'Go', createdAt: '2024-01-01 10:00:00', updatedAt: '2024-01-01 10:00:00' },
        { id: 2, name: 'React', createdAt: '2024-01-02 10:00:00', updatedAt: '2024-01-02 10:00:00' },
        { id: 3, name: 'TypeScript', createdAt: '2024-01-03 10:00:00', updatedAt: '2024-01-03 10:00:00' },
    ]

    const mockArticle = {
        id: 1,
        title: 'Go言語入門',
        url: 'https://example.com/go',
        summary: 'Go言語の基礎を学びます',
        tags: ['Go', 'TypeScript'],
        memo: '後で読む',
        createdAt: '2024-01-01 10:00:00',
        updatedAt: '2024-01-01 10:00:00',
    }

    beforeEach(() => {
        vi.clearAllMocks()
        vi.mocked(tagClient.getAll).mockResolvedValue(mockTags)
        vi.mocked(articleClient.getById).mockResolvedValue(mockArticle)
    })

    describe('初期表示', () => {
        it('フォームが正しく表示される', async () => {
            render(<ArticleEditForm />)
            expect(screen.getByText(/読み込み中.../)).toBeInTheDocument()

            await waitFor(() => {
                expect(screen.getByLabelText(/タイトル/)).toBeInTheDocument()
            })

            expect(screen.getByLabelText(/URL/)).toBeInTheDocument()
            expect(screen.getByLabelText(/要約/)).toBeInTheDocument()
            expect(screen.getByLabelText(/メモ/)).toBeInTheDocument()
            expect(screen.getByRole('button', { name: /更新/ })).toBeInTheDocument()
        })

        it('記事IDを使って記事データを取得する', async () => {
            render(<ArticleEditForm />)
            await waitFor(() => {
                expect(articleClient.getById).toHaveBeenCalledWith(1)
            })
        })

        it('取得した記事データがフォームに表示される', async () => {
            render(<ArticleEditForm />)

            await waitFor(() => {
                const titleInput = screen.getByLabelText(/タイトル/) as HTMLInputElement
                expect(titleInput.value).toBe('Go言語入門')
            })

            const urlInput = screen.getByLabelText(/URL/) as HTMLInputElement
            const summaryInput = screen.getByLabelText(/要約/) as HTMLTextAreaElement
            const memoInput = screen.getByLabelText(/メモ/) as HTMLTextAreaElement

            expect(urlInput.value).toBe('https://example.com/go')
            expect(summaryInput.value).toBe('Go言語の基礎を学びます')
            expect(memoInput.value).toBe('後で読む')
        })

        it('取得した記事のタグが選択状態で表示される', async () => {
            render(<ArticleEditForm />)

            await waitFor(() => {
                expect(screen.getByText('Go')).toBeInTheDocument()
            })

            const goTagButton = screen.getByText('Go').closest('button')
            const reactTagButton = screen.getByText('React').closest('button')
            const tsTagButton = screen.getByText('TypeScript').closest('button')

            expect(goTagButton).toHaveClass('bg-blue-600')
            expect(tsTagButton).toHaveClass('bg-blue-600')
            expect(reactTagButton).not.toHaveClass('bg-blue-600')
        })
    })

    describe('フォーム更新処理', () => {
        it('正常系：記事が正常に更新される', async () => {
            const updatedArticle = {
                ...mockArticle,
                title: 'Go言語応用編',
            }

            vi.mocked(articleClient.update).mockResolvedValue(updatedArticle)
            render(<ArticleEditForm />)

            await waitFor(() => {
                const titleInput = screen.getByLabelText(/タイトル/) as HTMLInputElement
                expect(titleInput.value).toBe('Go言語入門')
            })

            const titleInput = screen.getByLabelText(/タイトル/)
            fireEvent.change(titleInput, { target: { value: 'Go言語応用編' } })

            const submitButton = screen.getByRole('button', { name: /更新/ })
            fireEvent.click(submitButton)

            await waitFor(() => {
                expect(articleClient.update).toHaveBeenCalledWith(1, {
                    title: 'Go言語応用編',
                    url: 'https://example.com/go',
                    summary: 'Go言語の基礎を学びます',
                    tags: ['Go', 'TypeScript'],
                    memo: '後で読む',
                })
            })

            expect(mockPush).toHaveBeenCalledWith('/articles')
        })
    })
})
