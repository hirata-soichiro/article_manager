'use client'

import { useEffect, useState } from 'react'
import { useParams, useRouter } from 'next/navigation'
import Link from 'next/link'
import { articleClient } from '@/lib/api/articleClient'
import { Article } from '@/types/article'
import DeleteConfirmDialog from '@/components/DeleteConfirmDialog'

export default function ArticleDetailPage() {
    const router = useRouter()
    const params = useParams()
    const articleId = Number(params.id)

    const [article, setArticle] = useState<Article | null>(null)
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState<string | null>(null)

    // 削除ダイアログの状態
    const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false)
    const [isDeleting, setIsDeleting] = useState(false)
    const [deleteError, setDeleteError] = useState<string | null>(null)

    useEffect(() => {
        const fetchArticle = async () => {
            try {
                setLoading(true)
                const data = await articleClient.getById(articleId)
                setArticle(data)
                setError(null)
            } catch (err) {
                setError('記事の読み込みに失敗しました')
            } finally {
                setLoading(false)
            }
        }

        if (articleId) {
            fetchArticle()
        }
    }, [articleId])

    const handleDeleteClick = () => {
        setIsDeleteDialogOpen(true)
        setDeleteError(null)
    }

    const handleDeleteConfirm = async () => {
        try {
            setIsDeleting(true)
            setDeleteError(null)
            await articleClient.delete(articleId)
            router.push('/articles')
        } catch (err) {
            setDeleteError('記事の削除に失敗しました')
            setIsDeleting(false)
        }
    }

    const handleDeleteCancel = () => {
        setIsDeleteDialogOpen(false)
        setDeleteError(null)
    }

    if (loading) {
        return (
            <div className="max-w-4xl mx-auto px-4 py-8">
                <div className="flex justify-center items-center min-h-[400px]">
                    <div className="text-gray-600">
                        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
                        <p>記事を読み込んでいます...</p>
                    </div>
                </div>
            </div>
        )
    }

    if (error || !article) {
        return (
            <div className="max-w-4xl mx-auto px-4 py-8">
                <div className="bg-red-50 border border-red-200 rounded-lg p-4 text-red-600">
                    {error || '記事が見つかりません'}
                </div>
                <Link
                    href="/articles"
                    className="mt-4 inline-block px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition"
                >
                    記事一覧に戻る
                </Link>
            </div>
        )
    }

    return (
        <div className="min-h-screen bg-gray-50 pb-12">
            <div className="max-w-5xl mx-auto px-4 py-8">
                {/* ヘッダー */}
                <div className="mb-8 flex justify-between items-center">
                    <Link
                        href="/articles"
                        className="flex items-center gap-2 text-gray-600 hover:text-blue-600 transition-colors"
                    >
                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" />
                        </svg>
                        <span className="font-medium">記事一覧に戻る</span>
                    </Link>
                    <div className="flex gap-2">
                        <Link
                            href={`/articles/${article.id}/edit`}
                            className="p-3 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-all duration-200 flex items-center gap-2"
                            aria-label="記事を編集"
                        >
                            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
                            </svg>
                            <span className="text-sm font-medium hidden sm:inline">編集</span>
                        </Link>
                        <button
                            onClick={handleDeleteClick}
                            className="p-3 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-all duration-200 flex items-center gap-2"
                            aria-label="記事を削除"
                        >
                            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                            </svg>
                            <span className="text-sm font-medium hidden sm:inline">削除</span>
                        </button>
                    </div>
                </div>

                {deleteError && (
                    <div className="mb-6 bg-red-50 border border-red-200 rounded-lg p-4 text-red-600">
                        {deleteError}
                    </div>
                )}

                {/* 記事詳細 */}
                <div className="bg-white rounded-xl shadow-lg p-8 md:p-12 border border-gray-200">
                    {/* タイトル */}
                    <h1 className="text-4xl md:text-5xl font-bold text-gray-900 mb-6 leading-tight">
                        {article.title}
                    </h1>

                    {/* URL */}
                    <a
                        href={article.url}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="inline-flex items-center gap-2 text-blue-600 hover:text-blue-700 mb-8 group"
                    >
                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                        </svg>
                        <span className="text-sm group-hover:underline break-all">{article.url}</span>
                    </a>

                    {/* タグ */}
                    {article.tags.length > 0 && (
                        <div className="mb-8">
                            <div className="flex flex-wrap gap-2">
                                {article.tags.map((tag, index) => (
                                    <span
                                        key={index}
                                        className="px-3 py-1 bg-gray-100 text-gray-700 rounded-full text-sm font-medium"
                                    >
                                        {tag}
                                    </span>
                                ))}
                            </div>
                        </div>
                    )}

                    {/* 要約 */}
                    <div className="mb-8 bg-blue-50 border-l-4 border-blue-500 p-6 rounded-r-lg">
                        <h2 className="text-sm font-bold text-blue-900 uppercase tracking-wide mb-3">要約</h2>
                        <p className="text-gray-700 leading-relaxed whitespace-pre-wrap">{article.summary}</p>
                    </div>

                    {/* メモ */}
                    {article.memo && (
                        <div className="bg-amber-50 border-l-4 border-amber-500 p-6 rounded-r-lg">
                            <h2 className="text-sm font-bold text-amber-900 uppercase tracking-wide mb-3">メモ</h2>
                            <p className="text-gray-700 leading-relaxed whitespace-pre-wrap italic">{article.memo}</p>
                        </div>
                    )}
                </div>
            </div>

            {/* 削除確認ダイアログ */}
            <DeleteConfirmDialog
                isOpen={isDeleteDialogOpen}
                title="記事を削除しますか？"
                message={`「${article.title}」を削除します。この操作は取り消せません。`}
                onConfirm={handleDeleteConfirm}
                onCancel={handleDeleteCancel}
                isDeleting={isDeleting}
            />
        </div>
    )
}
