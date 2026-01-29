'use client'

import { useEffect, useState } from 'react'
import { useParams, useRouter } from 'next/navigation'
import Link from 'next/link'
import { articleClient } from '@/lib/api/articleClient'
import { Article } from '@/types/article'
import TagList from '@/components/TagList'
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
        <div className="max-w-4xl mx-auto px-4 py-8">
            {/* ヘッダー */}
            <div className="mb-6 flex justify-between items-start">
                <div>
                    <Link
                        href="/articles"
                        className="text-blue-600 hover:underline text-sm mb-2 inline-block"
                    >
                        ← 記事一覧に戻る
                    </Link>
                </div>
                <div className="flex gap-2">
                    <Link
                        href={`/articles/${article.id}/edit`}
                        className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition"
                    >
                        編集
                    </Link>
                    <button
                        onClick={handleDeleteClick}
                        className="px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700 transition"
                    >
                        削除
                    </button>
                </div>
            </div>

            {deleteError && (
                <div className="mb-6 bg-red-50 border border-red-200 rounded-lg p-4 text-red-600">
                    {deleteError}
                </div>
            )}

            {/* 記事詳細 */}
            <div className="bg-white rounded-lg shadow-md p-8">
                <h1 className="text-3xl font-bold text-gray-800 mb-4">
                    {article.title}
                </h1>

                <div className="mb-6">
                    <a
                        href={article.url}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-blue-600 hover:underline break-all"
                    >
                        {article.url}
                    </a>
                </div>

                <div className="mb-6">
                    <h2 className="text-lg font-semibold text-gray-700 mb-2">要約</h2>
                    <p className="text-gray-600 whitespace-pre-wrap">{article.summary}</p>
                </div>

                {article.tags.length > 0 && (
                    <div className="mb-6">
                        <h2 className="text-lg font-semibold text-gray-700 mb-2">タグ</h2>
                        <TagList tags={article.tags} />
                    </div>
                )}

                {article.memo && (
                    <div className="mb-6">
                        <h2 className="text-lg font-semibold text-gray-700 mb-2">メモ</h2>
                        <p className="text-gray-600 whitespace-pre-wrap">{article.memo}</p>
                    </div>
                )}

                <div className="border-t pt-4 mt-6">
                    <div className="text-sm text-gray-500">
                        <p>作成日時: {new Date(article.createdAt).toLocaleString('ja-JP')}</p>
                        <p>更新日時: {new Date(article.updatedAt).toLocaleString('ja-JP')}</p>
                    </div>
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
